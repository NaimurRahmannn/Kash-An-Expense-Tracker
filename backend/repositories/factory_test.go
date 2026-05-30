package repositories

import (
	"path/filepath"
	"runtime"
	"testing"

	"backend/config"
	csvrepo "backend/repositories/csv"
	postgresrepo "backend/repositories/postgres"
	"backend/storage"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath := filepath.Dir(filepath.Dir(file))
	beego.TestBeegoInit(appPath)
}

func TestNewRepositoriesReturnCSVImplementationsForCSVDriver(t *testing.T) {
	useStorageDriverConfig(t, "csv")

	userRepository := NewUserRepository()
	if _, ok := userRepository.(*csvrepo.UserRepository); !ok {
		t.Fatalf("expected CSV user repository, got %T", userRepository)
	}

	expenseRepository := NewExpenseRepository()
	if _, ok := expenseRepository.(*csvrepo.ExpenseRepository); !ok {
		t.Fatalf("expected CSV expense repository, got %T", expenseRepository)
	}
}

func TestNewRepositoriesReturnPostgresImplementationsForPostgresDriver(t *testing.T) {
	useStorageDriverConfig(t, "postgres")
	if err := storage.Close(); err != nil {
		t.Fatalf("expected storage close to succeed: %v", err)
	}
	if storage.PostgresDB() != nil {
		t.Fatal("expected Postgres DB to be nil when storage init was not called")
	}

	userRepository := NewUserRepository()
	if _, ok := userRepository.(*postgresrepo.UserRepository); !ok {
		t.Fatalf("expected Postgres user repository, got %T", userRepository)
	}

	expenseRepository := NewExpenseRepository()
	if _, ok := expenseRepository.(*postgresrepo.ExpenseRepository); !ok {
		t.Fatalf("expected Postgres expense repository, got %T", expenseRepository)
	}
}

func TestNewRepositoriesReturnCSVImplementationsForUnknownDriver(t *testing.T) {
	useFactoryStorageDriver(t, config.StorageDriver("sqlite"))

	userRepository := NewUserRepository()
	if _, ok := userRepository.(*csvrepo.UserRepository); !ok {
		t.Fatalf("expected fallback CSV user repository, got %T", userRepository)
	}

	expenseRepository := NewExpenseRepository()
	if _, ok := expenseRepository.(*csvrepo.ExpenseRepository); !ok {
		t.Fatalf("expected fallback CSV expense repository, got %T", expenseRepository)
	}
}

func useStorageDriverConfig(t *testing.T, driver string) {
	t.Helper()

	t.Setenv("STORAGE_DRIVER", "")

	previousDriver := beego.AppConfig.DefaultString("storage_driver", "csv")
	if err := beego.AppConfig.Set("storage_driver", driver); err != nil {
		t.Fatalf("expected storage_driver test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = storage.Close()
		_ = beego.AppConfig.Set("storage_driver", previousDriver)
	})
}

func useFactoryStorageDriver(t *testing.T, driver config.StorageDriver) {
	t.Helper()

	previousGetStorageDriver := getStorageDriver
	getStorageDriver = func() config.StorageDriver {
		return driver
	}

	t.Cleanup(func() {
		getStorageDriver = previousGetStorageDriver
		_ = storage.Close()
	})
}

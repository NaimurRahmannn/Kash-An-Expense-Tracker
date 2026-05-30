package repositories

import (
	"path/filepath"
	"runtime"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath := filepath.Dir(filepath.Dir(file))
	beego.TestBeegoInit(appPath)
}

func TestNewRepositoriesReturnNonNilForConfiguredDrivers(t *testing.T) {
	tests := []struct {
		name   string
		driver string
	}{
		{
			name:   "csv driver",
			driver: "csv",
		},
		{
			name:   "postgres driver falls back to csv adapter",
			driver: "postgres",
		},
		{
			name:   "unknown driver falls back to csv adapter",
			driver: "sqlite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useStorageDriverConfig(t, tt.driver)

			if repository := NewUserRepository(); repository == nil {
				t.Fatal("expected user repository to be non-nil")
			}
			if repository := NewExpenseRepository(); repository == nil {
				t.Fatal("expected expense repository to be non-nil")
			}
		})
	}
}

func useStorageDriverConfig(t *testing.T, driver string) {
	t.Helper()

	previousDriver := beego.AppConfig.DefaultString("storage_driver", "csv")
	if err := beego.AppConfig.Set("storage_driver", driver); err != nil {
		t.Fatalf("expected storage_driver test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("storage_driver", previousDriver)
	})
}

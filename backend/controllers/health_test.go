package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "backend/routers"

	beego "github.com/beego/beego/v2/server/web"
)

type healthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath := filepath.Dir(filepath.Dir(file))
	beego.TestBeegoInit(appPath)
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Fatal("expected success to be true")
	}

	if response.Message != "Server is running" {
		t.Fatalf("expected message %q, got %q", "Server is running", response.Message)
	}
}

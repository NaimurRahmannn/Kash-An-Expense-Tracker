package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestBaseControllerErrorResponse(t *testing.T) {
	beego.Router("/test/error-response", &testErrorController{})

	req := httptest.NewRequest(http.MethodGet, "/test/error-response", nil)
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Fatal("expected success to be false")
	}

	if response.Message != "Invalid request" {
		t.Fatalf("expected message %q, got %q", "Invalid request", response.Message)
	}
}

type testErrorController struct {
	BaseController
}

func (c *testErrorController) Get() {
	c.ErrorResponse(http.StatusBadRequest, "Invalid request")
}

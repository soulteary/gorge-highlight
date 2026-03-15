package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/soulteary/gorge-highlight/internal/highlight"

	"github.com/labstack/echo/v4"
)

func newTestDeps() *Deps {
	return &Deps{
		Highlighter: highlight.New(),
		Token:       "test-token",
		MaxBytes:    1048576,
	}
}

func TestHealthPing(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Error("expected ok in response")
	}
}

func TestRenderUnauthorized(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	body := `{"source":"print('hi')","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/highlight/render", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestRenderSuccess(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	body := `{"source":"x = 1","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/highlight/render", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", "test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "html") {
		t.Error("expected html in response")
	}
}

func TestRenderEmptySource(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	body := `{"source":"","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/highlight/render", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", "test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRenderTooLarge(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	deps.MaxBytes = 10
	RegisterRoutes(e, deps)

	body := `{"source":"this is a very long source string","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/highlight/render", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", "test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", rec.Code)
	}
}

func TestListLanguages(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	req := httptest.NewRequest(http.MethodGet, "/api/highlight/languages", nil)
	req.Header.Set("X-Service-Token", "test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "python") {
		t.Error("expected python in languages list")
	}
}

func TestTokenViaQueryParam(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	RegisterRoutes(e, deps)

	req := httptest.NewRequest(http.MethodGet, "/api/highlight/languages?token=test-token", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestNoTokenRequired(t *testing.T) {
	e := echo.New()
	deps := newTestDeps()
	deps.Token = ""
	RegisterRoutes(e, deps)

	body := `{"source":"x = 1","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/highlight/render", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

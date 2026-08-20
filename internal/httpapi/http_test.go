package httpapi_test

import (
	"emergency-claim-code/internal/httpapi"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/service"
	"emergency-claim-code/internal/store"
	"emergency-claim-code/internal/workflow"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "h.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	svc := service.New(repo)
	handler := httpapi.New(workflow.New(svc, repo), svc).Handler()
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("code %d", rec.Code)
	}
}

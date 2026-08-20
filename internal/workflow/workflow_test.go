package workflow_test

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/service"
	"emergency-claim-code/internal/store"
	"emergency-claim-code/internal/workflow"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "w.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	svc := service.New(repo)
	engine := workflow.New(svc, repo)
	got, err := engine.CreateReviewArchive(model.Record{ID: "r", BatchID: "b", Applicant: "A", Quantity: 1}, "reviewer", "fixed")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "archived" {
		t.Fatalf("got %s", got.Status)
	}
}

func TestWorkflowStages(t *testing.T) {
	if workflow.NextStage("register") != "review" {
		t.Fatal("stage")
	}
	if workflow.IsStage("bad") {
		t.Fatal("bad stage")
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	svc := service.New(repo)
	engine := workflow.New(svc, repo)
	item, err := engine.Register(model.Record{ID: "search-1", BatchID: "RB2972-01", Applicant: "Lin", Quantity: 1}, "clerk", "a")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Submit(item.ID, "clerk", "b"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Approve(item.ID, "reviewer", "c"); err != nil {
		t.Fatal(err)
	}
	found, code, err := engine.SearchUpdatePublish("RB2972-01", "approved", item.ID, "published", "reviewer", "d")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || code.State != "issued" {
		t.Fatal("publish workflow")
	}
}

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "i.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	svc := service.New(repo)
	engine := workflow.New(svc, repo)
	report, err := engine.ImportReport(model.Record{ID: "import-1", BatchID: "RB2972-01", Applicant: "Q", Quantity: 3}, "clerk", "fixed")
	if err != nil {
		t.Fatal(err)
	}
	if report == "" {
		t.Fatal("empty report")
	}
}

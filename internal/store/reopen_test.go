package store_test

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/store"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "claims.db")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	repo := repository.New(db)
	if err := repo.SaveRecord(model.Record{ID: "r1", BatchID: "RB2972-01", Applicant: "A", Quantity: 2, Status: "draft"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveAudit(model.AuditEvent{ID: "a1", RecordID: "r1", Action: "create"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveWorkflow(model.Workflow{ID: "w1", BatchID: "RB2972-01", Stage: "register"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveAttachment(model.Attachment{ID: "at1", RecordID: "r1", Name: "proof", Digest: "d"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveClaimCode(model.ClaimCode{Code: "c1", RecordID: "r1", BatchID: "RB2972-01", State: "issued"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo = repository.New(db)
	if _, err := repo.GetRecord("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.ListAudits("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetWorkflow("w1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.ListAttachments("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetClaimCode("c1"); err != nil {
		t.Fatal(err)
	}
}

func TestStoreMissingKey(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Get([]byte("records"), []byte("none")); err != store.ErrNotFound {
		t.Fatalf("got %v", err)
	}
}

package repository_test

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/store"
	"path/filepath"
	"testing"
)

func TestRepositoryRoundTrip(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "r.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	item := model.Record{ID: "r2", BatchID: "b", Applicant: "B", Quantity: 1, Status: "draft"}
	if err := repo.SaveRecord(item); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetRecord("r2")
	if err != nil {
		t.Fatal(err)
	}
	if got.Applicant != "B" {
		t.Fatal("wrong applicant")
	}
}

func TestRepositoryFiltering(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "r.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	_ = repo.SaveRecord(model.Record{ID: "r1", BatchID: "b", Applicant: "Alice", Quantity: 1, Status: "draft"})
	_ = repo.SaveRecord(model.Record{ID: "r2", BatchID: "c", Applicant: "Bob", Quantity: 1, Status: "approved"})
	items, err := repo.ListRecords()
	if err != nil || len(repository.FilterRecords(items, "alice", "")) != 1 {
		t.Fatal("filter failed")
	}
}

package service_test

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/service"
	"emergency-claim-code/internal/store"
	"path/filepath"
	"testing"
)

func newService(t *testing.T) *service.Service {
	db, err := store.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return service.New(repository.New(db))
}

func TestBusiness01Regression(t *testing.T) {
	svc := newService(t)
	item, err := svc.CreateRecord(model.Record{ID: "2972-01", BatchID: "RB2972-01", Applicant: "Lee", Quantity: 2}, "checker", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Submit(item.ID, "checker", "t2"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Approve(item.ID, "checker", "t3"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.RecallAndRetry(item.ID, "checker", "t4")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "submitted" {
		t.Fatalf("want submitted, got %s", got.Status)
	}
}

func TestServiceLifecycle(t *testing.T) {
	svc := newService(t)
	item, err := svc.CreateRecord(model.Record{ID: "r", BatchID: "b", Applicant: "A", Quantity: 1}, "a", "1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Submit(item.ID, "a", "2"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Approve(item.ID, "a", "3"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetRecord(item.ID)
	if err != nil || got.Status != "approved" {
		t.Fatal("not approved")
	}
}

func TestServiceCode(t *testing.T) {
	svc := newService(t)
	item, _ := svc.CreateRecord(model.Record{ID: "r", BatchID: "b", Applicant: "A", Quantity: 1}, "a", "1")
	_, _ = svc.Submit(item.ID, "a", "2")
	_, _ = svc.Approve(item.ID, "a", "3")
	code, err := svc.IssueCode(item.ID, "4")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RedeemCode(code.Code, "a", "5"); err != nil {
		t.Fatal(err)
	}
}

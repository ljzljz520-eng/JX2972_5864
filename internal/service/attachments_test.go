package service_test

import (
	"emergency-claim-code/internal/model"
	"testing"
)

func TestAttachmentLifecycle(t *testing.T) {
	svc := newService(t)
	item, err := svc.CreateRecord(model.Record{ID: "r", BatchID: "b", Applicant: "A", Quantity: 1}, "a", "1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddAttachment(model.Attachment{RecordID: item.ID, Name: "photo", Digest: "sha"}, "a", "2"); err != nil {
		t.Fatal(err)
	}
	ok, err := svc.HasAttachment(item.ID, "sha")
	if err != nil || !ok {
		t.Fatal("attachment missing")
	}
}

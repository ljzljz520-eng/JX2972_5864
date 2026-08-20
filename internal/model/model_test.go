package model_test

import (
	"emergency-claim-code/internal/model"
	"testing"
)

func TestRecordValidation(t *testing.T) {
	errs := model.ValidateRecord(model.Record{})
	if len(errs) != 4 {
		t.Fatalf("got %d", len(errs))
	}
}

func TestTransitions(t *testing.T) {
	if !model.CanTransition("draft", "submitted") {
		t.Fatal("draft should submit")
	}
	if model.CanTransition("archived", "submitted") {
		t.Fatal("archived should stop")
	}
	if !model.IsTerminal("archived") {
		t.Fatal("archived terminal")
	}
}

func TestEncoding(t *testing.T) {
	original := model.Record{ID: "r", Status: "draft"}
	data, err := model.Encode(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded model.Record
	if err := model.Decode(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != original.ID {
		t.Fatal("round trip")
	}
}

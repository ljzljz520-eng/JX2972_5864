package query_test

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/query"
	"testing"
)

func TestQueryApply(t *testing.T) {
	items := []model.Record{{ID: "2", BatchID: "b", Applicant: "Bob", Quantity: 4, Status: "approved"}, {ID: "1", BatchID: "a", Applicant: "Alice", Quantity: 2, Status: "draft"}}
	got := query.Apply(items, query.NewFilter("alice", "", "", 0))
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatal("query")
	}
}

func TestQueryReport(t *testing.T) {
	report := query.BuildReport([]model.Record{{Status: "draft", BatchID: "b"}, {Status: "draft", BatchID: "b"}})
	if report.Total != 2 || report.ByBatch["b"] != 2 {
		t.Fatal("report")
	}
}

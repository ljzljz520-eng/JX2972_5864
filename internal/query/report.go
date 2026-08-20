package query

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type Report struct {
	Total    int
	ByStatus map[string]int
	ByBatch  map[string]int
}

func BuildReport(items []model.Record) Report {
	report := Report{Total: len(items), ByStatus: make(map[string]int), ByBatch: make(map[string]int)}
	for _, item := range items {
		report.ByStatus[item.Status]++
		report.ByBatch[item.BatchID]++
	}
	return report
}

func (r Report) String() string {
	statuses := make([]string, 0, len(r.ByStatus))
	for key, value := range r.ByStatus {
		statuses = append(statuses, fmt.Sprintf("%s=%d", key, value))
	}
	return fmt.Sprintf("total=%d statuses=%s", r.Total, strings.Join(statuses, ","))
}

func SelectLatest(items []model.Record) (model.Record, bool) {
	if len(items) == 0 {
		return model.Record{}, false
	}
	latest := items[0]
	for _, item := range items[1:] {
		if item.Revision > latest.Revision {
			latest = item
		}
	}
	return latest, true
}

func EligibleForIssue(item model.Record) bool { return item.Status == "approved" && item.Quantity > 0 }

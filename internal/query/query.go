package query

import (
	"emergency-claim-code/internal/model"
	"sort"
	"strings"
)

type Filter struct {
	Term        string
	Status      string
	Batch       string
	MinQuantity int
}

func NewFilter(term, status, batch string, min int) Filter {
	return Filter{Term: strings.TrimSpace(term), Status: strings.TrimSpace(status), Batch: strings.TrimSpace(batch), MinQuantity: min}
}

func Matches(item model.Record, filter Filter) bool {
	if filter.Status != "" && item.Status != filter.Status {
		return false
	}
	if filter.Batch != "" && item.BatchID != filter.Batch {
		return false
	}
	if filter.MinQuantity > 0 && item.Quantity < filter.MinQuantity {
		return false
	}
	if filter.Term == "" {
		return true
	}
	term := strings.ToLower(filter.Term)
	return strings.Contains(strings.ToLower(item.ID), term) || strings.Contains(strings.ToLower(item.Applicant), term) || strings.Contains(strings.ToLower(item.BatchID), term)
}

func Apply(items []model.Record, filter Filter) []model.Record {
	result := make([]model.Record, 0, len(items))
	for _, item := range items {
		if Matches(item, filter) {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].BatchID == result[j].BatchID {
			return result[i].ID < result[j].ID
		}
		return result[i].BatchID < result[j].BatchID
	})
	return result
}

func GroupByBatch(items []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, item := range items {
		groups[item.BatchID] = append(groups[item.BatchID], item)
	}
	return groups
}

func StatusCounts(items []model.Record) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return counts
}

func RenderTable(items []model.Record) string {
	var b strings.Builder
	b.WriteString("ID|BATCH|APPLICANT|QTY|STATUS\n")
	for _, item := range items {
		b.WriteString(item.ID + "|" + item.BatchID + "|" + item.Applicant + "|")
		b.WriteString(strings.TrimSpace(string(rune('0' + item.Quantity%10))))
		b.WriteString("|" + item.Status + "\n")
	}
	return b.String()
}

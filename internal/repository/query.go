package repository

import (
	"emergency-claim-code/internal/model"
	"strings"
)

func MatchRecord(record model.Record, term string) bool {
	if term == "" {
		return true
	}
	needle := strings.ToLower(strings.TrimSpace(term))
	return strings.Contains(strings.ToLower(record.ID), needle) || strings.Contains(strings.ToLower(record.BatchID), needle) || strings.Contains(strings.ToLower(record.Applicant), needle) || strings.Contains(strings.ToLower(record.Status), needle)
}

func FilterRecords(records []model.Record, term, status string) []model.Record {
	result := make([]model.Record, 0, len(records))
	for _, record := range records {
		if MatchRecord(record, term) && (status == "" || record.Status == status) {
			result = append(result, record)
		}
	}
	return result
}

func SortRecords(records []model.Record) []model.Record {
	for i := 1; i < len(records); i++ {
		current := records[i]
		j := i - 1
		for j >= 0 && records[j].ID > current.ID {
			records[j+1] = records[j]
			j--
		}
		records[j+1] = current
	}
	return records
}

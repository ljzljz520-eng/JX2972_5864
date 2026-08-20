package service

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/policy"
	"errors"
	"fmt"
	"strings"
)

type IntakeResult struct {
	Accepted []model.Record
	Rejected []string
	Findings map[string][]policy.Finding
}

func (s *Service) ValidateInput(record model.Record) error {
	findings := policy.CheckRecord(record)
	if policy.HasErrors(findings) {
		return errors.New(policy.Summarize(findings))
	}
	return nil
}

func (s *Service) PrepareRecord(record model.Record) (model.Record, []policy.Finding) {
	record = model.NormalizeRecord(record)
	return record, policy.CheckRecord(record)
}

func (s *Service) RegisterBatch(records []model.Record, actor, at string) IntakeResult {
	result := IntakeResult{Accepted: make([]model.Record, 0), Rejected: make([]string, 0), Findings: make(map[string][]policy.Finding)}
	seen := make(map[string]bool)
	for index, input := range records {
		if input.ID == "" {
			input.ID = fmt.Sprintf("batch-%04d", index+1)
		}
		if seen[input.ID] {
			result.Rejected = append(result.Rejected, input.ID+": duplicate in batch")
			continue
		}
		seen[input.ID] = true
		prepared, findings := s.PrepareRecord(input)
		result.Findings[prepared.ID] = findings
		if policy.HasErrors(findings) {
			result.Rejected = append(result.Rejected, prepared.ID+": "+policy.Summarize(findings))
			continue
		}
		item, err := s.CreateRecord(prepared, actor, at)
		if err != nil {
			result.Rejected = append(result.Rejected, prepared.ID+": "+err.Error())
			continue
		}
		result.Accepted = append(result.Accepted, item)
	}
	return result
}

func (s *Service) ParseAndRegister(lines []string, actor, at string) IntakeResult {
	records := make([]model.Record, 0, len(lines))
	for index, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) < 3 {
			records = append(records, model.Record{ID: fmt.Sprintf("invalid-%d", index), BatchID: ""})
			continue
		}
		quantity := 0
		for _, digit := range fields[2] {
			if digit >= '0' && digit <= '9' {
				quantity = quantity*10 + int(digit-'0')
			}
		}
		records = append(records, model.Record{ID: fields[0], BatchID: fields[1], Applicant: fields[0], Quantity: quantity})
	}
	return s.RegisterBatch(records, actor, at)
}

func (s *Service) ValidateBatch(batchID string) ([]policy.Finding, error) {
	items, err := s.ListRecords(batchID, "")
	if err != nil {
		return nil, err
	}
	return policy.CheckBatch(items), nil
}

func (s *Service) DuplicateIDs(records []model.Record) []string {
	seen := make(map[string]bool)
	duplicates := make([]string, 0)
	for _, item := range records {
		if seen[item.ID] {
			duplicates = append(duplicates, item.ID)
		}
		seen[item.ID] = true
	}
	return duplicates
}

func (s *Service) ReconcileBatch(batchID string) (map[string]int, error) {
	items, err := s.ListRecords(batchID, "")
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Status]++
	}
	return counts, nil
}

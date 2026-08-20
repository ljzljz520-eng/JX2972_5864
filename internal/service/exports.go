package service

import (
	"emergency-claim-code/internal/model"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func (s *Service) ExportCSV(writer io.Writer, term, status string) error {
	items, err := s.ListRecords(term, status)
	if err != nil {
		return err
	}
	output := csv.NewWriter(writer)
	if err := output.Write([]string{"id", "batch_id", "applicant", "quantity", "status", "revision"}); err != nil {
		return err
	}
	for _, item := range items {
		if err := output.Write([]string{item.ID, item.BatchID, item.Applicant, fmt.Sprint(item.Quantity), item.Status, fmt.Sprint(item.Revision)}); err != nil {
			return err
		}
	}
	output.Flush()
	return output.Error()
}

func ParseCSV(input string) ([]model.Record, error) {
	reader := csv.NewReader(strings.NewReader(input))
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(header) < 3 {
		return nil, fmt.Errorf("csv header incomplete")
	}
	result := make([]model.Record, 0)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) < 4 {
			return nil, fmt.Errorf("csv row incomplete")
		}
		quantity := 0
		if _, err := fmt.Sscanf(row[3], "%d", &quantity); err != nil {
			return nil, err
		}
		result = append(result, model.Record{ID: row[0], BatchID: row[1], Applicant: row[2], Quantity: quantity})
	}
	return result, nil
}

func (s *Service) ImportCSV(input string, actor, at string) (IntakeResult, error) {
	records, err := ParseCSV(input)
	if err != nil {
		return IntakeResult{}, err
	}
	return s.RegisterBatch(records, actor, at), nil
}

func (s *Service) ExportSummary() (string, error) {
	metrics, err := s.Metrics()
	if err != nil {
		return "", err
	}
	batches, err := s.BatchMetrics()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s batches=%d", metrics.String(), len(batches)), nil
}

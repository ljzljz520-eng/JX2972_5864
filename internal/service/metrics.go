package service

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"sort"
)

type Metrics struct {
	Total     int
	Draft     int
	Submitted int
	Approved  int
	Recalled  int
	Archived  int
	Rejected  int
	Quantity  int
}

func (s *Service) Metrics() (Metrics, error) {
	items, err := s.ListRecords("", "")
	if err != nil {
		return Metrics{}, err
	}
	return CalculateMetrics(items), nil
}

func CalculateMetrics(items []model.Record) Metrics {
	metrics := Metrics{Total: len(items)}
	for _, item := range items {
		metrics.Quantity += item.Quantity
		switch item.Status {
		case "draft":
			metrics.Draft++
		case "submitted":
			metrics.Submitted++
		case "approved":
			metrics.Approved++
		case "recalled":
			metrics.Recalled++
		case "archived":
			metrics.Archived++
		case "rejected":
			metrics.Rejected++
		}
	}
	return metrics
}

func (m Metrics) CompletionRate() int {
	if m.Total == 0 {
		return 0
	}
	return (m.Archived * 100) / m.Total
}

func (m Metrics) Pending() int { return m.Draft + m.Submitted + m.Recalled + m.Approved }

func (m Metrics) String() string {
	return fmt.Sprintf("total=%d pending=%d archived=%d rejected=%d quantity=%d", m.Total, m.Pending(), m.Archived, m.Rejected, m.Quantity)
}

type BatchMetric struct {
	BatchID  string
	Total    int
	Approved int
	Quantity int
}

func (s *Service) BatchMetrics() ([]BatchMetric, error) {
	items, err := s.ListRecords("", "")
	if err != nil {
		return nil, err
	}
	byBatch := make(map[string]*BatchMetric)
	for _, item := range items {
		metric := byBatch[item.BatchID]
		if metric == nil {
			metric = &BatchMetric{BatchID: item.BatchID}
			byBatch[item.BatchID] = metric
		}
		metric.Total++
		metric.Quantity += item.Quantity
		if item.Status == "approved" {
			metric.Approved++
		}
	}
	result := make([]BatchMetric, 0, len(byBatch))
	for _, metric := range byBatch {
		result = append(result, *metric)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].BatchID < result[j].BatchID })
	return result, nil
}

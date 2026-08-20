package service

import (
	"emergency-claim-code/internal/model"
	"fmt"
)

func (s *Service) Snapshot(id string) (map[string]any, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return nil, err
	}
	audits, err := s.Audits(id)
	if err != nil {
		return nil, err
	}
	attachments, err := s.Attachments(id)
	if err != nil {
		return nil, err
	}
	return map[string]any{"record": record, "audits": audits, "attachments": attachments}, nil
}

func (s *Service) Export(id string) (string, error) {
	snapshot, err := s.Snapshot(id)
	if err != nil {
		return "", err
	}
	data, err := model.Encode(snapshot)
	return string(data), err
}

func (s *Service) ImportRecord(record model.Record, actor, at string) (model.Record, error) {
	if record.ID == "" {
		return model.Record{}, fmt.Errorf("import requires id")
	}
	record.Status = "draft"
	return s.CreateRecord(record, actor, at)
}

func (s *Service) ResetCache() { s.clearCache() }

func (s *Service) Health() string { return "ready" }

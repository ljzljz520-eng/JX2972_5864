package service

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/policy"
	"fmt"
	"strings"
)

type ArchiveEntry struct {
	Record          model.Record
	AuditCount      int
	AttachmentCount int
	Eligible        bool
	Reason          string
}

func (s *Service) ArchivePreview(id string) (ArchiveEntry, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return ArchiveEntry{}, err
	}
	audits, err := s.Audits(id)
	if err != nil {
		return ArchiveEntry{}, err
	}
	attachments, err := s.Attachments(id)
	if err != nil {
		return ArchiveEntry{}, err
	}
	entry := ArchiveEntry{Record: record, AuditCount: len(audits), AttachmentCount: len(attachments)}
	if err := policy.CanArchive(record, len(audits)); err != nil {
		entry.Reason = err.Error()
		return entry, nil
	}
	entry.Eligible = true
	return entry, nil
}

func (s *Service) ArchiveWithReason(id, actor, reason, at string) (model.Record, error) {
	if strings.TrimSpace(reason) == "" {
		return model.Record{}, fmt.Errorf("archive reason required")
	}
	record, err := s.Archive(id, actor, at)
	if err != nil {
		return model.Record{}, err
	}
	if err := s.audit(record, "archive-reason", actor, reason, at); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Restore(id, actor, at string) (model.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != "archived" {
		return model.Record{}, fmt.Errorf("only archived records restore")
	}
	record.Status = "approved"
	record.Revision++
	record.UpdatedAt = at
	if err := s.repo.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	s.invalidate(id)
	s.putCache(record)
	if err := s.audit(record, "restore", actor, "restored for correction", at); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Purge(id, actor, at string) error {
	record, err := s.GetRecord(id)
	if err != nil {
		return err
	}
	if !model.IsTerminal(record.Status) {
		return fmt.Errorf("record is active")
	}
	if err := s.repo.DeleteRecord(id); err != nil {
		return err
	}
	s.invalidate(id)
	return s.audit(record, "purge", actor, "record removed", at)
}

func (s *Service) ArchiveBatch(batchID, actor, at string) ([]model.Record, []error) {
	items, err := s.ListRecords(batchID, "approved")
	if err != nil {
		return nil, []error{err}
	}
	archived := make([]model.Record, 0)
	errs := make([]error, 0)
	for _, item := range items {
		got, err := s.Archive(item.ID, actor, at)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		archived = append(archived, got)
	}
	return archived, errs
}

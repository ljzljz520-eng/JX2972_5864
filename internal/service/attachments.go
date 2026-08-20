package service

import (
	"emergency-claim-code/internal/model"
	"errors"
	"strings"
)

func (s *Service) AddAttachment(input model.Attachment, actor, at string) (model.Attachment, error) {
	if input.ID == "" {
		input.ID = s.nextID("AT")
	}
	input.AddedAt = at
	if errs := model.ValidateAttachment(input); len(errs) > 0 {
		return model.Attachment{}, errors.New(strings.Join(errs, "; "))
	}
	if _, err := s.GetRecord(input.RecordID); err != nil {
		return model.Attachment{}, err
	}
	if err := s.repo.SaveAttachment(input); err != nil {
		return model.Attachment{}, err
	}
	item, err := s.GetRecord(input.RecordID)
	if err == nil {
		_ = s.audit(item, "attachment", actor, input.Name, at)
	}
	return input, err
}

func (s *Service) Attachments(recordID string) ([]model.Attachment, error) {
	return s.repo.ListAttachments(recordID)
}

func (s *Service) HasAttachment(recordID, digest string) (bool, error) {
	items, err := s.Attachments(recordID)
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.Digest == digest {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) VerifyAttachment(item model.Attachment) bool {
	return len(model.ValidateAttachment(item)) == 0 && strings.TrimSpace(item.Digest) != ""
}

package workflow

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
	"emergency-claim-code/internal/service"
	"errors"
	"fmt"
)

var ErrStep = errors.New("workflow step rejected")

type Engine struct {
	service *service.Service
	repo    *repository.Repository
}

func New(svc *service.Service, repo *repository.Repository) *Engine {
	return &Engine{service: svc, repo: repo}
}

func (e *Engine) Register(record model.Record, actor, at string) (model.Record, error) {
	return e.service.CreateRecord(record, actor, at)
}

func (e *Engine) Review(id, actor, at string, approve bool) (model.Record, error) {
	if approve {
		return e.service.Approve(id, actor, at)
	}
	return e.service.Reject(id, actor, at)
}

func (e *Engine) Confirm(id, actor, at string) (model.ClaimCode, error) {
	record, err := e.service.GetRecord(id)
	if err != nil {
		return model.ClaimCode{}, err
	}
	if record.Status != "approved" {
		return model.ClaimCode{}, fmt.Errorf("%w: record is %s", ErrStep, record.Status)
	}
	return e.service.IssueCode(id, at)
}

func (e *Engine) Archive(id, actor, at string) (model.Record, error) {
	return e.service.Archive(id, actor, at)
}

func (e *Engine) CreateReviewArchive(record model.Record, actor, at string) (model.Record, error) {
	item, err := e.Register(record, actor, at)
	if err != nil {
		return model.Record{}, err
	}
	if _, err = e.service.Submit(item.ID, actor, at); err != nil {
		return model.Record{}, err
	}
	if _, err = e.Review(item.ID, actor, at, true); err != nil {
		return model.Record{}, err
	}
	return e.Archive(item.ID, actor, at)
}

func (e *Engine) RecallRetry(id, actor, at string) (model.Record, error) {
	return e.service.RecallAndRetry(id, actor, at)
}

func (e *Engine) Publish(id, actor, at string) (model.ClaimCode, error) {
	if _, err := e.service.Approve(id, actor, at); err != nil {
		return model.ClaimCode{}, err
	}
	return e.Confirm(id, actor, at)
}

func (e *Engine) ImportReport(record model.Record, actor, at string) (string, error) {
	item, err := e.service.ImportRecord(record, actor, at)
	if err != nil {
		return "", err
	}
	return e.service.Export(item.ID)
}

func (e *Engine) SearchUpdatePublish(term, status, id, note, actor, at string) ([]model.Record, model.ClaimCode, error) {
	found, err := e.service.ListRecords(term, status)
	if err != nil {
		return nil, model.ClaimCode{}, err
	}
	if id == "" && len(found) > 0 {
		id = found[0].ID
	}
	if id == "" {
		return found, model.ClaimCode{}, ErrStep
	}
	item, err := e.service.GetRecord(id)
	if err != nil {
		return found, model.ClaimCode{}, err
	}
	item.Note = note
	item.UpdatedAt = at
	if err := e.repo.SaveRecord(item); err != nil {
		return found, model.ClaimCode{}, err
	}
	e.service.ResetCache()
	var code model.ClaimCode
	if item.Status == "approved" {
		code, err = e.Confirm(id, actor, at)
	} else {
		code, err = e.Publish(id, actor, at)
	}
	return found, code, err
}

func (e *Engine) AttachAndSnapshot(id string, attachment model.Attachment, actor, at string) (string, error) {
	if _, err := e.service.AddAttachment(attachment, actor, at); err != nil {
		return "", err
	}
	return e.service.Export(id)
}

package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/repository"
)

var ErrConflict = errors.New("record state conflict")
var ErrInvalidAction = errors.New("invalid action")

type Service struct {
	repo     *repository.Repository
	mu       sync.RWMutex
	cache    map[string]model.Record
	sequence int
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo, cache: make(map[string]model.Record)}
}

func (s *Service) nextID(prefix string) string {
	s.sequence++
	return fmt.Sprintf("%s-%04d", prefix, s.sequence)
}

func (s *Service) CreateRecord(input model.Record, actor, at string) (model.Record, error) {
	input = model.NormalizeRecord(input)
	if input.ID == "" {
		input.ID = s.nextID("RB")
	}
	if input.Status != "draft" {
		return model.Record{}, ErrInvalidAction
	}
	if errs := model.ValidateRecord(input); len(errs) > 0 {
		return model.Record{}, errors.New(strings.Join(errs, "; "))
	}
	input.CreatedAt = at
	input.UpdatedAt = at
	if err := s.repo.SaveRecord(input); err != nil {
		return model.Record{}, err
	}
	s.putCache(input)
	if err := s.audit(input, "create", actor, "record registered", at); err != nil {
		return model.Record{}, err
	}
	return input, nil
}

func (s *Service) GetRecord(id string) (model.Record, error) {
	s.mu.RLock()
	cached, ok := s.cache[id]
	s.mu.RUnlock()
	if ok {
		return model.CloneRecord(cached), nil
	}
	item, err := s.repo.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	s.putCache(item)
	return item, nil
}

func (s *Service) putCache(item model.Record) {
	s.mu.Lock()
	s.cache[item.ID] = model.CloneRecord(item)
	s.mu.Unlock()
}

func (s *Service) invalidate(id string) { s.mu.Lock(); delete(s.cache, id); s.mu.Unlock() }

func (s *Service) clearCache() { s.mu.Lock(); s.cache = make(map[string]model.Record); s.mu.Unlock() }

func (s *Service) Reload(id string) (model.Record, error) { s.invalidate(id); return s.GetRecord(id) }

func (s *Service) ListRecords(term, status string) ([]model.Record, error) {
	items, err := s.repo.ListRecords()
	if err != nil {
		return nil, err
	}
	return repository.SortRecords(repository.FilterRecords(items, term, status)), nil
}

func (s *Service) ChangeStatus(id, target, actor, at string) (model.Record, error) {
	item, err := s.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if !model.CanTransition(item.Status, target) {
		return model.Record{}, fmt.Errorf("%w: %s to %s", ErrConflict, item.Status, target)
	}
	item.Status = target
	item.Revision++
	item.UpdatedAt = at
	if err := s.repo.SaveRecord(item); err != nil {
		return model.Record{}, err
	}
	// Always invalidate the stale cached decision and repopulate from the
	// freshly persisted record. Recalling a batch must not leave the prior
	// approved decision live in the cache: a subsequent retry reads the
	// confirmed, independently persisted status rather than a stale copy,
	// so repeated recall-and-retry actions stay mutually independent.
	s.invalidate(id)
	s.putCache(item)
	if err := s.audit(item, "status", actor, target, at); err != nil {
		return model.Record{}, err
	}
	return item, nil
}

func (s *Service) RecallAndRetry(id, actor, at string) (model.Record, error) {
	if _, err := s.ChangeStatus(id, "recalled", actor, at); err != nil {
		return model.Record{}, err
	}
	return s.ChangeStatus(id, "submitted", actor, at)
}

func (s *Service) Approve(id, actor, at string) (model.Record, error) {
	return s.ChangeStatus(id, "approved", actor, at)
}

func (s *Service) Submit(id, actor, at string) (model.Record, error) {
	return s.ChangeStatus(id, "submitted", actor, at)
}

func (s *Service) Archive(id, actor, at string) (model.Record, error) {
	return s.ChangeStatus(id, "archived", actor, at)
}

func (s *Service) Reject(id, actor, at string) (model.Record, error) {
	return s.ChangeStatus(id, "rejected", actor, at)
}

func (s *Service) audit(item model.Record, action, actor, detail, at string) error {
	event := model.AuditEvent{ID: s.nextID("AU"), RecordID: item.ID, Action: action, Actor: actor, Detail: detail, At: at}
	return s.repo.SaveAudit(event)
}

func (s *Service) Audits(id string) ([]model.AuditEvent, error) { return s.repo.ListAudits(id) }

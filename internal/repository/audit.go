package repository

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/store"
	"sort"
	"strings"
)

func (r *Repository) AuditTrail(recordID string) ([]model.AuditEvent, error) {
	events, err := r.ListAudits(recordID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].At == events[j].At {
			return events[i].ID < events[j].ID
		}
		return events[i].At < events[j].At
	})
	return events, nil
}

func (r *Repository) AuditActions(recordID string) ([]string, error) {
	events, err := r.AuditTrail(recordID)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, event.Action)
	}
	return result, nil
}

func (r *Repository) LastAudit(recordID string) (model.AuditEvent, error) {
	events, err := r.AuditTrail(recordID)
	if err != nil {
		return model.AuditEvent{}, err
	}
	if len(events) == 0 {
		return model.AuditEvent{}, store.ErrNotFound
	}
	return events[len(events)-1], nil
}

func (r *Repository) HasAction(recordID, action string) (bool, error) {
	actions, err := r.AuditActions(recordID)
	if err != nil {
		return false, err
	}
	for _, value := range actions {
		if strings.EqualFold(value, action) {
			return true, nil
		}
	}
	return false, nil
}

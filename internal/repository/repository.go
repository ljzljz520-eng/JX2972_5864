package repository

import (
	"fmt"

	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/store"
)

var recordBucket = []byte("records")
var auditBucket = []byte("audits")
var workflowBucket = []byte("workflows")
var attachmentBucket = []byte("attachments")
var codeBucket = []byte("claim_codes")

type Repository struct{ db *store.Store }

func New(db *store.Store) *Repository { return &Repository{db: db} }

func (r *Repository) SaveRecord(record model.Record) error {
	data, err := model.Encode(record)
	if err != nil {
		return err
	}
	return r.db.Put(recordBucket, []byte(record.ID), data)
}

func (r *Repository) GetRecord(id string) (model.Record, error) {
	data, err := r.db.Get(recordBucket, []byte(id))
	if err != nil {
		return model.Record{}, err
	}
	var out model.Record
	err = model.Decode(data, &out)
	return out, err
}

func (r *Repository) DeleteRecord(id string) error { return r.db.Delete(recordBucket, []byte(id)) }

func (r *Repository) ListRecords() ([]model.Record, error) {
	entries, err := r.db.List(recordBucket)
	if err != nil {
		return nil, err
	}
	out := make([]model.Record, 0, len(entries))
	for _, entry := range store.SortEntries(entries) {
		var item model.Record
		if err := model.Decode(entry.Value, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *Repository) SaveAudit(event model.AuditEvent) error {
	data, err := model.Encode(event)
	if err != nil {
		return err
	}
	return r.db.Put(auditBucket, []byte(event.ID), data)
}

func (r *Repository) ListAudits(recordID string) ([]model.AuditEvent, error) {
	entries, err := r.db.List(auditBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.AuditEvent, 0)
	for _, entry := range store.SortEntries(entries) {
		var item model.AuditEvent
		if err := model.Decode(entry.Value, &item); err != nil {
			return nil, err
		}
		if recordID == "" || item.RecordID == recordID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *Repository) SaveWorkflow(item model.Workflow) error {
	data, err := model.Encode(item)
	if err != nil {
		return err
	}
	return r.db.Put(workflowBucket, []byte(item.ID), data)
}

func (r *Repository) GetWorkflow(id string) (model.Workflow, error) {
	data, err := r.db.Get(workflowBucket, []byte(id))
	if err != nil {
		return model.Workflow{}, err
	}
	var out model.Workflow
	err = model.Decode(data, &out)
	return out, err
}

func (r *Repository) ListWorkflows(batchID string) ([]model.Workflow, error) {
	entries, err := r.db.List(workflowBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.Workflow, 0)
	for _, entry := range store.SortEntries(entries) {
		var item model.Workflow
		if err := model.Decode(entry.Value, &item); err != nil {
			return nil, err
		}
		if batchID == "" || item.BatchID == batchID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *Repository) SaveAttachment(item model.Attachment) error {
	data, err := model.Encode(item)
	if err != nil {
		return err
	}
	return r.db.Put(attachmentBucket, []byte(item.ID), data)
}

func (r *Repository) ListAttachments(recordID string) ([]model.Attachment, error) {
	entries, err := r.db.List(attachmentBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.Attachment, 0)
	for _, entry := range store.SortEntries(entries) {
		var item model.Attachment
		if err := model.Decode(entry.Value, &item); err != nil {
			return nil, err
		}
		if recordID == "" || item.RecordID == recordID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *Repository) SaveClaimCode(item model.ClaimCode) error {
	data, err := model.Encode(item)
	if err != nil {
		return err
	}
	return r.db.Put(codeBucket, []byte(item.Code), data)
}

func (r *Repository) GetClaimCode(code string) (model.ClaimCode, error) {
	data, err := r.db.Get(codeBucket, []byte(code))
	if err != nil {
		return model.ClaimCode{}, err
	}
	var out model.ClaimCode
	err = model.Decode(data, &out)
	return out, err
}

func (r *Repository) Count(bucket []byte) (int, error) {
	entries, err := r.db.List(bucket)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

func (r *Repository) Summary() (string, error) {
	a, err := r.Count(recordBucket)
	if err != nil {
		return "", err
	}
	b, err := r.Count(auditBucket)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("records=%d audits=%d", a, b), nil
}

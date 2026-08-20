package repository

import (
	"emergency-claim-code/internal/model"
	"emergency-claim-code/internal/store"
	"fmt"
	"strings"
)

type BucketCounts struct {
	Records     int
	Audits      int
	Workflows   int
	Attachments int
	Codes       int
}

func (r *Repository) BucketCounts() (BucketCounts, error) {
	records, err := r.Count(recordBucket)
	if err != nil {
		return BucketCounts{}, err
	}
	audits, err := r.Count(auditBucket)
	if err != nil {
		return BucketCounts{}, err
	}
	workflows, err := r.Count(workflowBucket)
	if err != nil {
		return BucketCounts{}, err
	}
	attachments, err := r.Count(attachmentBucket)
	if err != nil {
		return BucketCounts{}, err
	}
	codes, err := r.Count(codeBucket)
	if err != nil {
		return BucketCounts{}, err
	}
	return BucketCounts{Records: records, Audits: audits, Workflows: workflows, Attachments: attachments, Codes: codes}, nil
}

func (c BucketCounts) Total() int {
	return c.Records + c.Audits + c.Workflows + c.Attachments + c.Codes
}

func (c BucketCounts) String() string {
	return fmt.Sprintf("records=%d audits=%d workflows=%d attachments=%d codes=%d", c.Records, c.Audits, c.Workflows, c.Attachments, c.Codes)
}

func (r *Repository) FindByBatch(batchID string) ([]model.Record, error) {
	records, err := r.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0)
	for _, record := range records {
		if record.BatchID == batchID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (r *Repository) FindByApplicant(applicant string) ([]model.Record, error) {
	records, err := r.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0)
	needle := strings.ToLower(applicant)
	for _, record := range records {
		if strings.Contains(strings.ToLower(record.Applicant), needle) {
			result = append(result, record)
		}
	}
	return result, nil
}

func (r *Repository) SaveMany(records []model.Record) error {
	return r.db.Update(func(tx *store.Txn) error {
		for _, record := range records {
			data, err := model.Encode(record)
			if err != nil {
				return err
			}
			if err := tx.Put(recordBucket, []byte(record.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteMany(ids []string) error {
	return r.db.Update(func(tx *store.Txn) error {
		for _, id := range ids {
			if err := tx.Delete(recordBucket, []byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) ReadRaw(bucket, key string) ([]byte, error) {
	return r.db.Get([]byte(bucket), []byte(key))
}

func (r *Repository) WriteRaw(bucket, key string, value []byte) error {
	return r.db.Put([]byte(bucket), []byte(key), value)
}

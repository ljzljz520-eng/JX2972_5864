package model

import "strings"

func ValidateRecord(r Record) []string {
	errs := make([]string, 0, 4)
	if strings.TrimSpace(r.ID) == "" {
		errs = append(errs, "id is required")
	}
	if strings.TrimSpace(r.BatchID) == "" {
		errs = append(errs, "batch_id is required")
	}
	if strings.TrimSpace(r.Applicant) == "" {
		errs = append(errs, "applicant is required")
	}
	if r.Quantity <= 0 {
		errs = append(errs, "quantity must be positive")
	}
	if len(r.BatchID) > 64 {
		errs = append(errs, "batch_id is too long")
	}
	if len(r.Applicant) > 120 {
		errs = append(errs, "applicant is too long")
	}
	return errs
}

func ValidateAttachment(a Attachment) []string {
	errs := make([]string, 0, 3)
	if strings.TrimSpace(a.ID) == "" {
		errs = append(errs, "attachment id is required")
	}
	if strings.TrimSpace(a.RecordID) == "" {
		errs = append(errs, "record id is required")
	}
	if strings.TrimSpace(a.Name) == "" {
		errs = append(errs, "attachment name is required")
	}
	if strings.TrimSpace(a.Digest) == "" {
		errs = append(errs, "attachment digest is required")
	}
	return errs
}

func ValidateWorkflow(w Workflow) []string {
	errs := make([]string, 0, 2)
	if w.ID == "" {
		errs = append(errs, "workflow id is required")
	}
	if w.BatchID == "" {
		errs = append(errs, "workflow batch is required")
	}
	if w.Stage == "" {
		errs = append(errs, "workflow stage is required")
	}
	return errs
}

func ValidateClaimCode(c ClaimCode) []string {
	errs := make([]string, 0, 2)
	if c.Code == "" {
		errs = append(errs, "code is required")
	}
	if c.RecordID == "" {
		errs = append(errs, "record id is required")
	}
	if c.BatchID == "" {
		errs = append(errs, "batch id is required")
	}
	return errs
}

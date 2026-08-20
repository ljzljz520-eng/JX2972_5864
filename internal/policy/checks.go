package policy

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type Finding struct {
	Code     string
	Field    string
	Message  string
	Severity string
}

func CheckRecord(record model.Record) []Finding {
	findings := make([]Finding, 0)
	for _, message := range model.ValidateRecord(record) {
		findings = append(findings, Finding{Code: "record.invalid", Field: "record", Message: message, Severity: "error"})
	}
	if record.Quantity > 100 {
		findings = append(findings, Finding{Code: "quantity.review", Field: "quantity", Message: "large quantity requires coordinator review", Severity: "warning"})
	}
	if strings.Contains(record.Applicant, "\n") {
		findings = append(findings, Finding{Code: "applicant.control", Field: "applicant", Message: "line breaks are not allowed", Severity: "error"})
	}
	return findings
}

func CheckBatch(records []model.Record) []Finding {
	findings := make([]Finding, 0)
	if len(records) == 0 {
		return append(findings, Finding{Code: "batch.empty", Field: "batch", Message: "batch has no records", Severity: "warning"})
	}
	seen := make(map[string]bool)
	for _, record := range records {
		if seen[record.ID] {
			findings = append(findings, Finding{Code: "record.duplicate", Field: "id", Message: record.ID, Severity: "error"})
		}
		seen[record.ID] = true
	}
	return findings
}

func CheckAttachment(attachment model.Attachment) []Finding {
	findings := make([]Finding, 0)
	if attachment.Kind == "" {
		findings = append(findings, Finding{Code: "attachment.kind", Field: "kind", Message: "attachment kind should be declared", Severity: "warning"})
	}
	if len(attachment.Digest) < 4 {
		findings = append(findings, Finding{Code: "attachment.digest", Field: "digest", Message: "digest is too short", Severity: "error"})
	}
	return findings
}

func HasErrors(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == "error" {
			return true
		}
	}
	return false
}

func Summarize(findings []Finding) string {
	if len(findings) == 0 {
		return "no findings"
	}
	parts := make([]string, 0, len(findings))
	for _, finding := range findings {
		parts = append(parts, fmt.Sprintf("%s:%s", finding.Code, finding.Message))
	}
	return strings.Join(parts, "; ")
}

func RequireEvidence(record model.Record, attachments []model.Attachment) error {
	if record.Quantity <= 20 {
		return nil
	}
	if len(attachments) == 0 {
		return fmt.Errorf("large claim requires evidence")
	}
	return nil
}

func RequireApproval(record model.Record) error {
	if record.Status != "approved" {
		return fmt.Errorf("status %s is not approved", record.Status)
	}
	return nil
}

func RequireArchive(record model.Record) error {
	if record.Status != "approved" {
		return fmt.Errorf("status %s cannot archive", record.Status)
	}
	if record.Revision < 2 {
		return fmt.Errorf("record needs review revision")
	}
	return nil
}

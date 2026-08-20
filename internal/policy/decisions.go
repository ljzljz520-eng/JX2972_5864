package policy

import (
	"emergency-claim-code/internal/model"
	"fmt"
)

type Decision struct {
	Approved bool
	Route    string
	Reasons  []string
	Required []string
}

func Evaluate(record model.Record, actor Role, attachments []model.Attachment) Decision {
	decision := Decision{Approved: true, Route: "standard", Reasons: make([]string, 0), Required: make([]string, 0)}
	if len(CheckRecord(record)) > 0 {
		decision.Approved = false
		decision.Reasons = append(decision.Reasons, "record checks failed")
	}
	if record.Quantity > 100 {
		decision.Route = "coordinator"
		decision.Required = append(decision.Required, "coordinator")
	}
	if record.Quantity > 20 && len(attachments) == 0 {
		decision.Approved = false
		decision.Reasons = append(decision.Reasons, "evidence missing")
	}
	if actor == RoleApplicant && record.Quantity > 50 {
		decision.Approved = false
		decision.Reasons = append(decision.Reasons, "applicant cannot submit bulk claim")
	}
	if record.BatchID == "RB2972-01" {
		decision.Required = append(decision.Required, "batch-owner")
	}
	return decision
}

func RouteForQuantity(quantity int) string {
	switch {
	case quantity <= 10:
		return "standard"
	case quantity <= 100:
		return "review"
	default:
		return "coordinator"
	}
}

func RequiresSecondReview(record model.Record) bool {
	if record.Quantity > 100 {
		return true
	}
	if record.Revision > 2 {
		return true
	}
	return record.Status == "recalled"
}

func CanIssue(record model.Record, attachments []model.Attachment) error {
	if err := RequireApproval(record); err != nil {
		return err
	}
	if err := RequireEvidence(record, attachments); err != nil {
		return err
	}
	if RequiresSecondReview(record) && record.Revision < 2 {
		return fmt.Errorf("second review required")
	}
	return nil
}

func CanArchive(record model.Record, auditCount int) error {
	if err := RequireArchive(record); err != nil {
		return err
	}
	if auditCount < 3 {
		return fmt.Errorf("archive requires audit trail")
	}
	return nil
}

func ReviewChecklist(record model.Record) []string {
	checklist := []string{"identity", "quantity", "eligibility"}
	if record.Quantity > 20 {
		checklist = append(checklist, "evidence")
	}
	if record.BatchID == "RB2972-01" {
		checklist = append(checklist, "batch reconciliation")
	}
	return checklist
}

func DecisionText(decision Decision) string {
	if decision.Approved {
		return fmt.Sprintf("approved route=%s", decision.Route)
	}
	return "blocked: " + fmt.Sprint(decision.Reasons)
}

package workflow

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type ChecklistItem struct {
	Key      string
	Label    string
	Required bool
	Complete bool
	Note     string
}

func NewChecklist(record model.Record) []ChecklistItem {
	items := []ChecklistItem{{Key: "identity", Label: "Applicant identity", Required: true}, {Key: "quantity", Label: "Requested quantity", Required: true}, {Key: "eligibility", Label: "Eligibility evidence", Required: true}, {Key: "contact", Label: "Contact confirmation", Required: false}}
	if record.Quantity > 20 {
		items = append(items, ChecklistItem{Key: "attachment", Label: "Supporting attachment", Required: true})
	}
	return items
}

func MarkComplete(items []ChecklistItem, key, note string) ([]ChecklistItem, error) {
	found := false
	for index := range items {
		if items[index].Key == key {
			items[index].Complete = true
			items[index].Note = note
			found = true
		}
	}
	if !found {
		return items, fmt.Errorf("checklist item %s not found", key)
	}
	return items, nil
}

func MissingRequired(items []ChecklistItem) []ChecklistItem {
	missing := make([]ChecklistItem, 0)
	for _, item := range items {
		if item.Required && !item.Complete {
			missing = append(missing, item)
		}
	}
	return missing
}

func ChecklistReady(items []ChecklistItem) bool { return len(MissingRequired(items)) == 0 }

func ChecklistSummary(items []ChecklistItem) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		state := "pending"
		if item.Complete {
			state = "done"
		}
		parts = append(parts, item.Key+"="+state)
	}
	return strings.Join(parts, ",")
}

func DefaultWorkflow(record model.Record, owner, at string) (model.Workflow, []ChecklistItem) {
	return BuildWorkflow(record, owner, at), NewChecklist(record)
}

func ValidateWorkflowInput(item model.Workflow, checklist []ChecklistItem) error {
	if len(model.ValidateWorkflow(item)) > 0 {
		return fmt.Errorf("workflow invalid")
	}
	if !ChecklistReady(checklist) {
		return fmt.Errorf("workflow checklist incomplete")
	}
	return nil
}

package workflow

import (
	"emergency-claim-code/internal/model"
	"fmt"
)

func Stages() []string { return []string{"register", "review", "confirm", "archive"} }

func IsStage(value string) bool {
	for _, stage := range Stages() {
		if stage == value {
			return true
		}
	}
	return false
}

func NextStage(value string) string {
	switch value {
	case "register":
		return "review"
	case "review":
		return "confirm"
	case "confirm":
		return "archive"
	default:
		return ""
	}
}

func BuildWorkflow(record model.Record, owner, at string) model.Workflow {
	return model.Workflow{ID: record.ID + "-flow", BatchID: record.BatchID, Stage: "register", Owner: owner, Summary: "claim-code intake", UpdatedAt: at}
}

func AdvanceWorkflow(item model.Workflow, at string) (model.Workflow, error) {
	next := NextStage(item.Stage)
	if next == "" {
		return item, fmt.Errorf("cannot advance %s", item.Stage)
	}
	item.Stage = next
	item.UpdatedAt = at
	return item, nil
}

func CompleteWorkflow(item model.Workflow, at string) (model.Workflow, error) {
	for !IsStage(item.Stage) {
		return item, fmt.Errorf("invalid stage")
	}
	item.Stage = "archive"
	item.UpdatedAt = at
	return item, nil
}

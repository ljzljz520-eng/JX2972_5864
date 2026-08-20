package query

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"sort"
	"strings"
)

type SortField string

const (
	SortID       SortField = "id"
	SortBatch    SortField = "batch"
	SortQuantity SortField = "quantity"
	SortRevision SortField = "revision"
)

type Plan struct {
	Filter Filter
	Sort   SortField
	Desc   bool
	Limit  int
	Offset int
}

func NewPlan(filter Filter, sortField SortField, desc bool, limit, offset int) Plan {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	return Plan{Filter: filter, Sort: sortField, Desc: desc, Limit: limit, Offset: offset}
}

func SortFieldValid(field SortField) bool {
	switch field {
	case SortID, SortBatch, SortQuantity, SortRevision:
		return true
	default:
		return false
	}
}

func ApplyPlan(items []model.Record, plan Plan) []model.Record {
	result := Apply(items, plan.Filter)
	sort.SliceStable(result, func(i, j int) bool {
		less := compare(result[i], result[j], plan.Sort)
		if plan.Desc {
			return !less && result[i].ID != result[j].ID
		}
		return less
	})
	start := plan.Offset
	if start > len(result) {
		start = len(result)
	}
	end := len(result)
	if plan.Limit > 0 && start+plan.Limit < end {
		end = start + plan.Limit
	}
	return result[start:end]
}

func compare(left, right model.Record, field SortField) bool {
	switch field {
	case SortBatch:
		if left.BatchID == right.BatchID {
			return left.ID < right.ID
		}
		return left.BatchID < right.BatchID
	case SortQuantity:
		if left.Quantity == right.Quantity {
			return left.ID < right.ID
		}
		return left.Quantity < right.Quantity
	case SortRevision:
		if left.Revision == right.Revision {
			return left.ID < right.ID
		}
		return left.Revision < right.Revision
	default:
		return left.ID < right.ID
	}
}

func ParsePlan(text string) (Plan, error) {
	plan := NewPlan(Filter{}, SortID, false, 0, 0)
	for _, part := range strings.Split(text, ",") {
		bits := strings.SplitN(part, "=", 2)
		if len(bits) != 2 {
			return plan, fmt.Errorf("invalid plan component")
		}
		switch bits[0] {
		case "q":
			plan.Filter.Term = bits[1]
		case "status":
			plan.Filter.Status = bits[1]
		case "batch":
			plan.Filter.Batch = bits[1]
		case "sort":
			plan.Sort = SortField(bits[1])
			if !SortFieldValid(plan.Sort) {
				return plan, fmt.Errorf("invalid sort")
			}
		case "desc":
			plan.Desc = bits[1] == "true"
		case "limit":
			if _, err := fmt.Sscanf(bits[1], "%d", &plan.Limit); err != nil {
				return plan, err
			}
		case "offset":
			if _, err := fmt.Sscanf(bits[1], "%d", &plan.Offset); err != nil {
				return plan, err
			}
		default:
			return plan, fmt.Errorf("unknown plan field %s", bits[0])
		}
	}
	return plan, nil
}

func Cursor(items []model.Record, after string, limit int) []model.Record {
	started := after == ""
	result := make([]model.Record, 0)
	for _, item := range items {
		if !started {
			if item.ID == after {
				started = true
			}
			continue
		}
		if limit > 0 && len(result) >= limit {
			break
		}
		result = append(result, item)
	}
	return result
}

func Facets(items []model.Record) map[string][]string {
	result := map[string][]string{"status": {}, "batch": {}}
	seenStatus := make(map[string]bool)
	seenBatch := make(map[string]bool)
	for _, item := range items {
		if !seenStatus[item.Status] {
			result["status"] = append(result["status"], item.Status)
			seenStatus[item.Status] = true
		}
		if !seenBatch[item.BatchID] {
			result["batch"] = append(result["batch"], item.BatchID)
			seenBatch[item.BatchID] = true
		}
	}
	sort.Strings(result["status"])
	sort.Strings(result["batch"])
	return result
}

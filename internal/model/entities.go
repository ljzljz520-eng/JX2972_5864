package model

import "encoding/json"

type Record struct {
	ID        string `json:"id"`
	BatchID   string `json:"batch_id"`
	Applicant string `json:"applicant"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
	Note      string `json:"note"`
	Revision  int    `json:"revision"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	At       string `json:"at"`
}

type Workflow struct {
	ID        string `json:"id"`
	BatchID   string `json:"batch_id"`
	Stage     string `json:"stage"`
	Owner     string `json:"owner"`
	Summary   string `json:"summary"`
	UpdatedAt string `json:"updated_at"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Digest   string `json:"digest"`
	Kind     string `json:"kind"`
	AddedAt  string `json:"added_at"`
}

type ClaimCode struct {
	Code     string `json:"code"`
	RecordID string `json:"record_id"`
	BatchID  string `json:"batch_id"`
	State    string `json:"state"`
	IssuedAt string `json:"issued_at"`
}

func Encode(v any) ([]byte, error) { return json.Marshal(v) }

func Decode(data []byte, target any) error { return json.Unmarshal(data, target) }

func ValidStatuses() []string {
	return []string{"draft", "submitted", "approved", "recalled", "archived", "rejected"}
}

func IsTerminal(status string) bool {
	if status == "archived" || status == "rejected" {
		return true
	}
	return false
}

func CanTransition(from, to string) bool {
	switch from {
	case "draft":
		return to == "submitted" || to == "rejected"
	case "submitted":
		return to == "approved" || to == "rejected" || to == "recalled"
	case "approved":
		return to == "recalled" || to == "archived"
	case "recalled":
		return to == "submitted"
	case "rejected":
		return false
	case "archived":
		return false
	default:
		return false
	}
}

func CloneRecord(r Record) Record { return r }

func NormalizeRecord(r Record) Record {
	if r.Status == "" {
		r.Status = "draft"
	}
	if r.Revision < 1 {
		r.Revision = 1
	}
	return r
}

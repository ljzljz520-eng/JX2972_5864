package policy

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type Role string

const (
	RoleApplicant   Role = "applicant"
	RoleClerk       Role = "clerk"
	RoleReviewer    Role = "reviewer"
	RoleAuditor     Role = "auditor"
	RoleCoordinator Role = "coordinator"
)

type Rule struct {
	Action      string
	Allowed     []Role
	From        []string
	To          string
	Description string
}

func DefaultRules() []Rule {
	return []Rule{
		{Action: "create", Allowed: []Role{RoleApplicant, RoleClerk}, From: []string{""}, To: "draft", Description: "register a claim"},
		{Action: "submit", Allowed: []Role{RoleApplicant, RoleClerk}, From: []string{"draft", "recalled"}, To: "submitted", Description: "send for review"},
		{Action: "approve", Allowed: []Role{RoleReviewer}, From: []string{"submitted"}, To: "approved", Description: "approve evidence"},
		{Action: "reject", Allowed: []Role{RoleReviewer}, From: []string{"draft", "submitted"}, To: "rejected", Description: "reject a claim"},
		{Action: "recall", Allowed: []Role{RoleClerk, RoleCoordinator}, From: []string{"submitted", "approved"}, To: "recalled", Description: "recall a claim"},
		{Action: "archive", Allowed: []Role{RoleAuditor, RoleCoordinator}, From: []string{"approved"}, To: "archived", Description: "archive a claim"},
	}
}

func FindRule(action string) (Rule, bool) {
	for _, rule := range DefaultRules() {
		if rule.Action == action {
			return rule, true
		}
	}
	return Rule{}, false
}

func HasRole(roles []Role, expected Role) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

func AllowedAction(actor Role, action, status string) bool {
	rule, ok := FindRule(action)
	if !ok {
		return false
	}
	if !HasRole(rule.Allowed, actor) {
		return false
	}
	for _, from := range rule.From {
		if from == status {
			return true
		}
	}
	return false
}

func Explain(action, status string) string {
	rule, ok := FindRule(action)
	if !ok {
		return "unknown action"
	}
	if !contains(rule.From, status) {
		return fmt.Sprintf("%s cannot run from %s", action, status)
	}
	return rule.Description
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func ParseRole(value string) (Role, error) {
	normalized := Role(strings.ToLower(strings.TrimSpace(value)))
	switch normalized {
	case RoleApplicant, RoleClerk, RoleReviewer, RoleAuditor, RoleCoordinator:
		return normalized, nil
	default:
		return "", fmt.Errorf("unknown role %q", value)
	}
}

func ValidateTransition(actor Role, action string, before, after string) error {
	rule, ok := FindRule(action)
	if !ok {
		return fmt.Errorf("action %s is undefined", action)
	}
	if !HasRole(rule.Allowed, actor) {
		return fmt.Errorf("role %s cannot %s", actor, action)
	}
	if !contains(rule.From, before) {
		return fmt.Errorf("action %s cannot start at %s", action, before)
	}
	if rule.To != after {
		return fmt.Errorf("action %s must end at %s", action, rule.To)
	}
	if !model.CanTransition(before, after) {
		return fmt.Errorf("model rejects %s to %s", before, after)
	}
	return nil
}

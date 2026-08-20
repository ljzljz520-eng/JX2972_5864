package httpapi

import (
	"net/http"
	"strings"
)

func Actor(r *http.Request) string {
	actor := strings.TrimSpace(r.Header.Get("X-Actor"))
	if actor == "" {
		return "anonymous"
	}
	return actor
}

func At(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-At"))
	if value == "" {
		return "static"
	}
	return value
}

func IsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

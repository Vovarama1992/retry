package utils

import (
	"encoding/json"
	"strings"

	"github.com/Vovarama1992/retry/pkg/domain"
)

type MentorArtstationRule struct{}

func (MentorArtstationRule) ShouldSkip(a domain.Action) bool {
	if a.ActionTypeName != "external_link_mentor_page" {
		return false
	}

	var m map[string]any
	if err := json.Unmarshal(a.Meta, &m); err != nil {
		return false
	}

	if name, ok := m["mentor_name"].(string); ok && strings.EqualFold(name, "artstation") {
		return true
	}
	return false
}

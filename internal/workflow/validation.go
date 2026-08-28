package workflow

import (
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"strings"
)

type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

func ValidateRecordForWorkflow(record model.Record, workflowName string) ValidationResult {
	result := ValidationResult{Valid: true, Errors: make([]string, 0), Warnings: make([]string, 0)}
	if record.ID == "" {
		result.Errors = append(result.Errors, "record id is required")
	}
	if workflowName == "intake" {
		if strings.TrimSpace(record.Title) == "" {
			result.Errors = append(result.Errors, "title is required")
		}
		if len([]rune(record.Body)) < 16 {
			result.Errors = append(result.Errors, "body is too short")
		}
	}
	if workflowName == "review" && record.Status == model.StatusReceived {
		result.Warnings = append(result.Warnings, "record has not been reviewed")
	}
	if workflowName == "publication" && record.Status != model.StatusApproved {
		result.Errors = append(result.Errors, "record is not approved")
	}
	result.Valid = len(result.Errors) == 0
	return result
}

func ValidateSteps(definition Definition) error {
	if definition.Name == "" {
		return errors.New("workflow name is required")
	}
	if len(definition.Steps) < 4 {
		return errors.New("workflow needs four steps")
	}
	seen := map[string]bool{}
	for _, step := range definition.Steps {
		if strings.TrimSpace(step) == "" {
			return errors.New("workflow step cannot be empty")
		}
		if seen[step] {
			return fmt.Errorf("workflow step %s is duplicated", step)
		}
		seen[step] = true
	}
	return nil
}

func DefaultDefinitions() []Definition {
	return []Definition{{Name: "intake", Steps: []string{"receive", "validate", "save", "display"}}, {Name: "review", Steps: []string{"register", "review", "archive", "query"}}, {Name: "publication", Steps: []string{"submit", "process", "notify", "track"}}}
}

func MatchWorkflow(record model.Record, definition Definition) bool {
	if definition.Name == "intake" {
		return record.Status == model.StatusReceived
	}
	if definition.Name == "review" {
		return record.Status == model.StatusReviewed || record.Status == model.StatusApproved
	}
	if definition.Name == "publication" {
		return record.Status == model.StatusPublished
	}
	return false
}

package model

import (
	"fmt"
	"strings"
)

func NormalizeCategory(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func ValidateTransition(from, to string) error {
	allowed := map[string][]string{
		StatusReceived:  {StatusReviewed, StatusApproved, StatusArchived},
		StatusReviewed:  {StatusApproved, StatusArchived},
		StatusApproved:  {StatusPublished, StatusArchived},
		StatusPublished: {StatusArchived},
		StatusArchived:  {},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return nil
		}
	}
	return fmt.Errorf("transition from %s to %s is not allowed", from, to)
}

func IsTerminal(status string) bool {
	return status == StatusPublished || status == StatusArchived
}

func IsVisible(status string) bool {
	return status == StatusApproved || status == StatusPublished
}

package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var categoryPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,31}$`)

type ReviewPolicy struct {
	MinimumBody     int
	AllowedRoles    map[string]bool
	RequireCategory bool
}

func DefaultReviewPolicy() ReviewPolicy {
	return ReviewPolicy{MinimumBody: 16, AllowedRoles: map[string]bool{"reviewer": true, "admin": true}, RequireCategory: true}
}

func (p ReviewPolicy) ValidateRecord(record Record) error {
	if len([]rune(record.Body)) < p.MinimumBody {
		return fmt.Errorf("body must contain %d characters", p.MinimumBody)
	}
	if p.RequireCategory && !categoryPattern.MatchString(record.Category) {
		return errorsForCategory(record.Category)
	}
	return nil
}

func errorsForCategory(category string) error {
	if strings.TrimSpace(category) == "" {
		return fmt.Errorf("category is required")
	}
	return fmt.Errorf("category %q must use lowercase letters, digits, or hyphens", category)
}

func (p ReviewPolicy) CanReview(user User) bool {
	if !user.Active {
		return false
	}
	return p.AllowedRoles[user.Role]
}

func (p ReviewPolicy) NeedsRevision(record Record) bool {
	if len(strings.TrimSpace(record.Body)) < p.MinimumBody {
		return true
	}
	return record.Status == StatusReviewed && record.Version > 5
}

type Schedule struct {
	Start    time.Time
	End      time.Time
	Timezone string
}

func (s Schedule) Valid() bool {
	if s.Start.IsZero() || s.End.IsZero() {
		return false
	}
	if !s.End.After(s.Start) {
		return false
	}
	return strings.TrimSpace(s.Timezone) != ""
}

func (s Schedule) Contains(value time.Time) bool {
	if !s.Valid() {
		return false
	}
	return !value.Before(s.Start) && !value.After(s.End)
}

func NormalizeRecord(record Record) Record {
	record.Title = strings.Join(strings.Fields(record.Title), " ")
	record.Body = strings.TrimSpace(record.Body)
	record.Category = NormalizeCategory(record.Category)
	if record.Version < 1 {
		record.Version = 1
	}
	return record
}

func StatusSequence() []string {
	return []string{StatusReceived, StatusReviewed, StatusApproved, StatusPublished, StatusArchived}
}

func CanDisplay(record Record, now time.Time) bool {
	if !IsVisible(record.Status) {
		return false
	}
	if record.PublishedAt.IsZero() {
		return false
	}
	return !record.PublishedAt.After(now)
}

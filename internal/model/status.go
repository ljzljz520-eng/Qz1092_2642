package model

import "time"

type StatusChange struct {
	From   string
	To     string
	At     time.Time
	Reason string
}

func (r *Record) ApplyStatus(change StatusChange) error {
	if err := ValidateTransition(r.Status, change.To); err != nil {
		return err
	}
	if change.At.IsZero() {
		return &timeError{"status change timestamp is required"}
	}
	r.Status = change.To
	r.UpdatedAt = change.At
	r.Version++
	if change.To == StatusPublished {
		r.PublishedAt = change.At
	}
	return nil
}

type timeError struct{ message string }

func (e *timeError) Error() string { return e.message }

func StatusRank(status string) int {
	switch status {
	case StatusReceived:
		return 1
	case StatusReviewed:
		return 2
	case StatusApproved:
		return 3
	case StatusPublished:
		return 4
	case StatusArchived:
		return 5
	default:
		return 0
	}
}

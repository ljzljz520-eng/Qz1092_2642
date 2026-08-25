package model

import (
	"testing"
	"time"
)

func TestRecordValidation(t *testing.T) {
	now := time.Now()
	if err := NewRecord("r1", "A story", "A sufficiently long community story", "food", "u1", now).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := NewRecord("", "", "short", "", "", now).Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestStatusTransition(t *testing.T) {
	r := NewRecord("r1", "A story", "A sufficiently long community story", "food", "u1", time.Now())
	if err := r.ApplyStatus(StatusChange{To: StatusReviewed, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if r.Status != StatusReviewed || r.Version != 2 {
		t.Fatalf("unexpected status %#v", r)
	}
}

package workflow

import (
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"sync"
	"time"
)

type Checkpoint struct {
	RecordID    string
	Step        string
	CompletedAt time.Time
	ActorID     string
	Note        string
}

type Tracker struct {
	mu    sync.RWMutex
	items map[string][]Checkpoint
}

func NewTracker() *Tracker { return &Tracker{items: map[string][]Checkpoint{}} }

func (t *Tracker) Complete(recordID, step, actorID, note string, at time.Time) (Checkpoint, error) {
	if t == nil {
		return Checkpoint{}, errors.New("tracker is not ready")
	}
	if recordID == "" || step == "" {
		return Checkpoint{}, errors.New("record and step are required")
	}
	if at.IsZero() {
		return Checkpoint{}, errors.New("checkpoint time is required")
	}
	checkpoint := Checkpoint{RecordID: recordID, Step: step, CompletedAt: at, ActorID: actorID, Note: note}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, prior := range t.items[recordID] {
		if prior.Step == step {
			return Checkpoint{}, fmt.Errorf("step %s already completed", step)
		}
	}
	t.items[recordID] = append(t.items[recordID], checkpoint)
	return checkpoint, nil
}

func (t *Tracker) Timeline(recordID string) []Checkpoint {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]Checkpoint(nil), t.items[recordID]...)
}

func (t *Tracker) CompleteCount(recordID string) int { return len(t.Timeline(recordID)) }

func (t *Tracker) ReadyFor(step string, record model.Record) bool {
	if step == "receive" {
		return record.Status == model.StatusReceived
	}
	if step == "review" {
		return record.Status == model.StatusReviewed || record.Status == model.StatusApproved
	}
	if step == "publish" {
		return record.Status == model.StatusPublished
	}
	return false
}

func (t *Tracker) Reset(recordID string) { t.mu.Lock(); defer t.mu.Unlock(); delete(t.items, recordID) }

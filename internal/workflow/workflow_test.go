package workflow

import (
	"github.com/charity/storydesk/internal/service"
	"github.com/charity/storydesk/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowDefinitions(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	registry := NewRegistry(service.New(st))
	if err := registry.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(registry.Names()) != 3 {
		t.Fatal("expected three workflows")
	}
}

func TestWorkflowRunner(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	runner := NewRunner(service.New(st))
	record, err := runner.RunIntake(service.Intake{Title: "Shelter", Body: "A sufficiently long shelter story", Category: "housing", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if !runner.CanAdvance(record.Status, "reviewed") {
		t.Fatal("received should advance")
	}
}

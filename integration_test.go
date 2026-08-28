package storydesk

import (
	"github.com/charity/storydesk/internal/service"
	"github.com/charity/storydesk/internal/store"
	"path/filepath"
	"testing"
)

func TestRecordFlow16(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	record, err := svc.ReceiveStory(service.Intake{Title: "Community aid", Body: "A sufficiently long community aid story", Category: "aid", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ReviewStory(record.ID, "reviewer", true, "approved"); err != nil {
		t.Fatal(err)
	}
	final, err := svc.PublishStory(record.ID, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if final.Status != "published" {
		t.Fatalf("expected published, got %s", final.Status)
	}
}

func TestWorkflowOne(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	record, err := service.New(st).ReceiveStory(service.Intake{Title: "Intake", Body: "A sufficiently long intake story", Category: "aid", AuthorID: "u1"})
	if err != nil || record.ID == "" {
		t.Fatalf("record=%#v err=%v", record, err)
	}
}

func TestWorkflowTwo(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := service.New(st).RegisterUser("Reviewer", "reviewer@example.com", "reviewer"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowThree(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.ListRecords(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishedTimestampIsFresh(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	record, err := svc.ReceiveStory(service.Intake{Title: "Fresh time", Body: "A sufficiently long story with fresh time", Category: "aid", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if !record.UpdatedAt.Equal(record.CreatedAt) {
		t.Fatalf("expected intake update time %s, got %s", record.CreatedAt, record.UpdatedAt)
	}
	if _, err = svc.ReviewStory(record.ID, "reviewer", true, "approved"); err != nil {
		t.Fatal(err)
	}
	published, err := svc.PublishStory(record.ID, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if published.UpdatedAt.Before(record.UpdatedAt) {
		t.Fatalf("expected new timestamp, got %s before %s", published.UpdatedAt, record.UpdatedAt)
	}
}

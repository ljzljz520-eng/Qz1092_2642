package service

import (
	"github.com/charity/storydesk/internal/store"
	"path/filepath"
	"testing"
)

func TestReviewAndPublish(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := New(st)
	record, err := svc.ReceiveStory(Intake{Title: "Books", Body: "A long story about a community library", Category: "Education", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ReviewStory(record.ID, "reviewer", true, "looks good"); err != nil {
		t.Fatal(err)
	}
	published, err := svc.PublishStory(record.ID, "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if published.Status != "published" {
		t.Fatalf("status=%s", published.Status)
	}
}

func TestArchiveStory(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := New(st)
	record, err := svc.ReceiveStory(Intake{Title: "Park", Body: "A long story about restoring a local park", Category: "Environment", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	archived, err := svc.ArchiveStory(record.ID, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != "archived" {
		t.Fatalf("status=%s", archived.Status)
	}
}

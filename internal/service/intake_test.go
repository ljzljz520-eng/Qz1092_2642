package service

import (
	"github.com/charity/storydesk/internal/store"
	"path/filepath"
	"testing"
)

func TestReceiveStory(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := New(st)
	record, err := svc.ReceiveStory(Intake{Title: "Warm meals", Body: "Volunteers delivered warm meals to neighbors", Category: "Food", AuthorID: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != "received" || record.Category != "food" {
		t.Fatalf("unexpected record %#v", record)
	}
}

func TestRegisterUser(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	user, err := New(st).RegisterUser("Ada", "ada@example.com", "author")
	if err != nil || user.ID == "" {
		t.Fatalf("user=%#v err=%v", user, err)
	}
}

package store

import (
	"github.com/charity/storydesk/internal/model"
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stories.db")
	now := time.Now().UTC()
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.NewRecord("r1", "Persisted", "A sufficiently long persisted story", "food", "u1", now)
	if err := first.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	got, err := second.GetRecord("r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != record.Title || got.Status != model.StatusReceived {
		t.Fatalf("unexpected record %#v", got)
	}
}

func TestStoreListsRecords(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for i := 0; i < 2; i++ {
		if err := st.SaveRecord(model.NewRecord(string(rune('a'+i)), "Story", "A sufficiently long story body", "food", "u1", time.Now())); err != nil {
			t.Fatal(err)
		}
	}
	items, err := st.ListRecords()
	if err != nil || len(items) != 2 {
		t.Fatalf("items=%d err=%v", len(items), err)
	}
}

package query

import (
	"github.com/charity/storydesk/internal/model"
	"github.com/charity/storydesk/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestFindFilters(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "stories.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if err := st.SaveRecord(model.NewRecord("r1", "Clean water", "A sufficiently long water story", "water", "u1", time.Now())); err != nil {
		t.Fatal(err)
	}
	items, err := New(st).Find(Filters{Search: "water"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%d err=%v", len(items), err)
	}
}

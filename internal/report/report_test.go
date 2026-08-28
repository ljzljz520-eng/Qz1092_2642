package report

import (
	"github.com/charity/storydesk/internal/model"
	"testing"
	"time"
)

func TestBuildReport(t *testing.T) {
	now := time.Now()
	records := []model.Record{model.NewRecord("r1", "One", "A sufficiently long story", "food", "u1", now)}
	snapshot := Build(records, now)
	if snapshot.Total != 1 || snapshot.Visible != 0 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
	if Line(snapshot) == "" {
		t.Fatal("expected line")
	}
}

func TestHealthSnapshot(t *testing.T) {
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	record := model.NewRecord("r-health", "Delayed", "A sufficiently long delayed story", "aid", "u1", now.Add(-49*time.Hour))
	snapshot := BuildHealthSnapshot([]model.Record{record}, now)
	if len(CriticalSignals(snapshot)) != 1 || snapshot.ByLevel[SignalCritical] != 1 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
	if !RequiresAttention(snapshot.Signals[0]) || HealthLine(snapshot) == "" || RenderHealth(snapshot) == "" {
		t.Fatalf("expected attention output")
	}
	if got := CategoriesNeedingAttention(snapshot, 1); len(got) != 1 || got[0] != "aid" {
		t.Fatalf("categories=%v", got)
	}
}

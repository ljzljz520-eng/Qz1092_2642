package report

import (
	"encoding/json"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"time"
)

type Snapshot struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Total       int            `json:"total"`
	Visible     int            `json:"visible"`
	ByStatus    map[string]int `json:"by_status"`
}

func Build(records []model.Record, now time.Time) Snapshot {
	snapshot := Snapshot{GeneratedAt: now, ByStatus: map[string]int{}}
	for _, record := range records {
		snapshot.Total++
		snapshot.ByStatus[record.Status]++
		if model.IsVisible(record.Status) {
			snapshot.Visible++
		}
	}
	return snapshot
}

func Render(snapshot Snapshot) (string, error) {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	return string(data), err
}

func Line(snapshot Snapshot) string {
	return fmt.Sprintf("stories=%d visible=%d", snapshot.Total, snapshot.Visible)
}

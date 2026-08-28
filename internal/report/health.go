package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charity/storydesk/internal/model"
)

type HealthSignal struct {
	RecordID string    `json:"record_id"`
	Status   string    `json:"status"`
	Level    string    `json:"level"`
	Reason   string    `json:"reason"`
	Age      int       `json:"age_hours"`
	At       time.Time `json:"updated_at"`
}

type HealthSnapshot struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Signals     []HealthSignal `json:"signals"`
	ByLevel     map[string]int `json:"by_level"`
	ByStatus    map[string]int `json:"by_status"`
	Oldest      time.Time      `json:"oldest_attention"`
	Categories  map[string]int `json:"attention_categories"`
}

const (
	SignalInfo     = "info"
	SignalWarn     = "warning"
	SignalCritical = "critical"
)

func BuildHealthSnapshot(records []model.Record, now time.Time) HealthSnapshot {
	snapshot := HealthSnapshot{
		GeneratedAt: now,
		Signals:     make([]HealthSignal, 0, len(records)),
		ByLevel:     map[string]int{},
		ByStatus:    map[string]int{},
		Categories:  map[string]int{},
	}
	for _, record := range records {
		snapshot.ByStatus[record.Status]++
		signal := SignalForRecord(record, now)
		if signal.Level == SignalInfo {
			continue
		}
		snapshot.Signals = append(snapshot.Signals, signal)
		snapshot.ByLevel[signal.Level]++
		snapshot.Categories[record.Category]++
		if snapshot.Oldest.IsZero() || signal.At.Before(snapshot.Oldest) {
			snapshot.Oldest = signal.At
		}
	}
	sort.SliceStable(snapshot.Signals, func(i, j int) bool {
		if snapshot.Signals[i].Level != snapshot.Signals[j].Level {
			return signalRank(snapshot.Signals[i].Level) > signalRank(snapshot.Signals[j].Level)
		}
		return snapshot.Signals[i].At.Before(snapshot.Signals[j].At)
	})
	return snapshot
}

func SignalForRecord(record model.Record, now time.Time) HealthSignal {
	age := AgeHours(record.UpdatedAt, now)
	signal := HealthSignal{RecordID: record.ID, Status: record.Status, Level: SignalInfo, Reason: "within service target", Age: age, At: record.UpdatedAt}
	if record.Status == model.StatusReceived && age >= 48 {
		signal.Level, signal.Reason = SignalCritical, "received story has waited two days"
	} else if record.Status == model.StatusReceived && age >= 24 {
		signal.Level, signal.Reason = SignalWarn, "received story needs review"
	} else if record.Status == model.StatusReviewed && age >= 48 {
		signal.Level, signal.Reason = SignalCritical, "reviewed story has no approval decision"
	} else if record.Status == model.StatusReviewed && age >= 24 {
		signal.Level, signal.Reason = SignalWarn, "reviewed story is awaiting approval"
	} else if record.Status == model.StatusApproved && age >= 24 {
		signal.Level, signal.Reason = SignalWarn, "approved story is awaiting publication"
	} else if record.Status == model.StatusPublished && record.PublishedAt.IsZero() {
		signal.Level, signal.Reason = SignalCritical, "published story is missing publication time"
	} else if record.Status == model.StatusArchived && !record.PublishedAt.IsZero() && record.PublishedAt.After(record.UpdatedAt) {
		signal.Level, signal.Reason = SignalCritical, "archive time predates publication time"
	}
	return signal
}

func AgeHours(updated, now time.Time) int {
	if updated.IsZero() || now.Before(updated) {
		return 0
	}
	return int(now.Sub(updated) / time.Hour)
}

func signalRank(level string) int {
	switch level {
	case SignalCritical:
		return 3
	case SignalWarn:
		return 2
	case SignalInfo:
		return 1
	default:
		return 0
	}
}

func RequiresAttention(signal HealthSignal) bool {
	return signal.Level == SignalWarn || signal.Level == SignalCritical
}

func CriticalSignals(snapshot HealthSnapshot) []HealthSignal {
	items := make([]HealthSignal, 0)
	for _, signal := range snapshot.Signals {
		if signal.Level == SignalCritical {
			items = append(items, signal)
		}
	}
	return items
}

func WarningSignals(snapshot HealthSnapshot) []HealthSignal {
	items := make([]HealthSignal, 0)
	for _, signal := range snapshot.Signals {
		if signal.Level == SignalWarn {
			items = append(items, signal)
		}
	}
	return items
}

func CategoriesNeedingAttention(snapshot HealthSnapshot, minimum int) []string {
	if minimum < 1 {
		minimum = 1
	}
	values := make([]string, 0)
	for category, count := range snapshot.Categories {
		if count >= minimum {
			values = append(values, category)
		}
	}
	sort.Strings(values)
	return values
}

func HealthLine(snapshot HealthSnapshot) string {
	return fmt.Sprintf("attention=%d critical=%d warning=%d", len(snapshot.Signals), snapshot.ByLevel[SignalCritical], snapshot.ByLevel[SignalWarn])
}

func RenderHealth(snapshot HealthSnapshot) string {
	lines := []string{HealthLine(snapshot)}
	for _, signal := range snapshot.Signals {
		lines = append(lines, fmt.Sprintf("%s %s %s (%dh): %s", signal.Level, signal.RecordID, signal.Status, signal.Age, signal.Reason))
	}
	return strings.Join(lines, "\n")
}

func StatusTargetHours(status string) int {
	switch status {
	case model.StatusReceived:
		return 24
	case model.StatusReviewed:
		return 24
	case model.StatusApproved:
		return 24
	case model.StatusPublished, model.StatusArchived:
		return 0
	default:
		return 12
	}
}

func IsPastTarget(record model.Record, now time.Time) bool {
	target := StatusTargetHours(record.Status)
	return target > 0 && AgeHours(record.UpdatedAt, now) >= target
}

func AttentionCountByStatus(snapshot HealthSnapshot) map[string]int {
	counts := make(map[string]int)
	for _, signal := range snapshot.Signals {
		counts[signal.Status]++
	}
	return counts
}

func AttentionSummary(snapshot HealthSnapshot) []string {
	counts := AttentionCountByStatus(snapshot)
	statuses := make([]string, 0, len(counts))
	for status := range counts {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	lines := make([]string, 0, len(statuses))
	for _, status := range statuses {
		lines = append(lines, fmt.Sprintf("%s=%d", status, counts[status]))
	}
	return lines
}

func SignalsForStatus(snapshot HealthSnapshot, status string) []HealthSignal {
	items := make([]HealthSignal, 0)
	for _, signal := range snapshot.Signals {
		if status == "" || strings.EqualFold(status, signal.Status) {
			items = append(items, signal)
		}
	}
	return items
}

func SLAReport(records []model.Record, now time.Time) map[string]int {
	result := map[string]int{"within_target": 0, "past_target": 0, "not_tracked": 0}
	for _, record := range records {
		if StatusTargetHours(record.Status) == 0 {
			result["not_tracked"]++
			continue
		}
		if IsPastTarget(record, now) {
			result["past_target"]++
		} else {
			result["within_target"]++
		}
	}
	return result
}

func HasCritical(snapshot HealthSnapshot) bool {
	return snapshot.ByLevel[SignalCritical] > 0
}

func SnapshotTotal(snapshot HealthSnapshot) int {
	total := 0
	for _, count := range snapshot.ByStatus {
		total += count
	}
	return total
}

func SnapshotGeneratedAt(snapshot HealthSnapshot) time.Time {
	return snapshot.GeneratedAt.UTC()
}

func MergeHealthSnapshots(left, right HealthSnapshot) HealthSnapshot {
	merged := HealthSnapshot{GeneratedAt: left.GeneratedAt, Signals: make([]HealthSignal, 0, len(left.Signals)+len(right.Signals)), ByLevel: map[string]int{}, ByStatus: map[string]int{}, Categories: map[string]int{}, Oldest: left.Oldest}
	if right.GeneratedAt.After(merged.GeneratedAt) {
		merged.GeneratedAt = right.GeneratedAt
	}
	merged.Signals = append(merged.Signals, left.Signals...)
	merged.Signals = append(merged.Signals, right.Signals...)
	for _, snapshot := range []HealthSnapshot{left, right} {
		for key, value := range snapshot.ByLevel {
			merged.ByLevel[key] += value
		}
		for key, value := range snapshot.ByStatus {
			merged.ByStatus[key] += value
		}
		for key, value := range snapshot.Categories {
			merged.Categories[key] += value
		}
		if !snapshot.Oldest.IsZero() && (merged.Oldest.IsZero() || snapshot.Oldest.Before(merged.Oldest)) {
			merged.Oldest = snapshot.Oldest
		}
	}
	sort.SliceStable(merged.Signals, func(i, j int) bool { return merged.Signals[i].At.Before(merged.Signals[j].At) })
	return merged
}

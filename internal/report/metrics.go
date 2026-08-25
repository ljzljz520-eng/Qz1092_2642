package report

import (
	"github.com/charity/storydesk/internal/model"
	"sort"
)

type Metric struct {
	Label string
	Value int
}

func StatusMetrics(records []model.Record) []Metric {
	counts := map[string]int{}
	for _, record := range records {
		counts[record.Status]++
	}
	labels := []string{model.StatusReceived, model.StatusReviewed, model.StatusApproved, model.StatusPublished, model.StatusArchived}
	metrics := make([]Metric, 0, len(labels))
	for _, label := range labels {
		metrics = append(metrics, Metric{Label: label, Value: counts[label]})
	}
	return metrics
}

func TopCategories(records []model.Record) []string {
	counts := map[string]int{}
	for _, record := range records {
		counts[record.Category]++
	}
	values := make([]string, 0, len(counts))
	for category := range counts {
		values = append(values, category)
	}
	sort.Slice(values, func(i, j int) bool {
		if counts[values[i]] == counts[values[j]] {
			return values[i] < values[j]
		}
		return counts[values[i]] > counts[values[j]]
	})
	return values
}

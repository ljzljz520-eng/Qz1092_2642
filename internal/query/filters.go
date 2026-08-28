package query

import (
	"github.com/charity/storydesk/internal/model"
	"sort"
	"strings"
)

func SortByRecent(records []model.Record) []model.Record {
	result := append([]model.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	return result
}

func Summarize(records []model.Record) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}

func Categories(records []model.Record) []string {
	seen := map[string]bool{}
	values := make([]string, 0)
	for _, record := range records {
		category := strings.TrimSpace(record.Category)
		if category != "" && !seen[category] {
			seen[category] = true
			values = append(values, category)
		}
	}
	sort.Strings(values)
	return values
}

func StatusLabel(status string) string {
	switch status {
	case model.StatusReceived:
		return "Received"
	case model.StatusReviewed:
		return "Reviewed"
	case model.StatusApproved:
		return "Approved"
	case model.StatusPublished:
		return "Published"
	case model.StatusArchived:
		return "Archived"
	default:
		return "Unknown"
	}
}

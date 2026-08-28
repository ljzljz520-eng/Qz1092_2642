package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"sort"
	"strconv"
)

type CategorySummary struct {
	Category string
	Total    int
	Visible  int
}

func SummarizeCategories(records []model.Record) []CategorySummary {
	groups := map[string]*CategorySummary{}
	for _, record := range records {
		entry := groups[record.Category]
		if entry == nil {
			entry = &CategorySummary{Category: record.Category}
			groups[record.Category] = entry
		}
		entry.Total++
		if model.IsVisible(record.Status) {
			entry.Visible++
		}
	}
	result := make([]CategorySummary, 0, len(groups))
	for _, entry := range groups {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Category < result[j].Category })
	return result
}

func CSV(snapshot Snapshot, summaries []CategorySummary) (string, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if err := writer.Write([]string{"metric", "value"}); err != nil {
		return "", err
	}
	if err := writer.Write([]string{"total", strconv.Itoa(snapshot.Total)}); err != nil {
		return "", err
	}
	if err := writer.Write([]string{"visible", strconv.Itoa(snapshot.Visible)}); err != nil {
		return "", err
	}
	if err := writer.Write([]string{"category", "total", "visible"}); err != nil {
		return "", err
	}
	for _, summary := range summaries {
		if err := writer.Write([]string{summary.Category, strconv.Itoa(summary.Total), strconv.Itoa(summary.Visible)}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func Describe(snapshot Snapshot) string {
	return fmt.Sprintf("%d stories tracked; %d are visible; generated %s", snapshot.Total, snapshot.Visible, snapshot.GeneratedAt.Format("2006-01-02"))
}

func SortMetrics(metrics []Metric) []Metric {
	result := append([]Metric(nil), metrics...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Value == result[j].Value {
			return result[i].Label < result[j].Label
		}
		return result[i].Value > result[j].Value
	})
	return result
}

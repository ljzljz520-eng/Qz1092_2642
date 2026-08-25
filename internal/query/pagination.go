package query

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"sort"
	"time"
)

type Page struct {
	Items    []model.Record
	Page     int
	PageSize int
	Total    int
	HasNext  bool
}

func Paginate(records []model.Record, page, pageSize int) (Page, error) {
	if page < 1 {
		return Page{}, errors.New("page must be positive")
	}
	if pageSize < 1 || pageSize > 100 {
		return Page{}, errors.New("page size must be between 1 and 100")
	}
	ordered := SortByRecent(records)
	start := (page - 1) * pageSize
	if start >= len(ordered) {
		return Page{Items: []model.Record{}, Page: page, PageSize: pageSize, Total: len(ordered)}, nil
	}
	end := start + pageSize
	if end > len(ordered) {
		end = len(ordered)
	}
	return Page{Items: ordered[start:end], Page: page, PageSize: pageSize, Total: len(ordered), HasNext: end < len(ordered)}, nil
}

func UpdatedBetween(records []model.Record, start, end time.Time) []model.Record {
	filtered := make([]model.Record, 0)
	for _, record := range records {
		if (start.IsZero() || !record.UpdatedAt.Before(start)) && (end.IsZero() || record.UpdatedAt.Before(end)) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func GroupByCategory(records []model.Record) map[string][]model.Record {
	groups := map[string][]model.Record{}
	for _, record := range records {
		groups[record.Category] = append(groups[record.Category], record)
	}
	for category := range groups {
		sort.Slice(groups[category], func(i, j int) bool { return groups[category][i].UpdatedAt.After(groups[category][j].UpdatedAt) })
	}
	return groups
}

func VisibleAt(records []model.Record, now time.Time) []model.Record {
	visible := make([]model.Record, 0)
	for _, record := range records {
		if model.CanDisplay(record, now) {
			visible = append(visible, record)
		}
	}
	return SortByRecent(visible)
}

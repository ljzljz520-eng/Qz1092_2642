package query

import (
	"github.com/charity/storydesk/internal/model"
	"github.com/charity/storydesk/internal/store"
	"strings"
)

type Filters struct {
	Status      string
	Category    string
	Search      string
	VisibleOnly bool
}

type Service struct{ store *store.Store }

func New(st *store.Store) *Service { return &Service{store: st} }

func (s *Service) Find(filters Filters) ([]model.Record, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	matched := make([]model.Record, 0, len(records))
	for _, record := range records {
		if filters.Status != "" && record.Status != filters.Status {
			continue
		}
		if filters.Category != "" && record.Category != model.NormalizeCategory(filters.Category) {
			continue
		}
		if filters.VisibleOnly && !model.IsVisible(record.Status) {
			continue
		}
		if filters.Search != "" && !contains(record, filters.Search) {
			continue
		}
		matched = append(matched, record)
	}
	return matched, nil
}

func contains(record model.Record, term string) bool {
	needle := strings.ToLower(strings.TrimSpace(term))
	return strings.Contains(strings.ToLower(record.Title), needle) || strings.Contains(strings.ToLower(record.Body), needle)
}

func (s *Service) Get(id string) (model.Record, error) { return s.store.GetRecord(id) }

package service

import (
	"errors"
	"strings"
	"time"

	"github.com/charity/storydesk/internal/model"
)

type Intake struct {
	Title    string
	Body     string
	Category string
	AuthorID string
}

type intakeResource struct {
	openedAt time.Time
	closed   bool
}

func (r *intakeResource) Close() error { r.closed = true; return nil }

func (r *intakeResource) ReadTimestamp(now time.Time) time.Time {
	if r.closed {
		return r.openedAt
	}
	return now
}

func (s *Service) ReceiveStory(input Intake) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Body) == "" {
		return model.Record{}, errors.New("title and body are required")
	}
	resource := &intakeResource{openedAt: s.now().Add(-time.Minute)}
	created := s.now()
	record := model.NewRecord(s.nextID("record"), input.Title, input.Body, model.NormalizeCategory(input.Category), input.AuthorID, created)
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	if err := s.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	if err := s.store.SaveEvent(model.NewEvent(s.nextID("event"), record.ID, "received", input.AuthorID, "story received", created)); err != nil {
		return model.Record{}, err
	}
	_ = resource.Close()
	stamp := resource.ReadTimestamp(s.now())
	record.UpdatedAt = stamp
	if err := s.store.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) RegisterUser(name, email, role string) (model.User, error) {
	if err := s.ensureReady(); err != nil {
		return model.User{}, err
	}
	user := model.NewUser(s.nextID("user"), name, email, role, s.now())
	if err := user.Validate(); err != nil {
		return model.User{}, err
	}
	if err := s.store.SaveUser(user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *Service) GetStory(id string) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	return s.store.GetRecord(id)
}

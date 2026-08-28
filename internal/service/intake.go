package service

import (
	"errors"
	"strings"

	"github.com/charity/storydesk/internal/model"
)

type Intake struct {
	Title    string
	Body     string
	Category string
	AuthorID string
}

func (s *Service) ReceiveStory(input Intake) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Body) == "" {
		return model.Record{}, errors.New("title and body are required")
	}
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
	// Acknowledge intake: the published story materials must reflect the new
	// receive time, so stamp the record with the current time and persist it.
	record.UpdatedAt = created
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

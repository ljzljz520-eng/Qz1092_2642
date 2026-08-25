package service

import (
	"errors"
	"fmt"

	"github.com/charity/storydesk/internal/model"
)

func (s *Service) ReviewStory(id, reviewerID string, approve bool, reason string) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	if reviewerID == "" {
		return model.Record{}, errors.New("reviewer is required")
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	target := model.StatusReviewed
	if approve {
		target = model.StatusApproved
	}
	change := model.StatusChange{From: record.Status, To: target, At: s.now(), Reason: reason}
	if err := record.ApplyStatus(change); err != nil {
		return model.Record{}, err
	}
	if err := s.store.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	message := fmt.Sprintf("story moved to %s", target)
	if err := s.store.SaveEvent(model.NewEvent(s.nextID("event"), id, "review", reviewerID, message, change.At)); err != nil {
		return model.Record{}, err
	}
	return record, s.store.SaveAudit(model.NewAudit(s.nextID("audit"), "review", id, reviewerID, reason, change.At))
}

func (s *Service) ArchiveStory(id, actorID string) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if model.IsTerminal(record.Status) && record.Status != model.StatusPublished {
		return record, errors.New("story is already archived")
	}
	if err := record.ApplyStatus(model.StatusChange{From: record.Status, To: model.StatusArchived, At: s.now(), Reason: "archive request"}); err != nil {
		return model.Record{}, err
	}
	if err := s.store.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, s.store.SaveAudit(model.NewAudit(s.nextID("audit"), "archive", id, actorID, "archive request", record.UpdatedAt))
}

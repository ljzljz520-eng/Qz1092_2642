package service

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
)

func (s *Service) PublishStory(id, actorID string) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusApproved {
		return model.Record{}, errors.New("only approved stories can be published")
	}
	change := model.StatusChange{From: record.Status, To: model.StatusPublished, At: s.now(), Reason: "publish request"}
	if err := record.ApplyStatus(change); err != nil {
		return model.Record{}, err
	}
	if err := s.store.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	if err := s.store.SaveEvent(model.NewEvent(s.nextID("event"), id, "published", actorID, "story published", change.At)); err != nil {
		return model.Record{}, err
	}
	return record, s.store.SaveAudit(model.NewAudit(s.nextID("audit"), "publish", id, actorID, "publish request", change.At))
}

func (s *Service) ChangeStatus(id, actorID, target string) (model.Record, error) {
	if err := s.ensureReady(); err != nil {
		return model.Record{}, err
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if target == model.StatusArchived {
		return s.ArchiveStory(id, actorID)
	}
	if target == model.StatusPublished {
		return s.PublishStory(id, actorID)
	}
	if target != model.StatusReviewed && target != model.StatusApproved {
		return model.Record{}, errors.New("unsupported target status")
	}
	if err := record.ApplyStatus(model.StatusChange{From: record.Status, To: target, At: s.now(), Reason: "status request"}); err != nil {
		return model.Record{}, err
	}
	if err := s.store.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, s.store.SaveAudit(model.NewAudit(s.nextID("audit"), "status", id, actorID, target, record.UpdatedAt))
}

package store

import (
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
	"sort"
	"time"
)

type Transaction struct {
	Records []model.Record
	Events  []model.Event
	Audits  []model.Audit
}

func (s *Store) SaveTransaction(transaction Transaction) error {
	if len(transaction.Records) == 0 {
		return errors.New("transaction requires a record")
	}
	return s.withWrite(func(tx *bbolt.Tx) error {
		for _, record := range transaction.Records {
			if err := record.Validate(); err != nil {
				return err
			}
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := tx.Bucket(recordBucket).Put([]byte(record.ID), data); err != nil {
				return err
			}
		}
		for _, event := range transaction.Events {
			if event.ID == "" {
				return errors.New("event id is required")
			}
			data, err := encode(event)
			if err != nil {
				return err
			}
			if err := tx.Bucket(eventBucket).Put([]byte(event.ID), data); err != nil {
				return err
			}
		}
		for _, audit := range transaction.Audits {
			if audit.ID == "" {
				return errors.New("audit id is required")
			}
			data, err := encode(audit)
			if err != nil {
				return err
			}
			if err := tx.Bucket(auditBucket).Put([]byte(audit.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Timeline(recordID string) ([]model.Event, []model.Audit, error) {
	events, err := s.ListEvents(recordID)
	if err != nil {
		return nil, nil, err
	}
	audits, err := s.ListAudits(recordID)
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt.Before(events[j].CreatedAt) })
	sort.Slice(audits, func(i, j int) bool { return audits[i].CreatedAt.Before(audits[j].CreatedAt) })
	return events, audits, nil
}

func (s *Store) Snapshot(recordID string, exportedAt time.Time) (model.RecordEnvelope, error) {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return model.RecordEnvelope{}, err
	}
	events, audits, err := s.Timeline(recordID)
	if err != nil {
		return model.RecordEnvelope{}, err
	}
	if exportedAt.IsZero() {
		exportedAt = time.Now().UTC()
	}
	return model.RecordEnvelope{Record: record, Events: events, Audits: audits, ExportedAt: exportedAt}, nil
}

func (s *Store) RequireVersion(id string, version int) (model.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Version != version {
		return model.Record{}, fmt.Errorf("record version mismatch: want %d got %d", version, record.Version)
	}
	return record, nil
}

func (s *Store) UpdateIfVersion(record model.Record, expected int) error {
	current, err := s.RequireVersion(record.ID, expected)
	if err != nil {
		return err
	}
	if record.UpdatedAt.Before(current.UpdatedAt) {
		return errors.New("updated timestamp cannot move backwards")
	}
	return s.UpdateRecord(record)
}

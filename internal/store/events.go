package store

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveEvent(event model.Event) error {
	if event.ID == "" || event.RecordID == "" {
		return errors.New("event identity is required")
	}
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(eventBucket).Put([]byte(event.ID), data) })
}

func (s *Store) ListEvents(recordID string) ([]model.Event, error) {
	events := make([]model.Event, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(eventBucket).ForEach(func(_, value []byte) error {
			var event model.Event
			if err := decode(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	return events, err
}

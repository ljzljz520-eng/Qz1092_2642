package store

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
	"sort"
	"time"
)

type Health struct {
	Open    bool
	Records int
	Users   int
	Events  int
	Audits  int
}

func (s *Store) Health() (Health, error) {
	health := Health{}
	err := s.withRead(func(tx *bbolt.Tx) error {
		health.Open = true
		health.Records = tx.Bucket(recordBucket).Stats().KeyN
		health.Users = tx.Bucket(userBucket).Stats().KeyN
		health.Events = tx.Bucket(eventBucket).Stats().KeyN
		health.Audits = tx.Bucket(auditBucket).Stats().KeyN
		return nil
	})
	return health, err
}

func (s *Store) PurgeBefore(cutoff time.Time) (int, error) {
	if cutoff.IsZero() {
		return 0, errors.New("cutoff is required")
	}
	removed := 0
	err := s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, value []byte) error {
			var record model.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.UpdatedAt.Before(cutoff) && model.IsTerminal(record.Status) {
				keys = append(keys, append([]byte(nil), key...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
			removed++
		}
		return nil
	})
	return removed, err
}

func (s *Store) RecordsByStatus(status string) ([]model.Record, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	filtered := records[:0]
	for _, record := range records {
		if record.Status == status {
			filtered = append(filtered, record)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt) })
	return filtered, nil
}

func (s *Store) TouchRecord(id string, when time.Time) (model.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if when.IsZero() {
		return model.Record{}, errors.New("touch time is required")
	}
	record.UpdatedAt = when
	record.Version++
	if err := s.UpdateRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

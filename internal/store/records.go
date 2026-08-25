package store

import (
	"errors"
	"sort"

	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(record model.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(recordBucket).Put([]byte(record.ID), data) })
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.withRead(func(tx *bbolt.Tx) error {
		value := tx.Bucket(recordBucket).Get([]byte(id))
		if value == nil {
			return errors.New("record not found")
		}
		return decode(value, &record)
	})
	return record, err
}

func (s *Store) ListRecords() ([]model.Record, error) {
	items := make([]model.Record, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordBucket).ForEach(func(_, value []byte) error {
			var record model.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			items = append(items, record)
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.Before(items[j].UpdatedAt) })
	return items, err
}

func (s *Store) UpdateRecord(record model.Record) error {
	if record.ID == "" {
		return errors.New("record id is required")
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		if bucket.Get([]byte(record.ID)) == nil {
			return errors.New("record not found")
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (s *Store) DeleteRecord(id string) error {
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(recordBucket).Delete([]byte(id)) })
}

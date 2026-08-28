package store

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveAudit(audit model.Audit) error {
	if audit.ID == "" || audit.Action == "" {
		return errors.New("audit identity is required")
	}
	data, err := encode(audit)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put([]byte(audit.ID), data) })
}

func (s *Store) ListAudits(recordID string) ([]model.Audit, error) {
	audits := make([]model.Audit, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, value []byte) error {
			var audit model.Audit
			if err := decode(value, &audit); err != nil {
				return err
			}
			if recordID == "" || audit.RecordID == recordID {
				audits = append(audits, audit)
			}
			return nil
		})
	})
	return audits, err
}

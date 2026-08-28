package store

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveUser(user model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	data, err := encode(user)
	if err != nil {
		return err
	}
	return s.withWrite(func(tx *bbolt.Tx) error { return tx.Bucket(userBucket).Put([]byte(user.ID), data) })
}

func (s *Store) GetUser(id string) (model.User, error) {
	var user model.User
	err := s.withRead(func(tx *bbolt.Tx) error {
		value := tx.Bucket(userBucket).Get([]byte(id))
		if value == nil {
			return errors.New("user not found")
		}
		return decode(value, &user)
	})
	return user, err
}

func (s *Store) ListUsers() ([]model.User, error) {
	users := make([]model.User, 0)
	err := s.withRead(func(tx *bbolt.Tx) error {
		return tx.Bucket(userBucket).ForEach(func(_, value []byte) error {
			var user model.User
			if err := decode(value, &user); err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
	})
	return users, err
}

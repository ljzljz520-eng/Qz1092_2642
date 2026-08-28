package service

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/charity/storydesk/internal/store"
)

type Service struct {
	store    *store.Store
	clock    func() time.Time
	sequence uint64
}

func New(st *store.Store) *Service {
	return &Service{store: st, clock: time.Now}
}

func (s *Service) withClock(clock func() time.Time) {
	if clock != nil {
		s.clock = clock
	}
}

func (s *Service) nextID(prefix string) string {
	n := atomic.AddUint64(&s.sequence, 1)
	return fmt.Sprintf("%s-%d", prefix, n)
}

func (s *Service) ensureReady() error {
	if s == nil || s.store == nil {
		return errors.New("service is not ready")
	}
	return nil
}

func (s *Service) Store() *store.Store { return s.store }

func (s *Service) now() time.Time { return s.clock().UTC() }

package service

import (
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"strings"
	"sync"
	"time"
)

type Notification struct {
	ID        string
	RecordID  string
	Recipient string
	Channel   string
	Subject   string
	Body      string
	SentAt    time.Time
}

type Notifier struct {
	mu     sync.Mutex
	outbox []Notification
}

func NewNotifier() *Notifier { return &Notifier{outbox: make([]Notification, 0)} }

func (n *Notifier) Queue(record model.Record, recipient, channel string, now time.Time) (Notification, error) {
	if n == nil {
		return Notification{}, errors.New("notifier is not ready")
	}
	if strings.TrimSpace(recipient) == "" {
		return Notification{}, errors.New("recipient is required")
	}
	if channel != "email" && channel != "sms" && channel != "dashboard" {
		return Notification{}, errors.New("unsupported notification channel")
	}
	if now.IsZero() {
		return Notification{}, errors.New("notification time is required")
	}
	notification := Notification{ID: fmt.Sprintf("notice-%d", now.UnixNano()), RecordID: record.ID, Recipient: recipient, Channel: channel, Subject: "公益故事状态更新", Body: fmt.Sprintf("故事 %s 当前状态为 %s", record.Title, record.Status), SentAt: now}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.outbox = append(n.outbox, notification)
	return notification, nil
}

func (n *Notifier) Pending() []Notification {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]Notification(nil), n.outbox...)
}

func (n *Notifier) CountFor(recordID string) int {
	count := 0
	for _, notification := range n.Pending() {
		if notification.RecordID == recordID {
			count++
		}
	}
	return count
}

func (n *Notifier) Flush(now time.Time) []Notification {
	n.mu.Lock()
	defer n.mu.Unlock()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	flushed := make([]Notification, len(n.outbox))
	copy(flushed, n.outbox)
	n.outbox = n.outbox[:0]
	for i := range flushed {
		flushed[i].SentAt = now
	}
	return flushed
}

func (s *Service) PublishWithNotice(id, actorID, recipient, channel string, notifier *Notifier) (model.Record, Notification, error) {
	record, err := s.PublishStory(id, actorID)
	if err != nil {
		return model.Record{}, Notification{}, err
	}
	if notifier == nil {
		return model.Record{}, Notification{}, errors.New("notifier is required")
	}
	notification, err := notifier.Queue(record, recipient, channel, s.now())
	if err != nil {
		return model.Record{}, Notification{}, err
	}
	return record, notification, nil
}

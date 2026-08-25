package service

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"sync"
	"time"
)

type WorkItem struct {
	ID        string
	RecordID  string
	Kind      string
	Priority  int
	CreatedAt time.Time
	Attempts  int
	index     int
}

type workHeap []*WorkItem

func (h workHeap) Len() int { return len(h) }
func (h workHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].CreatedAt.Before(h[j].CreatedAt)
	}
	return h[i].Priority > h[j].Priority
}
func (h workHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *workHeap) Push(value any) {
	item := value.(*WorkItem)
	item.index = len(*h)
	*h = append(*h, item)
}
func (h *workHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[:n-1]
	return item
}

type Queue struct {
	mu     sync.Mutex
	items  workHeap
	closed bool
}

func NewQueue() *Queue { return &Queue{items: make(workHeap, 0)} }

func (q *Queue) Enqueue(record model.Record, kind string, priority int, now time.Time) (WorkItem, error) {
	if q == nil {
		return WorkItem{}, errors.New("queue is not ready")
	}
	if record.ID == "" || kind == "" {
		return WorkItem{}, errors.New("record and kind are required")
	}
	if now.IsZero() {
		return WorkItem{}, errors.New("queue time is required")
	}
	if priority < 0 || priority > 10 {
		return WorkItem{}, errors.New("priority must be between 0 and 10")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return WorkItem{}, errors.New("queue is closed")
	}
	item := &WorkItem{ID: fmt.Sprintf("work-%d", now.UnixNano()), RecordID: record.ID, Kind: kind, Priority: priority, CreatedAt: now}
	heap.Push(&q.items, item)
	return *item, nil
}

func (q *Queue) Next() (WorkItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed || len(q.items) == 0 {
		return WorkItem{}, false
	}
	item := heap.Pop(&q.items).(*WorkItem)
	item.Attempts++
	return *item, true
}

func (q *Queue) Retry(item WorkItem, now time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return errors.New("queue is closed")
	}
	if item.Attempts >= 3 {
		return errors.New("work item exceeded retry limit")
	}
	item.CreatedAt = now
	heap.Push(&q.items, &item)
	return nil
}

func (q *Queue) Close() { q.mu.Lock(); defer q.mu.Unlock(); q.closed = true; q.items = nil }

func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }

func (q *Queue) Drain() []WorkItem {
	items := make([]WorkItem, 0)
	for {
		item, ok := q.Next()
		if !ok {
			return items
		}
		items = append(items, item)
	}
}

func (s *Service) EnqueuePublication(queue *Queue, record model.Record) (WorkItem, error) {
	if queue == nil {
		return WorkItem{}, errors.New("queue is required")
	}
	if record.Status != model.StatusApproved {
		return WorkItem{}, errors.New("only approved stories enter publication queue")
	}
	return queue.Enqueue(record, "publish", 5, s.now())
}

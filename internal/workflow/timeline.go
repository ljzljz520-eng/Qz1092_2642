package workflow

import (
	"github.com/charity/storydesk/internal/model"
	"sort"
	"time"
)

func SortEvents(events []model.Event) []model.Event {
	copyEvents := append([]model.Event(nil), events...)
	sort.SliceStable(copyEvents, func(i, j int) bool { return copyEvents[i].CreatedAt.Before(copyEvents[j].CreatedAt) })
	return copyEvents
}

func LatestEvent(events []model.Event) (model.Event, bool) {
	if len(events) == 0 {
		return model.Event{}, false
	}
	ordered := SortEvents(events)
	return ordered[len(ordered)-1], true
}

func IsRecent(value time.Time, now time.Time) bool {
	if value.IsZero() {
		return false
	}
	return now.Sub(value) <= 24*time.Hour && !value.After(now)
}

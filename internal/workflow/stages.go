package workflow

import (
	"errors"
	"github.com/charity/storydesk/internal/model"
	"github.com/charity/storydesk/internal/service"
)

type Runner struct{ service *service.Service }

func NewRunner(svc *service.Service) *Runner { return &Runner{service: svc} }

func (r *Runner) RunIntake(input service.Intake) (model.Record, error) {
	if r == nil || r.service == nil {
		return model.Record{}, errors.New("runner is not ready")
	}
	return r.service.ReceiveStory(input)
}

func (r *Runner) RunReview(id, reviewerID string, approve bool) (model.Record, error) {
	if r == nil || r.service == nil {
		return model.Record{}, errors.New("runner is not ready")
	}
	return r.service.ReviewStory(id, reviewerID, approve, "workflow review")
}

func (r *Runner) RunPublication(id, actorID string) (model.Record, error) {
	if r == nil || r.service == nil {
		return model.Record{}, errors.New("runner is not ready")
	}
	return r.service.PublishStory(id, actorID)
}

func (r *Runner) CanAdvance(status, target string) bool {
	for _, value := range AllowedForStatus(status) {
		if value == target {
			return true
		}
	}
	return false
}

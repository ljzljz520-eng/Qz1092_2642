package workflow

import (
	"errors"
	"fmt"
	"github.com/charity/storydesk/internal/model"
	"github.com/charity/storydesk/internal/service"
)

type Definition struct {
	Name  string
	Steps []string
}

type Registry struct {
	service     *service.Service
	definitions map[string]Definition
}

func NewRegistry(svc *service.Service) *Registry {
	definitions := map[string]Definition{
		"intake":      {Name: "intake", Steps: []string{"receive", "validate", "save", "display"}},
		"review":      {Name: "review", Steps: []string{"register", "review", "archive", "query"}},
		"publication": {Name: "publication", Steps: []string{"submit", "process", "notify", "track"}},
	}
	return &Registry{service: svc, definitions: definitions}
}

func (r *Registry) Definition(name string) (Definition, error) {
	definition, ok := r.definitions[name]
	if !ok {
		return Definition{}, fmt.Errorf("workflow %s not found", name)
	}
	return definition, nil
}

func (r *Registry) Names() []string {
	return []string{"intake", "review", "publication"}
}

func (r *Registry) Service() *service.Service { return r.service }

func (r *Registry) Validate() error {
	if r == nil || r.service == nil {
		return errors.New("workflow registry is not ready")
	}
	for _, name := range r.Names() {
		definition, err := r.Definition(name)
		if err != nil {
			return err
		}
		if len(definition.Steps) < 4 {
			return errors.New("workflow must have four steps")
		}
	}
	return nil
}

func AllowedForStatus(status string) []string {
	switch status {
	case model.StatusReceived:
		return []string{model.StatusReviewed, model.StatusArchived}
	case model.StatusReviewed:
		return []string{model.StatusApproved, model.StatusArchived}
	case model.StatusApproved:
		return []string{model.StatusPublished, model.StatusArchived}
	default:
		return []string{model.StatusArchived}
	}
}

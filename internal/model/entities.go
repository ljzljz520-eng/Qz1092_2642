package model

import (
	"errors"
	"strings"
	"time"
)

type Record struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	AuthorID    string    `json:"author_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	Version     int       `json:"version"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Kind      string    `json:"kind"`
	ActorID   string    `json:"actor_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Audit struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	RecordID  string    `json:"record_id"`
	ActorID   string    `json:"actor_id"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	StatusReceived  = "received"
	StatusReviewed  = "reviewed"
	StatusApproved  = "approved"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

func NewRecord(id, title, body, category, authorID string, now time.Time) Record {
	return Record{ID: id, Title: strings.TrimSpace(title), Body: strings.TrimSpace(body), Category: strings.TrimSpace(category), Status: StatusReceived, AuthorID: authorID, CreatedAt: now, UpdatedAt: now, Version: 1}
}

func NewUser(id, name, email, role string, now time.Time) User {
	return User{ID: id, Name: strings.TrimSpace(name), Email: strings.TrimSpace(email), Role: role, Active: true, CreatedAt: now}
}

func NewEvent(id, recordID, kind, actorID, message string, now time.Time) Event {
	return Event{ID: id, RecordID: recordID, Kind: kind, ActorID: actorID, Message: message, CreatedAt: now}
}

func NewAudit(id, action, recordID, actorID, details string, now time.Time) Audit {
	return Audit{ID: id, Action: action, RecordID: recordID, ActorID: actorID, Details: details, CreatedAt: now}
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("record title is required")
	}
	if len([]rune(r.Body)) < 16 {
		return errors.New("record body must contain at least 16 characters")
	}
	if strings.TrimSpace(r.Category) == "" {
		return errors.New("record category is required")
	}
	if r.AuthorID == "" {
		return errors.New("record author is required")
	}
	return nil
}

func (u User) Validate() error {
	if u.ID == "" || u.Name == "" {
		return errors.New("user identity is required")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("user email is invalid")
	}
	if u.Role != "author" && u.Role != "reviewer" && u.Role != "admin" {
		return errors.New("user role is invalid")
	}
	return nil
}

package model

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"
)

type RecordEnvelope struct {
	Record     Record    `json:"record"`
	Events     []Event   `json:"events"`
	Audits     []Audit   `json:"audits"`
	ExportedAt time.Time `json:"exported_at"`
}

func MarshalRecord(record Record) ([]byte, error) {
	if err := record.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(record)
}

func UnmarshalRecord(data []byte) (Record, error) {
	var record Record
	if len(data) == 0 {
		return record, errors.New("record data is empty")
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return record, err
	}
	return NormalizeRecord(record), nil
}

func MarshalEnvelope(envelope RecordEnvelope) ([]byte, error) {
	if envelope.Record.ID == "" {
		return nil, errors.New("envelope record is required")
	}
	ordered := append([]Event(nil), envelope.Events...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].CreatedAt.Before(ordered[j].CreatedAt) })
	envelope.Events = ordered
	return json.MarshalIndent(envelope, "", "  ")
}

func UnmarshalEnvelope(data []byte) (RecordEnvelope, error) {
	var envelope RecordEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return envelope, err
	}
	if envelope.Record.ID == "" {
		return envelope, errors.New("envelope record is required")
	}
	return envelope, nil
}

func CanonicalStatus(value string) string {
	status := strings.ToLower(strings.TrimSpace(value))
	for _, candidate := range StatusSequence() {
		if status == candidate {
			return candidate
		}
	}
	return ""
}

func ParseTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, errors.New("time is required")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func FormatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func CloneRecord(record Record) Record { return record }

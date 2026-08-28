package service

import (
	"bufio"
	"encoding/csv"
	"errors"
	"github.com/charity/storydesk/internal/model"
	"io"
	"strings"
)

type ImportResult struct {
	Accepted int
	Rejected int
	Errors   []string
}

func (s *Service) ImportCSV(reader io.Reader, actorID string) (ImportResult, error) {
	if err := s.ensureReady(); err != nil {
		return ImportResult{}, err
	}
	if reader == nil {
		return ImportResult{}, errors.New("csv reader is required")
	}
	parser := csv.NewReader(bufio.NewReader(reader))
	parser.FieldsPerRecord = -1
	result := ImportResult{Errors: make([]string, 0)}
	line := 0
	for {
		line++
		row, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		if line == 1 && len(row) > 0 && strings.EqualFold(strings.TrimSpace(row[0]), "title") {
			continue
		}
		if len(row) < 3 {
			result.Rejected++
			result.Errors = append(result.Errors, "line has fewer than three columns")
			continue
		}
		input := Intake{Title: row[0], Body: row[1], Category: row[2], AuthorID: actorID}
		if _, err := s.ReceiveStory(input); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.Accepted++
	}
	return result, nil
}

func ExportRecords(records []model.Record) string {
	var builder strings.Builder
	builder.WriteString("id,title,category,status,updated_at\n")
	for _, record := range records {
		builder.WriteString(strings.Join([]string{record.ID, csvCell(record.Title), csvCell(record.Category), record.Status, record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}, ","))
		builder.WriteByte('\n')
	}
	return builder.String()
}

func csvCell(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}

func ValidateImportRow(row []string) error {
	if len(row) != 3 {
		return errors.New("row must have title, body, and category")
	}
	if strings.TrimSpace(row[0]) == "" || strings.TrimSpace(row[1]) == "" {
		return errors.New("title and body are required")
	}
	return nil
}

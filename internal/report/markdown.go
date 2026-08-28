package report

import (
	"github.com/charity/storydesk/internal/model"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TimelineLine struct {
	At     time.Time
	Label  string
	Detail string
}

func BuildTimeline(record model.Record, events []model.Event, audits []model.Audit) []TimelineLine {
	lines := make([]TimelineLine, 0, 1+len(events)+len(audits))
	lines = append(lines, TimelineLine{At: record.CreatedAt, Label: "received", Detail: record.Title})
	for _, event := range events {
		lines = append(lines, TimelineLine{At: event.CreatedAt, Label: event.Kind, Detail: event.Message})
	}
	for _, audit := range audits {
		lines = append(lines, TimelineLine{At: audit.CreatedAt, Label: audit.Action, Detail: audit.Details})
	}
	sort.SliceStable(lines, func(i, j int) bool { return lines[i].At.Before(lines[j].At) })
	return lines
}

func RenderTimeline(lines []TimelineLine) string {
	var builder strings.Builder
	for _, line := range lines {
		builder.WriteString("- ")
		builder.WriteString(line.At.UTC().Format(time.RFC3339))
		builder.WriteString(" **")
		builder.WriteString(line.Label)
		builder.WriteString("** ")
		builder.WriteString(strings.ReplaceAll(line.Detail, "\n", " "))
		builder.WriteByte('\n')
	}
	return builder.String()
}

func RenderRecord(record model.Record, lines []TimelineLine) string {
	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(record.Title)
	builder.WriteByte('\n')
	builder.WriteString("\n")
	builder.WriteString(record.Body)
	builder.WriteString("\n\n")
	builder.WriteString("Status: ")
	builder.WriteString(record.Status)
	builder.WriteString("\n\n")
	builder.WriteString(RenderTimeline(lines))
	return builder.String()
}

func StatusTable(metrics []Metric) string {
	var builder strings.Builder
	builder.WriteString("| status | count |\n| --- | ---: |\n")
	for _, metric := range SortMetrics(metrics) {
		builder.WriteString("| ")
		builder.WriteString(metric.Label)
		builder.WriteString(" | ")
		builder.WriteString(strconv.Itoa(metric.Value))
		builder.WriteString(" |\n")
	}
	return builder.String()
}

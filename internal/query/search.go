package query

import (
	"github.com/charity/storydesk/internal/model"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type SearchTerm struct {
	Raw     string
	Words   []string
	Pattern *regexp.Regexp
}

func ParseSearch(value string) SearchTerm {
	trimmed := strings.TrimSpace(value)
	words := strings.FieldsFunc(strings.ToLower(trimmed), func(r rune) bool { return unicode.IsSpace(r) || r == ',' || r == ';' })
	pattern, _ := regexp.Compile(strings.Join(words, ".*"))
	return SearchTerm{Raw: trimmed, Words: words, Pattern: pattern}
}

func MatchTerm(record model.Record, term SearchTerm) bool {
	if len(term.Words) == 0 {
		return true
	}
	text := strings.ToLower(record.Title + " " + record.Body + " " + record.Category)
	if term.Pattern != nil && term.Pattern.MatchString(text) {
		return true
	}
	for _, word := range term.Words {
		if !strings.Contains(text, word) {
			return false
		}
	}
	return true
}

func Score(record model.Record, term SearchTerm) int {
	if !MatchTerm(record, term) {
		return 0
	}
	text := strings.ToLower(record.Title + " " + record.Body)
	score := 1
	for _, word := range term.Words {
		if strings.Contains(strings.ToLower(record.Title), word) {
			score += 5
		}
		if strings.Contains(text, word) {
			score++
		}
	}
	if model.IsVisible(record.Status) {
		score += 2
	}
	return score
}

func Rank(records []model.Record, value string) []model.Record {
	term := ParseSearch(value)
	result := make([]model.Record, 0)
	for _, record := range records {
		if MatchTerm(record, term) {
			result = append(result, record)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		left, right := Score(result[i], term), Score(result[j], term)
		if left == right {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		}
		return left > right
	})
	return result
}

func ExtractHighlights(record model.Record, value string) []string {
	term := ParseSearch(value)
	highlights := make([]string, 0)
	for _, word := range term.Words {
		if strings.Contains(strings.ToLower(record.Title), word) {
			highlights = append(highlights, record.Title)
		}
		if strings.Contains(strings.ToLower(record.Body), word) {
			highlights = append(highlights, word)
		}
	}
	return highlights
}

func Deduplicate(records []model.Record) []model.Record {
	seen := map[string]bool{}
	result := make([]model.Record, 0, len(records))
	for _, record := range records {
		if record.ID != "" && !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

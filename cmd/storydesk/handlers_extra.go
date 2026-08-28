package main

import (
	"encoding/json"
	"github.com/charity/storydesk/internal/query"
	"github.com/charity/storydesk/internal/report"
	"net/http"
	"strconv"
	"time"
)

type listResponse struct {
	Items    any  `json:"items"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasNext  bool `json:"has_next"`
}

func pagedStoriesHandler(queries *query.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := parseInt(r.URL.Query().Get("page"), 1)
		size := parseInt(r.URL.Query().Get("page_size"), 20)
		items, err := queries.Find(query.Filters{Status: r.URL.Query().Get("status"), Category: r.URL.Query().Get("category"), Search: r.URL.Query().Get("q"), VisibleOnly: r.URL.Query().Get("visible") == "1"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err := query.Paginate(items, page, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, listResponse{Items: result.Items, Page: result.Page, PageSize: result.PageSize, Total: result.Total, HasNext: result.HasNext})
	}
}

func summaryHandler(queries *query.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		items, err := queries.Find(query.Filters{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		snapshot := report.Build(items, nowUTC())
		writeJSON(w, snapshot)
	}
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func nowUTC() time.Time { return time.Now().UTC() }

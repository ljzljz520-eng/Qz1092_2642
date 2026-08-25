package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/charity/storydesk/internal/query"
	"github.com/charity/storydesk/internal/service"
	"github.com/charity/storydesk/internal/store"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	path := flag.String("db", filepath.Join(os.TempDir(), "storydesk.db"), "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	st, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	queries := query.New(st)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/stories", storiesHandler(queries))
	mux.HandleFunc("/stories/page", pagedStoriesHandler(queries))
	mux.HandleFunc("/summary", summaryHandler(queries))
	mux.HandleFunc("/stories/", storyHandler(svc))
	log.Printf("storydesk listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func storiesHandler(queries *query.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := query.Filters{Status: r.URL.Query().Get("status"), Category: r.URL.Query().Get("category"), Search: r.URL.Query().Get("q"), VisibleOnly: r.URL.Query().Get("visible") == "1"}
		items, err := queries.Find(filters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	}
}

func storyHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)
		if id == "" || id == "." {
			http.Error(w, "story id is required", http.StatusBadRequest)
			return
		}
		record, err := svc.GetStory(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(record); err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

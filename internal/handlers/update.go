package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"web/internal/db"
)

type UpdateResponse struct {
	ID       int     `json:"id"`
	Current  int     `json:"current"`
	Total    int     `json:"total"`
	Progress float64 `json:"progress"` // 0..100
	Done     bool    `json:"done"`
}

func Update(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Must be POST (required)
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id <= 0 {
			http.Error(w, "id inválido", 400)
			return
		}

		s, err := db.IncrementEpisode(database, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		progress := 0.0
		if s.TotalEpisodes > 0 {
			progress = (float64(s.CurrentEpisode) / float64(s.TotalEpisodes)) * 100.0
		}

		resp := UpdateResponse{
			ID:       s.ID,
			Current:  s.CurrentEpisode,
			Total:    s.TotalEpisodes,
			Progress: progress,
			Done:     s.CurrentEpisode >= s.TotalEpisodes,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
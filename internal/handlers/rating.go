package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"web/internal/db"
)

type RatingBody struct {
	Rating int `json:"rating"`
}

func Rating(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use POST for rating updates
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id <= 0 {
			http.Error(w, "id inválido", 400)
			return
		}

		var body RatingBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}

		if body.Rating < 0 || body.Rating > 10 {
			http.Error(w, "rating debe ser 0..10", 400)
			return
		}

		if err := db.SetRating(database, id, body.Rating); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}
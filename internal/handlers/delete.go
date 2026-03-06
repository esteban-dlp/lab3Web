package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"web/internal/db"
)

func Delete(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Must be DELETE (challenge)
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id <= 0 {
			http.Error(w, "id inválido", 400)
			return
		}

		if err := db.DeleteSeries(database, id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}
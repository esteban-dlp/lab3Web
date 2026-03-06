package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"web/internal/db"
	"web/internal/views"
)

type EditViewData struct {
	Series db.Series
	Error  string
}

func Edit(database *sql.DB) http.HandlerFunc {
	renderer, err := views.NewRenderer()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id <= 0 {
			http.Error(w, "id inválido", 400)
			return
		}

		switch r.Method {
		case "GET":
			s, err := db.GetSeriesByID(database, id)
			if err != nil {
				http.Error(w, "Serie no encontrada", 404)
				return
			}
			renderer.Render(w, "edit.html", EditViewData{Series: s})
			return

		case "POST":
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "No se pudo leer el body", 400)
				return
			}
			values, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				http.Error(w, "Body inválido", 400)
				return
			}

			name := strings.TrimSpace(values.Get("series_name"))
			currentStr := strings.TrimSpace(values.Get("current_episode"))
			totalStr := strings.TrimSpace(values.Get("total_episodes"))

			current, err1 := strconv.Atoi(currentStr)
			total, err2 := strconv.Atoi(totalStr)

			// validation
			if name == "" {
				s, _ := db.GetSeriesByID(database, id)
				renderer.Render(w, "edit.html", EditViewData{Series: s, Error: "El nombre es obligatorio."})
				return
			}
			if err1 != nil || err2 != nil || current < 1 || total < 1 {
				s, _ := db.GetSeriesByID(database, id)
				renderer.Render(w, "edit.html", EditViewData{Series: s, Error: "Episodios deben ser números >= 1."})
				return
			}
			if current > total {
				s, _ := db.GetSeriesByID(database, id)
				renderer.Render(w, "edit.html", EditViewData{Series: s, Error: fmt.Sprintf("Actual (%d) no puede ser mayor al total (%d).", current, total)})
				return
			}

			// SQL UPDATE (challenge)
			if err := db.UpdateSeries(database, id, name, current, total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
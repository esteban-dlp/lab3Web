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

type CreateViewData struct {
	Error string
}

func Create(database *sql.DB) http.HandlerFunc {
	renderer, err := views.NewRenderer()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			renderer.Render(w, "create.html", CreateViewData{})
			return

		case "POST":
			// Parse form body manually as requested (application/x-www-form-urlencoded)
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

			// Server-side validation
			if name == "" {
				renderer.Render(w, "create.html", CreateViewData{Error: "El nombre es obligatorio."})
				return
			}
			if err1 != nil || err2 != nil {
				renderer.Render(w, "create.html", CreateViewData{Error: "Episodios deben ser números válidos."})
				return
			}
			if current < 1 || total < 1 {
				renderer.Render(w, "create.html", CreateViewData{Error: "Los episodios deben ser >= 1."})
				return
			}
			if current > total {
				renderer.Render(w, "create.html", CreateViewData{Error: fmt.Sprintf("El episodio actual (%d) no puede ser mayor al total (%d).", current, total)})
				return
			}

			if err := db.InsertSeries(database, name, current, total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// POST/Redirect/GET with 303
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
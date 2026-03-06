package handlers

import (
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"web/internal/db"
	"web/internal/views"
)

type HomeViewData struct {
	Series     []db.Series
	Query      string
	Sort       string
	Dir        string
	Page       int
	PageSize   int
	TotalCount int
	TotalPages int
}

func Home(database *sql.DB) http.HandlerFunc {
	renderer, err := views.NewRenderer()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query().Get("q")
		sort := r.URL.Query().Get("sort")
		dir := r.URL.Query().Get("dir")

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		params := db.ListParams{
			Query:    q,
			Sort:     sort,
			Dir:      dir,
			Page:     page,
			PageSize: pageSize,
		}

		series, totalCount, err := db.ListSeries(database, params)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if params.PageSize <= 0 {
			params.PageSize = 10
		}
		totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))
		if totalPages == 0 {
			totalPages = 1
		}
		if params.Page <= 0 {
			params.Page = 1
		}

		renderer.Render(w, "home.html", HomeViewData{
			Series:     series,
			Query:      params.Query,
			Sort:       params.Sort,
			Dir:        params.Dir,
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		})
	}
}
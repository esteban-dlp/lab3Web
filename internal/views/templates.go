package views

import (
	"html/template"
	"math"
	"net/http"
	"path/filepath"
)

func progress(current, total int) int {
	if total <= 0 {
		return 0
	}
	p := int(math.Round((float64(current) / float64(total)) * 100))
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

func inc(x int) int { return x + 1 }
func dec(x int) int {
	if x <= 1 {
		return 1
	}
	return x - 1
}

// If you click the same column, toggle dir, else default asc
func toggleDir(currentSort, currentDir, clicked string) string {
	if currentSort == clicked {
		if currentDir == "asc" {
			return "desc"
		}
		return "asc"
	}
	return "asc"
}

type Renderer struct {
	t *template.Template
}

func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"progress":  progress,
		"inc":       inc,
		"dec":       dec,
		"toggleDir": toggleDir,
	}

	t := template.New("root").Funcs(funcs)

	parsed, err := t.ParseFiles(
		filepath.Join("templates", "home.html"),
		filepath.Join("templates", "create.html"),
		filepath.Join("templates", "edit.html"),
	)
	if err != nil {
		return nil, err
	}
	return &Renderer{t: parsed}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = r.t.ExecuteTemplate(w, name, data)
}
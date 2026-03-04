package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	var err error

	db, err = sql.Open("sqlite3", "./series.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS series (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		current_episode INTEGER,
		total_episodes INTEGER
	)
	`)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)

	fmt.Println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, name, current_episode, total_episodes FROM series")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	fmt.Fprintln(w, "<h1>Series Tracker</h1>")
	fmt.Fprintln(w, `<a href="/create">Agregar serie</a><br><br>`)

	fmt.Fprintln(w, "<table border='1'>")
	fmt.Fprintln(w, "<tr><th>ID</th><th>Name</th><th>Current</th><th>Total</th><th>Action</th></tr>")

	for rows.Next() {

		var id int
		var name string
		var current int
		var total int

		rows.Scan(&id, &name, &current, &total)

		fmt.Fprintf(w,
			"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td>"+
				"<td><button onclick='nextEpisode(%d)'>+1</button></td></tr>",
			id, name, current, total, id)

	}

	fmt.Fprintln(w, "</table>")

	fmt.Fprintln(w, `
<script>
async function nextEpisode(id){
	await fetch("/update?id=" + id, {method:"POST"})
	location.reload()
}
</script>
`)
}

func createHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		fmt.Fprintln(w, `
		<h1>Agregar Serie</h1>

		<form method="POST" action="/create">

		Nombre:
		<input type="text" name="series_name" required><br>

		Episodio actual:
		<input type="number" name="current_episode" min="1" value="1"><br>

		Total episodios:
		<input type="number" name="total_episodes" min="1"><br>

		<button type="submit">Enviar</button>

		</form>

		<a href="/">Volver</a>
		`)

		return
	}

	if r.Method == "POST" {

		name := r.FormValue("series_name")
		current := r.FormValue("current_episode")
		total := r.FormValue("total_episodes")

		_, err := db.Exec(
			"INSERT INTO series(name,current_episode,total_episodes) VALUES(?,?,?)",
			name, current, total,
		)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	_, err := db.Exec(`
	UPDATE series
	SET current_episode = current_episode + 1
	WHERE id = ? AND current_episode < total_episodes
	`, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintln(w, "ok")
}
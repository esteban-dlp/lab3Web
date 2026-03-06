package main

import (
	"fmt"
	"log"
	"net/http"

	"web/internal/db"
	"web/internal/handlers"
)

func main() {
	database, err := db.Open("./series.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.InitSchema(database); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Static files (CSS/JS/Images)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	mux.HandleFunc("/", handlers.Home(database))
	mux.HandleFunc("/create", handlers.Create(database))
	mux.HandleFunc("/edit", handlers.Edit(database))

	// Fetch endpoints
	mux.HandleFunc("/update", handlers.Update(database))   // POST /update?id=1
	mux.HandleFunc("/series", handlers.Delete(database))  // DELETE /series?id=1
	mux.HandleFunc("/rating", handlers.Rating(database))  // POST /rating?id=1 (JSON)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
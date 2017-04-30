package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rigellute/note-server/httpMethods"
	"log"
	"net/http"
)

func main() {
	// Create the db connection
	db := setUpDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleNotes(w, r, db)
	}) // set router
	err := http.ListenAndServe(":3001", nil) // set listen port
	if err != nil {
		panic(err)
	}
}

func handleNotes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseForm() // parse arguments, you have to call this by yourself

	fmt.Println(r.Method)

	if r.Method == `GET` {
		httpMethods.GetNotes(w, r, db)
		return
	} else if r.Method == `POST` {
		httpMethods.PostNotes(w, r, db)
		return
	}

	http.Error(w, http.StatusText(400), 400)
}

func setUpDB() *sql.DB {
	const dbUrl = `postgresql://localhost/go_notes?sslmode=disable`

	// Update global db variable (this does not create db connection)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Force connection through Ping(), check if error
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

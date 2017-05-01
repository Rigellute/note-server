package httpMethods

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	// for accessing postgres
	_ "github.com/lib/pq"
)

// PostNotes used in main.go HandleFunc
func PostNotes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	book := r.FormValue("book")

	if len(book) == 0 {
		http.Error(w, "Must provide a book value", 400)
		return
	}

	note := r.FormValue("note")

	if len(note) == 0 {
		http.Error(w, "The note cannot be empty", 400)
		return
	}

	// Create book if not exists and create a note for it
	result, err := db.Exec(`
    WITH
      check_book AS (
        INSERT INTO books (title)
        SELECT $2
        WHERE NOT EXISTS (SELECT * FROM books WHERE title = $2)
        RETURNING *
      )
      , book as (
        SELECT id FROM check_book
        UNION
        SELECT id FROM books WHERE title = $2
      )
      INSERT INTO notes (note_content, book_id)
      VALUES ($1, (select id FROM book))
  `, note, book)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	fmt.Println("Rows affected", rowsAffected)

	form, _ := url.ParseQuery(r.URL.RawQuery)
	form.Add("resType", "html")
	r.URL.RawQuery = form.Encode()

	// success so return the get notes query
	GetNotes(w, r, db)
}

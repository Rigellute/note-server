package httpMethods

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	// Required for postgres query
	_ "github.com/lib/pq"
)

type book struct {
	Title string `json:"title"`
}

// Homepage is used int main.go
func Homepage(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	rows, err := db.Query(`
    SELECT title
    FROM books
  `)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Create slice of notes from postgres
	books := make([]*book, 0)

	for rows.Next() {
		// Create new Note struct
		book := new(book)

		// Check the rows for Note struct property
		if err = rows.Scan(&book.Title); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Add the full Note struct to Notes slice
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	absHTMLPath, _ := filepath.Abs("../note-server/html/book-list.html")

	var t *template.Template
	t, err = template.ParseFiles(absHTMLPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, books)

}

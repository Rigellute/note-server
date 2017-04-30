package httpMethods

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
	"time"
)

/**
 * Struct properties need to be uppercase to be exported
 * To turn a struct into json, the json.Marshal method
 * needs these exported methods to be visible
 */

type Note struct {
	Note_content string    `json:"note_content"`
	Created_at   time.Time `json:"created_at"`
}

func GetNotes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get query params for book
	book := strings.Join(r.Form["book"], "")

	fmt.Println("Querying", book)

	rows, err := db.Query(`
    SELECT note_content, notes.created_at
    FROM books
    JOIN notes ON notes.book_id = books.id
    WHERE title ILIKE $1;`, book)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Create slice of Notes
	notes := make([]*Note, 0)

	for rows.Next() {
		// Create new Note struct
		note := new(Note)

		// Check the rows for Note struct property
		err := rows.Scan(&note.Note_content, &note.Created_at)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// Add the full Note struct to Notes slice
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsNotes, err := json.Marshal(notes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsNotes)

}

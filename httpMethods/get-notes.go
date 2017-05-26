package httpMethods

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
	// Required for postgres query
	_ "github.com/lib/pq"
)

/**
 * Struct properties need to be uppercase to be exported
 * To turn a struct into json, the json.Marshal method
 * needs these exported methods to be visible
 */

type note struct {
	ID          int       `json:"id"`
	NoteContent string    `json:"note_content"`
	CreatedAt   time.Time `json:"created_at"`
	Book        string    `json:"book"`
}

type bookStruct struct {
	Title string  `json:"title"`
	Notes []*note `json:"notes"`
}

// GetNotes is used int main.go
func GetNotes(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// HACK: On the homepage there are two inputs named book so get one or the other
	book := r.Form["book"][0]

	if book == "" {
		book = r.Form["book"][1]
	}

	rows, err := db.Query(`
    SELECT notes.id, note_content, notes.created_at
    FROM books
    JOIN notes ON notes.book_id = books.id
    WHERE title ILIKE $1
    ORDER BY notes.created_at
    DESC
    ;`, book)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Create slice of notes from postgres
	notes := make([]*note, 0)

	for rows.Next() {
		// Create new Note struct
		note := new(note)

		// Check the rows for Note struct properties
		if err = rows.Scan(&note.ID, &note.NoteContent, &note.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Construct the note struct
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookWithData := bookStruct{Title: book, Notes: notes}

	// default return type is html, but json is allowed
	if r.FormValue("resType") == "json" {
		var jsNotes []byte
		jsNotes, err = json.Marshal(bookWithData)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsNotes)
		return
	}

	absHTMLPath, _ := filepath.Abs("../note-server/html/note-list.html")

	var t *template.Template
	t, err = template.ParseFiles(absHTMLPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, bookWithData)

}

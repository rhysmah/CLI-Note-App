package models

import "time"

// Note represents a single note entry.
// It contains simple data: a title, content, and tags.
// It contains metadata: an identifier, creation, and modification timestamps.
// The Note struct implements JSON serialization through struct tags.
type Note struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Tags       []string  `json:"tags"`
}

type NoteTitle struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

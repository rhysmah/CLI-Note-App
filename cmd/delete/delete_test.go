package delete

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/testutil"
	bolt "go.etcd.io/bbolt"
)

const (
	noteTitle   = "valid_note"
	noteContent = "valid_content"
)

// Test conditions
// > delete an existing note
// > delete an non-existing note

func createTestNote() models.Note {
	return models.Note{
		ID:         uuid.New().String(),
		Title:      noteTitle,
		Content:    noteContent,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Tags:       []string{},
	}
}

func storeNoteContent(tx *bolt.Tx, note models.Note) error {
	bucket := tx.Bucket([]byte(db.NotesBucket))
	if bucket == nil {
		return fmt.Errorf("bucket %s does not exist", db.NotesBucket)
	}

	noteJSON, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("failed to marshal note as JSON: %w", err)
	}

	err = bucket.Put([]byte(note.ID), noteJSON)
	if err != nil {
		return fmt.Errorf("failed to store note in database %q", db.NotesBucket)
	}

	return nil
}

// storeNoteInDB persists the given note in the BoltDB database.
// It marshals the note to JSON and stores it using the note's ID as the key.
func storeNoteInDB(note models.Note, database *bolt.DB) error {
	return database.Update(func(tx *bolt.Tx) error {
		if err := storeNoteContent(tx, note); err != nil {
			return fmt.Errorf("error storing note %q in database: %w", note.Title, err)
		}

		if err := storeNoteTitle(tx, note); err != nil {
			return fmt.Errorf("error storing note %q in database: %w", note.Title, err)
		}

		fmt.Printf("Added note %q to database\n", note.Title)
		return nil
	})
}

func storeNoteTitle(tx *bolt.Tx, note models.Note) error {
	bucket := tx.Bucket([]byte(db.NotesTitleBucket))

	if bucket == nil {
		return fmt.Errorf("failed to marshal note as JSON: %s", db.NotesTitleBucket)
	}

	err := bucket.Put([]byte(note.Title), []byte(note.ID))
	if err != nil {
		return fmt.Errorf("failed to store title in database %q", db.NotesBucket)
	}

	return nil
}

func TestDeleteExistingNote(t *testing.T) {
	testDb, _, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	note := createTestNote()

	err := storeNoteInDB(note, testDb)
	if err != nil {
		t.Errorf("Error adding note to database: %v", err)
	}

	testNoteContentSaved(t, note, testDb)
	testNoteTitleSaved(t, note, testDb)

	err = deleteNote(note.Title, testDb)
	if err != nil {
		t.Errorf("Error deleting note from database: %v", err)
	}

	testNoteContentDeleted(t, note, testDb)
	testNoteTitleDeleted(t, note, testDb)
}

func testNoteContentSaved(t *testing.T, note models.Note, database *bolt.DB) {
	var retrievedNoteContent models.Note

	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesBucket)
			return fmt.Errorf("bucket %s not found", db.NotesBucket)
		}

		retrievedNote := bucket.Get([]byte(note.ID))
		if retrievedNote == nil {
			t.Errorf("Note not found; expected Note with ID %s", note.ID)
			return fmt.Errorf("Note not found; expected Note with ID %s", note.ID)
		}

		return json.Unmarshal(retrievedNote, &retrievedNoteContent)
	})

	if err != nil {
		t.Errorf("Could not unmarshal note")
	}
	if retrievedNoteContent.ID != note.ID {
		t.Errorf("Incorrect Note ID retrieved; expected %s", note.ID)
	}
	if retrievedNoteContent.Content != note.Content {
		t.Errorf("Incorrect Note content; expected %s", note.Content)
	}
}

func testNoteTitleSaved(t *testing.T, note models.Note, database *bolt.DB) {
	var retrievedNoteID string

	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesTitleBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
			return fmt.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
		}

		// retrieving the Note ID associated with the Note Title
		retrievedNote := bucket.Get([]byte(note.Title))
		if retrievedNote == nil {
			t.Errorf("Note not found; expected %s", note.Title)
			return fmt.Errorf("Note not found; expected %s", note.Title)
		}

		retrievedNoteID = string(retrievedNote)
		return nil
	})

	if err != nil {
		t.Errorf("Note not found; expected %s", note.Title)
	}
	if retrievedNoteID != note.ID {
		t.Errorf("Incorrect ID mapping for title %q: got %q, want %q", note.Title, retrievedNoteID, note.ID)
	}
}

func testNoteTitleDeleted(t *testing.T, note models.Note, database *bolt.DB) {
	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesTitleBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
			return fmt.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
		}

		retrievedNote := bucket.Get([]byte(note.Title))
		if retrievedNote != nil {
			t.Errorf("Note title found; expected nothing")
			return fmt.Errorf("Note title found; expected nothing")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error accessing bucket %q: %v", db.NotesTitleBucket, err)
	}
}

func testNoteContentDeleted(t *testing.T, note models.Note, database *bolt.DB) {
	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesBucket)
			return fmt.Errorf("bucket %s not found", db.NotesBucket)
		}

		retrievedNote := bucket.Get([]byte(note.ID))
		if retrievedNote != nil {
			t.Errorf("Note found; expected nothing")
			return fmt.Errorf("Note found; expected nothing")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error accessing bucket %q: %v", db.NotesBucket, err)
	}
}

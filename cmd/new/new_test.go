package new

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/testutil"

	bolt "go.etcd.io/bbolt"
)

// TODO: Add tests to ensure note was actually saved to database

const (
	testNoteTitle        = "new_note"
	testNoteContent      = "sample text"
	testInvalidNoteTitle = "new:note"
)

func TestNewNote(t *testing.T) {
	note, err := createNote(testNoteTitle)
	if err != nil {
		t.Errorf("Couldn't create note: %v", err)
	}

	if note.ID == "" {
		t.Errorf("Note ID not properly set; got %s", note.ID)
	}
	if note.Title != testNoteTitle {
		t.Errorf("Note title not correct; got %s", note.Title)
	}
	if note.Content != "" {
		t.Errorf("Note content not empty; got %s", note.Content)
	}
	if note.CreatedAt.IsZero() {
		t.Errorf("Created timestamp not set")
	}
	if note.ModifiedAt.IsZero() {
		t.Errorf("Modified timestamp not set")
	}
	if len(note.Tags) != 0 {
		t.Errorf("Should be no tags; got %v", note.Tags)
	}
}

func TestInvalidNote(t *testing.T) {
	_, err := createNote(testInvalidNoteTitle)
	if err == nil {
		t.Errorf("Note title invalid; should have thrown error; got nil")
	}
}

func createTestNote() models.Note {
	return models.Note{
		ID:         uuid.New().String(),
		Title:      testNoteTitle,
		Content:    testNoteContent,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Tags:       []string{},
	}
}

func TestStoreNoteInDB(t *testing.T) {
	testDb, _, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	note := createTestNote()

	err := storeNoteInDB(note, testDb)
	if err != nil {
		t.Errorf("Error adding note to database: %v", err)
	}

	testNoteContentSaved(t, note, testDb)
	testNoteTitleSaved(t, note, testDb)
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

func TestStoreNoteInDB_Error(t *testing.T) {
	// Set up test DB without required buckets
	tempDir, err := os.MkdirTemp("", "notes-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	testDb, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	defer testDb.Close()

	// Should fail since buckets don't exist
	note := createTestNote()
	err = storeNoteInDB(note, testDb)

	if err == nil {
		t.Error("Expected error when buckets don't exist, got nil")
	}
}

package testutil

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
	bolt "go.etcd.io/bbolt"
)

const (
	TestValidNoteTitle   = "new_note"
	TestNoteContent      = "sample text"
	TestInvalidNoteTitle = "new:note"
)

func CreateTestNote() models.Note {
	return models.Note{
		ID:         uuid.New().String(),
		Title:      TestValidNoteTitle,
		Content:    TestNoteContent,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Tags:       []string{},
	}
}

func TestNoteContentSaved(t *testing.T, note models.Note, database *bolt.DB) {
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
			return fmt.Errorf("note not found; expected Note with ID %s", note.ID)
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

func TestNoteTitleSaved(t *testing.T, note models.Note, database *bolt.DB) {
	var retrievedNoteID string

	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesTitleBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
			return fmt.Errorf("bucket not found; expected %s", db.NotesTitleBucket)
		}

		// retrieving the Note ID associated with the Note Title
		retrievedNote := bucket.Get([]byte(note.Title))
		if retrievedNote == nil {
			t.Errorf("Note not found; expected %s", note.Title)
			return fmt.Errorf("note not found; expected %s", note.Title)
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

// TestSetupDB creates a temporary BoltDB database for testing purposes.
// It returns the database connection, the temporary directory path, and a cleanup function.
// The cleanup function should be deferred by the caller to ensure proper cleanup of resources.
// If any setup step fails, it will call t.Fatalf() and clean up any partially created resources.
func SetupTestDB(t *testing.T) (*bolt.DB, string, func()) {
	testTempDir, err := os.MkdirTemp("", "notes-test-*")
	if err != nil {
		t.Fatalf("Couldn't create temp directory: %v", err)
	}

	testDBPath := filepath.Join(testTempDir, "test.db")

	testDB, err := bolt.Open(testDBPath, db.DbReadWritePermissions, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		os.RemoveAll(testTempDir)
		t.Fatalf("couldn't create test database for testing: %v", err)
	}

	err = testDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.NotesBucket))
		return err
	})
	if err != nil {
		testDB.Close()
		os.RemoveAll(testTempDir)
		t.Fatalf("Couldn't create %v bucket: %v", db.NotesBucket, err)
	}

	err = testDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.NotesTitleBucket))
		return err
	})
	if err != nil {
		testDB.Close()
		os.RemoveAll(testTempDir)
		t.Fatalf("Couldn't create %v bucket: %v", db.NotesBucket, err)
	}

	cleanup := func() {
		testDB.Close()
		os.RemoveAll(testTempDir)
	}

	return testDB, testTempDir, cleanup
}

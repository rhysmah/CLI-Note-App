package new

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/testutil"

	bolt "go.etcd.io/bbolt"
)

// TODO: Add tests to ensure note was actually saved to database

const (
	testNoteTitle        = "new_note"
	testInvalidNoteTitle = "new:note"
)

func TestNewNote(t *testing.T) {
	note, err := createNote(testNoteTitle)
	if err != nil {
		t.Errorf("Couldn't create note: %v", err)
	}

	// Verify all fields
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
		ID: uuid.New().String(),
		Title: "test_note",
		Content: "",
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
		Tags: []string{},
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
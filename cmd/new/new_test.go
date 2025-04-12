package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/rhysmah/CLI-Note-App/testutil"

	bolt "go.etcd.io/bbolt"
)

func TestNewNote(t *testing.T) {
	note, err := createNote(testutil.TestValidNoteTitle)
	if err != nil {
		t.Errorf("Couldn't create note: %v", err)
	}

	if note.ID == "" {
		t.Errorf("Note ID not properly set; got %s", note.ID)
	}
	if note.Title != testutil.TestValidNoteTitle {
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
	_, err := createNote(testutil.TestInvalidNoteTitle)
	if err == nil {
		t.Errorf("Note title invalid; should have thrown error; got nil")
	}
}

func TestStoreNoteInDB(t *testing.T) {
	testDb, _, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	note := testutil.CreateTestNote()

	err := StoreNoteInDB(note, testDb)
	if err != nil {
		t.Errorf("Error adding note to database: %v", err)
	}

	testutil.TestNoteContentSaved(t, note, testDb)
	testutil.TestNoteTitleSaved(t, note, testDb)
}

func TestStoreNoteInDB_Error(t *testing.T) {
	// Set up test DB without required buckets
	tempDir, err := os.MkdirTemp("", "notes-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			// Need logging to properly capture this
			fmt.Println(err)
		}
	}()

	dbPath := filepath.Join(tempDir, "test.db")
	testDb, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	defer func() {
		err := testDb.Close()
		if err != nil {
			// Need logging to properly capture this
			fmt.Println(err)
		}
	}()

	defer func() {
		err := testDb.Close()
		if err != nil {
			// Need logging to properly capture this
			fmt.Println(err)
		}
	}()

	// Should fail since buckets don't exist
	note := testutil.CreateTestNote()
	err = StoreNoteInDB(note, testDb)

	if err == nil {
		t.Error("Expected error when buckets don't exist, got nil")
	}
}

func TestPreventDuplicateNoteInsertion(t *testing.T) {
	// Set up test database
	testDB, _, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	originalDB := root.NotesDB
	root.NotesDB = testDB
	defer func() { root.NotesDB = originalDB }()

	testNoteArg := []string{"test_note"}

	// Simulation: execute the command with a title
	newCmd := NewCommand()
	newCmd.SetArgs(testNoteArg)
	err := newCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// Create new command that should fail
	newCmd = NewCommand()
	newCmd.SetArgs(testNoteArg)
	err = newCmd.Execute()

	// Should receive note about duplicate note name
	if err == nil {
		t.Error("Should have received duplicate note error; received no error")
	}

	// Check that correct error was thrown
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		fmt.Println(err.Error())
		t.Errorf("Should have mentioned error about note existing: %v", err)
	}
}

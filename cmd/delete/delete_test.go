package delete

import (
	"fmt"
	"testing"

	"github.com/rhysmah/CLI-Note-App/cmd/new"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/rhysmah/CLI-Note-App/testutil"

	bolt "go.etcd.io/bbolt"
)

func TestDeleteExistingNote(t *testing.T) {
	testDb, _ := testutil.SetupTestDB(t)

	note := testutil.CreateTestNote()

	err := new.StoreNoteInDB(note, testDb)
	if err != nil {
		t.Errorf("Error adding note to database: %v", err)
	}

	testutil.TestNoteContentSaved(t, note, testDb)
	testutil.TestNoteTitleSaved(t, note, testDb)

	err = deleteNote(note.Title, testDb)
	if err != nil {
		t.Errorf("Error deleting note from database: %v", err)
	}

	testNoteContentNotInDB(t, note, testDb)
	testNoteTitleNotInDB(t, note, testDb)
}

func TestDeleteNonExistingNote(t *testing.T) {
	testDb, _ := testutil.SetupTestDB(t)

	note := testutil.CreateTestNote()

	testNoteContentNotInDB(t, note, testDb)
	testNoteTitleNotInDB(t, note, testDb)

	err := deleteNote(note.Title, testDb)
	if err == nil {
		t.Errorf("Expected error deleting non-existing note: %v", err)
	}
}

func testNoteTitleNotInDB(t *testing.T, note models.Note, database *bolt.DB) {
	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesTitleBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesTitleBucket)
			return fmt.Errorf("bucket not found; expected %s", db.NotesTitleBucket)
		}

		retrievedNote := bucket.Get([]byte(note.Title))
		if retrievedNote != nil {
			t.Errorf("Note title found; expected nothing")
			return fmt.Errorf("note title found; expected nothing")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error accessing bucket %q: %v", db.NotesTitleBucket, err)
	}
}

func testNoteContentNotInDB(t *testing.T, note models.Note, database *bolt.DB) {
	err := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesBucket))

		if bucket == nil {
			t.Errorf("Bucket not found; expected %s", db.NotesBucket)
			return fmt.Errorf("bucket %s not found", db.NotesBucket)
		}

		retrievedNote := bucket.Get([]byte(note.ID))
		if retrievedNote != nil {
			t.Errorf("Note found; expected nothing")
			return fmt.Errorf("note found; expected nothing")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error accessing bucket %q: %v", db.NotesBucket, err)
	}
}

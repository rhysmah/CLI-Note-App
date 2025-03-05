package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rhysmah/CLI-Note-App/db"
	bolt "go.etcd.io/bbolt"
)

func TestSetupDB(t *testing.T) (*bolt.DB, string, func()) {
	testTempDir, err := os.MkdirTemp("", "notes-test-*")
	if err != nil {
		t.Fatalf("Couldn't create temp directory: %v", err)
	}

	testDBPath := filepath.Join(testTempDir, "test.db")
	testDB, err := bolt.Open(testDBPath,
		db.DbReadWritePermissions,
		&bolt.Options{Timeout: 1 * time.Second})
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

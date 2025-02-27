package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	standardNotesDir       = ".notes"
	notesDBFile            = "notes.db"
	readWriteAccess        = 0755
	dbReadWritePermissions = 0600
	NotesBucket            = "Notes"
	NotesTitleBucket       = "NotesTitle"
)

// Initialize sets up and returns a new BoltDB instance for storing notes.
// It creates the directory structure and database file if they don't exist.
// If userPath is empty, it defaults to the user's home directory.
func Initialize(userPath string) (*bolt.DB, error) {
	notesDirectory, err := defineNotesDirectory(userPath)
	if err != nil {
		return nil, fmt.Errorf("error creating notes directory: %w", err)
	}

	if err := os.MkdirAll(notesDirectory, readWriteAccess); err != nil {
		return nil, fmt.Errorf("error creating default notes directory: %w", err)
	}

	dbFile := filepath.Join(notesDirectory, notesDBFile)
	db, err := setupNotesDB(dbFile)
	if err != nil {
		return nil, fmt.Errorf("error with database: %w", err)
	}
	return db, nil
}

// defineNotesDirectory determines the directory path where notes will be stored.
// If userPath is empty, it uses the user's home directory with a '.notes' subdirectory.
// Otherwise, it creates a '.notes' subdirectory in the specified userPath.
func defineNotesDirectory(userPath string) (string, error) {
	if userPath == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("user home directory not found: %w", err)
		}
		return filepath.Join(userHomeDir, standardNotesDir), nil
	}
	return filepath.Join(userPath, standardNotesDir), nil
}

// setupNotesDB opens or creates a BoltDB database file at the specified path.
// It configures the database with appropriate permissions and timeout settings.
func setupNotesDB(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, dbReadWritePermissions, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening / creating database: %w", err)
	}
	if err := createNoteBucket(db); err != nil {
		return nil, err
	}
	if err := createNoteTitleBucket(db); err != nil {
		return nil, err
	}
	return db, nil
}

// createNoteBucket ensures that the notes bucket exists in the database.
// If the bucket doesn't exist, it creates it.
func createNoteBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(NotesBucket))
		if err != nil {
			return fmt.Errorf("error creating %q bucket: %w", NotesBucket, err)
		}
		return nil
	})
}

func createNoteTitleBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(NotesTitleBucket))
		if err != nil {
			return fmt.Errorf("error creating %q bucket: %w", NotesTitleBucket, err)
		}
		return nil
	})
}

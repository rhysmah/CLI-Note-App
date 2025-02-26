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
	notesBucket            = "Notes"
)

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

func setupNotesDB(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, dbReadWritePermissions, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening / creating database: %w", err)
	}

	if err := createNoteBucket(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createNoteBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(notesBucket))
		if err != nil {
			return fmt.Errorf("error creating %q bucket: %w", notesBucket, err)
		}
		return nil
	})
}

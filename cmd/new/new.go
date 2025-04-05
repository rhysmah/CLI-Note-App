package new

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/spf13/cobra"

	bolt "go.etcd.io/bbolt"
)

const (
	createCmdFull  = "new"
	createCmdShort = "Create a new note"
	createCmdDesc  = `Create a new note with the specified name.
The note will be saved as '[note-name]_[date].txt' in your notes directory.
Note names cannot contain special characters or exceed 50 characters.`
)

// init registers the new note command with the root command.
func init() {
	newCommand := NewCommand()
	root.RootCmd.AddCommand(newCommand)
}

// NewCommand creates and returns a cobra.Command for creating new notes.
// The command requires exactly one argument: the note title.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   createCmdFull,
		Short: createCmdShort,
		Long:  createCmdDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			noteTitle := args[0]

			exists, err := checkIfNoteExists(noteTitle, root.NotesDB)
			if err != nil {
				return fmt.Errorf("error checking if note already exists: %w", err)
			}
			if exists {
				return fmt.Errorf("note %q already exists!\nPlease choose another name for your note", noteTitle)
			}

			note, err := createNote(args[0])
			if err != nil {
				return fmt.Errorf("error creating note: %w", err)
			}

			if err = StoreNoteInDB(note, root.NotesDB); err != nil {
				return fmt.Errorf("error saving note to database: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func checkIfNoteExists(title string, database *bolt.DB) (bool, error) {
	var exists bool
	err := database.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(db.NotesTitleBucket))
		if bucket == nil {
			return fmt.Errorf("error finding bucket %q", db.NotesTitleBucket)
		}

		val := bucket.Get([]byte(title))
		exists = val != nil
		return nil

	})
	return exists, err
}

// createNote instantiates a new Note with the given title and validates it.
func createNote(title string) (models.Note, error) {
	newNote := models.Note{
		ID:         uuid.New().String(),
		Title:      title,
		Content:    "",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	validator := newValidator()
	if err := validator.Run(newNote); err != nil {
		return models.Note{}, fmt.Errorf("invalid note name: %w", err)
	}
	return newNote, nil
}

// storeNoteInDB persists the given note in the BoltDB database.
// It marshals the note to JSON and stores it using the note's ID as the key.
func StoreNoteInDB(note models.Note, database *bolt.DB) error {
	return database.Update(func(tx *bolt.Tx) error {
		if err := StoreNoteContent(tx, note); err != nil {
			return fmt.Errorf("error storing note %q in database: %w", note.Title, err)
		}
		if err := StoreNoteTitle(tx, note); err != nil {
			return fmt.Errorf("error storing note %q in database: %w", note.Title, err)
		}
		fmt.Printf("Note %q successfully added to database!\nUse 'cli-note edit %s' to open your default text editor and start writing!'\n", note.Title, note.Title)
		return nil
	})
}

// StoreNoteContent saves a note to the database within an existing transaction.
// It marshals the note to JSON and stores it using the note's ID as the key.
// Returns an error if the bucket doesn't exist, JSON marshaling fails, or the database operation fails.
func StoreNoteContent(tx *bolt.Tx, note models.Note) error {
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

// StoreNoteTitle stores a mapping from note title to note ID in the titles bucket.
// It allows notes to be looked up by their titles.
// The function expects to be called within an existing bolt transaction.
func StoreNoteTitle(tx *bolt.Tx, note models.Note) error {
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

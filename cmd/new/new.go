package new

import (
	"encoding/json"
	"errors"
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
	createCmdFull  = "create"
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

			note, err := createNote(args[0])
			if err != nil {
				return fmt.Errorf("error creating note: %w", err)
			}

			if err = storeNoteInDB(note, root.NotesDB); err != nil {
				return fmt.Errorf("error saving note to database: %w", err)
			}

			return nil

		},
	}
	return cmd
}

// createNote instantiates a new Note with the given title and validates it.
// It generates a UUID, sets creation and modification times, and initializes 
// an empty content and tags slice.
func createNote(title string) (models.Note, error) {
	newNote := models.Note{
		ID:         uuid.New().String(),
		Title:      title,
		Content:    "",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Tags:       []string{},
	}

	validator := newValidator()
	if err := validator.Run(newNote); err != nil {
		return models.Note{}, fmt.Errorf("invalid note name: %w", err)
	}
	return newNote, nil
}

// storeNoteInDB persists the given note in the BoltDB database.
// It marshals the note to JSON and stores it using the note's ID as the key.
func storeNoteInDB(note models.Note, database *bolt.DB) error {

	return database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesBucket))
		if bucket == nil {
			return errors.New("bucket does not exist")
		}

		noteJSON, err := json.Marshal(note)
		if err != nil {
			return fmt.Errorf("failed to marshal note as JSON: %w", err)
		}

		err = bucket.Put([]byte(note.ID), noteJSON)
		if err != nil {
			return fmt.Errorf("failed to store note in database %q", db.NotesBucket)
		}

		fmt.Printf("Added note %q to database", note.Title)
		return nil
	})
}

package edit

import (
	"encoding/json"
	"fmt"

	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/rhysmah/CLI-Note-App/models"
	"github.com/spf13/cobra"

	bolt "go.etcd.io/bbolt"
)

const (
	editCmdFull  = "edit"
	editCmdShort = "Edit a note"
	editCmdLong  = "Edit a note by opening it in your default text editor"
)

// init registers the edit note command with the root command.
func init() {
	editCommand := EditCommand()
	root.RootCmd.AddCommand(editCommand)
}

func EditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   editCmdFull,
		Short: editCmdShort,
		Long:  editCmdLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			noteTitle := args[0]
			note, err := retrieveNote(noteTitle, root.NotesDB)
			if err != nil {
				return fmt.Errorf("error retrieving note %q: %w", noteTitle, err)
			}

			fmt.Println(note.Title)

			return nil
		},
	}
	return cmd
}

// TODO:
// User runs: notes edit "My Note Title"
// We retrieve the note from the database
// We create a temporary file with the note's content
// We determine which editor to use
// We open the editor with the temporary file
// User makes changes and closes the editor
// We read the modified content
// If changes were made, we update the note in the database
// We clean up the temporary file

func retrieveNote(noteTitle string, database *bolt.DB) (models.Note, error) {
	noteID, err := getNoteIDByTitle(noteTitle, database)
	if err != nil {
		return models.Note{}, err
	}
	return getNoteContent(noteID, database)
}

func getNoteIDByTitle(noteTitle string, database *bolt.DB) (string, error) {
	var retrievedNoteID string

	err := database.View(func(tx *bolt.Tx) error {

		notesTitleBucket := tx.Bucket([]byte(db.NotesTitleBucket))
		if notesTitleBucket == nil {
			return fmt.Errorf("bucket %s does not exist", db.NotesTitleBucket)
		}

		noteID := notesTitleBucket.Get([]byte(noteTitle))
		if noteID == nil {
			return fmt.Errorf("note %q does not exist", noteTitle)
		}

		retrievedNoteID = string(noteID)

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error retrieving note ID %w", err)
	}

	return retrievedNoteID, nil
}

func getNoteContent(noteID string, database *bolt.DB) (models.Note, error) {
	var retrievedNote models.Note

	err := database.View(func(tx *bolt.Tx) error {

		notesBucket := tx.Bucket([]byte(db.NotesBucket))
		if notesBucket == nil {
			return fmt.Errorf("bucket %s does not exist", db.NotesBucket)
		}

		note := notesBucket.Get([]byte(noteID))
		if note == nil {
			return fmt.Errorf("noteID %s does not exist", noteID)
		}

		return json.Unmarshal(note, &retrievedNote)
	})

	if err != nil {
		return models.Note{}, fmt.Errorf("error retrieving note: %w", err)
	}

	return retrievedNote, nil
}

// func getNote(note string, database *bolt.DB) (models.Note, error) {

// 	err := database.Update(func(tx *bolt.Tx) error {

// 		notesTitleBucket := tx.Bucket([]byte(db.NotesTitleBucket))
// 		if notesTitleBucket == nil {
// 			return fmt.Errorf("bucket %s does not exist", db.NotesTitleBucket)
// 		}

// 		var noteTitle string

// 		notesBucket := tx.Bucket([]byte(db.NotesBucket))
// 		if notesBucket == nil {
// 			return fmt.Errorf("bucket %s does not exist", db.NotesBucket)
// 		}

// 	})
// }

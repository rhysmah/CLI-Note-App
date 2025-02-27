package delete

import (
	"errors"
	"fmt"
	"os"

	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/spf13/cobra"

	bolt "go.etcd.io/bbolt"
)

const (
	deleteCmdFull  = "delete"
	deleteCmdShort = "Delete a note"
	deleteCmdDesc  = `Delete an existing note from your notes database.

Usage:
  notes delete <note-id>

The note-id is required and must match the ID of an existing note.
This action cannot be undone.`
)

func init() {
	deleteCommand := DeleteCommand()
	root.RootCmd.AddCommand(deleteCommand)
}

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   deleteCmdFull,
		Short: deleteCmdShort,
		Long:  deleteCmdFull,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := deleteNote(args[0], root.NotesDB); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
		},
	}
	return cmd
}

func deleteNote(title string, database *bolt.DB) error {
	return database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(db.NotesBucket))
		if bucket == nil {
			return errors.New("NotesBucket does not exist")
		}

		noteID, err := extractNoteID(title, tx)
		if err != nil {
			return fmt.Errorf("error finding note title %q: %w", title, err)
		}

		if err := bucket.Delete([]byte(noteID)); err != nil {
			return fmt.Errorf("error deleting note %q: %w", title, err)
		}

		// Clean up the title mapping
		titlesBucket := tx.Bucket([]byte(db.NotesTitleBucket))
		if titlesBucket == nil {
			return errors.New("NotesTitleBucket does not exist")
		}
		if err := titlesBucket.Delete([]byte(title)); err != nil {
			return fmt.Errorf("error removing title mapping for %q: %w", title, err)
		}

		fmt.Printf("Successfully deleted note %q from database", title)
		return nil
	})
}

func extractNoteID(noteTitle string, tx *bolt.Tx) (string, error) {
	bucket := tx.Bucket([]byte(db.NotesTitleBucket))
	if bucket == nil {
		return "", errors.New("NotesTitleBucket does not exist")
	}
	noteID := bucket.Get([]byte(noteTitle))
	if noteID == nil {
		return "", fmt.Errorf("error finding note %q", noteTitle)
	}
	return string(noteID), nil
}

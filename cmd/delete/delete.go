package delete

import (
	"errors"
	"fmt"

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
  notes delete <note-title>

The note-title is required and must match the title of an existing note.
This action cannot be undone.`
)

// init registers the delete command with the root command.
func init() {
	deleteCommand := DeleteCommand()
	root.RootCmd.AddCommand(deleteCommand)
}

// DeleteCommand creates and returns a cobra.Command for deleting notes.
// The command requires exactly one argument: the note title to delete.
func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   deleteCmdFull,
		Short: deleteCmdShort,
		Long:  deleteCmdDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := deleteNote(args[0], root.NotesDB); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

// deleteNote removes a note from the database using its title.
// It deletes both the note content and its title mapping.
// Returns an error if the note doesn't exist or if deletion fails.
func deleteNote(title string, database *bolt.DB) error {
	return database.Update(func(tx *bolt.Tx) error {
		if err := deleteNoteContent(title, tx); err != nil {
			return fmt.Errorf("failed to delete note content: %w", err)
		}
		if err := deleteNoteTitle(title, tx); err != nil {
			return fmt.Errorf("failed to delete note title mapping: %w", err)
		}
		fmt.Printf("Successfully deleted note %q from database\n", title)
		return nil
	})
}

// deleteNoteContent removes the note content from the NotesBucket using the note's ID.
// The ID is retrieved from the title mapping before deletion.
func deleteNoteContent(title string, tx *bolt.Tx) error {
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
	return nil
}

// deleteNoteTitle removes the title-to-ID mapping from the NotesTitleBucket.
func deleteNoteTitle(title string, tx *bolt.Tx) error {
	titlesBucket := tx.Bucket([]byte(db.NotesTitleBucket))
	if titlesBucket == nil {
		return errors.New("NotesTitleBucket does not exist")
	}
	if err := titlesBucket.Delete([]byte(title)); err != nil {
		return fmt.Errorf("error removing title mapping for %q: %w", title, err)
	}
	return nil
}

// extractNoteID retrieves the note's ID from the title mapping bucket.
// Returns the note ID as a string if found, or an error if the title
// doesn't exist or the bucket is missing.
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

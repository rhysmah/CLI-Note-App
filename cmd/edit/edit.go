package edit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

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

			// Create temporary file to write data
			tempFile, err := os.CreateTemp("", "temp-file-*.txt")
			if err != nil {
				return fmt.Errorf("error creating temp file: %w", err)
			}

			defer func() {
				if err := tempFile.Close(); err != nil {
					log.Printf("error closing temp file: %v", err)
				}
				if err = os.Remove(tempFile.Name()); err != nil {
					log.Printf("error removing temp file: %v", err)
				}
			}()

			// Copy data from current note to temp file
			if _, err := tempFile.WriteString(note.Content); err != nil {
				return fmt.Errorf("error writing to temp file: %w", err)
			}

			editor := determineEditor()
			command := exec.Command(editor, tempFile.Name())
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				return fmt.Errorf("error running editor: %w", err)
			}

			// Read back the edited file
			editedContent, err := os.ReadFile(tempFile.Name())
			if err != nil {
				return fmt.Errorf("error reading edited file: %w", err)
			}

			// Check if content changed
			if string(editedContent) != note.Content {
				note.Content = string(editedContent)
				note.ModifiedAt = time.Now()

				// Save the updated note
				if err := updateNote(note, root.NotesDB); err != nil {
					return fmt.Errorf("error saving updated note: %w", err)
				}

				fmt.Println("Note updated successfully.")
			} else {
				fmt.Println("No changes made to note.")
			}

			return nil
		},
	}
	return cmd
}

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

func updateNote(note models.Note, database *bolt.DB) error {
	return database.Update(func(tx *bolt.Tx) error {
		// Get the notes bucket
		notesBucket := tx.Bucket([]byte(db.NotesBucket))
		if notesBucket == nil {
			return fmt.Errorf("bucket %s does not exist", db.NotesBucket)
		}

		// Marshal the updated note to JSON
		noteJSON, err := json.Marshal(note)
		if err != nil {
			return fmt.Errorf("failed to marshal note: %w", err)
		}

		// Store the updated note
		if err := notesBucket.Put([]byte(note.ID), noteJSON); err != nil {
			return fmt.Errorf("failed to store note: %w", err)
		}

		return nil
	})
}

func determineEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin": // macOS
		for _, editor := range []string{"nano", "vim", "vi"} {
			if _, err := exec.LookPath(editor); err == nil {
				return editor
			}
		}

		// Fall back to TextEdit if available
		if _, err := exec.LookPath("TextEdit"); err == nil {
			return "open -a TextEdit"
		}

		// Last resort
		return "nano"
	default: // Linux and others
		// Try common editors
		for _, editor := range []string{"nano", "vim", "vi", "emacs"} {
			if _, err := exec.LookPath(editor); err == nil {
				return editor
			}
		}
		return "nano" // Default
	}
}

package list

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
	listCmdFull  = "list"
	listCmdShort = "List all notes"
	listCmdDesc  = "Display a list of all your notes."
)

func init() {
	listCmd := ListCommand()
	root.RootCmd.AddCommand(listCmd)
}

// NewCommand creates and returns a cobra.Command for creating new notes.
// The command requires exactly one argument: the note title.
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   listCmdFull,
		Short: listCmdShort,
		Long:  listCmdDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Get flags
			sort, _ := cmd.Flags().GetString("sort-by")
			reverse, _ := cmd.Flags().GetBool("reverse")

			// Convert flag values to SortType, OrderType
			sortBy := convertToSortBy(sort)
			orderBy := convertToSortOrder(reverse)

			// Sort and displays values based on flags

			notes, err := getNotes(root.NotesDB)
			if err != nil {
				return fmt.Errorf("error opening database")
			}

			if len(notes) == 0 {
				fmt.Println("You have no notes")
				return nil
			}

			sortNotes(notes, sortBy, orderBy)

			for _, note := range notes {
				fmt.Println(note.Title)
			}

			return nil
		},
	}

	// TODO: create constants for arguments
	cmd.Flags().StringP("sort-by", "s", "modified", "Sort by: title, created, modified")
	cmd.Flags().BoolP("reverse", "r", false, "Reverse the sort order")

	return cmd
}


// TODO: create a note print function.


func convertToSortBy(sort string) SortBy {
	switch sort {
	case "modified":
		return SortByModified
	case "created":
		return SortByCreated
	case "title":
		return SortByTitle
	default:
		return SortByTitle
	}
}

func convertToSortOrder(reverse bool) SortOrder {
	switch reverse {
	case false:
		return SortOrderAscending
	case true:
		return SortOrderDescending
	default:
		return SortOrderAscending
	}
}

func getNotes(database *bolt.DB) ([]models.Note, error) {
	var notes []models.Note

	err := database.View(func(tx *bolt.Tx) error {
		
		notesBucket := tx.Bucket([]byte(db.NotesBucket))
		if notesBucket == nil {
			return fmt.Errorf("bucket %s does not exist", db.NotesBucket)
		}

		return notesBucket.ForEach(func(k, v []byte) error {
			var note models.Note
			if err := json.Unmarshal(v, &note); err != nil {
				return fmt.Errorf("error reading note data: %w", err)
			}

			notes = append(notes, note)
			return nil
		})
	})

	if err != nil {
		return nil, fmt.Errorf("error retrieving note titles: %w", err)
	}

	return notes, nil
}

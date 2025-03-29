package list

import (
	"fmt"

	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/spf13/cobra"
)

const (
	listCmdFull  = "list"
	listCmdShort = "List all notes"
	listCmdDesc  = "Display a list of all your notes."

	sortFlag  = "sort-by"
	orderFlag = "reverse"
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

			sort, _ := cmd.Flags().GetString(sortFlag)
			sortBy := convertToSortBy(sort)

			order, _ := cmd.Flags().GetBool(orderFlag)
			orderBy := convertToSortOrder(order)

			notes, err := db.GetNotes(root.NotesDB)
			if err != nil {
				return fmt.Errorf("error opening database")
			}

			if len(notes) == 0 {
				fmt.Println("You have no notes")
				return nil
			}

			sortNotes(notes, sortBy, orderBy)
			DisplayNotes(notes, sortBy, orderBy)

			return nil
		},
	}

	// TODO: create constants for arguments
	cmd.Flags().StringP(sortFlag, "s", "modified", "Sort by: title, created, modified")
	cmd.Flags().BoolP(orderFlag, "r", false, "Reverse the sort order")

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

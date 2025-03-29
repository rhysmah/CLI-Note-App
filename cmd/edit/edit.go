package edit

import (
	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/spf13/cobra"
)


const (
	editCmdFull = "edit"
	editCmdShort = "Edit a note"
	editCmdLong = "Edit a note by opening it in your default text editor"
)

// init registers the edit note command with the root command.
func init() {
	editCommand :=  EditCommand() 
	root.RootCmd.AddCommand(editCommand)
}

func EditCommand() {
	cmd := &cobra.Command{
		Use: editCmdFull,
		Short: editCmdShort,
		Long: editCmdLong,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

		}
	}
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

func editNote() {
	
}

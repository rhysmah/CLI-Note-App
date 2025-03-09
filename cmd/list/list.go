package list

import (
	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/spf13/cobra"
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	return cmd
}

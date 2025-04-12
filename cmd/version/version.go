package version

import (
	"fmt"
	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/spf13/cobra"
)

func init() {
	VersionCmd := VersionCmd
	root.RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cli-notes v0.1")
	},
}

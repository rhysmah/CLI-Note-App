package version

import (
	"fmt"
	"github.com/rhysmah/CLI-Note-App/cmd/root"
	"github.com/spf13/cobra"
)

const (
	AppName    = "cli-notes"
	AppVersion = "0.1.1" // Semantic versioning (Major.Minor.Patch)
)

func init() {
	VersionCmd := VersionCmd
	root.RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s  |  %s\n", AppName, AppVersion)
	},
}

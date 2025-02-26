package root

import (
	"fmt"
	"os"

	"github.com/rhysmah/CLI-Note-App/db"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	NotesDB *bolt.DB
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "CLI-Note-App",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Include logic to setup notes directory the first time user starts program

		db, err := db.Initialize("")
		if err != nil {
			fmt.Printf("error initializing database: %s", err)
			os.Exit(1)
		}
		NotesDB = db
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		NotesDB.Close()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

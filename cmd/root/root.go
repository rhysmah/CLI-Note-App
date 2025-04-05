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
	Use:   "cli-note",
	Short: "A simple CLI-based note-taking app",
	Long: `CLI Note is exactly as the name implies: a dead simple CLI-based note-taking app.

You can create, edit, delete, and list your notes.

Basic Commands:
	new         Create a new .txt file
	edit        Open a file using your OS's default text editor
	delete      Delete a file via filename
	list        List all notes (name, creation date, modified date)

When you run CLI Notes for the first time, a small database is created locally on your machine.
This database is located in your home directory at ~/.cli-notes/

Examples:
  # Create a new note
  cli-note new "Shopping List"
  
  # Edit an existing note
  cli-note edit "Shopping List"
  
  # Delete a note
  cli-note delete "Shopping List"
  
  # List all your notes
  cli-note list`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// TODO
		// Include logic to setup notes directory the first time user starts program
		// For the above, check for a flag
		// Update command description
		// Improve error messages
		// Consider global logger

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

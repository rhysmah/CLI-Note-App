package new

import (
	"github.com/rhysmah/CLI-Note-App/cmd/root"
)

const (
	createCmdFull  = "create"
	createCmdShort = "Create a new note"
	createCmdDesc  = `Create a new note with the specified name.
The note will be saved as '[note-name]_[date].txt' in your notes directory.
Note names cannot contain special characters or exceed 50 characters.`
)

func init() {
	root.RootCmd.AddCommand()
}

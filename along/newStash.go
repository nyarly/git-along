package along

import (
	"github.com/spf13/cobra"
)

// newStashCmd represents the newStash command
var newStashCmd = &cobra.Command{
	Use:   "new-stash [name]",
	Short: "Create a new stash branch in this workspace",
	Long: longUsage(
		`Create a configstash branch called 'name' (default: along). You'll still
		need to add a remote and push to it, but new-stash takes care of the
		fiddly process of creating an empty branch without ancestors.`),
	RunE: runNewStash,
	Args: cobra.MaximumNArgs(1),
}

func runNewStash(cmd *cobra.Command, args []string) error {
	name := "along"
	if len(args) > 0 {
		name = args[0]
	}
	return emptyBranch(name)
}

func init() {
	alongCmd.AddCommand(newStashCmd)
}

package along

import (
	"github.com/spf13/cobra"
)

var (
	alongCmd = &cobra.Command{
		Use:   "git-along",
		Short: "git along gives config files a branch of their own to live in.",
	}
)

func init() {
}

// Execute runs the along command.
func Execute() error {
	return alongCmd.Execute()
}

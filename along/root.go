package along

import (
	"github.com/spf13/cobra"
)

var (
	alongCmd = &cobra.Command{
		Use:   "git-along",
		Short: "git along gives config files a branch of their own to live in.",
		Long: longUsage(
			`creates raw branches for specific files in your repo, so that you can use them
			without polluting the repo and let them evolve separately from the code.
			git-along acts like a specialized porcelain,
			so it ignores git’s normal rules (e.g. .gitignore),
			and you’re safe to stash things there.

			Rough example of use:

			Setup:
			git branch nixsupport
			git config --bool branch.nixsupport.configstash true
			# optionally, add files...
			git along add nixsupport shell.nix
			git along add nixsupport .envrc
			git push origin nixsupport

			Normal use:
			got along diff # to check for changes
			git along store nixsupport
			git push nixsupport

			git pull nixsupport
			git along retrieve nixsupport`),
	}
)

func init() {
}

// Execute runs the along command.
func Execute() error {
	return alongCmd.Execute()
}

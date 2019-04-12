package along

import (
	"github.com/spf13/cobra"
)

var (
	alongCmd = &cobra.Command{
		Use:   "git-along",
		Short: "git along gives config files a branch of their own to live in.",
		Long: longUsage(
			`git-along:
      creates raw branches for specific files in your repo, so that you can use them
			without polluting the repo and let them evolve separately from the code.
			git-along acts like a specialized porcelain,
			so it ignores git’s normal rules (e.g. .gitignore),
			and you’re safe to stash things there.

			Rough example of use:

			Setup:
			git along new-stash
			git along add  shell.nix
			git remote add along git@github.com:me/myconfigs.git
			git push -u along along:nix-thisproject

			Normal use:
			git along diff # to check for changes
			git along store
			git along add  .envrc
			git push along

			git pull along
			git along retrieve`),
	}

	stashBranchName = "along"
)

func init() {
	alongCmd.PersistentFlags().StringVar(&stashBranchName, "branchname", "along", "the branch to work against (default: along)")
}

// Execute runs the along command.
func Execute() error {
	return alongCmd.Execute()
}

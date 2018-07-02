package along

import (
	"fmt"
	"os"

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
func Execute() {
	if err := alongCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
branch=${1:?Usage: git stashconfig <branch>}

if [ $(git config branch.$branch.configstash) != true ]; then
  echo "Branch $branch is not listed as a configstash branch!"
  echo "  (try: git config --bool branch.$branch.configstash true)"
  exit 2
fi

paths=$(git config --get-all branch.$branch.stashedfile)

if [[ -z $paths ]]; then
  echo "Empty branch.$branch.stashedfile!"
  echo "  (try: git config --add --path branch.$branch.stashedfile <path>)"
  exit 3
fi

treeid=$((for p in $paths; do
  echo "100644 blob $(git hash-object -w $p)	$p"
done) | git mktree)
commitid=$(git commit-tree -p "$branch" -m "Stashconfig" "${treeid:?Tree id is empty!}")

git update-ref "refs/heads/$branch" "${commitid?Commit id is empty!}"
*/

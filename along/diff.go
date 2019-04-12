// Copyright © 2018 Judson Lester <nyarly@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package along

import (
	"fmt"

	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [flags]",
	Short: "Prints a diff of workspace files and their branch version.",
	Long: `Prints out a difference of the files in your current workspace
and the version recorded at the head of their along branch.

Example:
> git along diff
No difference.`,
	RunE: runDiff,
}

func init() {
	alongCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	pathlist, err := stashedfiles(stashBranchName)
	if err != nil {
		return err
	}

	hasDiff := false
	for _, path := range pathlist {
		diff, err := git("diff", branchpath(stashBranchName, path), path)
		if err != nil {
			return err
		}
		if len(diff) != 0 {
			fmt.Printf("%s:\n%s\n\n", path, diff)
			hasDiff = true
		}
	}
	if !hasDiff {
		fmt.Printf("No differences\n")
	}

	return nil
}

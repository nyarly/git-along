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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <branch> <file>",
	Short: "remove a file to a stashbranch",
	Long:  `Removes a file to a particular stashbranch.`,
	RunE:  runRemove,
	Args:  cobra.ExactArgs(2),
}

func init() {
	alongCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	branch := args[0]
	path := args[1]

	pathlist, err := stashedfiles(branch)
	if err != nil {
		return err
	}
	n := -1
	for i, p := range pathlist {
		if p == path {
			n = i
			break
		}
	}
	if n == -1 {
		return fmt.Errorf("%s not in %s", path, branch)
	}
	pathlist[n] = pathlist[len(pathlist)-1]
	pathlist = pathlist[:len(pathlist)-1]

	return storePaths(branch, pathlist)
}

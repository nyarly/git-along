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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "add a file to a stashbranch",
	Long:  `Adds a file to a particular stashbranch.`,
	RunE:  runAdd,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	alongCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	pathlist, err := stashedfiles(stashBranchName)

	for _, path := range args {
		if err != nil {
			return err
		}
		for _, p := range pathlist {
			if p == path {
				return errors.Errorf("%q already recorded in %q", path, stashBranchName)
			}
		}
		pathlist = append(pathlist, path)
	}

	return storePaths(stashBranchName, pathlist)
}

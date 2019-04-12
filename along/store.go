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
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "save stashed files",
	Long: longUsage(
		`stores the files tracked by git along into their stash stashBranchName. You'll use
		this commonly to commit changes to those files.`),
	RunE: runStore,
}

func init() {
	alongCmd.AddCommand(storeCmd)
}

func runStore(cmd *cobra.Command, args []string) error {
	pathlist, err := stashedfiles(stashBranchName)
	if err != nil {
		return err
	}

	return storePaths(stashBranchName, pathlist)
}

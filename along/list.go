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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list paths in a stashbranch",
	Long:  `Returns a list of the paths stored in the chosen branch`,
	RunE:  runList,
	Args:  cobra.ExactArgs(1),
}

func init() {
	alongCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	branch := args[0]
	pathlist, err := stashedfiles(branch)
	if err != nil {
		return err
	}

	for _, path := range pathlist {
		fmt.Println(path)
	}
	return nil
}

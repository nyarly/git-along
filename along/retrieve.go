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
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// retrieveCmd represents the retrieve command
var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "get files out of the stashbranch",
	Long: longUsage(
		`Updates files from the stashbranch. Checks out all the files recorded
		there. Will clobber the local versions of any files that exist. Note that
		old-but-removed stashed files will not be deleted.`),
	RunE: runRetrieve,
}

func init() {
	alongCmd.AddCommand(retrieveCmd)
}

func runRetrieve(cmd *cobra.Command, args []string) error {
	pathlist, err := stashedfiles(stashBranchName)
	if err != nil {
		return err
	}

	for _, path := range pathlist {
		contents, err := git("show", branchpath(stashBranchName, path))
		if err != nil {
			return err
		}
		dst, err := os.Create(path)
		if err != nil {
			return err
		}
		defer dst.Close()
		src := bytes.NewBuffer(contents)
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}
	return nil
}

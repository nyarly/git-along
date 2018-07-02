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
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runStore,
	Args: cobra.ExactArgs(1),
}

func init() {
	alongCmd.AddCommand(storeCmd)
}

type fileObject struct {
	path, hash string
}

func (obj fileObject) treeEntry() string {
	return fmt.Sprintf("100644 blob %s\t%s", obj.hash, obj.path)
}

func runStore(cmd *cobra.Command, args []string) error {
	branch := args[0]
	pathlist, err := stashedfiles(branch)
	if err != nil {
		return err
	}

	objects := []fileObject{}
	for _, path := range pathlist {
		hash, err := git("hash-object", "-w", path)
		if err != nil {
			return err
		}

		objects = append(objects, fileObject{
			path: path,
			hash: strings.TrimSpace(string(hash)),
		})
	}

	treeList := []string{}
	for _, obj := range objects {
		treeList = append(treeList, obj.treeEntry())
	}

	lstree := strings.Join(treeList, "\n")

	mktree := makeGit("mktree")
	mktree.Stdin = bytes.NewBufferString(lstree)
	treeid, err := runGit(mktree)
	if err != nil {
		return err
	}
	treeid = bytes.TrimSpace(treeid)
	if len(treeid) == 0 {
		return errors.Errorf("treeid is empty")
	}

	commitid, err := git("commit-tree", "-p", branch, "-m", "git-along", string(treeid))
	if err != nil {
		return err
	}
	commitid = bytes.TrimSpace(commitid)
	if len(commitid) == 0 {
		return errors.Errorf("commitid is empty")
	}

	git("update-ref", branchhead(branch), string(commitid))
	return nil
}

func branchhead(branch string) string {
	return fmt.Sprintf("refs/heads/%s", branch)
}

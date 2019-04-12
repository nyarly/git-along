package along

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type (
	fileObject struct {
		path, hash string
	}

	fileObjects []fileObject
)

func (obj fileObject) treeEntry() string {
	return fmt.Sprintf("100644 blob %s\t%s", obj.hash, obj.baseName())
}

func (obj fileObject) baseName() string {
	return filepath.Base(obj.path)
}

func (obj fileObject) dirName() string {
	return filepath.Dir(obj.path)
}

// omg hacky. suit for purpose, but watch out
func (obj fileObject) nextSegment(under int) string {
	return strings.Split(obj.path, string(os.PathSeparator))[under]
}

func (os fileObjects) Less(i, j int) bool {
	return os[i].path < os[j].path
}

func (os fileObjects) Swap(i, j int) {
	os[i], os[j] = os[j], os[i]
}

func (os fileObjects) Len() int {
	return len(os)
}

func pathsIntoBranch(pathlist []string, branch string) error {
	lstree, err := pathsTree(pathlist)
	if err != nil {
		return err
	}

	return storeTree(branch, lstree)
}

func pathsTree(pathlist []string) (fileObjects, error) {
	objects := fileObjects{}
	for _, path := range pathlist {
		hash, err := git("hash-object", "-w", path)
		if err != nil {
			return objects, err
		}

		objects = append(objects, fileObject{
			path: path,
			hash: strings.TrimSpace(string(hash)),
		})
	}
	return objects, nil
}

func storeTree(branch string, lstree fileObjects) error {
	sort.Sort(lstree)

	treeid, err := storeSubtree(".", 0, lstree)
	if err != nil {
		return err
	}
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

	_, err = git("update-ref", branchhead(branch), string(commitid))
	return err
}

func storeSubtree(path string, deep int, lstree fileObjects) (string, error) {
	lines := []string{}
	var curSub *string
	var subStart int
	var err error

	for i, fo := range lstree {
		switch {
		case fo.dirName() == path:
			if curSub != nil {
				lines, err = insertTree(lines, deep, lstree[subStart:i])
				if err != nil {
					return "", err
				}
				curSub = nil
			}
			lines = append(lines, fo.treeEntry())
		case curSub == nil:
			seg := fo.nextSegment(deep)
			curSub = &seg
			subStart = i
		case *curSub != fo.nextSegment(deep):
			lines, err = insertTree(lines, deep, lstree[subStart:i])
			seg := fo.nextSegment(deep)
			curSub = &seg
		}
	}
	if curSub != nil {
		lines, err = insertTree(lines, deep, lstree[subStart:])
		if err != nil {
			return "", err
		}
		curSub = nil
	}

	b, err := makeTree(strings.Join(lines, "\n"))
	return string(b), err
}

func insertTree(lines []string, deep int, ls fileObjects) ([]string, error) {
	fo := ls[0]
	h, err := storeSubtree(fo.dirName(), deep+1, ls)
	if err != nil {
		return lines, err
	}
	treeline := fmt.Sprintf("040000 tree %s\t%s", h, fo.nextSegment(deep))
	return append(lines, treeline), nil
}

func makeTree(lstree string) ([]byte, error) {
	mktree := makeGit("mktree")
	mktree.Stdin = bytes.NewBufferString(lstree)
	treeid, err := runGit(mktree)
	if err != nil {
		return []byte{}, err
	}
	treeid = bytes.TrimSpace(treeid)
	if len(treeid) == 0 {
		return []byte{}, errors.Errorf("treeid is empty")
	}
	return treeid, nil
}

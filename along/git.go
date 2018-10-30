package along

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type fileObject struct {
	path, hash string
}

func (obj fileObject) treeEntry() string {
	return fmt.Sprintf("100644 blob %s\t%s", obj.hash, obj.path)
}

func git(args ...string) ([]byte, error) {
	return runGit(makeGit(args...))
}

func makeGit(args ...string) *exec.Cmd {
	return exec.Command("git", args...)
}

func runGit(git *exec.Cmd) ([]byte, error) {
	out, err := git.CombinedOutput()
	if ee, is := err.(*exec.ExitError); is {
		return nil, errors.Wrapf(err, "%s:\n\t%s%s", strings.Join(git.Args, " "), strings.TrimSpace(string(out)), string(ee.Stderr))
	}
	return out, errors.Wrapf(err, "%s", strings.Join(git.Args, " "))
}

func nonZeroExit(err error) bool {
	_, is := errors.Cause(err).(*exec.ExitError)
	return is
}

func stashedfiles(branch string) ([]string, error) {
	stashkey := configstash(branch)
	configstash, err := git("config", stashkey)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(configstash)) != "true" {
		return nil, errors.Errorf("Branch $branch is not listed as a configstash branch!\n (try: `git config --bool %s true`, is currently %q)", stashkey, configstash)
	}

	paths, err := git("ls-tree", "--name-only", branch)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(paths)), "\n"), nil
}

func configstash(branch string) string {
	return fmt.Sprintf("branch.%s.configstash", branch)
}

func stashedfile(branch string) string {
	return fmt.Sprintf("branch.%s.stashedfile", branch)
}

func branchhead(branch string) string {
	return fmt.Sprintf("refs/heads/%s", branch)
}

func branchpath(branch, path string) string {
	return fmt.Sprintf("%s:%s", branch, path)
}

func storePaths(branch string, pathlist []string) error {
	if err := excludePaths(pathlist); err != nil {
		return err
	}

	lstree, err := pathsTree(pathlist)
	if err != nil {
		return err
	}

	return storeTree(branch, lstree)
}

func excludePaths(pathlist []string) error {
	for _, path := range pathlist {
		if _, fail := git("check-ignore", "-q", path); fail != nil {
			if !nonZeroExit(fail) {
				return fail
			}
			if err := excludePath(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func excludePath(path string) error {
	excludeFile, err := os.OpenFile(".git/info/exclude", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer excludeFile.Close()
	if _, err := excludeFile.Write(append([]byte(path), '\n')); err != nil {
		return err
	}
	return nil
}

func pathsTree(pathlist []string) (string, error) {
	objects := []fileObject{}
	for _, path := range pathlist {
		hash, err := git("hash-object", "-w", path)
		if err != nil {
			return "", err
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

	return strings.Join(treeList, "\n"), nil
}

func storeTree(branch, lstree string) error {
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

	_, err = git("update-ref", branchhead(branch), string(commitid))
	return err
}

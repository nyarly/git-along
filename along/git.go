package along

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func git(args ...string) ([]byte, error) {
	return runGit(makeGit(args...))
}

func makeGit(args ...string) *exec.Cmd {
	return exec.Command("git", args...)
}

func runGit(git *exec.Cmd) ([]byte, error) {
	watch := &bytes.Buffer{}
	if git.Stdin != nil {
		git.Stdin = io.TeeReader(git.Stdin, watch)
	}
	out, err := git.CombinedOutput()
	if ee, is := err.(*exec.ExitError); is {
		return nil, errors.Wrapf(err, "%s:\n\t%s%s\nconsumed stdin: \n%s", strings.Join(git.Args, " "), strings.TrimSpace(string(out)), string(ee.Stderr), watch.String())
	}
	return out, errors.Wrapf(err, "%s\nconsumed stdin: \n%s", strings.Join(git.Args, " "), watch.String())
}

func nonZeroExit(err error) bool {
	_, is := errors.Cause(err).(*exec.ExitError)
	return is
}

func emptyBranch(name string) error {
	treeid, err := makeTree("")
	if err != nil {
		return err
	}
	commitid, err := git("commit-tree", "-m", "git-along", string(treeid))
	if err != nil {
		return err
	}
	commitid = bytes.TrimSpace(commitid)
	if len(commitid) == 0 {
		return errors.Errorf("commitid is empty")
	}

	if _, err := git("branch", "--no-track", name, string(commitid)); err != nil {
		return err
	}

	stashkey := configstash(name)
	_, err = git("config", "--bool", stashkey, "true")
	return err
}

func stashedfiles(branch string) ([]string, error) {
	var err error
	stashkey := configstash(branch)
	configstash := []byte("<unset>")

	configstash, err = git("config", stashkey)
	if nonZeroExit(err) || strings.TrimSpace(string(configstash)) != "true" {
		return nil, errors.Errorf("Branch %q is not listed as a configstash branch!\n (try: `git config --bool %s true`, is currently %q)", branch, stashkey, configstash)
	}
	if err != nil {
		return nil, err
	}

	paths, err := git("ls-tree", "--name-only", "-r", branch)
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return []string{}, nil
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

	return pathsIntoBranch(pathlist, branch)
}

func excludePaths(pathlist []string) error {
	for _, path := range pathlist {
    if err := excludePath(path); err != nil {
      return err
    }
	}
	return nil
}

func excludePath(path string) error {
  if _, fail := git("check-ignore", "-q", path); fail != nil {
    if !nonZeroExit(fail) { // not ignored
      return fail
    }
  } else {
    return nil // already ignored
  }
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

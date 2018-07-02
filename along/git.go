package along

import (
	"fmt"
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
	out, err := git.CombinedOutput()
	if ee, is := err.(*exec.ExitError); is {
		return nil, errors.Wrapf(err, "%s:\n\t%s%s", strings.Join(git.Args, " "), strings.TrimSpace(string(out)), string(ee.Stderr))
	}
	return out, errors.Wrapf(err, "%s", strings.Join(git.Args, " "))
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

	fileskey := stashedfile(branch)
	paths, err := git("config", "--get-all", fileskey)
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

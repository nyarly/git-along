package along

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func git(args []string) ([]byte, error) {
	git := exec.Command("git", args...)
	out, err := git.CombinedOutput()
	return out, errors.Wrapf(err, "%s", git)
}

func stashedfiles(branch string) ([]string, error) {
	stashkey := configstash(branch)
	configstash, err := git("config", stashkey)
	if err != nil {
		return nil, err
	}
	if configstash != "true" {
		return nil, errors.Errof("Branch $branch is not listed as a configstash branch!\n (try: `git config --bool %s true`, is currently %q)", stashkey, configstash)
	}

	fileskey := stashedfile(branch)
	paths, err = git("config", "--get-all", fileskey)
	if err != nil {
		return nil, err
	}
	return strings.Split(paths, "\n"), nil
}

func configstash(branch string) string {
	return fmt.Sprintf("branch.%s.configstash", branch)
}

func stashedfile(branch string) string {
	return fmt.Sprintf("branch.%s.stashedfile", branch)
}

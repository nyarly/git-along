# git-along

Basic concept: creates raw branches for specific files in your repo, so that you can use them
* without polluting the repo
* and let them evolve separately from the code.

Install: `go get github.com/nyarly/git-along/along`

Rough example of use:

### Setup
```
git along new-stash nixsupport
git along add nixsupport shell.nix
git along add nixsupport .envrc
git remote add along github.com/me/myconfigs
git push -u along nixsupport:nix-thisproject
```

### Day to day
```
git along diff # to check for changes
git along store nixsupport
git push nixsupport

git pull nixsupport
git along retrieve nixsupport
```

`git-along` acts like a specialized porcelain,
so it ignores git’s normal rules (e.g. .gitignore),
and you’re safe to stash things there.
There’s still a little bit of finagling regarding
managing the remote for a stash branch,
(as well as the open question of naming that concept better…)
and maybe automating e.g. `.git/info/exclude` based on the stash branch contents.

But all in all, it works well for what it’s supposed to do: allow for the control of idiosyncratic project config files.

# Future work

* Branching from stash branches
* Figure out what to call "stash branches"
* History
* Remotes

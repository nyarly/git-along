package along

type fileObject struct {
	path, hash string
}

func (obj fileObject) treeEntry() string {
	return fmt.Sprintf("100644 blob %s\t%s", obj.hash, obj.path)
}

func pathsIntoBranch(pathlist []string, branch string) error {
	lstree, err := pathsTree(pathlist)
	if err != nil {
		return err
	}

	return storeTree(branch, lstree)
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
	treeid, err := makeTree(lstree)
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

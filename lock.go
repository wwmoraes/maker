package maker

// Lock stores the md5 hash sum for each managed snippet
type Lock map[string]map[string]string

func (lock Lock) Set(repo, name, version string) {
	if lock[repo] == nil {
		lock[repo] = make(map[string]string)
	}

	lock[repo][name] = version
}

func (lock Lock) Get(repo, name string) string {
	repoLock, exists := lock[repo]
	if !exists {
		return ""
	}

	return repoLock[name]
}

func (lock Lock) Unset(repo, name string) {
	repoLock, exists := lock[repo]
	if !exists {
		return
	}

	delete(repoLock, name)
}

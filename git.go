package gitmono

import (
	"github.com/gogs/git-module"
)

func OpenCurrentRepo() (*git.Repository, error) {
	gitDirPath := "./"
	return git.Open(gitDirPath)
}

package gitmono

import (
	"github.com/gogs/git-module"
)

// GitMono contains repository instance and command parameters
type GitMono struct {
	repo     *git.Repository
	projects []string
}

type Config struct {
	Projects []string
}

func OpenCurrentRepo(config *Config) (*GitMono, error) {
	repo, err := git.Open("./")
	if err != nil {
		return nil, err
	}

	monorepo := GitMono{
		repo:     repo,
		projects: config.Projects,
	}

	return &monorepo, nil
}

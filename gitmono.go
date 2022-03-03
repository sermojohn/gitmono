package gitmono

import (
	"github.com/gogs/git-module"
)

// GitMono contains repository instance and command parameters
type GitMono struct {
	repo   *git.Repository
	config *Config
}

// Config defines generic configuration applying to multiple commands
type Config struct {
	DryRun        bool
	CommitScheme  string
	VersionPrefix string
	Project       string
}

// OpenRepo open a git repository and returns the monorepo wrapper
func OpenRepo(path string) (*GitMono, error) {
	repo, err := git.Open("./")
	if err != nil {
		return nil, err
	}

	monorepo := GitMono{
		repo: repo,
	}

	return &monorepo, nil
}

func (gm *GitMono) Config(config *Config) {
	gm.config = config
}

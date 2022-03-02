package gitmono

import (
	"github.com/gogs/git-module"
)

// GitMono contains repository instance and command parameters
type GitMono struct {
	repo   *git.Repository
	config Config
}

// Config defines generic configuration applying to multiple commands
type Config struct {
	DryRun        bool
	CommitScheme  string
	VersionPrefix string
	PrintTag      bool
}

// OpenRepo open a git repository and returns the monorepo wrapper
func OpenRepo(path string, config *Config) (*GitMono, error) {
	repo, err := git.Open("./")
	if err != nil {
		return nil, err
	}

	monorepo := GitMono{
		repo:   repo,
		config: *config,
	}

	return &monorepo, nil
}

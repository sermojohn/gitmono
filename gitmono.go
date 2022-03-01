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
	Projects      []string
	DryRun        bool
	CommitScheme  string
	VersionPrefix string
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

// SetConfig configures the tool instance with configuration options
func (m *GitMono) SetConfig(config *Config) {
	m.config = *config
}

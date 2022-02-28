package gitmono

import (
	"github.com/gogs/git-module"
)

// GitMono contains repository instance and command parameters
type GitMono struct {
	repo          *git.Repository
	projects      []string
	dryRun        bool
	commitScheme  string
	versionPrefix string
}

type Config struct {
	Projects      []string
	DryRun        bool
	CommitScheme  string
	VersionPrefix string
}

func OpenCurrentRepo(config *Config) (*GitMono, error) {
	repo, err := git.Open("./")
	if err != nil {
		return nil, err
	}

	monorepo := GitMono{
		repo:          repo,
		projects:      config.Projects,
		dryRun:        config.DryRun,
		commitScheme:  config.CommitScheme,
		versionPrefix: config.VersionPrefix,
	}

	return &monorepo, nil
}

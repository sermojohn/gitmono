package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Log performs log operation for the specifed monorepo project
type Log struct {
	repo   *gitmono.GitRepository
	config *gitmono.Config
}

// NewLog creates a new logger instance
func NewLog(repo *gitmono.GitRepository, config *gitmono.Config) *Log {
	return &Log{
		repo:   repo,
		config: config,
	}
}

// Log performs log operation for the provided git references range and monorepo project
func (l *Log) Log(from, to string) ([]*git.Commit, error) {
	logOption := git.LogOptions{
		Path: l.config.Project,
	}

	return l.repo.Log(fmt.Sprintf("%s..%s", from, to), logOption)
}

// CommitHashByRevision lookup the commit hash for a revision/reference.
func (l *Log) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.repo.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

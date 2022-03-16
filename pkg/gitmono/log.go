package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Log performs log operation for the specifed monorepo project
type Log struct {
	logger gitmono.GitLogger
	config *gitmono.Config
}

// NewLog creates a new logger instance
func NewLog(logger gitmono.GitLogger, config *gitmono.Config) *Log {
	return &Log{
		logger: logger,
		config: config,
	}
}

// Log performs log operation for the provided git references range and monorepo project
func (l *Log) Log(from, to string) ([]*git.Commit, error) {
	logOption := git.LogOptions{
		Path: l.config.Project,
	}

	return l.logger.Log(fmt.Sprintf("%s..%s", from, to), logOption)
}

// CommitHashByRevision lookup the commit hash for a revision/reference.
func (l *Log) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.logger.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

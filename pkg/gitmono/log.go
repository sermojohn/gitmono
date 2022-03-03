package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Log performs log operation for the specifed monorepo project
type Log struct {
	monorepo *gitmono.MonoRepo
}

// NewLog creates a new logger instance
func NewLog(monorepo *gitmono.MonoRepo) *Log {
	return &Log{monorepo: monorepo}
}

// Log performs log operation for the provided git references range and monorepo project
func (l *Log) Log(from, to string) ([]*git.Commit, error) {
	logOption := git.LogOptions{
		Path: l.monorepo.GetConfig().Project,
	}

	return l.monorepo.Log(fmt.Sprintf("%s..%s", from, to), logOption)
}

// CommitHashByRevision lookup the commit hash for a revision/reference.
func (l *Log) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.monorepo.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

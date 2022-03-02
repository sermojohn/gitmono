package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
)

// Logger performs log operation for the specifed monorepo project
type Logger struct {
	mono *GitMono
}

// NewLogger creates a new logger instance
func NewLogger(mono *GitMono) *Logger {
	return &Logger{mono: mono}
}

// Log performs log operation for the provided git references range and monorepo project
func (l *Logger) Log(from, to string) ([]*git.Commit, error) {
	logOption := git.LogOptions{
		Path: l.mono.config.Project,
	}

	return l.mono.repo.Log(fmt.Sprintf("%s..%s", from, to), logOption)
}

// CommitHashByRevision lookup the commit hash for a revision/reference.
func (l *Logger) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.mono.repo.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

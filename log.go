package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
)

// Logger performs log operation for the monorepo projects provided
type Logger struct {
	mono *GitMono
}

// NewLogger creates a new logger instance
func NewLogger(mono *GitMono) *Logger {
	return &Logger{mono: mono}
}

// Log performs log operation for the provided git references range and monorepo projects
func (l *Logger) Log(from, to string, projects ...string) ([]*git.Commit, error) {
	logOptions := make([]git.LogOptions, 0, len(projects))
	for _, project := range projects {
		logOptions = append(logOptions, git.LogOptions{
			Path: project,
		})
	}

	return l.mono.repo.Log(fmt.Sprintf("%s..%s", from, to), logOptions...)
}

// CommitHashByRevision lookup the commit hash for a revision/reference.
func (l *Logger) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.mono.repo.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

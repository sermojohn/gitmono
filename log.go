package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
)

type Logger struct {
	mono *GitMono
}

func NewLogger(mono *GitMono) *Logger {
	return &Logger{mono: mono}
}

func (l *Logger) Log(from, to string) ([]*git.Commit, error) {
	if len(l.mono.projects) != 1 {
		return nil, fmt.Errorf("expected single project")
	}

	return l.mono.repo.Log(fmt.Sprintf("%s..%s", from, to), git.LogOptions{
		Path: l.mono.projects[0],
	})
}

func (l *Logger) CommitHashByRevision(rev string) (string, error) {
	commit, err := l.mono.repo.CommitByRevision(rev)
	if err != nil {
		return "", err
	}

	return commit.ID.String(), nil

}

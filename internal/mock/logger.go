package mock

import "github.com/gogs/git-module"

// Logger mock logger
type Logger struct {
	LogOutput                  []*git.Commit
	LogError                   error
	DidLog                     bool
	CommitHashByRevisionOutput string
	CommitHashByRevisionError  error
	DidCommitHashByRevision    bool
}

// Log mock log implementation
func (l *Logger) Log(from, to string) ([]*git.Commit, error) {
	if l.LogError != nil {
		return nil, l.LogError
	}

	l.DidLog = true
	return l.LogOutput, nil
}

// CommitHashByRevision mock commit from revision implementation
func (l *Logger) CommitHashByRevision(rev string) (string, error) {
	if l.CommitHashByRevisionError != nil {
		return "", l.CommitHashByRevisionError
	}

	l.DidCommitHashByRevision = true
	return l.CommitHashByRevisionOutput, nil
}

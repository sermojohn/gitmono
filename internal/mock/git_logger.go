package mock

import "github.com/gogs/git-module"

type GitLogger struct {
	LogInputs              []LogInput
	LogOutput              []*git.Commit
	LogError               error
	CommitByRevisionInputs []CommitByRevisionInput
	CommitByRevisionOutput *git.Commit
	CommitByRevisionError  error
}

func (gl *GitLogger) Log(rev string, opts ...git.LogOptions) ([]*git.Commit, error) {
	if gl.LogError != nil {
		return nil, gl.LogError
	}

	gl.LogInputs = append(gl.LogInputs, LogInput{rev: rev, opts: opts})
	return gl.LogOutput, nil
}

func (gl *GitLogger) CommitByRevision(rev string, opts ...git.CommitByRevisionOptions) (*git.Commit, error) {
	if gl.CommitByRevisionError != nil {
		return nil, gl.CommitByRevisionError
	}

	gl.CommitByRevisionInputs = append(gl.CommitByRevisionInputs, CommitByRevisionInput{rev: rev, opts: opts})
	return gl.CommitByRevisionOutput, nil
}

type LogInput struct {
	rev  string
	opts []git.LogOptions
}

type CommitByRevisionInput struct {
	rev  string
	opts []git.CommitByRevisionOptions
}

package mock

import (
	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// CommitParser mock commit parser
type CommitParser struct {
	DidGetBumperFromCommit       bool
	GetBumperFromCommitOutput    gitmono.Bumper
	GetBumperFromCommitOutputMap map[string]gitmono.Bumper
}

// GetBumperFromCommit mock bumper from commit implementation
func (cp *CommitParser) GetBumperFromCommit(c *git.Commit) gitmono.Bumper {
	cp.DidGetBumperFromCommit = true
	bumper, ok := cp.GetBumperFromCommitOutputMap[c.Message]
	if ok {
		return bumper
	}
	return cp.GetBumperFromCommitOutput
}

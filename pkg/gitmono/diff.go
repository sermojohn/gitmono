package gitmono

import (
	"strings"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Diff performs diff operation for a monorepo.
type Diff struct {
	differ gitmono.GitDiffer
	config *gitmono.Config
}

// NewDiff creates a new differ instance.
func NewDiff(differ gitmono.GitDiffer, config *gitmono.Config) *Diff {
	diff := Diff{
		differ: differ,
		config: config,
	}

	return &diff
}

// Diff performs diff for the provided git references range
// Matches changed files to the provided monorepo project and return the list of files
func (d *Diff) Diff(from, to string) ([]string, error) {
	diffRes, err := d.differ.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}

	var (
		changedFiles = []string{}
		project      = d.config.Project
	)
	for _, file := range diffRes.Files {
		if project == "." || strings.HasPrefix(file.Name, project) {
			changedFiles = append(changedFiles, file.Name)
		}
	}
	return changedFiles, nil
}

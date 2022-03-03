package gitmono

import (
	"strings"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Diff performs diff operation for a monorepo.
type Diff struct {
	monorepo *gitmono.MonoRepo
}

// NewDiff creates a new differ instance.
func NewDiff(monorepo *gitmono.MonoRepo) *Diff {
	diff := Diff{
		monorepo: monorepo,
	}

	return &diff
}

// Diff performs diff for the provided git references range
// Matches changed files to the provided monorepo project and return the list of files
func (d *Diff) Diff(from, to string) ([]string, error) {
	diffRes, err := d.monorepo.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}

	var (
		changedFiles = []string{}
		project      = d.monorepo.GetConfig().Project
	)
	for _, file := range diffRes.Files {
		if project == "." || strings.HasPrefix(file.Name, project) {
			changedFiles = append(changedFiles, file.Name)
		}
	}
	return changedFiles, nil
}

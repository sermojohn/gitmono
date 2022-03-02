package gitmono

import (
	"strings"

	"github.com/gogs/git-module"
)

// Differ performs diff operation for a monorepo.
type Differ struct {
	mono *GitMono
}

// NewDiffer creates a new differ instance.
func NewDiffer(mono *GitMono) *Differ {
	differ := Differ{
		mono: mono,
	}

	return &differ
}

// Diff performs diff for the provided git references range
// Matches changed files to the provided monorepo project and return the list of files
func (d *Differ) Diff(from, to string) ([]string, error) {
	diffRes, err := d.mono.repo.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}

	var (
		changedFiles = []string{}
		project      = d.mono.config.Project
	)
	for _, file := range diffRes.Files {
		if project == "." || strings.HasPrefix(file.Name, project) {
			changedFiles = append(changedFiles, file.Name)
		}
	}
	return changedFiles, nil
}

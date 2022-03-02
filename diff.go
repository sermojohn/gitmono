package gitmono

import (
	"fmt"
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
func (d *Differ) Diff(from, to, project string) ([]string, error) {
	diffRes, err := d.mono.repo.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}

	changedFiles := []string{}
	for _, file := range diffRes.Files {
		if project == "." || strings.HasPrefix(file.Name, project) {
			changedFiles = append(changedFiles, file.Name)
		}
	}
	if len(changedFiles) == 0 {
		return nil, fmt.Errorf("no diff")
	}
	return changedFiles, nil
}

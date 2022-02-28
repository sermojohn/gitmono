package gitmono

import (
	"strings"

	"github.com/gogs/git-module"
)

type Differ struct {
	mono *GitMono
}

func NewDiffer(mono *GitMono) *Differ {
	differ := Differ{
		mono: mono,
	}

	return &differ
}

func (d *Differ) Diff(from, to string) ([]string, error) {
	diffRes, err := d.mono.repo.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}

	diffedProjects := make([]string, 0, len(d.mono.config.Projects))
	diffedProjectsIndex := make(map[string]struct{}, len(d.mono.config.Projects))
	for _, file := range diffRes.Files {
		if project, matched := d.matchFile(file.Name); matched {
			if _, ok := diffedProjectsIndex[project]; !ok {
				diffedProjectsIndex[project] = struct{}{}
				diffedProjects = append(diffedProjects, project)
			}
		}
	}
	return diffedProjects, nil
}

func (d *Differ) matchFile(name string) (string, bool) {
	for _, project := range d.mono.config.Projects {
		if project == "." || strings.HasPrefix(name, project) {
			return project, true
		}
	}
	return "", false
}

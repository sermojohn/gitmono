package gitmono

import (
	"strings"

	"github.com/gogs/git-module"
)

type Differ struct {
	FromRef  string
	ToRef    string
	Projects []string
}

func NewDiffer(from string, to string, projects []string) *Differ {
	differ := Differ{
		FromRef:  from,
		ToRef:    to,
		Projects: projects,
	}

	return &differ
}

func (d *Differ) Diff() ([]string, error) {
	repo, err := OpenCurrentRepo()
	if err != nil {
		return nil, err
	}

	diffRes, err := repo.Diff(d.ToRef, 0, 0, 0, git.DiffOptions{
		Base: d.FromRef,
	})
	if err != nil {
		return nil, err
	}

	diffedProjects := make([]string, 0, len(d.Projects))
	diffedProjectsIndex := make(map[string]struct{}, len(d.Projects))
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
	for _, project := range d.Projects {
		if strings.HasPrefix(name, project) {
			return project, true
		}
	}
	return "", false
}

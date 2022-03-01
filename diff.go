package gitmono

import (
	"log"
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
// Matches changed files to the list of monorepo projects
func (d *Differ) Diff(from, to string, projects ...string) ([]string, error) {
	diffRes, err := d.mono.repo.Diff(to, 0, 0, 0, git.DiffOptions{
		Base: from,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("diff found %d files", len(diffRes.Files))

	diffedProjects := make([]string, 0, len(projects))
	diffedProjectsIndex := make(map[string]struct{}, len(projects))
	for _, file := range diffRes.Files {
		if project, matched := d.matchFile(file.Name, projects); matched {
			log.Printf("project %s matches changed file: %s\n", project, file.Name)
			if _, ok := diffedProjectsIndex[project]; !ok {
				diffedProjectsIndex[project] = struct{}{}
				diffedProjects = append(diffedProjects, project)
			}
		}
	}
	return diffedProjects, nil
}

func (d *Differ) matchFile(name string, projects []string) (string, bool) {
	for _, project := range projects {
		if project == "." || strings.HasPrefix(name, project) {
			return project, true
		}
	}
	return "", false
}

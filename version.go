package gitmono

import (
	"fmt"
	"strings"

	"github.com/gogs/git-module"
)

type Versioner struct {
	mono *GitMono
}

type VersionedCommit struct {
	CommitID string
	Tag      string
	Version  string
	Project  string
}

func NewVersioner(mono *GitMono) *Versioner {
	return &Versioner{
		mono: mono,
	}
}

func (v *Versioner) GetCurrentVersion() (*VersionedCommit, error) {
	if len(v.mono.projects) != 1 {
		return nil, fmt.Errorf("expected single project")
	}
	tags, err := v.mono.repo.Tags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		project, version := v.parseProjectVersion(tag)
		if !strings.EqualFold(project, v.mono.projects[0]) {
			continue
		}

		logger := &Logger{v.mono}
		commitHash, err := logger.CommitHashByRevision(tag)
		if err != nil {
			return nil, err
		}

		currentVersion := VersionedCommit{
			Tag:      tag,
			Version:  version,
			Project:  project,
			CommitID: commitHash,
		}

		return &currentVersion, nil
	}

	return nil, nil
}

func (v *Versioner) parseProjectVersion(tag string) (string, string) {
	idx := strings.LastIndex(tag, "/")
	if idx == -1 {
		return "", tag
	}

	return tag[0:idx], tag[idx+1:]
}

func (v *Versioner) GetNewCommits() ([]*git.Commit, error) {
	currentVersion, err := v.GetCurrentVersion()
	if err != nil {
		return nil, err
	}

	logger := NewLogger(v.mono)
	newCommits, err := logger.Log(currentVersion.CommitID, "HEAD")
	if err != nil {
		return nil, err
	}

	return newCommits, nil
}

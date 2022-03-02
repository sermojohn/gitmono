package gitmono

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-version"
)

// Versioner combines git commands to read and write releases
type Versioner struct {
	mono *GitMono
}

// VersionedCommit points a commit that is assigned a version
type VersionedCommit struct {
	CommitID      string
	Project       string
	VersionPrefix string
	Version       *version.Version
}

// GetTag returns the tag to version a commit with
func (vc *VersionedCommit) GetTag() string {
	var projectPrefix string
	if vc.Project != "." {
		projectPrefix = fmt.Sprintf("%s/", vc.Project)
	}
	return fmt.Sprintf("%s%s", projectPrefix, vc.GetVersion())
}

// GetVersion returns the version part of the tag
func (vc *VersionedCommit) GetVersion() string {
	return fmt.Sprintf("%s%s", vc.VersionPrefix, vc.Version.String())
}

// NewVersioner creates a new versioner instance
func NewVersioner(mono *GitMono) *Versioner {
	return &Versioner{
		mono: mono,
	}
}

// GetCurrentVersion retrieves the current version for the specified project
func (v *Versioner) GetCurrentVersion(project string) (*VersionedCommit, error) {
	tagger := &Tagger{mono: v.mono}
	tags, err := tagger.Tags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		parsedProject, version := v.parseProjectVersion(tag)
		if !strings.EqualFold(parsedProject, project) {
			continue
		}

		logger := &Logger{v.mono}
		commitHash, err := logger.CommitHashByRevision(tag)
		if err != nil {
			return nil, err
		}

		parsedVersion, err := v.parseVersion(version)
		if err != nil {
			return nil, err
		}

		currentVersion := VersionedCommit{
			Version:       parsedVersion,
			VersionPrefix: v.mono.config.VersionPrefix,
			Project:       parsedProject,
			CommitID:      commitHash,
		}

		log.Printf("current version: %v\n", currentVersion)
		return &currentVersion, nil
	}

	return nil, nil
}

func (v *Versioner) parseProjectVersion(tag string) (string, string) {
	idx := strings.LastIndex(tag, "/")
	if idx == -1 {
		return ".", tag
	}

	return tag[0:idx], tag[idx+1:]
}

func (v *Versioner) parseVersion(vv string) (*version.Version, error) {
	var (
		versionPrefix = v.mono.config.VersionPrefix
		versionValue  = vv
	)
	if versionPrefix != "" && strings.HasPrefix(vv, versionPrefix) {
		versionValue = vv[len(versionPrefix):]
	}

	parsedVersion, err := version.NewSemver(versionValue)
	if err != nil {
		return nil, err
	}

	return parsedVersion, nil
}

// ReleaseNewVersion calculates the new version for the provided project and performs release
//
// Returns an error if there are no new commits for the provided project
func (v *Versioner) ReleaseNewVersion(commitID string, project string) (*VersionedCommit, error) {
	if commitID == "" {
		commitID = "HEAD"
	}
	currentVersion, err := v.GetCurrentVersion(project)
	if err != nil {
		return nil, err
	}
	if currentVersion == nil {
		return nil, nil
	}

	logger := NewLogger(v.mono)
	newCommits, err := logger.Log(currentVersion.CommitID, commitID, project)
	if err != nil {
		return nil, err
	}
	if len(newCommits) == 0 {
		return nil, fmt.Errorf("no new commits were found")
	}

	var (
		commitParser = commitParser{scheme: v.mono.config.CommitScheme}
		bump         bumper
		lastCommitID string
	)
	for i, cm := range newCommits {
		if i == 0 {
			lastCommitID = cm.ID.String()
		}
		bump = commitParser.parseCommit(cm)
		if bump != nil {
			break
		}
	}

	if bump == nil {
		bump = patchBumper
	}
	newVersion, err := bump.bump(currentVersion.Version)
	if err != nil {
		return nil, err
	}

	newVersionedCommit := VersionedCommit{
		CommitID:      lastCommitID,
		Version:       newVersion,
		VersionPrefix: currentVersion.VersionPrefix,
		Project:       currentVersion.Project,
	}

	err = v.createReleaseTag(&newVersionedCommit)
	if err != nil {
		return nil, err
	}

	return &newVersionedCommit, nil
}

// InitVersion identifies checks if project has version and releases the initial version
func (v *Versioner) InitVersion(commitID string, project string) (*VersionedCommit, error) {
	if commitID == "" {
		commitID = "HEAD"
	}
	currentVersion, err := v.GetCurrentVersion(project)
	if err != nil {
		return nil, err
	}
	if currentVersion != nil {
		return nil, nil
	}

	initVersion, err := version.NewSemver("0.1.0")
	if err != nil {
		return nil, err
	}

	newVersionedCommit := VersionedCommit{
		CommitID:      commitID,
		Project:       project,
		Version:       initVersion,
		VersionPrefix: v.mono.config.VersionPrefix,
	}
	err = v.createReleaseTag(&newVersionedCommit)
	if err != nil {
		return nil, err
	}

	return &newVersionedCommit, nil
}

func (v *Versioner) createReleaseTag(vc *VersionedCommit) error {
	log.Printf("about to create tag: %s to commit: %s\n", vc.GetTag(), vc.CommitID)
	if v.mono.config.DryRun {
		return nil
	}

	tagger := NewTagger(v.mono)
	return tagger.CreateTag(vc)
}

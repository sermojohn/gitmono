package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/hashicorp/go-version"
)

// EnvVars contains the accepted environment variables
type EnvVars struct {
	CommitterName  string
	CommitterEmail string
}

// Config defines generic configuration applying to multiple commands
type Config struct {
	DryRun        bool
	CommitScheme  string
	VersionPrefix string
	Project       string
}

// Logger performs log commands on the repo
//
// Log returns the commits for the specified reference range in reverse chronological order
// CommitHashByRevision returns the commit hash matching to a revision
type Logger interface {
	Log(from, to string) ([]*git.Commit, error)
	CommitHashByRevision(rev string) (string, error)
}

// Tagger performs tag commands on the repo
//
// CreateTag writes the given tag to the given commit
// ListProjectVersionTags retrieves all tags for a project using the tag list pattern
type Tagger interface {
	CreateTag(versionedCommit *VersionedCommit) error
	ListProjectVersionTags() ([]string, error)
}

// Versioner maintains version using tags
//
// GetCurrentVersion retrieves the latest version
// ReleaseNewVersion creates and persists a new version
// InitVersion creates and persists the initial version
type Versioner interface {
	GetCurrentVersion() (*VersionedCommit, error)
	ReleaseNewVersion(commitID string) (*VersionedCommit, error)
	InitVersion(commitID string) (*VersionedCommit, error)
}

// Differ performs diff on the repo
type Differ interface {
	Diff(from, to string) ([]string, error)
}

// CommitParser parses the provided commit
//
// GetBumperFromCommit parses a commit message and decide the bumper to use for version
type CommitParser interface {
	GetBumperFromCommit(*git.Commit) Bumper
}

// Bumper bumps a version to a new version
type Bumper interface {
	Bump(*version.Version) (*version.Version, error)
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

// GitTagger abstracts git tag operations
type GitTagger interface {
	Tags(opts ...git.TagsOptions) ([]string, error)
	CreateTag(name, rev string, opts ...git.CreateTagOptions) error
}

// GitLogger abstracts git log operations
type GitLogger interface {
	Log(rev string, opts ...git.LogOptions) ([]*git.Commit, error)
	CommitByRevision(rev string, opts ...git.CommitByRevisionOptions) (*git.Commit, error)
}

// GitDiffer abstracts git diff operations
type GitDiffer interface {
	Diff(rev string, maxFiles, maxFileLines, maxLineChars int, opts ...git.DiffOptions) (*git.Diff, error)
}

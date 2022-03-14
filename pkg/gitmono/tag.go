package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Tag performs tag operation for a monorepo
type Tag struct {
	monorepo *gitmono.MonoRepo
}

// NewTag creates a new tagger instance
func NewTag(monorepo *gitmono.MonoRepo) *Tag {
	return &Tag{
		monorepo: monorepo,
	}
}

// Tags retrieves all repository tags ordered by descending creation date
func (t *Tag) Tags() ([]string, error) {
	return t.monorepo.Tags()
}

// ListProjectTags retrieves all project tags ordered by descending version value
func (t *Tag) ListProjectTags() ([]string, error) {
	return t.monorepo.Tags(git.TagsOptions{
		SortKey: "-version:refname",
		Pattern: fmt.Sprintf("%s/v*", t.monorepo.GetConfig().Project),
	})
}

// CreateTag create a tag on the provided commit
func (t *Tag) CreateTag(versionedCommit *gitmono.VersionedCommit) error {
	return t.monorepo.CreateTag(versionedCommit.GetTag(), versionedCommit.CommitID)
}

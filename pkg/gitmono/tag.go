package gitmono

import "github.com/sermojohn/gitmono"

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

// Tags retrieves all repository tags
func (t *Tag) Tags() ([]string, error) {
	return t.monorepo.Tags()
}

// CreateTag create a tag on the provided commit
func (t *Tag) CreateTag(versionedCommit *gitmono.VersionedCommit) error {
	return t.monorepo.CreateTag(versionedCommit.GetTag(), versionedCommit.CommitID)
}

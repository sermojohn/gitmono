package gitmono

import (
	"fmt"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// Tag performs tag operation for a monorepo
type Tag struct {
	tagger  gitmono.GitTagger
	config  *gitmono.Config
	envVars *gitmono.EnvVars
}

// NewTag creates a new tagger instance
func NewTag(tagger gitmono.GitTagger, config *gitmono.Config, envVars *gitmono.EnvVars) *Tag {
	return &Tag{
		tagger:  tagger,
		config:  config,
		envVars: envVars,
	}
}

// ListProjectVersionTags retrieves all project tags ordered by descending version value
func (t *Tag) ListProjectVersionTags() ([]string, error) {
	return t.tagger.Tags(git.TagsOptions{
		SortKey: "-version:refname",
		Pattern: fmt.Sprintf("%s%s*", gitmono.GetProjectTagPrefix(t.config.Project), t.config.VersionPrefix),
	})
}

// CreateTag create an annotated tag on the provided commit
func (t *Tag) CreateTag(versionedCommit *gitmono.VersionedCommit) error {
	var committer *git.Signature

	if t.envVars.CommitterName != "" && t.envVars.CommitterEmail != "" {
		committer = &git.Signature{
			Name:  t.envVars.CommitterName,
			Email: t.envVars.CommitterEmail,
		}
	}

	return t.tagger.CreateTag(versionedCommit.GetTag(), versionedCommit.CommitID, git.CreateTagOptions{
		Annotated: true,
		Author:    committer,
	})
}

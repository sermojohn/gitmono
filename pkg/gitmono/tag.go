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
func NewTag(repo *gitmono.GitRepository, config *gitmono.Config, envVars *gitmono.EnvVars) *Tag {
	return &Tag{
		tagger:  repo,
		config:  config,
		envVars: envVars,
	}
}

// Tags retrieves all repository tags ordered by descending creation date
func (t *Tag) Tags() ([]string, error) {
	return t.tagger.Tags()
}

// ListProjectTags retrieves all project tags ordered by descending version value
func (t *Tag) ListProjectTags() ([]string, error) {
	return t.tagger.Tags(git.TagsOptions{
		SortKey: "-version:refname",
		Pattern: fmt.Sprintf("%s/v*", t.config.Project),
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

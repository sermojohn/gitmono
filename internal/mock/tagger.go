package mock

import "github.com/sermojohn/gitmono"

type Tagger struct {
	ListProjectVersionTagsOutput []string
	ListProjectVersionTagsError  error
	CreateTagError               error
}

func (t *Tagger) ListProjectVersionTags() ([]string, error) {
	if t.ListProjectVersionTagsError != nil {
		return nil, t.ListProjectVersionTagsError
	}

	return t.ListProjectVersionTagsOutput, nil
}

func (t *Tagger) CreateTag(versionedCommit *gitmono.VersionedCommit) error {
	return t.CreateTagError
}

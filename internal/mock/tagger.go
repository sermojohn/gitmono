package mock

import "github.com/sermojohn/gitmono"

type Tagger struct {
	ListProjectTagsOutput []string
	ListProjectTagsError  error
	TagsOutput            []string
	TagsError             error
	CreateTagError        error
}

func (t *Tagger) ListProjectTags() ([]string, error) {
	if t.ListProjectTagsError != nil {
		return nil, t.ListProjectTagsError
	}

	return t.ListProjectTagsOutput, nil
}

func (t *Tagger) Tags() ([]string, error) {
	return t.TagsOutput, t.TagsError
}

func (t *Tagger) CreateTag(versionedCommit *gitmono.VersionedCommit) error {
	return t.CreateTagError
}

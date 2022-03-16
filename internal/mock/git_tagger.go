package mock

import "github.com/gogs/git-module"

type GitTagger struct {
	TagsOutput      []string
	TagsError       error
	TagsInputs      []*TagsInput
	CreateTagError  error
	CreateTagInputs []*CreateTagInput
}

func (gt *GitTagger) Tags(opts ...git.TagsOptions) ([]string, error) {
	if gt.TagsError != nil {
		return nil, gt.TagsError
	}
	gt.TagsInputs = append(gt.TagsInputs, &TagsInput{Opts: opts})
	return gt.TagsOutput, nil
}

func (gt *GitTagger) CreateTag(name, rev string, opts ...git.CreateTagOptions) error {
	if gt.CreateTagError != nil {
		return gt.CreateTagError
	}
	gt.CreateTagInputs = append(gt.CreateTagInputs, &CreateTagInput{
		Name: name,
		Rev:  rev,
		Opts: opts,
	})
	return nil
}

type TagsInput struct {
	Opts []git.TagsOptions
}

type CreateTagInput struct {
	Name string
	Rev  string
	Opts []git.CreateTagOptions
}

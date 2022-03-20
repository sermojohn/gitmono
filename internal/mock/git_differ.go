package mock

import "github.com/gogs/git-module"

type GitDiffer struct {
	DiffInputs []DiffInput
	DiffOutput *git.Diff
	DiffError  error
}

func (gd *GitDiffer) Diff(rev string, maxFiles, maxFileLines, maxLineChars int, opts ...git.DiffOptions) (*git.Diff, error) {
	if gd.DiffError != nil {
		return nil, gd.DiffError
	}

	gd.DiffInputs = append(gd.DiffInputs, DiffInput{
		rev:          rev,
		maxFiles:     maxFiles,
		maxFileLines: maxFileLines,
		maxLineChars: maxLineChars,
		opts:         opts,
	})
	return gd.DiffOutput, nil
}

type DiffInput struct {
	rev          string
	maxFiles     int
	maxFileLines int
	maxLineChars int
	opts         []git.DiffOptions
}

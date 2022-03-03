package gitmono

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

var (
	majorBumper major
	minorBumper minor
	patchBumper patch
)

type major struct{}

func (m major) Bump(cv *version.Version) (*version.Version, error) {
	segments := cv.Segments()
	if len(segments) < 1 {
		return nil, fmt.Errorf("expected 1 version segment, got %d", len(segments))
	}

	return version.NewVersion(fmt.Sprintf("%d.0.0", segments[0]+1))
}

type minor struct{}

func (m minor) Bump(cv *version.Version) (*version.Version, error) {
	segments := cv.Segments()
	if len(segments) < 2 {
		return nil, fmt.Errorf("expected 2 version segment, got %d", len(segments))
	}

	return version.NewVersion(fmt.Sprintf("%d.%d.0", segments[0], segments[1]+1))
}

type patch struct{}

func (m patch) Bump(cv *version.Version) (*version.Version, error) {
	segments := cv.Segments()
	if len(segments) < 3 {
		return nil, fmt.Errorf("expected 3 version segment, got %d", len(segments))
	}

	return version.NewVersion(fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2]+1))
}

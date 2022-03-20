package gitmono

import (
	"fmt"

	"github.com/hashicorp/go-version"
	ctx "github.com/sermojohn/gitmono"
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

// compareBumpers compares two bumpers.
// Returns -1, 0, or 1 if bumper A is smaller, equal,
// or larger than the bumper B, respectively.
func compareBumpers(bumperA, bumperB ctx.Bumper) (int, error) {
	if bumperA == nil && bumperB == nil {
		return 0, nil
	}
	if bumperA == nil {
		return -1, nil
	}
	if bumperB == nil {
		return 1, nil
	}

	versionOne, err := version.NewVersion("1.0.0")
	if err != nil {
		return 0, err
	}
	versionA, err := bumperA.Bump(versionOne)
	if err != nil {
		return 0, err
	}
	versionB, err := bumperB.Bump(versionOne)
	if err != nil {
		return 0, err
	}
	return versionA.Compare(versionB), nil
}

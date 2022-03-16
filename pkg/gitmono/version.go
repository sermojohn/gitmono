package gitmono

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-version"
	ctx "github.com/sermojohn/gitmono"
)

var (
	// ErrNoCommitsRelease signals that no commit exist so failing the release
	ErrNoCommitsRelease = fmt.Errorf("no commits to release")
)

// Version combines git commands to read and write releases
type Version struct {
	config       *ctx.Config
	logger       ctx.Logger
	tagger       ctx.Tagger
	commitParser ctx.CommitParser
}

// NewVersion creates a new version instance
func NewVersion(config *ctx.Config, logger ctx.Logger, tagger ctx.Tagger, commitParser ctx.CommitParser) *Version {
	return &Version{
		config:       config,
		logger:       logger,
		tagger:       tagger,
		commitParser: commitParser,
	}
}

// GetCurrentVersion retrieves the current version
func (v *Version) GetCurrentVersion() (*ctx.VersionedCommit, error) {
	tags, err := v.tagger.ListProjectVersionTags()
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		log.Printf("no tags found\n")
		return nil, nil
	}

	latestTag := tags[0]
	parsedVersion := v.parseProjectVersion(latestTag)

	commitHash, err := v.logger.CommitHashByRevision(latestTag)
	if err != nil {
		return nil, err
	}

	version, err := v.parseVersion(parsedVersion)
	if err != nil {
		return nil, err
	}

	currentVersion := ctx.VersionedCommit{
		Version:       version,
		Project:       v.config.Project,
		VersionPrefix: v.config.VersionPrefix,
		CommitID:      commitHash,
	}

	log.Printf("current version: '%s'\n", currentVersion.Version)
	return &currentVersion, nil
}

func (v *Version) parseProjectVersion(tag string) string {
	idx := strings.LastIndex(tag, "/")
	if idx == -1 {
		return tag
	}

	return tag[idx+1:]
}

func (v *Version) parseVersion(vv string) (*version.Version, error) {
	var (
		versionPrefix = v.config.VersionPrefix
		versionValue  = vv
	)
	if versionPrefix != "" && strings.HasPrefix(vv, versionPrefix) {
		versionValue = vv[len(versionPrefix):]
	}

	parsedVersion, err := version.NewSemver(versionValue)
	if err != nil {
		return nil, err
	}

	return parsedVersion, nil
}

// ReleaseNewVersion calculates the new version and performs release
//
// Returns an error if there are no new commits for the provided project
func (v *Version) ReleaseNewVersion(commitID string) (*ctx.VersionedCommit, error) {
	currentVersion, err := v.GetCurrentVersion()
	if err != nil {
		return nil, err
	}
	if currentVersion == nil {
		return nil, nil
	}

	newVersionedCommit, err := v.bumpVersion(currentVersion, commitID)
	if err != nil {
		return nil, err
	}

	err = v.createReleaseTag(newVersionedCommit)
	if err != nil {
		return nil, err
	}

	return newVersionedCommit, nil
}

// InitVersion identifies checks if project has version and releases the initial version
func (v *Version) InitVersion(commitID string) (*ctx.VersionedCommit, error) {
	currentVersion, err := v.GetCurrentVersion()
	if err != nil {
		return nil, err
	}
	if currentVersion != nil {
		return nil, nil
	}

	initVersion, err := version.NewSemver("0.1.0")
	if err != nil {
		return nil, err
	}

	newVersionedCommit := ctx.VersionedCommit{
		CommitID:      commitID,
		Project:       v.config.Project,
		Version:       initVersion,
		VersionPrefix: v.config.VersionPrefix,
	}
	err = v.createReleaseTag(&newVersionedCommit)
	if err != nil {
		return nil, err
	}

	return &newVersionedCommit, nil
}

func (v *Version) createReleaseTag(vc *ctx.VersionedCommit) error {
	log.Printf("new tag: %s at: %s\n", vc.GetTag(), vc.CommitID)
	if v.config.DryRun {
		return nil
	}

	return v.tagger.CreateTag(vc)
}

func (v *Version) bumpVersion(currentVersion *ctx.VersionedCommit, commitID string) (*ctx.VersionedCommit, error) {
	newCommits, err := v.logger.Log(currentVersion.CommitID, commitID)
	if err != nil {
		return nil, err
	}
	if len(newCommits) == 0 {
		return nil, ErrNoCommitsRelease
	}

	var bump ctx.Bumper
	for _, commit := range newCommits {
		commitBump := v.commitParser.GetBumperFromCommit(commit)

		res, err := compareBumpers(bump, commitBump)
		if err != nil {
			return nil, err
		}

		// bumper is lower than the commit bumper
		if res == -1 {
			bump = commitBump
		}

		if bump == majorBumper {
			break
		}
	}

	if bump == nil {
		bump = patchBumper
	}
	newVersion, err := bump.Bump(currentVersion.Version)
	if err != nil {
		return nil, err
	}

	newVersionedCommit := ctx.VersionedCommit{
		CommitID:      commitID,
		Version:       newVersion,
		VersionPrefix: currentVersion.VersionPrefix,
		Project:       currentVersion.Project,
	}
	return &newVersionedCommit, nil
}

// compareBumpers compares two bumpers.
// Returns -1, 0, or 1 if bumper A is smaller, equal,
// or larger than the bumper B, respectively.
func compareBumpers(bumperA, bumperB ctx.Bumper) (int, error) {
	if bumperA == nil {
		return -1, nil
	}
	if bumperA == bumperB {
		return 0, nil
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

package gitmono

import (
	"log"
	"strings"

	"github.com/hashicorp/go-version"
	ctx "github.com/sermojohn/gitmono"
)

// Version combines git commands to read and write releases
type Version struct {
	config       *ctx.Config
	logger       ctx.Logger
	tagger       ctx.Tagger
	commitParser ctx.CommitParser
}

// NewVersion creates a new version instance
func NewVersion(monorepo *ctx.MonoRepo, logger ctx.Logger, tagger ctx.Tagger, commitParser ctx.CommitParser) *Version {
	return &Version{
		config:       monorepo.GetConfig(),
		logger:       logger,
		tagger:       tagger,
		commitParser: commitParser,
	}
}

// GetCurrentVersion retrieves the current version
func (v *Version) GetCurrentVersion() (*ctx.VersionedCommit, error) {
	tags, err := v.tagger.Tags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		parsedProject, version := v.parseProjectVersion(tag)
		if !strings.EqualFold(parsedProject, v.config.Project) {
			continue
		}

		commitHash, err := v.logger.CommitHashByRevision(tag)
		if err != nil {
			return nil, err
		}

		parsedVersion, err := v.parseVersion(version)
		if err != nil {
			return nil, err
		}

		currentVersion := ctx.VersionedCommit{
			Version:       parsedVersion,
			VersionPrefix: v.config.VersionPrefix,
			Project:       parsedProject,
			CommitID:      commitHash,
		}

		log.Printf("current version: %v\n", currentVersion)
		return &currentVersion, nil
	}

	return nil, nil
}

func (v *Version) parseProjectVersion(tag string) (string, string) {
	idx := strings.LastIndex(tag, "/")
	if idx == -1 {
		return ".", tag
	}

	return tag[0:idx], tag[idx+1:]
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
	log.Printf("about to create tag: %s to commit: %s\n", vc.GetTag(), vc.CommitID)
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
		return nil, nil
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

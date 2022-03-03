package main

import (
	"github.com/sermojohn/gitmono"
)

type versionReleaseOptions struct {
	CommitID string `short:"c" description:"The commit ID to release versions on"`
	PrintTag bool   `long:"print-tag" description:"Print tag instead of version"`
}
type versionReleaseCommand struct {
	mono    *gitmono.GitMono
	cmdOpts versionReleaseOptions
}

func (vrc *versionReleaseCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vrc.mono)
	newVersion, err := versioner.ReleaseNewVersion(vrc.cmdOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if vrc.cmdOpts.PrintTag {
			printTag(newVersion)
			return nil
		}
		printVersion(newVersion)
	}
	return nil
}

func (vrc *versionReleaseCommand) name() string {
	return "release"
}

func (vrc *versionReleaseCommand) options() interface{} {
	return &vrc.cmdOpts
}

package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

// ReleaseOptions contains options applying to the release command
type ReleaseOptions struct {
	CommitID string `short:"c" description:"The commit ID to release versions on"`
}
type versionReleaseCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vrc *versionReleaseCommand) Execute(args []string) error {
	var releaseOpts ReleaseOptions
	_, err := flags.NewParser(&releaseOpts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	versioner := gitmono.NewVersioner(vrc.mono)
	newVersion, err := versioner.ReleaseNewVersion(releaseOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if vrc.options.PrintTag {
			printTag(newVersion)
			return nil
		}
		printVersion(newVersion)
	}
	return nil
}

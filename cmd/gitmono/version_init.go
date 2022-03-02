package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

// VersionOptions contains options applying to the version command
type VersionOptions struct {
	CommitID string `short:"c" description:"The commit ID to release initial versions on"`
}
type versionInitCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vic *versionInitCommand) Execute(args []string) error {
	var versionOpts VersionOptions
	_, err := flags.NewParser(&versionOpts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	versioner := gitmono.NewVersioner(vic.mono)
	newVersion, err := versioner.InitVersion(versionOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if vic.options.PrintTag {
			printTag(newVersion)
			return nil
		}
		printVersion(newVersion)
	}
	return nil
}

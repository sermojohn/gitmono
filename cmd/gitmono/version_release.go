package main

import (
	"github.com/sermojohn/gitmono"
)

type releaseOptions struct {
	CommitID string `short:"c" description:"The commit ID to release versions on"`
	PrintTag bool   `long:"print-tag" description:"Print tag instead of version"`
}
type releaseCommand struct {
	versioner gitmono.Versioner
	cmdOpts   releaseOptions
}

func newReleaseCommand(versioner gitmono.Versioner) *releaseCommand {
	return &releaseCommand{
		versioner: versioner,
	}
}

func (rc *releaseCommand) Execute(args []string) error {
	newVersion, err := rc.versioner.ReleaseNewVersion(rc.cmdOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if rc.cmdOpts.PrintTag {
			printTag(newVersion)
			return nil
		}
		printVersion(newVersion)
	}
	return nil
}

func (rc *releaseCommand) name() string {
	return "release"
}

func (rc *releaseCommand) options() interface{} {
	return &rc.cmdOpts
}

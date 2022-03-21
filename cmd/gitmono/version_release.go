package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type releaseOptions struct {
	CommitID string `short:"c" default:"HEAD" description:"The commit ID to release versions on"`
	PrintTag bool   `long:"print-tag" description:"Print tag instead of version"`
}
type releaseCommand struct {
	versioner    gitmono.Versioner
	outputWriter io.Writer
	cmdOpts      releaseOptions
}

func newReleaseCommand(versioner gitmono.Versioner, w io.Writer) *releaseCommand {
	return &releaseCommand{
		versioner:    versioner,
		outputWriter: w,
	}
}

// Execute trigger the release command
func (rc *releaseCommand) Execute(args []string) error {
	newVersion, err := rc.versioner.ReleaseNewVersion(rc.cmdOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if rc.cmdOpts.PrintTag {
			printTag(rc.outputWriter, newVersion)
			return nil
		}
		printVersion(rc.outputWriter, newVersion)
	}
	return nil
}

func (rc *releaseCommand) name() string {
	return "release"
}

func (rc *releaseCommand) options() interface{} {
	return &rc.cmdOpts
}

package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type initOptions struct {
	CommitID string `short:"c" description:"The commit ID to release initial versions on"`
	PrintTag bool   `long:"print-tag" description:"Print tag instead of version"`
}
type initCommand struct {
	versioner    gitmono.Versioner
	outputWriter io.Writer
	cmdOpts      initOptions
}

func newInitCommand(versioner gitmono.Versioner, w io.Writer) *initCommand {
	return &initCommand{
		versioner:    versioner,
		outputWriter: w,
	}
}

func (ic *initCommand) Execute(args []string) error {
	newVersion, err := ic.versioner.InitVersion(ic.cmdOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if ic.cmdOpts.PrintTag {
			printTag(ic.outputWriter, newVersion)
			return nil
		}
		printVersion(ic.outputWriter, newVersion)
	}
	return nil
}

func (ic *initCommand) name() string {
	return "init"
}

func (ic *initCommand) options() interface{} {
	return &ic.cmdOpts
}

package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type versionOptions struct {
	PrintTag    bool `long:"print-tag" description:"Print the tag of the current version"`
	PrintCommit bool `long:"print-commit" description:"Print the commit of the current version"`
}
type versionCommand struct {
	versioner gitmono.Versioner
	w         io.Writer
	cmdOpts   versionOptions
}

func newVersionCommand(versioner gitmono.Versioner, w io.Writer) *versionCommand {
	return &versionCommand{
		versioner: versioner,
		w:         w,
	}
}

func (vc *versionCommand) Execute(args []string) error {
	currentVersion, err := vc.versioner.GetCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion != nil {
		if vc.cmdOpts.PrintTag {
			printTag(vc.w, currentVersion)
			return nil
		}
		if vc.cmdOpts.PrintCommit {
			printCommit(vc.w, currentVersion)
			return nil
		}
		printVersion(vc.w, currentVersion)
	}

	return nil
}

func (vc *versionCommand) name() string {
	return "version"
}

func (vc *versionCommand) options() interface{} {
	return &vc.cmdOpts
}

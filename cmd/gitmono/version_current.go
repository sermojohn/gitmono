package main

import (
	"github.com/sermojohn/gitmono"
)

type versionCurrentOptions struct {
	PrintTag bool `long:"print-tag" description:"Print tag instead of version"`
}
type versionCurrentCommand struct {
	mono    *gitmono.GitMono
	cmdOpts versionCurrentOptions
}

func (vcc *versionCurrentCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vcc.mono)
	currentVersion, err := versioner.GetCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion != nil {
		if vcc.cmdOpts.PrintTag {
			printTag(currentVersion)
			return nil
		}
		printVersion(currentVersion)
	}

	return nil
}

func (vcc *versionCurrentCommand) name() string {
	return "version"
}

func (vcc *versionCurrentCommand) options() interface{} {
	return &vcc.cmdOpts
}

package main

import (
	"github.com/sermojohn/gitmono"
)

type versionCurrentCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vcc *versionCurrentCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vcc.mono)
	currentVersion, err := versioner.GetCurrentVersion(vcc.options.Project)
	if err != nil {
		return err
	}

	if currentVersion != nil {
		if vcc.options.PrintTag {
			printTag(currentVersion)
			return nil
		}
		printVersion(currentVersion)
	}

	return nil
}

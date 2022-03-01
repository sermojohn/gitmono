package main

import (
	"github.com/sermojohn/gitmono"
)

type versionCurrentCommand struct {
	mono *gitmono.GitMono
}

func (vcc *versionCurrentCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vcc.mono)
	currentVersion, err := versioner.CurrentVersion()
	checkError(err)

	if currentVersion != nil {
		printVersion(currentVersion)
	}

	return nil
}

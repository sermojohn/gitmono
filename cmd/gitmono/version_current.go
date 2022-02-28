package main

import (
	"github.com/sermojohn/gitmono"
)

type VersionCurrentCommander struct {
	mono *gitmono.GitMono
}

func (vcc *VersionCurrentCommander) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vcc.mono)
	currentVersion, err := versioner.CurrentVersion()
	checkError(err)

	if currentVersion != nil {
		printVersion(currentVersion)
	}

	return nil
}

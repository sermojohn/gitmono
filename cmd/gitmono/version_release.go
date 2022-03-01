package main

import (
	"github.com/sermojohn/gitmono"
)

type versionReleaseCommand struct {
	mono *gitmono.GitMono
}

func (vrc *versionReleaseCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vrc.mono)
	newVersion, err := versioner.NewVersion()
	checkError(err)

	if newVersion != nil {
		printVersion(newVersion)
	}
	return nil
}

package main

import "github.com/sermojohn/gitmono"

type VersionInitCommander struct {
	mono *gitmono.GitMono
}

func (vic *VersionInitCommander) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vic.mono)
	newVersions, err := versioner.InitVersion()
	checkError(err)

	for _, newVersion := range newVersions {
		printVersion(newVersion)
	}
	return nil
}

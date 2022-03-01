package main

import "github.com/sermojohn/gitmono"

type versionInitCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vic *versionInitCommand) Execute(args []string) error {
	versioner := gitmono.NewVersioner(vic.mono)
	newVersions, err := versioner.InitVersion(vic.options.Projects)
	checkError(err)

	for _, newVersion := range newVersions {
		printVersion(newVersion)
	}
	return nil
}

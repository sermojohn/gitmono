package main

import (
	"fmt"

	"github.com/sermojohn/gitmono"
)

type versionReleaseCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vrc *versionReleaseCommand) Execute(args []string) error {
	if len(vrc.options.Projects) != 1 {
		return fmt.Errorf("expected single project")
	}

	versioner := gitmono.NewVersioner(vrc.mono)
	newVersion, err := versioner.ReleaseNewVersion(vrc.options.Projects[0])
	checkError(err)

	if newVersion != nil {
		printVersion(newVersion)
	}
	return nil
}

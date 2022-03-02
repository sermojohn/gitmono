package main

import (
	"fmt"

	"github.com/sermojohn/gitmono"
)

type versionCurrentCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

func (vcc *versionCurrentCommand) Execute(args []string) error {
	if len(vcc.options.Projects) != 1 {
		return fmt.Errorf("expected single project")
	}

	versioner := gitmono.NewVersioner(vcc.mono)
	currentVersion, err := versioner.GetCurrentVersion(vcc.options.Projects[0])
	if err != nil {
		return err
	}

	if currentVersion != nil {
		printVersion(currentVersion)
	}

	return nil
}

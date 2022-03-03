package main

import (
	"fmt"

	"github.com/sermojohn/gitmono"
)

type versionInitOptions struct {
	CommitID string `short:"c" description:"The commit ID to release initial versions on"`
	PrintTag bool   `long:"print-tag" description:"Print tag instead of version"`
}
type versionInitCommand struct {
	mono    *gitmono.GitMono
	cmdOpts versionInitOptions
}

func (vic *versionInitCommand) Execute(args []string) error {
	fmt.Printf("diff called with: %v, opts: %v\n", args, vic.cmdOpts)

	versioner := gitmono.NewVersioner(vic.mono)
	newVersion, err := versioner.InitVersion(vic.cmdOpts.CommitID)
	if err != nil {
		return err
	}

	if newVersion != nil {
		if vic.cmdOpts.PrintTag {
			printTag(newVersion)
			return nil
		}
		printVersion(newVersion)
	}
	return nil
}

func (vic *versionInitCommand) name() string {
	return "init"
}

func (vic *versionInitCommand) options() interface{} {
	return &vic.cmdOpts
}

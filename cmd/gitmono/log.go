package main

import (
	"fmt"

	"github.com/sermojohn/gitmono"
)

type logCommand struct {
	mono    *gitmono.GitMono
	cmdOpts logOptions
}
type logOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

// Execute trigger the log command
func (lc *logCommand) Execute(args []string) error {
	fmt.Printf("diff called with: %v, opts: %v\n", args, lc.cmdOpts)

	logger := gitmono.NewLogger(lc.mono)
	commits, err := logger.Log(lc.cmdOpts.FromRef, lc.cmdOpts.ToRef)
	if err != nil {
		return err
	}

	printCommits(commits)
	return nil
}

func (lc *logCommand) name() string {
	return "log"
}

func (lc *logCommand) options() interface{} {
	return &lc.cmdOpts
}

package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type logCommand struct {
	logger  gitmono.Logger
	w       io.Writer
	cmdOpts logOptions
}
type logOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func newLogCommand(logger gitmono.Logger, w io.Writer) *logCommand {
	return &logCommand{
		logger: logger,
		w:      w,
	}
}

// Execute trigger the log command
func (lc *logCommand) Execute(args []string) error {
	commits, err := lc.logger.Log(lc.cmdOpts.FromRef, lc.cmdOpts.ToRef)
	if err != nil {
		return err
	}

	printCommits(lc.w, commits)
	return nil
}

func (lc *logCommand) name() string {
	return "log"
}

func (lc *logCommand) options() interface{} {
	return &lc.cmdOpts
}

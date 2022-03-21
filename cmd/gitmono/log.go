package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type logCommand struct {
	logger       gitmono.Logger
	outputWriter io.Writer
	cmdOpts      logOptions
}
type logOptions struct {
	FromRef string `short:"f" required:"1" description:"The starting point of reference range"`
	ToRef   string `short:"t" default:"HEAD" description:"The ending point of reference range"`
}

func newLogCommand(logger gitmono.Logger, w io.Writer) *logCommand {
	return &logCommand{
		logger:       logger,
		outputWriter: w,
	}
}

// Execute trigger the log command
func (lc *logCommand) Execute(args []string) error {
	commits, err := lc.logger.Log(lc.cmdOpts.FromRef, lc.cmdOpts.ToRef)
	if err != nil {
		return err
	}

	printCommits(lc.outputWriter, commits)
	return nil
}

func (lc *logCommand) name() string {
	return "log"
}

func (lc *logCommand) options() interface{} {
	return &lc.cmdOpts
}

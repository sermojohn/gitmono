package main

import "github.com/sermojohn/gitmono"

type logCommand struct {
	logger  gitmono.Logger
	cmdOpts logOptions
}
type logOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func newLogCommand(logger gitmono.Logger) *logCommand {
	return &logCommand{
		logger: logger,
	}
}

// Execute trigger the log command
func (lc *logCommand) Execute(args []string) error {
	commits, err := lc.logger.Log(lc.cmdOpts.FromRef, lc.cmdOpts.ToRef)
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

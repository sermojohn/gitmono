package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

type LogCommander struct {
	mono *gitmono.GitMono
}

type LogOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func (lc *LogCommander) Execute(args []string) error {
	var opts LogOptions
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	checkError(err)

	logger := gitmono.NewLogger(lc.mono)
	commits, err := logger.Log(opts.FromRef, opts.ToRef)
	checkError(err)

	printCommits(commits)
	return nil
}

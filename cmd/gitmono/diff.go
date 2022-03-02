package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

type diffCommand struct {
	mono *gitmono.GitMono
}

// DiffOptions contains the options applying to the diff command
type DiffOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func (dc *diffCommand) Execute(args []string) error {
	var diffOpts DiffOptions
	_, err := flags.NewParser(&diffOpts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return err
	}

	differ := gitmono.NewDiffer(dc.mono)
	changedFiles, err := differ.Diff(diffOpts.FromRef, diffOpts.ToRef)
	if err != nil {
		return err
	}

	printFiles(changedFiles)
	return nil
}

func printFiles(files []string) {
	for _, file := range files {
		fmt.Printf("%s\n", file)
	}
}

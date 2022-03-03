package main

import (
	"fmt"

	"github.com/sermojohn/gitmono"
)

type diffOptions struct {
	FromRef string `short:"f" required:"1" description:"The starting point of reference range"`
	ToRef   string `short:"t" required:"1" description:"The ending point of reference range"`
}

type diffCommand struct {
	differ  gitmono.Differ
	cmdOpts diffOptions
}

func newDiffCommand(differ gitmono.Differ) *diffCommand {
	return &diffCommand{
		differ: differ,
	}
}

func (dc *diffCommand) Execute(args []string) error {
	changedFiles, err := dc.differ.Diff(dc.cmdOpts.FromRef, dc.cmdOpts.ToRef)
	if err != nil {
		return err
	}

	printFiles(changedFiles)
	return nil
}

func (dc *diffCommand) name() string {
	return "diff"
}

func (dc *diffCommand) options() interface{} {
	return &dc.cmdOpts
}

func printFiles(files []string) {
	for _, file := range files {
		fmt.Printf("%s\n", file)
	}
}

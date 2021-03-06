package main

import (
	"io"

	"github.com/sermojohn/gitmono"
)

type diffOptions struct {
	FromRef string `short:"f" required:"1" description:"The starting point of reference range"`
	ToRef   string `short:"t" default:"HEAD" description:"The ending point of reference range"`
}

type diffCommand struct {
	differ       gitmono.Differ
	outputWriter io.Writer
	cmdOpts      diffOptions
}

func newDiffCommand(differ gitmono.Differ, w io.Writer) *diffCommand {
	return &diffCommand{
		differ:       differ,
		outputWriter: w,
	}
}

// Execute trigger the diff command
func (dc *diffCommand) Execute(args []string) error {
	changedFiles, err := dc.differ.Diff(dc.cmdOpts.FromRef, dc.cmdOpts.ToRef)
	if err != nil {
		return err
	}

	printFiles(dc.outputWriter, changedFiles)
	return nil
}

func (dc *diffCommand) name() string {
	return "diff"
}

func (dc *diffCommand) options() interface{} {
	return &dc.cmdOpts
}

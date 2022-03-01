package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

type diffCommand struct {
	mono    *gitmono.GitMono
	options *Options
}

type DiffOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func (dc *diffCommand) Execute(args []string) error {
	var diffOpts DiffOptions
	_, err := flags.NewParser(&diffOpts, flags.IgnoreUnknown).Parse()
	checkError(err)

	differ := gitmono.NewDiffer(dc.mono)
	projects, err := differ.Diff(diffOpts.FromRef, diffOpts.ToRef)
	checkError(err)
	printProjects(projects)

	return nil
}

func printProjects(projects []string) {
	for _, project := range projects {
		fmt.Printf("%s\n", project)
	}
}

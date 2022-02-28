package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

type DiffCommander struct {
	mono *gitmono.GitMono
}

type DiffOptions struct {
	FromRef string `short:"f" description:"The starting point of reference range"`
	ToRef   string `short:"t" description:"The ending point of reference range"`
}

func (dc *DiffCommander) Execute(args []string) error {
	var opts DiffOptions
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	checkError(err)

	differ := gitmono.NewDiffer(dc.mono)
	projects, err := differ.Diff(opts.FromRef, opts.ToRef)
	checkError(err)
	printProjects(projects)

	return nil
}

func printProjects(projects []string) {
	for _, project := range projects {
		fmt.Printf("%s\n", project)
	}
}

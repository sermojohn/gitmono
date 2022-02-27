package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

// Options holds the CLI args
type Options struct {
	DiffSubcommand bool     `short:"d" description:"Print modified projects for a reference range"`
	FromRef        string   `short:"f" description:"The starting point of reference range"`
	ToRef          string   `short:"t" description:"The ending point of reference range"`
	Projects       []string `short:"p" description:"The list of project directories to account"`
}

func main() {
	var opts Options

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if opts.DiffSubcommand {
		differ := gitmono.Differ{
			FromRef:  opts.FromRef,
			ToRef:    opts.ToRef,
			Projects: opts.Projects,
		}

		projects, err := differ.Diff()
		if err != nil {
			fmt.Printf("failed with error: %v\n", err)
			os.Exit(1)
		}

		printProjects(projects)

		os.Exit(0)
	}

}

func printProjects(projects []string) {
	for _, project := range projects {
		fmt.Printf("%s\n", project)
	}
}

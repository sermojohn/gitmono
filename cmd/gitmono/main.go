package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogs/git-module"
	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

// Options holds the CLI args
type Options struct {
	DiffSubcommand    bool     `short:"d" description:"Subcommand to print changed projects for a reference range"`
	VersionSubcommand bool     `short:"v" description:"Subcommand to print the current version of a project"`
	LogSubcommand     bool     `short:"l" description:"Subcommand to print the log of commits of a project"`
	FromRef           string   `short:"f" description:"The starting point of reference range"`
	ToRef             string   `short:"t" description:"The ending point of reference range"`
	Projects          []string `short:"p" description:"The list of project directories to account"`
}

func main() {
	var opts Options

	_, err := flags.Parse(&opts)
	checkError(err)

	mono, err := gitmono.OpenCurrentRepo(&gitmono.Config{
		Projects: opts.Projects,
	})
	checkError(err)

	if opts.DiffSubcommand {
		differ := gitmono.NewDiffer(mono)

		projects, err := differ.Diff(opts.FromRef, opts.ToRef)
		checkError(err)

		printProjects(projects)
		os.Exit(0)
	}

	if opts.VersionSubcommand {
		if len(opts.Projects) != 1 {
			fmt.Printf("expected single project to be provided")
			os.Exit(1)
		}

		versioner := gitmono.NewVersioner(mono)

		currentVersion, err := versioner.GetCurrentVersion()
		checkError(err)
		printVersion(currentVersion)

		// newCommits, err := versioner.GetNewCommits()
		// checkError(err)
		// if len(newCommits) != 0 {
		// 	printCommits(newCommits)
		// 	os.Exit(1)
		// }

		os.Exit(0)
	}

	if opts.LogSubcommand {
		if len(opts.Projects) != 1 {
			fmt.Printf("expected single project to be provided")
			os.Exit(1)
		}

		logger := gitmono.NewLogger(mono)

		commits, err := logger.Log(opts.FromRef, opts.ToRef)
		checkError(err)

		printCommits(commits)
		os.Exit(0)
	}
}

func printProjects(projects []string) {
	for _, project := range projects {
		fmt.Printf("%s\n", project)
	}
}

func printCommits(commits []*git.Commit) {
	for _, commit := range commits {
		fmt.Printf("%s %s\n", commit.ID.String(), strings.Trim(commit.Message, "\n"))
	}
}

func printVersion(version *gitmono.VersionedCommit) {
	fmt.Printf("%v\n", version)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

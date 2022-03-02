package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gogs/git-module"
	"github.com/jessevdk/go-flags"
	"github.com/sermojohn/gitmono"
)

// Commands contains the CLI commands
type Commands struct {
	DiffCommand           diffCommand           `command:"diff"`
	LogCommand            logCommand            `command:"log"`
	VersionCurrentCommand versionCurrentCommand `command:"version"`
	VersionReleaseCommand versionReleaseCommand `command:"release"`
	VersionInitCommand    versionInitCommand    `command:"init"`
}

// Options contains the generic options applying to all commands
type Options struct {
	Project       string `short:"p" default:"." description:"The project directory to operate on. Defaults to '.' to operate on the whole repo"`
	Verbose       bool   `short:"v" description:"Enable verbose loggging"`
	DryRun        bool   `long:"dry-run" description:"Do not persist any write action"`
	CommitScheme  string `long:"commit-scheme" default:"common" description:"The scheme parse commit messages with (common, conventional)"`
	VersionPrefix string `long:"version-prefix" default:"" description:"The prefix to prepend to version"`
	PrintTag      bool   `long:"print-tag" description:"Print tag instead of version"`
}

// Config creates the tool configuration from the provided options
func (opts *Options) Config() *gitmono.Config {
	return &gitmono.Config{
		DryRun:        opts.DryRun,
		CommitScheme:  opts.CommitScheme,
		VersionPrefix: opts.VersionPrefix,
		PrintTag:      opts.PrintTag,
		Project:       opts.Project,
	}
}

func main() {
	var opts Options
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	checkError(err, true)

	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	mono, err := gitmono.OpenRepo("./", opts.Config())
	checkError(err, opts.Verbose)

	var commands = Commands{
		DiffCommand:           diffCommand{mono: mono},
		LogCommand:            logCommand{mono: mono},
		VersionCurrentCommand: versionCurrentCommand{mono: mono, options: &opts},
		VersionReleaseCommand: versionReleaseCommand{mono: mono, options: &opts},
		VersionInitCommand:    versionInitCommand{mono: mono, options: &opts},
	}
	_, err = flags.NewParser(&commands, flags.IgnoreUnknown).Parse()
	checkError(err, opts.Verbose)
}

func printCommits(commits []*git.Commit) {
	for _, commit := range commits {
		fmt.Printf("%s %s\n", commit.ID.String(), strings.Trim(commit.Message, "\n"))
	}
}

func printVersion(version *gitmono.VersionedCommit) {
	fmt.Printf("%s\n", version.GetVersion())
}

func printTag(version *gitmono.VersionedCommit) {
	fmt.Printf("%s\n", version.GetTag())
}

func checkError(err error, verbose bool) {
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

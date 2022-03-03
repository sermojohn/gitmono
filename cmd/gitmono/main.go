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

// Options contains the generic options applying to all commands
type Options struct {
	Project       string `short:"p" default:"." description:"The project directory to operate on"`
	Verbose       bool   `short:"v" description:"Enable verbose loggging"`
	DryRun        bool   `long:"dry-run" description:"Do not persist any write action"`
	CommitScheme  string `long:"commit-scheme" default:"common" description:"The scheme parse commit messages with (common, conventional)"`
	VersionPrefix string `long:"version-prefix" default:"" description:"The prefix to prepend to version"`
}

// command internal interface of a command-line command
type command interface {
	flags.Commander
	name() string
	options() interface{}
}

// Config creates the tool configuration from the provided options
func (opts *Options) Config() *gitmono.Config {
	return &gitmono.Config{
		Project:       opts.Project,
		DryRun:        opts.DryRun,
		CommitScheme:  opts.CommitScheme,
		VersionPrefix: opts.VersionPrefix,
	}
}

func main() {
	mono, err := gitmono.OpenRepo("./")
	checkError(err, true)

	var (
		opts     Options
		commands = []command{
			&diffCommand{mono: mono},
			&logCommand{mono: mono},
			&versionInitCommand{mono: mono},
			&versionReleaseCommand{mono: mono},
			&versionCurrentCommand{mono: mono},
		}
		flagsParser = flags.NewParser(&opts, flags.IgnoreUnknown|flags.HelpFlag)
	)

	// inject options to global component
	flagsParser.CommandHandler = func(command flags.Commander, args []string) error {
		mono.Config(opts.Config())
		return command.Execute(args)
	}
	for _, command := range commands {
		cmd, err := flagsParser.AddCommand(command.name(), "", "", command)
		checkError(err, true)

		_, err = cmd.AddGroup(command.name(), "", command.options())
		checkError(err, true)
	}

	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	// parse options and trigger command
	_, err = flagsParser.Parse()
	checkError(err, true)
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

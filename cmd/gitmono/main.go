package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gogs/git-module"
	"github.com/jessevdk/go-flags"
	ctx "github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/pkg/gitmono"
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
func (opts *Options) Config() *ctx.Config {
	return &ctx.Config{
		Project:       opts.Project,
		DryRun:        opts.DryRun,
		CommitScheme:  opts.CommitScheme,
		VersionPrefix: opts.VersionPrefix,
	}
}

func main() {
	monorepo, err := ctx.OpenRepo("./")
	checkError(err)

	ctx := newContext(monorepo)
	var (
		opts     Options
		commands = []command{
			newDiffCommand(ctx.differ),
			newLogCommand(ctx.logger),
			newInitCommand(ctx.versioner),
			newReleaseCommand(ctx.versioner),
			newVersionCommand(ctx.versioner),
		}
		flagsParser = flags.NewParser(&opts, flags.IgnoreUnknown|flags.HelpFlag)
	)

	// inject options to global component
	flagsParser.CommandHandler = func(command flags.Commander, args []string) error {
		monorepo.SetConfig(opts.Config())
		return command.Execute(args)
	}
	for _, command := range commands {
		cmd, err := flagsParser.AddCommand(command.name(), "", "", command)
		checkError(err)

		_, err = cmd.AddGroup(command.name(), "", command.options())
		checkError(err)
	}

	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	// parse options and trigger command
	_, err = flagsParser.Parse()
	checkError(err)
}

func printCommits(commits []*git.Commit) {
	for _, commit := range commits {
		fmt.Printf("%s %s\n", commit.ID.String(), strings.Trim(commit.Message, "\n"))
	}
}

func printVersion(version *ctx.VersionedCommit) {
	fmt.Printf("%s\n", version.GetVersion())
}

func printTag(version *ctx.VersionedCommit) {
	fmt.Printf("%s\n", version.GetTag())
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr)
		os.Exit(1)
	}
}

type context struct {
	versioner ctx.Versioner
	tagger    ctx.Tagger
	differ    ctx.Differ
	logger    ctx.Logger
}

func newContext(monorepo *ctx.MonoRepo) *context {
	log := gitmono.NewLog(monorepo)
	tag := gitmono.NewTag(monorepo)
	diff := gitmono.NewDiff(monorepo)
	commitParse := gitmono.NewCommitParse(monorepo)
	version := gitmono.NewVersion(monorepo, log, tag, commitParse)
	return &context{
		versioner: version,
		tagger:    tag,
		differ:    diff,
		logger:    log,
	}
}

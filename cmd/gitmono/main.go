package main

import (
	"fmt"
	"io"
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
	ctx, err := newContext()
	checkError(err)

	var (
		opts        Options
		flagsParser = flags.NewParser(&opts, flags.IgnoreUnknown|flags.HelpFlag|flags.PrintErrors)
		commands    = []command{
			newDiffCommand(ctx.differ, ctx.logWriter),
			newLogCommand(ctx.logger, ctx.logWriter),
			newInitCommand(ctx.versioner, ctx.logWriter),
			newReleaseCommand(ctx.versioner, ctx.logWriter),
			newVersionCommand(ctx.versioner, ctx.logWriter),
		}
	)

	// inject options to global component
	flagsParser.CommandHandler = func(command flags.Commander, args []string) error {
		*ctx.config = *opts.Config()
		return command.Execute(args)
	}
	for _, command := range commands {
		cmd, err := flagsParser.AddCommand(command.name(), "", "", command)
		checkError(err)

		_, err = cmd.AddGroup(command.name(), "", command.options())
		checkError(err)
	}

	// parse options and trigger command
	_, err = flagsParser.Parse()
	checkError(err)

	// options are ready to use
	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}
}

func printCommits(outputWriter io.Writer, commits []*git.Commit) {
	for _, commit := range commits {
		fmt.Fprintf(outputWriter, "%s %s\n", commit.ID.String(), strings.Trim(commit.Message, "\n"))
	}
}

func printVersion(outputWriter io.Writer, version *ctx.VersionedCommit) {
	fmt.Fprintf(outputWriter, "%s\n", version.GetVersion())
}

func printTag(outputWriter io.Writer, version *ctx.VersionedCommit) {
	fmt.Fprintf(outputWriter, "%s\n", version.GetTag())
}

func printCommit(outputWriter io.Writer, version *ctx.VersionedCommit) {
	fmt.Fprintf(outputWriter, "%s\n", version.CommitID)
}

func printFiles(outputWriter io.Writer, files []string) {
	for _, file := range files {
		fmt.Fprintf(outputWriter, "%s\n", file)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr)
		os.Exit(1)
	}
}

type context struct {
	// components
	versioner ctx.Versioner
	tagger    ctx.Tagger
	differ    ctx.Differ
	logger    ctx.Logger
	logWriter io.Writer
	// state
	config  *ctx.Config
	envVars *ctx.EnvVars
}

func newContext() (*context, error) {
	repo, err := git.Open("./")
	if err != nil {
		return nil, err
	}

	config := &ctx.Config{}
	envVars := loadEnvVars(os.LookupEnv)
	logger := gitmono.NewLog(repo, config)
	tagger := gitmono.NewTag(repo, config, envVars)
	differ := gitmono.NewDiff(repo, config)
	commitParse := gitmono.NewCommitParse(config)
	versioner := gitmono.NewVersion(config, logger, tagger, commitParse)

	return &context{
		config:    config,
		envVars:   envVars,
		logger:    logger,
		tagger:    tagger,
		differ:    differ,
		versioner: versioner,
		logWriter: os.Stdout,
	}, nil
}

func loadEnvVars(loaderFunc func(string) (string, bool)) *ctx.EnvVars {
	envVars := ctx.EnvVars{}
	if value, found := loaderFunc("GIT_COMMITTER_NAME"); found {
		envVars.CommitterName = value
	}

	if value, found := loaderFunc("GIT_COMMITTER_EMAIL"); found {
		envVars.CommitterEmail = value
	}

	return &envVars
}

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

// Commands holds the CLI args
type Commands struct {
	DiffCommand           diffCommand           `command:"diff"`
	LogCommand            logCommand            `command:"log"`
	VersionCurrentCommand versionCurrentCommand `command:"version"`
	VersionReleaseCommand versionReleaseCommand `command:"release"`
	VersionInitCommand    versionInitCommand    `command:"init"`
}
type Options struct {
	Projects      []string `short:"p" description:"The list of project directories to account"`
	Verbose       bool     `short:"v" description:"Enable verbose loggging"`
	DryRun        bool     `long:"dry-run" description:"Do not persist any write action"`
	CommitScheme  string   `long:"commit-scheme" description:"The scheme parse commit messages with"`
	VersionPrefix string   `long:"version-prefix" description:"The prefix to prepend to version"`
}

func (opts *Options) Config() *gitmono.Config {
	return &gitmono.Config{
		Projects:      opts.Projects,
		DryRun:        opts.DryRun,
		CommitScheme:  opts.CommitScheme,
		VersionPrefix: opts.VersionPrefix,
	}
}

func main() {
	mono, err := gitmono.OpenRepo("./")
	checkError(err)

	var opts Options
	_, err = flags.NewParser(&opts, flags.IgnoreUnknown).Parse()
	checkError(err)
	if len(opts.Projects) == 0 {
		opts.Projects = []string{"."}
	}

	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}
	mono.SetConfig(opts.Config())

	var commands = Commands{
		DiffCommand:           diffCommand{mono: mono},
		LogCommand:            logCommand{mono: mono},
		VersionCurrentCommand: versionCurrentCommand{mono: mono},
		VersionReleaseCommand: versionReleaseCommand{mono: mono},
		VersionInitCommand:    versionInitCommand{mono: mono},
	}
	_, err = flags.NewParser(&commands, flags.IgnoreUnknown).Parse()
	checkError(err)
}

func printCommits(commits []*git.Commit) {
	for _, commit := range commits {
		fmt.Printf("%s %s\n", commit.ID.String(), strings.Trim(commit.Message, "\n"))
	}
}

func printVersion(version *gitmono.VersionedCommit) {
	fmt.Printf("%s\n", version.GetTag())
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

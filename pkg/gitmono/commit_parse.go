package gitmono

import (
	"regexp"
	"strings"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
)

// CommitParse implements commit parsing
type CommitParse struct {
	config *gitmono.Config
}

// NewCommitParse creates a new commit parser
func NewCommitParse(config *gitmono.Config) *CommitParse {
	return &CommitParse{
		config: config,
	}
}

// GetBumperFromCommit parses commit message based on scheme
//
// Returns the bumper to use for this commit or nil
func (cp *CommitParse) GetBumperFromCommit(commit *git.Commit) gitmono.Bumper {
	var b gitmono.Bumper
	msg := commit.Message

	switch cp.config.CommitScheme {
	case "conventional":
		b = conventionalCommitParse(msg)
	case "common":
		b = commonCommitParse(msg)
	}

	return b
}

// parseConventionalCommit implements the Conventional Commit scheme. Given a commit message
// it will return the correct version bumper. In the case of non-confirming conventional commit
// it will return nil and the caller will decide what action to take.
// https://www.conventionalcommits.org/en/v1.0.0/#summary
func conventionalCommitParse(msg string) gitmono.Bumper {
	matches := findNamedMatches(conventionalCommitRex, msg)

	// If the commit contains a footer with 'BREAKING CHANGE:' it is always a major bump
	if strings.Contains(msg, "\nBREAKING CHANGE:") {
		return majorBumper
	}

	// if the type/scope in the header includes a trailing '!' this is a breaking change
	if breaking, ok := matches["breaking"]; ok && breaking == "!" {
		return majorBumper
	}

	// if the type in the header is 'feat' it is a minor change
	if typ, ok := matches["type"]; ok && typ == "feat" {
		return minorBumper
	}

	return nil
}

// findNamedMatches is a helper function for use with regexes containing named capture groups.
// It takes a regex and a string and returns a map with keys corresponding to the named captures
// in the regex. If there are no matches the map will be empty.
// https://play.golang.org/p/GR_6YHaEvef
func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

var (
	// autotag commit message scheme:
	majorRex = regexp.MustCompile(`(?i)\[major\]|\#major`)
	minorRex = regexp.MustCompile(`(?i)\[minor\]|\#minor`)
	patchRex = regexp.MustCompile(`(?i)\[patch\]|\#patch`)

	// conventional commit message scheme:
	// https://regex101.com/r/XciTmT/2
	conventionalCommitRex = regexp.MustCompile(`^\s*(?P<type>\w+)(?P<scope>(?:\([^()\r\n]*\)|\()?(?P<breaking>!)?)(?P<subject>:.*)?`)
)

// parseAutotagCommit implements the autotag (default) commit scheme.
// A git commit message header containing:
//  - [major] or #major: major version bump
//  - [minor] or #minor: minor version bump
//  - [patch] or #patch: patch version bump
// If no action is present nil is returned and the caller must decide what action to take.
func commonCommitParse(msg string) gitmono.Bumper {
	if majorRex.MatchString(msg) {
		return majorBumper
	}

	if minorRex.MatchString(msg) {
		return minorBumper
	}

	if patchRex.MatchString(msg) {
		return patchBumper
	}

	return nil
}

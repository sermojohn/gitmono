package gitmono

import (
	"log"
	"regexp"
	"strings"

	"github.com/gogs/git-module"
)

type CommitParser struct {
	scheme string
}

func (cp *CommitParser) parseCommit(commit *git.Commit) bumper {
	var b bumper
	msg := commit.Message

	switch cp.scheme {
	case "conventional":
		b = parseConventionalCommit(msg)
	case "", "autotag":
		b = parseAutotagCommit(msg)
	}

	return b
}

// parseConventionalCommit implements the Conventional Commit scheme. Given a commit message
// it will return the correct version bumper. In the case of non-confirming conventional commit
// it will return nil and the caller will decide what action to take.
// https://www.conventionalcommits.org/en/v1.0.0/#summary
func parseConventionalCommit(msg string) bumper {
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
func parseAutotagCommit(msg string) bumper {
	if majorRex.MatchString(msg) {
		log.Println("major bump")
		return majorBumper
	}

	if minorRex.MatchString(msg) {
		log.Println("minor bump")
		return minorBumper
	}

	if patchRex.MatchString(msg) {
		log.Println("patch bump")
		return patchBumper
	}

	return nil
}
[![CI Status](https://github.com/sermojohn/gitmono/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/sermojohn/gitmono/actions/workflows/ci.yaml?query=branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/sermojohn/gitmono)](https://goreportcard.com/report/github.com/sermojohn/gitmono)
[![Coverage Status](https://coveralls.io/repos/github/sermojohn/gitmono/badge.svg?branch=main)](https://coveralls.io/github/sermojohn/gitmono?branch=main)
![Latest Release](https://img.shields.io/github/v/release/sermojohn/gitmono)

# gitmono
CLI tool that wraps git to manage versioning of projects in monorepos.

All commands accept the monorepo project to operate on, using the `-p` flag. When none is specified, the flag defaults to `.` that enables the tool for any repository type (also non-monorepo). The specified project should match the subdirectory of the project from the root path.

## Commands

`init` command releases and prints the initial version (0.1.0) for the specified project, if it has never been released.

`version` command prints the latest version of the specified project.

`release` command releases and prints a new version of the specified project, after looking up for a version bump command in the commits that reference files of that project.

`diff` command prints out the modified files for the specified git reference range that match the project subdirectory.

`log` command prints out the commits for the specified git reference range that match project subdirectory.

## Examples

1. Initialise versioning of a project:
```
gitmono init -p mod1

0.1.0
```

2. Get latest version of a project:
```
gitmono version -p mod1

0.1.0
```

3. Release new version for a project:
```
gitmono release -p mod1

mod1/0.2.0
```


4. Get the modified files of a project in the provided reference range:
```
gitmono diff -p mod1 -f head~1 -t head

mod1/go.mod
```

5. Get commit log of a project in the provided reference range:
```
gitmono log -p mod1 -f head~2 -t head

7fd4cd0d6141b3cbc3c4f0a2206090140d2d7722 fix: mod1 modified
```

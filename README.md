[![Go Report Card](https://goreportcard.com/badge/github.com/sermojohn/gitmono)](https://goreportcard.com/report/github.com/sermojohn/gitmono)

# gitmono
Git wrapper for monorepos

This tool provides commands to manage versioning of projects in monorepos.

All commands accept the monorepo project to operate on using the `-p` flag. When no specified, it defaults to `.` that enables the tool for non-monorepos.
The project value is expected to match the subdirectory of the project from the root path. 

## Commands

`diff` command returns the modified files for the specified git range belong to the specified project.

`log` command prints out the commits of the specified project.

`init` command releases the initial version (0.1.0) for the specified project, if it has never been released.

`version` command extracts the latest version of the specified project.

`release` command releases a new version of the specified project, after looking up for a version bump command in the commits that reference files of that project.


## Examples

1. Initialise the version of all monorepo projects:
```
gitmono init --version-prefix=v -p mod1

v0.1.0
```

2. Get the latest version of a monorepo project:
```
gitmono version -p mod1 --version-prefix=v

v0.1.0
```

3. Get the subset of modified monorepo projects:
```
gitmono diff -f head~1 -t head -p mod1

mod1/go.mod
```

4. Release new version for a modified monorepo project:
```
gitmono release -p mod1 --commit-scheme=conventional --version-prefix=v

mod1/v0.2.0
```

5. Get log of commits relevant to a monorepo project:
```
gitmono log -p mod1 -f head~2 -t head

7fd4cd0d6141b3cbc3c4f0a2206090140d2d7722 fix: mod1 modified
```


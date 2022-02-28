# gitmono
Git wrapper for monorepos

The following commands take one or more project paths as input (using the `-p` argument).
The projects are any of the monorepo subdirectories and are considered separate projects residing in the monorepo.
This tool provides commands to manage versioning of projects in monorepos.

## Commands

`diff` command prints out the projects whose scope matches to at least one file included in the diff of the specified reference range.

`version` command extracts the latest version of the specified project.

`release` command releases a new version of the specified project, after looking up for a version bump command in the commits that reference files of that project.

`init` command releases a new version (0.1.0) for the specified projects that have never been released.

`log` command prints out the commits of the specified project.

## Examples

1. Initialise the version of all monorepo projects:
```
gitmono init -version-prefix=v -p mod1 -p mod2

mod2/0.1.0
mod1/0.1.0
```

2. Get the latest version of a monorepo project:
```
gitmono version -p mod1

mod1/0.1.0
```

3. Get the subset of modified monorepo projects:
```
gitmono diff -f head~1 -t head -p mod2 -p mod1

mod1
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


# gitmono
Git wrapper for monorepos

The following commands take one or more project paths as input, that serve as the monorepo subdirectories that are considered separate projects residing in the monorepo.
This tool provides commands to manage versioning of projects in monorepos.

## Diff
`diff` command prints out the projects whose scope matches to at least one file included in the diff of the specified reference range.

## Version
`version` command extracts the latest version of the specified project.
`release` command releases a new version of the specified project, after looking up for a version bump command in the commits that reference files of that project.
`init` command releases a new version (0.1.0) for the specified projects that have never been released.
## Log
`log` command prints out the commits of the specified project.

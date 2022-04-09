package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gogs/git-module"
	"github.com/google/go-cmdtest"
	"github.com/stretchr/testify/assert"
)

// TestCommand uses google/go-cmdtest to run commands and verify output
func TestCommand(t *testing.T) {
	t.Parallel()

	cleanup := setupRepo(t)
	defer cleanup()

	tests := []struct {
		name  string
		suite string
	}{
		{
			name:  "init command",
			suite: "testdata/init",
		},
		{
			name:  "version command",
			suite: "testdata/version",
		},
		{
			name:  "release command",
			suite: "testdata/release",
		},
		{
			name:  "diff command",
			suite: "testdata/diff",
		},
		{
			name:  "log command",
			suite: "testdata/log",
		},
		{
			name:  "verbose flag",
			suite: "testdata/verbose",
		},
		{
			name:  "help command",
			suite: "testdata/help",
		},
	}

	for _, tt := range tests {
		ts, err := cmdtest.Read(tt.suite)
		assert.Nil(t, err)

		ts.Commands["gitmono"] = cmdtest.InProcessProgram("gitmono", run)
		ts.Run(t, true)
	}
}

func TestCommandFailures(t *testing.T) {
	ts, err := cmdtest.Read("testdata/repofailure")
	assert.Nil(t, err)

	err = os.Setenv("GIT_REPO_PATH", "invalid_repo_path")
	assert.Nil(t, err)

	ts.Commands["gitmono"] = cmdtest.InProcessProgram("gitmono", run)
	ts.Run(t, true)
}

func setupRepo(t *testing.T) func() {
	// repo path lookup
	wd, err := os.Getwd()
	assert.Nil(t, err)

	repoPath := fmt.Sprintf("%s/%s", wd, "testrepo")
	err = os.Setenv("GIT_REPO_PATH", repoPath)
	assert.Nil(t, err)

	// fetch test repository
	if !isExist(repoPath) {
		err := os.Mkdir(repoPath, os.ModePerm)
		assert.Nil(t, err)

		err = git.Clone("https://github.com/sermojohn/gitmono-testrepo.git", repoPath, git.CloneOptions{})
		assert.Nil(t, err)
	}

	return func() {
		err := os.RemoveAll(repoPath)
		if err != nil {
			fmt.Printf("cleaned clone with err: %v", err)
		}
	}
}

// isExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

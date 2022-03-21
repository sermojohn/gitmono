package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
	"github.com/stretchr/testify/assert"
)

func Test_versionCommand_Execute(t *testing.T) {
	t.Parallel()

	type fields struct {
		versioner    gitmono.Versioner
		outputWriter io.Writer
		cmdOpts      versionOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "get current version tag",
			fields: fields{
				versioner: &mock.Versioner{
					GetCurrentVersionOutput: &gitmono.VersionedCommit{
						Project:       "test",
						CommitID:      "commit-id",
						VersionPrefix: "v",
						Version:       newVersion(t, "1.0.0"),
					},
				},
				outputWriter: &bytes.Buffer{},
				cmdOpts: versionOptions{
					PrintTag: true,
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "test/v1.0.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "get current version commit",
			fields: fields{
				versioner: &mock.Versioner{
					GetCurrentVersionOutput: &gitmono.VersionedCommit{
						Project:       "test",
						CommitID:      "commit-id",
						VersionPrefix: "v",
						Version:       newVersion(t, "1.0.0"),
					},
				},
				outputWriter: &bytes.Buffer{},
				cmdOpts: versionOptions{
					PrintCommit: true,
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "commit-id\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "get current version default output",
			fields: fields{
				versioner: &mock.Versioner{
					GetCurrentVersionOutput: &gitmono.VersionedCommit{
						Project:       "test",
						CommitID:      "commit-id",
						VersionPrefix: "v",
						Version:       newVersion(t, "1.0.0"),
					},
				},
				outputWriter: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "v1.0.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "get current version failure",
			fields: fields{
				versioner: &mock.Versioner{
					GetCurrentVersionError: fmt.Errorf("failed to get current version"),
				},
				outputWriter: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "", f.outputWriter.(*bytes.Buffer).String())
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := newVersionCommand(tt.fields.versioner, tt.fields.outputWriter)
			vc.cmdOpts = tt.fields.cmdOpts

			if err := vc.Execute([]string{}); (err != nil) != tt.wantErr {
				t.Errorf("versionCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

func newVersion(t *testing.T, v string) *version.Version {
	semVer, err := version.NewSemver(v)
	assert.Nil(t, err)
	return semVer
}

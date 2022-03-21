package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
	"github.com/stretchr/testify/assert"
)

func Test_releaseCommand_name(t *testing.T) {
	rc := &releaseCommand{}
	assert.Equal(t, "release", rc.name())
}

func Test_releaseCommand_Execute(t *testing.T) {
	t.Parallel()

	type fields struct {
		versioner    gitmono.Versioner
		outputWriter io.Writer
		cmdOpts      releaseOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "release version and print default output",
			fields: fields{
				versioner: &mock.Versioner{ReleaseNewVersionOutput: &gitmono.VersionedCommit{
					Project:       "test",
					CommitID:      "commit-id",
					VersionPrefix: "v",
					Version:       newVersion(t, "2.0.0"),
				}},
				outputWriter: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "v2.0.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "release version and print tag",
			fields: fields{
				versioner: &mock.Versioner{ReleaseNewVersionOutput: &gitmono.VersionedCommit{
					Project:       "test",
					CommitID:      "commit-id",
					VersionPrefix: "v",
					Version:       newVersion(t, "2.0.0"),
				}},
				outputWriter: &bytes.Buffer{},
				cmdOpts: releaseOptions{
					PrintTag: true,
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "test/v2.0.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "release version failure",
			fields: fields{
				versioner:    &mock.Versioner{ReleaseNewVersionError: fmt.Errorf("failed to release version")},
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
			rc := newReleaseCommand(tt.fields.versioner, tt.fields.outputWriter)
			rc.cmdOpts = tt.fields.cmdOpts

			if err := rc.Execute([]string{}); (err != nil) != tt.wantErr {
				t.Errorf("releaseCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

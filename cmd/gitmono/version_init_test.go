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

func Test_initCommand_Execute(t *testing.T) {
	t.Parallel()

	type fields struct {
		versioner    gitmono.Versioner
		outputWriter io.Writer
		cmdOpts      initOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "inititate version and print default",
			fields: fields{
				versioner: &mock.Versioner{
					InitVersionOutput: &gitmono.VersionedCommit{
						Project:       "test",
						CommitID:      "commit-id",
						VersionPrefix: "v",
						Version:       newVersion(t, "0.1.0"),
					},
				},
				outputWriter: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "v0.1.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "inititate version and print tag",
			fields: fields{
				versioner: &mock.Versioner{
					InitVersionOutput: &gitmono.VersionedCommit{
						Project:       "test",
						CommitID:      "commit-id",
						VersionPrefix: "v",
						Version:       newVersion(t, "0.1.0"),
					},
				},
				outputWriter: &bytes.Buffer{},
				cmdOpts: initOptions{
					PrintTag: true,
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "test/v0.1.0\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "inititate version failure",
			fields: fields{
				versioner: &mock.Versioner{
					InitVersionError: fmt.Errorf("failed to init version"),
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
			ic := newInitCommand(tt.fields.versioner, tt.fields.outputWriter)
			ic.cmdOpts = tt.fields.cmdOpts

			if err := ic.Execute([]string{}); (err != nil) != tt.wantErr {
				t.Errorf("initCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

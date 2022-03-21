package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
	"github.com/stretchr/testify/assert"
)

func Test_logCommand_Execute(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger       gitmono.Logger
		outputWriter io.Writer
		cmdOpts      logOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "log command success",
			fields: fields{
				logger: &mock.Logger{
					LogOutput: []*git.Commit{
						{ID: git.MustIDFromString("2852c690001da9fc36db2eacf9e7ee77585fd7aa"), Message: "commitA"},
						{ID: git.MustIDFromString("2852c690001da9fc36db2eacf9e7ee77585fd7ab"), Message: "commitB"},
					},
				},
				outputWriter: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "2852c690001da9fc36db2eacf9e7ee77585fd7aa commitA\n2852c690001da9fc36db2eacf9e7ee77585fd7ab commitB\n", f.outputWriter.(*bytes.Buffer).String())
			},
		},
		{
			name: "log command failure",
			fields: fields{
				logger: &mock.Logger{
					LogError: fmt.Errorf("failed to perform git log"),
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
			lc := newLogCommand(tt.fields.logger, tt.fields.outputWriter)
			lc.cmdOpts = tt.fields.cmdOpts

			if err := lc.Execute([]string{}); (err != nil) != tt.wantErr {
				t.Errorf("logCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

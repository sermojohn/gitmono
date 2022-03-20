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

func Test_diffCommand_Execute(t *testing.T) {
	t.Parallel()

	type fields struct {
		differ  gitmono.Differ
		w       io.Writer
		cmdOpts diffOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "diff command success",
			fields: fields{
				differ: &mock.Differ{DiffOutput: []string{"test/file1", "test/file2"}},
				cmdOpts: diffOptions{
					FromRef: "ref-1",
					ToRef:   "ref-2",
				},
				w: &bytes.Buffer{},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Equal(t, "test/file1\ntest/file2\n", f.w.(*bytes.Buffer).String())
			},
		},
		{
			name: "diff command failure",
			fields: fields{
				differ: &mock.Differ{DiffError: fmt.Errorf("failed to diff")},
				cmdOpts: diffOptions{
					FromRef: "ref-1",
					ToRef:   "ref-2",
				},
				w: &bytes.Buffer{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &diffCommand{
				differ:  tt.fields.differ,
				w:       tt.fields.w,
				cmdOpts: tt.fields.cmdOpts,
			}
			if err := dc.Execute([]string{}); (err != nil) != tt.wantErr {
				t.Errorf("diffCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

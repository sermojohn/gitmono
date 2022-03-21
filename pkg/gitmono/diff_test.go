package gitmono

import (
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
)

func TestDiff_Diff(t *testing.T) {
	t.Parallel()

	type fields struct {
		config *gitmono.Config
		differ gitmono.GitDiffer
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "fetch project diff",
			fields: fields{
				config: &gitmono.Config{Project: "test"},
				differ: &mock.GitDiffer{DiffOutput: &git.Diff{Files: []*git.DiffFile{{Name: "test/file1"}, {Name: "test2/file2"}}}},
			},
			args: args{
				from: "commit1",
				to:   "commit2",
			},
			want: []string{"test/file1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Diff{
				differ: tt.fields.differ,
				config: tt.fields.config,
			}
			got, err := d.Diff(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff.Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff.Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

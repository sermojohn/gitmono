package gitmono

import (
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	ctx "github.com/sermojohn/gitmono"
)

func Test_CommitParse_GetBumperFromCommit(t *testing.T) {
	t.Parallel()

	type fields struct {
		config *ctx.Config
	}
	type args struct {
		commit *git.Commit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ctx.Bumper
	}{
		{
			name: "no commit message scheme",
			fields: fields{
				config: &ctx.Config{},
			},
			args: args{
				commit: &git.Commit{Message: "any message #major or breaking!"},
			},
			want: nil,
		},
		{
			name: "conventional message scheme",
			fields: fields{
				config: &ctx.Config{
					CommitScheme: "conventional",
				},
			},
			args: args{
				commit: &git.Commit{Message: "breaking!"},
			},
			want: majorBumper,
		},
		{
			name: "common message scheme",
			fields: fields{
				config: &ctx.Config{
					CommitScheme: "common",
				},
			},
			args: args{
				commit: &git.Commit{Message: "#minor"},
			},
			want: minorBumper,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &CommitParse{
				config: tt.fields.config,
			}
			if got := cp.GetBumperFromCommit(tt.args.commit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commitParser.parseCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}

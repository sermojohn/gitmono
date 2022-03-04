package gitmono

import (
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
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

func Test_conventionalCommitParse(t *testing.T) {
	t.Parallel()

	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want gitmono.Bumper
	}{
		{
			name: "use exclamation mark",
			args: args{
				msg: "fix!: bug that introduces breaking API change",
			},
			want: majorBumper,
		},
		{
			name: "use breaking change footer",
			args: args{
				msg: "fix: bug that introduces breaking API change\nBREAKING CHANGE: API change",
			},
			want: majorBumper,
		},
		{
			name: "feature implmenentation",
			args: args{
				msg: "feat: application feature",
			},
			want: minorBumper,
		},
		{
			name: "fix implmenentation",
			args: args{
				msg: "fix: application fix",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conventionalCommitParse(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("conventionalCommitParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonCommitParse(t *testing.T) {
	t.Parallel()

	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want gitmono.Bumper
	}{
		{
			name: "introduce breaking change using #major",
			args: args{
				msg: "introduces breaking API change #major",
			},
			want: majorBumper,
		},
		{
			name: "introduce breaking change using [major]",
			args: args{
				msg: "introduces breaking API change [major]",
			},
			want: majorBumper,
		},
		{
			name: "trigger minor bump using #minor",
			args: args{
				msg: "#minor",
			},
			want: minorBumper,
		},
		{
			name: "trigger minor bump using [minor]",
			args: args{
				msg: "[minor]",
			},
			want: minorBumper,
		},
		{
			name: "trigger patch bump using #patch",
			args: args{
				msg: "#patch",
			},
			want: patchBumper,
		},
		{
			name: "trigger patch bump using [patch]",
			args: args{
				msg: "[patch]",
			},
			want: patchBumper,
		},
		{
			name: "no commit message bump",
			args: args{
				msg: "test message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commonCommitParse(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commonCommitParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

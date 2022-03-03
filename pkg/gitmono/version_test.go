package gitmono

import (
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	"github.com/hashicorp/go-version"
	ctx "github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"

	"github.com/stretchr/testify/assert"
)

func TestVersion_bumpVersion(t *testing.T) {
	t.Parallel()
	type fields struct {
		config       *ctx.Config
		logger       ctx.Logger
		tagger       ctx.Tagger
		commitParser ctx.CommitParser
	}
	type args struct {
		currentVersion *ctx.VersionedCommit
		commitID       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ctx.VersionedCommit
		wantErr bool
	}{
		{
			name: "no commits",
			fields: fields{
				logger: &mock.Logger{
					LogOutput: []*git.Commit{},
				},
			},
			args: args{
				currentVersion: &ctx.VersionedCommit{
					CommitID: "1",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "one commit",
			fields: fields{
				config: &ctx.Config{},
				logger: &mock.Logger{
					LogOutput: []*git.Commit{
						{Message: "test"},
					},
				},
				commitParser: &mock.CommitParser{
					GetBumperFromCommitOutput: majorBumper,
				},
			},
			args: args{
				currentVersion: &ctx.VersionedCommit{
					CommitID: "1",
					Version:  newVersion(t, "1.0.0"),
				},
				commitID: "2",
			},
			want: &ctx.VersionedCommit{
				CommitID: "2",
				Version:  newVersion(t, "2.0.0"),
			},
			wantErr: false,
		},
		{
			name: "multiple commits",
			fields: fields{
				config: &ctx.Config{},
				logger: &mock.Logger{
					LogOutput: []*git.Commit{
						{Message: "patch"},
						{Message: "minor"},
						{Message: "major"},
					},
				},
				commitParser: &mock.CommitParser{
					GetBumperFromCommitOutputMap: map[string]ctx.Bumper{
						"patch": patchBumper,
						"minor": minorBumper,
						"major": majorBumper,
					},
				},
			},
			args: args{
				currentVersion: &ctx.VersionedCommit{
					CommitID: "1",
					Version:  newVersion(t, "1.0.0"),
				},
				commitID: "2",
			},
			want: &ctx.VersionedCommit{
				CommitID: "2",
				Version:  newVersion(t, "2.0.0"),
			},
			wantErr: false,
		},
		{
			name: "higher commit bump wins",
			fields: fields{
				config: &ctx.Config{},
				logger: &mock.Logger{
					LogOutput: []*git.Commit{
						{Message: "minor"},
						{Message: "patch"},
						{Message: "major"},
						{Message: "minor"},
					},
				},
				commitParser: &mock.CommitParser{
					GetBumperFromCommitOutputMap: map[string]ctx.Bumper{
						"patch": patchBumper,
						"minor": minorBumper,
						"major": majorBumper,
					},
				},
			},
			args: args{
				currentVersion: &ctx.VersionedCommit{
					CommitID: "1",
					Version:  newVersion(t, "1.0.0"),
				},
				commitID: "2",
			},
			want: &ctx.VersionedCommit{
				CommitID: "2",
				Version:  newVersion(t, "2.0.0"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				config:       tt.fields.config,
				logger:       tt.fields.logger,
				tagger:       tt.fields.tagger,
				commitParser: tt.fields.commitParser,
			}
			got, err := v.bumpVersion(tt.args.currentVersion, tt.args.commitID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.bumpVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.bumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newVersion(t *testing.T, v string) *version.Version {
	semVer, err := version.NewSemver(v)
	assert.Nil(t, err)
	return semVer
}

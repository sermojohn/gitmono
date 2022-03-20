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
			wantErr: true,
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
		},
		{
			name: "default bump on no commit bumper",
			fields: fields{
				config: &ctx.Config{},
				logger: &mock.Logger{
					LogOutput: []*git.Commit{
						{Message: "test"},
					},
				},
				commitParser: &mock.CommitParser{},
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
				Version:  newVersion(t, "1.0.1"),
			},
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

func TestVersion_GetCurrentVersion(t *testing.T) {
	t.Parallel()

	type fields struct {
		config *ctx.Config
		logger ctx.Logger
		tagger ctx.Tagger
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ctx.VersionedCommit
		wantErr bool
	}{
		{
			name: "current version using tag default ordering",
			fields: fields{
				config: &ctx.Config{Project: "test2", VersionPrefix: "v"},
				tagger: &mock.Tagger{ListProjectVersionTagsOutput: []string{"test2/v1.0.0", "test2/v0.1.0"}},
				logger: &mock.Logger{},
			},
			want: &ctx.VersionedCommit{Project: "test2", VersionPrefix: "v", Version: newVersion(t, "1.0.0")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				config: tt.fields.config,
				logger: tt.fields.logger,
				tagger: tt.fields.tagger,
			}
			got, err := v.GetCurrentVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.GetCurrentVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Version, tt.want.Version) {
				t.Errorf("Version.GetCurrentVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_ReleaseNewVersion(t *testing.T) {
	t.Parallel()

	type fields struct {
		config       *ctx.Config
		logger       ctx.Logger
		tagger       ctx.Tagger
		commitParser ctx.CommitParser
	}
	type args struct {
		commitID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ctx.VersionedCommit
		wantErr bool
	}{
		{
			name: "release new project version",
			fields: fields{
				config:       &ctx.Config{Project: "test2", VersionPrefix: "v", CommitScheme: "conventional"},
				tagger:       &mock.Tagger{ListProjectVersionTagsOutput: []string{"test2/v1.0.0", "test2/v0.1.0"}},
				logger:       &mock.Logger{CommitHashByRevisionOutput: "test2-version-commit-id", LogOutput: []*git.Commit{{Message: "feat: test"}}},
				commitParser: &mock.CommitParser{GetBumperFromCommitOutput: minorBumper},
			},
			args: args{
				commitID: "test2-new-commit-id",
			},
			want: &ctx.VersionedCommit{
				Project:       "test2",
				VersionPrefix: "v",
				Version:       newVersion(t, "1.1.0"),
				CommitID:      "test2-new-commit-id",
			},
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
			got, err := v.ReleaseNewVersion(tt.args.commitID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.ReleaseNewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.ReleaseNewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_createReleaseTag(t *testing.T) {
	t.Parallel()

	type fields struct {
		config       *ctx.Config
		logger       ctx.Logger
		tagger       ctx.Tagger
		commitParser ctx.CommitParser
	}
	type args struct {
		vc *ctx.VersionedCommit
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "prevent tagging on dry-run enabled",
			fields: fields{
				config: &ctx.Config{
					DryRun: true,
				},
				tagger: &mock.Tagger{CreateTagInputs: []*ctx.VersionedCommit{}},
			},
			args: args{
				vc: &ctx.VersionedCommit{
					Version: newVersion(t, "2.0.0"),
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Len(t, f.tagger.(*mock.Tagger).CreateTagInputs, 0)
			},
			wantErr: false,
		},
		{
			name: "perform tagging on dry-run disabled",
			fields: fields{
				config: &ctx.Config{},
				tagger: &mock.Tagger{CreateTagInputs: []*ctx.VersionedCommit{}},
			},
			args: args{
				vc: &ctx.VersionedCommit{
					Version: newVersion(t, "2.0.0"),
				},
			},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Len(t, f.tagger.(*mock.Tagger).CreateTagInputs, 1)
				assert.Equal(t, &ctx.VersionedCommit{Version: newVersion(t, "2.0.0")}, f.tagger.(*mock.Tagger).CreateTagInputs[0])
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
			if err := v.createReleaseTag(tt.args.vc); (err != nil) != tt.wantErr {
				t.Errorf("Version.createReleaseTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

func TestVersion_InitVersion(t *testing.T) {
	t.Parallel()

	type fields struct {
		config *ctx.Config
		tagger ctx.Tagger
		logger ctx.Logger
	}
	type args struct {
		commitID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ctx.VersionedCommit
		wantErr bool
	}{
		{
			name: "perform init version for new project",
			fields: fields{
				config: &ctx.Config{Project: "test"},
				tagger: &mock.Tagger{ListProjectVersionTagsOutput: []string{}},
			},
			args: args{
				commitID: "tag-commit-id",
			},
			want: &ctx.VersionedCommit{
				CommitID: "tag-commit-id",
				Project:  "test",
				Version:  newVersion(t, "0.1.0"),
			},
		},
		{
			name: "no action for existing project",
			fields: fields{
				config: &ctx.Config{Project: "test"},
				tagger: &mock.Tagger{ListProjectVersionTagsOutput: []string{"test/0.1.0"}},
				logger: &mock.Logger{},
			},
			args: args{
				commitID: "tag-commit-id",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				config: tt.fields.config,
				tagger: tt.fields.tagger,
				logger: tt.fields.logger,
			}
			got, err := v.InitVersion(tt.args.commitID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.InitVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.InitVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

package gitmono

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectTagPrefix(t *testing.T) {
	t.Parallel()

	type args struct {
		project string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single project",
			args: args{
				project: ".",
			},
			want: "",
		},
		{
			name: "monorepo project",
			args: args{
				project: "test",
			},
			want: "test/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectTagPrefix(tt.args.project); got != tt.want {
				t.Errorf("GetProjectTagPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionedCommit_GetTag(t *testing.T) {
	t.Parallel()

	type fields struct {
		CommitID      string
		Project       string
		VersionPrefix string
		Version       *version.Version
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "single project tag",
			fields: fields{
				Project: ".",
				Version: newVersion(t, "1.0.0"),
			},
			want: "1.0.0",
		},
		{
			name: "monorepo project tag",
			fields: fields{
				Project: "test",
				Version: newVersion(t, "1.0.0"),
			},
			want: "test/1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &VersionedCommit{
				CommitID:      tt.fields.CommitID,
				Project:       tt.fields.Project,
				VersionPrefix: tt.fields.VersionPrefix,
				Version:       tt.fields.Version,
			}
			if got := vc.GetTag(); got != tt.want {
				t.Errorf("VersionedCommit.GetTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newVersion(t *testing.T, v string) *version.Version {
	semVer, err := version.NewSemver(v)
	assert.Nil(t, err)
	return semVer
}

func TestVersionedCommit_GetVersion(t *testing.T) {
	type fields struct {
		CommitID      string
		Project       string
		VersionPrefix string
		Version       *version.Version
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no prefix",
			fields: fields{
				Version: newVersion(t, "1.0.0"),
			},
			want: "1.0.0",
		},
		{
			name: "with prefix",
			fields: fields{
				Version:       newVersion(t, "1.0.0"),
				VersionPrefix: "v",
			},
			want: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &VersionedCommit{
				CommitID:      tt.fields.CommitID,
				Project:       tt.fields.Project,
				VersionPrefix: tt.fields.VersionPrefix,
				Version:       tt.fields.Version,
			}
			if got := vc.GetVersion(); got != tt.want {
				t.Errorf("VersionedCommit.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

package gitmono

import (
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestTag_CreateTag(t *testing.T) {
	t.Parallel()

	type fields struct {
		tagger  gitmono.GitTagger
		config  *gitmono.Config
		envVars *gitmono.EnvVars
	}
	type args struct {
		versionedCommit *gitmono.VersionedCommit
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "create annotated tag",
			fields: fields{
				tagger:  &mock.GitTagger{},
				config:  &gitmono.Config{},
				envVars: &gitmono.EnvVars{},
			},
			args: args{
				versionedCommit: &gitmono.VersionedCommit{
					Version: newVersion(t, "1.0.0"),
					Project: ".",
				},
			},
			wantErr: false,
			assertFunc: func(t *testing.T, f *fields) {
				assert.Len(t, f.tagger.(*mock.GitTagger).CreateTagInputs, 1)
				tag := f.tagger.(*mock.GitTagger).CreateTagInputs[0]
				assert.Equal(t, "1.0.0", tag.Name)
			},
		},
		{
			name: "create annotated tag with committer",
			fields: fields{
				tagger: &mock.GitTagger{},
				config: &gitmono.Config{},
				envVars: &gitmono.EnvVars{
					CommitterName:  "alice",
					CommitterEmail: "alice@example.com",
				},
			},
			args: args{
				versionedCommit: &gitmono.VersionedCommit{
					Version: newVersion(t, "1.0.0"),
					Project: ".",
				},
			},
			wantErr: false,
			assertFunc: func(t *testing.T, f *fields) {
				assert.Len(t, f.tagger.(*mock.GitTagger).CreateTagInputs, 1)
				tag := f.tagger.(*mock.GitTagger).CreateTagInputs[0]
				assert.Equal(t, "1.0.0", tag.Name)
				assert.Len(t, tag.Opts, 1)
				assert.Equal(t, &git.Signature{
					Name:  "alice",
					Email: "alice@example.com",
				}, tag.Opts[0].Author)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tag{
				tagger:  tt.fields.tagger,
				config:  tt.fields.config,
				envVars: tt.fields.envVars,
			}
			if err := tr.CreateTag(tt.args.versionedCommit); (err != nil) != tt.wantErr {
				t.Errorf("Tag.CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

func TestTag_ListProjectTags(t *testing.T) {
	t.Parallel()

	type fields struct {
		tagger gitmono.GitTagger
		config *gitmono.Config
	}
	tests := []struct {
		name       string
		fields     fields
		want       []string
		wantErr    bool
		assertFunc func(*testing.T, *fields)
	}{
		{
			name: "list tags",
			fields: fields{
				tagger: &mock.GitTagger{
					TagsOutput: []string{"project1/tag1", "project1/tag2"},
				},
				config: &gitmono.Config{Project: "project1"},
			},
			want: []string{"project1/tag1", "project1/tag2"},
			assertFunc: func(t *testing.T, f *fields) {
				assert.Len(t, f.tagger.(*mock.GitTagger).TagsInputs, 1)
				assert.Len(t, f.tagger.(*mock.GitTagger).TagsInputs[0].Opts, 1)
				assert.Equal(t, "project1/v*", f.tagger.(*mock.GitTagger).TagsInputs[0].Opts[0].Pattern)
				assert.Equal(t, "-version:refname", f.tagger.(*mock.GitTagger).TagsInputs[0].Opts[0].SortKey)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tag{
				tagger: tt.fields.tagger,
				config: tt.fields.config,
			}
			got, err := tr.ListProjectTags()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tag.ListProjectTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tag.ListProjectTags() = %v, want %v", got, tt.want)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t, &tt.fields)
			}
		})
	}
}

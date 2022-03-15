package main

import (
	"reflect"
	"testing"

	ctx "github.com/sermojohn/gitmono"
)

func Test_loadEnvVars(t *testing.T) {
	t.Parallel()

	type args struct {
		loaderFunc func(string) (string, bool)
	}
	tests := []struct {
		name string
		args args
		want *ctx.EnvVars
	}{
		{
			name: "no vars",
			args: args{
				loaderFunc: func(s string) (string, bool) {
					return "", false
				},
			},
			want: &ctx.EnvVars{},
		},
		{
			name: "git committer vars",
			args: args{
				loaderFunc: func(s string) (string, bool) {
					if s == "GIT_COMMITTER_NAME" {
						return "alice", true
					}
					if s == "GIT_COMMITTER_EMAIL" {
						return "alice@example.com", true
					}
					return "", false
				},
			},
			want: &ctx.EnvVars{
				CommitterName:  "alice",
				CommitterEmail: "alice@example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadEnvVars(tt.args.loaderFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

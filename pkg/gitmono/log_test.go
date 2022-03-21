package gitmono

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gogs/git-module"
	"github.com/sermojohn/gitmono"
	"github.com/sermojohn/gitmono/internal/mock"
)

func TestLog_CommitHashByRevision(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger gitmono.GitLogger
	}
	type args struct {
		rev string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "fetch commit hash",
			fields: fields{
				logger: &mock.GitLogger{
					CommitByRevisionOutput: &git.Commit{ID: git.MustIDFromString("2852c690001da9fc36db2eacf9e7ee77585fd7aa")},
				},
			},
			want: "2852c690001da9fc36db2eacf9e7ee77585fd7aa",
		},
		{
			name: "fetch commit hash failure",
			fields: fields{
				logger: &mock.GitLogger{
					CommitByRevisionError: fmt.Errorf("failed to fetch commit"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			got, err := l.CommitHashByRevision(tt.args.rev)
			if (err != nil) != tt.wantErr {
				t.Errorf("Log.CommitHashByRevision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Log.CommitHashByRevision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLog_Log(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger gitmono.GitLogger
		config *gitmono.Config
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*git.Commit
		wantErr bool
	}{
		{
			name: "get project commit log",
			fields: fields{
				config: &gitmono.Config{
					Project: "test",
				},
				logger: &mock.GitLogger{
					LogOutput: []*git.Commit{{Message: "test"}},
				},
			},
			want: []*git.Commit{{Message: "test"}},
		},
		{
			name: "get project commit log error",
			fields: fields{
				config: &gitmono.Config{
					Project: "test",
				},
				logger: &mock.GitLogger{
					LogError: fmt.Errorf("failed to fetch commit log"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
				config: tt.fields.config,
			}
			got, err := l.Log(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Log.Log() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Log.Log() = %v, want %v", got, tt.want)
			}
		})
	}
}

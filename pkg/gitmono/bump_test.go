package gitmono

import (
	"runtime"
	"testing"

	"github.com/hashicorp/go-version"
	ctx "github.com/sermojohn/gitmono"
)

// whitebox testing for autotag bump interface
func TestMinorBumper(t *testing.T) {
	for k, v := range map[string]string{
		"1":                  "1.1.0",
		"1.0":                "1.1.0",
		"1.0.0":              "1.1.0",
		"1.0.12":             "1.1.0",
		"1.0.0-patch":        "1.1.0",
		"1.0.0+build123":     "1.1.0",
		"1.0.0+build123.foo": "1.1.0",
		"1.0.0.0":            "1.1.0",
	} {
		tv, err := version.NewVersion(k)
		checkFatal(t, err)

		nv, err := minorBumper.Bump(tv)
		checkFatal(t, err)

		if nv.String() != v {
			t.Fatalf("Expected '%s' got '%s'", v, nv.String())
		}
	}
}

func TestPatchBumper(t *testing.T) {
	// in retro this didn't have to be a map, but w/e
	for k, v := range map[string]string{
		"1":                      "1.0.1",
		"1.0":                    "1.0.1",
		"1.0.0":                  "1.0.1",
		"1.0.0-patch":            "1.0.1",
		"1.0.0+build123":         "1.0.1",
		"1.0.0+build123.foo.bar": "1.0.1",
	} {
		tv, err := version.NewVersion(k)
		checkFatal(t, err)

		nv, err := patchBumper.Bump(tv)
		checkFatal(t, err)

		if nv.String() != v {
			t.Fatalf("Expected '%s' got '%s'", v, nv.String())
		}
	}
}

func TestMajorBumper(t *testing.T) {
	for k, v := range map[string]string{
		"1":                  "2.0.0",
		"1.0":                "2.0.0",
		"1.1":                "2.0.0",
		"1.0.0":              "2.0.0",
		"1.1.0":              "2.0.0",
		"1.0.0-patch":        "2.0.0",
		"1.0.0+build123":     "2.0.0",
		"1.0.0+build123.foo": "2.0.0",
		"1.0.12":             "2.0.0",
	} {
		tv, err := version.NewVersion(k)
		checkFatal(t, err)

		nv, err := majorBumper.Bump(tv)
		checkFatal(t, err)

		if nv.String() != v {
			t.Fatalf("Expected '%s' got '%s'", v, nv.String())
		}
	}
}

func checkFatal(t *testing.T, err error) {
	if err == nil {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Unable to get caller")
	}
	t.Fatalf("Fail at %v:%v; %v", file, line, err)
}

func TestCompareBumpers(t *testing.T) {
	t.Parallel()

	type args struct {
		bumperA ctx.Bumper
		bumperB ctx.Bumper
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "empty bumpers",
			args: args{},
			want: 0,
		},
		{
			name: "one higher feat bumper",
			args: args{
				bumperA: patchBumper,
				bumperB: majorBumper,
			},
			want: -1,
		},
		{
			name: "one higher feat bumper",
			args: args{
				bumperA: majorBumper,
				bumperB: minorBumper,
			},
			want: 1,
		},
		{
			name: "equal bumpers",
			args: args{
				bumperA: minorBumper,
				bumperB: minorBumper,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compareBumpers(tt.args.bumperA, tt.args.bumperB)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareBumpers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareBumpers() = %v, want %v", got, tt.want)
			}
		})
	}
}

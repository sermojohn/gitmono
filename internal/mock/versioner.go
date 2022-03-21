package mock

import "github.com/sermojohn/gitmono"

type Versioner struct {
	GetCurrentVersionError  error
	GetCurrentVersionOutput *gitmono.VersionedCommit
	ReleaseNewVersionError  error
	ReleaseNewVersionOutput *gitmono.VersionedCommit
	ReleaseNewVersionInputs []string
	InitVersionError        error
	InitVersionOutput       *gitmono.VersionedCommit
	InitVersionInputs       []string
}

func (v *Versioner) GetCurrentVersion() (*gitmono.VersionedCommit, error) {
	if v.GetCurrentVersionError != nil {
		return nil, v.GetCurrentVersionError
	}

	return v.GetCurrentVersionOutput, nil
}

func (v *Versioner) ReleaseNewVersion(commitID string) (*gitmono.VersionedCommit, error) {
	if v.ReleaseNewVersionError != nil {
		return nil, v.ReleaseNewVersionError
	}

	v.ReleaseNewVersionInputs = append(v.ReleaseNewVersionInputs, commitID)
	return v.ReleaseNewVersionOutput, nil
}

func (v *Versioner) InitVersion(commitID string) (*gitmono.VersionedCommit, error) {
	if v.InitVersionError != nil {
		return nil, v.InitVersionError
	}

	v.InitVersionInputs = append(v.InitVersionInputs, commitID)
	return v.InitVersionOutput, nil
}

package gitmono

type Tagger struct {
	mono *GitMono
}

func (t *Tagger) Tags() ([]string, error) {
	return t.mono.repo.Tags()
}

func (t *Tagger) WriteTag(tagCommit *VersionedCommit) error {
	return t.mono.repo.CreateTag(tagCommit.GetTag(), tagCommit.CommitID)
}

package gitmono

// Tagger performs tag operation for a monorepo
type Tagger struct {
	mono *GitMono
}

// Tags retrieves all repository tags
func (t *Tagger) Tags() ([]string, error) {
	return t.mono.repo.Tags()
}

// CreateTag create a tag on the provided commit
func (t *Tagger) CreateTag(versionedCommit *VersionedCommit) error {
	return t.mono.repo.CreateTag(versionedCommit.GetTag(), versionedCommit.CommitID)
}

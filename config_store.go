package poller

import (
	"io/ioutil"
	"sync"
)

// A Store defines a place where configuration can be loaded/persisted.
type Store interface {
	// Read configuration and returns the complete list of checks.
	Load() (*CheckList, error)

	// Persist `checks`.
	Persist(*CheckList) error
}

// The InMemoryStore keeps checks in memory.
type inMemoryStore struct {
	checks *CheckList
}

type fileStore struct {
	filename string
	checks   *CheckList
	mu       sync.Mutex
}

func NewInMemoryStore() Store {
	return &inMemoryStore{NewCheckList()}
}

func (ims *inMemoryStore) Load() (*CheckList, error) {
	return ims.checks, nil
}

func (ims *inMemoryStore) Persist(checks *CheckList) error {
	ims.checks = checks
	return nil
}

// NewFileStore instantiates a Store where data are saved
// on disk. Configuration is stored as a unique json file.
func NewFileStore(filename string) Store {
	return &fileStore{filename: filename, checks: NewCheckList()}
}

func (fs *fileStore) Load() (*CheckList, error) {
	data, err := ioutil.ReadFile(fs.filename)
	if err != nil {
		return nil, err
	}

	return NewCheckListFromJSON(data)
}

func (fs *fileStore) Persist(cl *CheckList) error {
	data, err := cl.JSON()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fs.filename, data, 0644)
}

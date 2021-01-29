package maker

import (
	"errors"
	"fmt"
	"io"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Config stores the snippets required and any extra Maker setting
type Config struct {
	Repositories map[string]*Repository `yaml:"repositories"`
}

type Repository struct {
	*git.Repository `yaml:"-"`

	Snippets map[string]string `yaml:"snippets"`
	URL      string            `yaml:"url"`

	sync.Mutex `yaml:"-"`
}

func (repository *Repository) Init() error {
	if repository.Repository != nil {
		return nil
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repository.URL,
	})
	if err != nil {
		return err
	}

	repository.Repository = repo
	return nil
}

func (repository *Repository) changeReference(reference string) error {
	_, err := repository.ResolveRevision(plumbing.Revision(reference))

	return err
}

func (repository *Repository) isOnReference(reference string) (bool, error) {
	currentReference, err := repository.Head()
	if err != nil {
		return false, err
	}

	if currentReference.Target().String() == reference {
		return true, nil
	}

	if currentReference.Hash().String() == reference {
		return true, nil
	}

	return false, err
}

// Get returns the snippet file contents reader. The caller is responsible for
// closing it after usage.
func (repository *Repository) Get(reference, name string) (io.ReadCloser, error) {
	// lock to ensure we retrieve the file from the proper reference
	repository.Lock()
	defer repository.Unlock()

	is, err := repository.isOnReference(reference)
	if err != nil {
		return nil, err
	}

	if !is {
		err = repository.changeReference(reference)
		if err != nil {
			return nil, err
		}
	}

	trees, err := repository.TreeObjects()
	if err != nil {
		return nil, err
	}

	var file *object.File
	err = trees.ForEach(func(t *object.Tree) error {
		file, err = t.File(fmt.Sprintf("snippets/%s.mk", name))
		if errors.Is(err, object.ErrFileNotFound) {
			return nil
		}

		if err == nil && file != nil {
			return storer.ErrStop
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return file.Reader()
}

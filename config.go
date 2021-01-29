package maker

import "fmt"

// Config stores the snippets required and any extra Maker setting
type Config struct {
	Repositories []*Repository `yaml:"repositories"`
}

func (config *Config) AddRepository(repo *Repository) error {
	if repo == nil {
		return fmt.Errorf("no repository provided")
	}

	err := repo.Init()
	if err != nil {
		return err
	}

	// TODO lock before appending
	config.Repositories = append(config.Repositories, repo)

	return nil
}

// GetRepository returns a repository for the given reference, which can be
// either an alias or an URL. An empty reference returns the first
// (i.e. default) entry found in the configuration.
func (config *Config) GetRepository(reference string) (*Repository, error) {
	if len(config.Repositories) == 0 {
		return nil, fmt.Errorf("no repositories configured")
	}

	if reference == "" {
		return config.Repositories[0], nil
	}

	for _, repository := range config.Repositories {
		if repository.Alias == reference {
			return repository, nil
		}

		if repository.URL == reference {
			return repository, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}

package maker

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/wwmoraes/maker/pkg/semver"
)

const (
	SnippetsDirectory = ".make"
	ConfFilename      = "maker.yaml"
	LockFilename      = "maker.lock"
)

// Maker manages the configuration, lock data and snippet files on a directory
type Maker struct {
	configFile File
	lockFile   File
	directory  billy.Filesystem
	conf       Config
	lock       Lock
}

func closeDescriptor(fd UnlockCloser) error {
	err := fd.Unlock()
	if err != nil {
		return err
	}

	return fd.Close()
}

// NewDefault creates a standard Maker instance using the OS filesystem and the
// current directory as the repository root.
func NewDefault() (mk *Maker, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root := osfs.New(pwd)

	// always try to make the directory
	err = root.MkdirAll(SnippetsDirectory, 0750)
	if err != nil && !os.IsExist(err) {
		return mk, err
	}

	var info fs.FileInfo
	info, err = root.Stat(SnippetsDirectory)
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsDir() {
		return nil, fmt.Errorf("%s is not a valid directory", SnippetsDirectory)
	}

	// make sure we can RWX on the target directory
	if info.Mode()&0700 == 0 {
		return nil, fmt.Errorf("%s directory must be readable, writable and executable", SnippetsDirectory)
	}

	snippetsFS, err := root.Chroot(SnippetsDirectory)
	if err != nil {
		return nil, err
	}

	// check if the maker config file is valid
	info, err = root.Stat(ConfFilename)
	if err != nil && !os.IsNotExist(err) {
		return mk, err
	}

	if err == nil && !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a valid Maker file", ConfFilename)
	}

	// check if the maker lock file is valid
	info, err = root.Stat(LockFilename)
	if err != nil && !os.IsNotExist(err) {
		return mk, err
	}

	if err == nil && !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a valid Maker lock file", LockFilename)
	}

	// open the maker file for usage
	confFD, err := root.OpenFile(ConfFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return mk, err
	}
	confFD.Lock()
	runtime.SetFinalizer(confFD, closeDescriptor)

	// open the locker file for usage
	lockFD, err := root.OpenFile(LockFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return mk, err
	}
	lockFD.Lock()
	runtime.SetFinalizer(lockFD, closeDescriptor)

	return New(confFD, lockFD, snippetsFS)
}

// New returns an instance of Maker using the provided file descriptors to read
// and write data from, and a target directory to manage the snippets within.
//
// The caller is responsible for closing both file descriptors
func New(conf, lock File, directory billy.Filesystem) (mk *Maker, err error) {
	mk = &Maker{
		configFile: conf,
		lockFile:   lock,
		directory:  directory,
		lock:       make(Lock),
	}

	err = unmarshalInto(conf, &mk.conf)
	if err != nil {
		return nil, err
	}

	err = unmarshalInto(lock, &mk.lock)
	if err != nil {
		return nil, err
	}

	for _, repository := range mk.conf.Repositories {
		err = repository.Init()
		if err != nil {
			return nil, err
		}
	}

	return mk, nil
}

// Init creates an empty configuration data with the default repository
func (mk *Maker) Init() error {
	if len(mk.conf.Repositories) > 0 {
		return fmt.Errorf("maker is already initialized")
	}

	err := mk.conf.AddRepository(&Repository{
		Alias:    "wwmoraes",
		URL:      "https://github.com/wwmoraes/maker-snippets.git",
		Snippets: make(map[string]string),
	})
	if err != nil {
		return err
	}

	return mk.Sync()
}

// Add fetches a snippet file and adds its info into the config and lock files
func (mk *Maker) Add(name string) error {
	name, versionStr, found := strings.Cut(name, "@")
	if !found {
		versionStr = "*"
	}

	alias, name, found := strings.Cut(name, ":")
	if !found {
		name = alias
		alias = ""
	}

	repository, err := mk.conf.GetRepository(alias)
	if err != nil {
		return err
	}

	if repository.HasSnippet(name) {
		return fmt.Errorf("snippet %s already added", name)
	}

	constraint, err := semver.NewConstraint(versionStr)
	if err != nil && !errors.Is(err, semver.ErrInvalidVersion) {
		return err
	}

	versions := make([]string, 0)
	// no constraint found, passthrough branch/tag name directly
	if constraint == nil {
		fmt.Printf("DEBUG: %s isn't a valid constraint - using as tag/branch name directly\n", versionStr)
		versions = append(versions, versionStr)
	} else {
		fmt.Printf("DEBUG: %s is a valid constraint - checking matching versions\n", versionStr)

		refs, err := repository.References()
		if err != nil {
			return err
		}

		err = refs.ForEach(func(r *plumbing.Reference) error {
			name := r.Name()
			if !(name.IsBranch() || name.IsTag()) {
				return nil
			}

			fmt.Printf("DEBUG: checking ref %s = %s\n", name.Short(), versionStr)

			version, err := semver.NewVersion(name.Short())
			if err != nil {
				fmt.Printf("DEBUG: %s not a semver\n", name.Short())
				return nil
			}

			if constraint.Match(version, false) {
				versions = append(versions, version.String())
				return nil
				// return storer.ErrStop
			}

			return nil
		})
		if err != nil && errors.Is(err, storer.ErrStop) {
			return err
		}
	}

	if len(versions) == 0 {
		return fmt.Errorf("no matching version or branch found")
	}

	sort.Strings(versions)
	hash, err := repository.ResolveRevision(plumbing.Revision(versions[0]))
	if err != nil {
		return err
	}

	repository.SetSnippet(name, versionStr)
	mk.lock.Set(repository.URL, name, hash.String())

	fmt.Println("installing", color.MagentaString(name))
	// TODO install newly added snippet

	return mk.Sync()
}

// Remove removes a snippet file and its info from the config and lock files
func (mk *Maker) Remove(name string) error {
	for _, repository := range mk.conf.Repositories {
		delete(repository.Snippets, name)
		mk.lock.Unset(repository.URL, name)
	}

	fmt.Println("removed", color.MagentaString(name))

	return mk.Sync()
}

// Install fetches the snippets if they're not present, or if there's any local
// changes
func (mk *Maker) Install(force bool) (err error) {
	for _, repository := range mk.conf.Repositories {
		for name, version := range repository.Snippets {
			lockVersion := mk.lock.Get(repository.URL, name)
			if lockVersion == "" {
				lockVersion = version
				// TODO generate lock on install
			}

			file, err := repository.Get(version, name)
			if err != nil {
				return err
			}

			installed, err := mk.install(name, file)
			if err != nil {
				return err
			}

			if installed {
				fmt.Println("updated ", color.MagentaString(name))
			} else {
				fmt.Println("skipped ", color.MagentaString(name))
			}
		}
	}

	return mk.Sync()
}

// Sync marshals the current configuration and lock data back into their
// respective file handlers. Flushing/syncing to an underlying persistent media
// is the caller's responsibility.
func (mk *Maker) Sync() (err error) {
	err = mk.configFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = mk.configFile.Seek(0, 0)
	if err != nil {
		return err
	}

	err = marshalInto(&mk.conf, mk.configFile)
	if err != nil {
		return err
	}

	err = mk.lockFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = mk.lockFile.Seek(0, 0)
	if err != nil {
		return err
	}

	err = marshalInto(&mk.lock, mk.lockFile)
	if err != nil {
		return err
	}

	return nil
}

func (mk *Maker) install(name string, file FileReader) (bool, error) {
	fd, err := mk.directory.OpenFile(fmt.Sprintf("%s.mk", name), os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return false, err
	}
	defer fd.Close()

	currentData, err := io.ReadAll(fd)
	if err != nil {
		return false, err
	}

	hash := plumbing.ComputeHash(plumbing.BlobObject, currentData)
	if hash.String() == file.ID().String() {
		return false, nil
	}

	reader, err := file.Reader()
	if err != nil {
		return false, err
	}
	defer reader.Close()

	err = fd.Truncate(0)
	if err != nil {
		return false, err
	}

	_, err = fd.Seek(0, 0)
	if err != nil {
		return false, err
	}

	_, err = io.Copy(fd, reader)
	if err != nil {
		return false, err
	}

	return true, nil
}

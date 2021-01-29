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

// Truncable represents an IO object that supports truncating its size to a
// specific value
type Truncable interface {
	// Truncate changes the size of the object. It does not change the I/O offset
	Truncate(size int64) error
}

// Lockable represents an IO object that supports locking its access to the
// current process
type Lockable interface {
	// Lock applies an advisory lock e.g. flock. It protects against access from
	// other processes
	Lock() error
}

// Unlockable represents an IO object that supports unlocking its access to be
// accessible by any process
type Unlockable interface {
	// Unlock removes the advisory lock to enable other processes to access it
	Unlock() error
}

// File represents a descriptor object that supports read, write, seek, truncate
// and lock operations.
type File interface {
	io.Reader
	io.Writer
	io.Seeker

	Truncable
	Lockable
	Unlockable
}

// UnlockCloser represents an IO object that supports unlocking and closing
type UnlockCloser interface {
	io.Closer
	Unlockable
}

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

// GetRepository returns a repository for the given reference, which can be
// either an alias or an URL. An empty reference returns the first
// (i.e. default) entry found in the configuration.
func (mk *Maker) GetRepository(reference string) (*Repository, error) {
	if len(mk.conf.Repositories) == 0 {
		return nil, fmt.Errorf("no repositories configured")
	}

	if reference == "" {
		return mk.conf.Repositories[0], nil
	}

	for _, repository := range mk.conf.Repositories {
		if repository.Alias == reference {
			return repository, nil
		}

		if repository.URL == reference {
			return repository, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
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

	repository, err := mk.GetRepository(alias)
	if err != nil {
		return err
	}

	if _, exists := repository.Snippets[name]; exists {
		return fmt.Errorf("snippet %s already added", name)
	}

	constraint, err := semver.NewConstraint(versionStr)
	if err != nil && !errors.Is(err, semver.ErrInvalidVersion) {
		return err
	}

	versions := make([]string, 0)
	// no constraint found, passthrough branch/tag name directly
	if constraint == nil {
		versions = append(versions, versionStr)
	} else {
		refs, err := repository.References()
		if err != nil {
			return err
		}

		err = refs.ForEach(func(r *plumbing.Reference) error {
			name := r.Name()
			if !(name.IsBranch() || name.IsTag()) {
				return nil
			}

			fmt.Printf("DEBUG: checking ref %s\n", name.Short())

			version, err := semver.NewVersion(name.Short())
			if err != nil {
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

	repository.Snippets[name] = versionStr
	mk.lock[repository.URL][name] = hash.String()

	fmt.Println("installing", color.MagentaString(name))
	// TODO install newly added snippet

	return mk.Sync()
}

// Remove removes a snippet file and its info from the config and lock files
func (mk *Maker) Remove(name string) error {
	for _, repository := range mk.conf.Repositories {
		delete(repository.Snippets, name)
		delete(mk.lock[repository.URL], name)
	}

	fmt.Println("removed", color.MagentaString(name))

	return mk.Sync()
}

// Install fetches the snippets if they're not present, or if there's any local
// changes
func (mk *Maker) Install(force bool) (err error) {
	for _, repository := range mk.conf.Repositories {
		for name, version := range repository.Snippets {
			file, err := repository.Get(version, name)
			if err != nil {
				return err
			}

			fd, err := mk.directory.OpenFile(fmt.Sprintf("%s.mk", name), os.O_RDWR|os.O_CREATE, 0640)
			if err != nil {
				return err
			}
			defer fd.Close()

			currentData, err := io.ReadAll(fd)
			if err != nil {
				return err
			}

			hash := plumbing.ComputeHash(plumbing.BlobObject, currentData)
			if hash.String() == file.ID().String() {
				fmt.Println("skipped:", color.MagentaString(name))
				continue
			}

			reader, err := file.Reader()
			if err != nil {
				return err
			}
			defer reader.Close()

			err = fd.Truncate(0)
			if err != nil {
				return err
			}

			_, err = fd.Seek(0, 0)
			if err != nil {
				return err
			}

			_, err = io.Copy(fd, reader)
			if err != nil {
				return err
			}

			fmt.Println("updated:", color.MagentaString(name))
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

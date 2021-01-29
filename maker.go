package maker

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	SnippetsDirectory = ".make"
	ConfFilename      = "maker.yaml"
	LockFilename      = "maker.lock"
)

type Truncable interface {
	// Truncate changes the size of the file. It does not change the I/O offset.
	Truncate(size int64) error
}

type Lockable interface {
	// Lock locks the file like e.g. flock. It protects against access from
	// other processes.
	Lock() error
}

type Unlockable interface {
	// Unlock removes the advisory lock.
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
// and write data from, and a target directory for the snippets.
//
// The caller is responsible for closing both file descriptors, and
// closing before Maker is done will cause IO errors
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

	return mk, nil
}

// Add fetches a snippet file and adds its info into the config and lock files
func (mk *Maker) Add(name string) error {
	versionStr := "*"
	data := strings.SplitAfterN(name, "@", 2)
	name = data[0]
	if len(data) == 2 {
		versionStr = data[1]
	}

	// TODO manage default repository name somehow
	repository := mk.conf.Repositories["wwmoraes"]

	if _, exists := repository.Snippets[name]; exists {
		return fmt.Errorf("snippet %s already added", name)
	}

	// constraint, err := semver.NewConstraint(versionStr)
	// if err != nil {
	//   return err
	// }

	hash, err := repository.ResolveRevision(plumbing.Revision(versionStr))
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
	// repo, err := local.New("../maker-repository")
	// if err != nil {
	//   return err
	// }

	for _, repository := range mk.conf.Repositories {
		err = repository.Init()
		if err != nil {
			return err
		}

		for name, version := range repository.Snippets {
			reader, err := repository.Get(version, name)
			if err != nil {
				return err
			}
			defer reader.Close()

			fd, err := mk.directory.OpenFile(fmt.Sprintf("%s.mk", name), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
			if err != nil {
				return err
			}
			defer fd.Close()

			_, err = io.Copy(fd, reader)
			if err != nil {
				return err
			}

			fmt.Println("skipped:", color.MagentaString(name))
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

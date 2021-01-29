package maker

import (
	"io"
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

package timee

import (
	"os"
	"time"
)

// Timespec provides access to file times.
type Timespec interface {
	BirthTime() time.Time
	HasBirthTime() bool
}

type btime struct {
	v time.Time
}
type statFunc func(string) (os.FileInfo, error)

// Get returns the Timespec for the given FileInfo
func Get(fi os.FileInfo) Timespec {
	return getTimespec(fi)
}
func stat(name string, sf statFunc) (Timespec, error) {
	fi, err := sf(name)
	if err != nil {
		return nil, err
	}
	return getTimespec(fi), nil
}



func (btime) HasBirthTime() bool { return true }

func (b btime) BirthTime() time.Time { return b.v }
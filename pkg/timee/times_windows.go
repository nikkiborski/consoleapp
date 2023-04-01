package timee

import (
	"os"
	"syscall"
	"time"
)

// HasChangeTime and HasBirthTime are true if and only if
// the target OS supports them.
const (
	HasChangeTime = false
	HasBirthTime  = true
)

type timespec struct {
	btime
}

func getTimespec(fi os.FileInfo) Timespec {
	var t timespec
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	t.btime.v = time.Unix(0, stat.CreationTime.Nanoseconds())
	return t
}

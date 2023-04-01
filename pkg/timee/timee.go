package timee

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)
var (
	findProcErr                      error
	procGetFileInformationByHandleEx *syscall.Proc
)
type Timespec interface {
	BirthTime() time.Time
	HasBirthTime() bool
}
// Stat returns the Timespec for the given filename.
func Stat(name string) (Timespec, error) {
	ts, err := platformSpecficStat(name)
	if err == nil {
		return ts, err
	}

	return stat(name, os.Stat)
}
func stat(name string, sf statFunc) (Timespec, error) {
	fi, err := sf(name)
	if err != nil {
		return nil, err
	}
	return getTimespec(fi), nil
}
type timespec struct {
	atime
	mtime
	noctime
	btime
}
type noctime struct{}

func (noctime) HasChangeTime() bool { return false }
func getTimespec(fi os.FileInfo) Timespec {
	var t timespec
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	t.atime.v = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	t.mtime.v = time.Unix(0, stat.LastWriteTime.Nanoseconds())
	t.btime.v = time.Unix(0, stat.CreationTime.Nanoseconds())
	return t
}
type statFunc func(string) (os.FileInfo, error)
func platformSpecficStat(name string) (Timespec, error) {
	if findProcErr != nil {
		return nil, findProcErr
	}

	return openHandleAndStat(name, syscall.FILE_FLAG_BACKUP_SEMANTICS)
}

func openHandleAndStat(name string, attrs uint32) (Timespec, error) {
	pathp, e := syscall.UTF16PtrFromString(name)
	if e != nil {
		return nil, e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, attrs, 0)
	if e != nil {
		return nil, e
	}
	defer syscall.Close(h)

	return statFile(h)
}
type atime struct {
	v time.Time
}

type mtime struct {
	v time.Time
}

type ctime struct {
	v time.Time
}
func (ctime) HasChangeTime() bool { return true }

func (c ctime) ChangeTime() time.Time { return c.v }


type btime struct {
	v time.Time
}
type timespecEx struct {
    atime
    mtime
    ctime
    btime
}
func (btime) HasBirthTime() bool { return true }

func (b btime) BirthTime() time.Time { return b.v }
func (a atime) AccessTime() time.Time { return a.v }
func statFile(h syscall.Handle) (Timespec, error) {
	var fileInfo fileBasicInfo
	if err := getFileInformationByHandleEx(h, &fileInfo); err != nil {
		return nil, err
	}

	var t timespecEx
	t.atime.v = time.Unix(0, fileInfo.LastAccessTime.Nanoseconds())
	t.mtime.v = time.Unix(0, fileInfo.LastWriteTime.Nanoseconds())
	t.ctime.v = time.Unix(0, fileInfo.ChangeTime.Nanoseconds())
	t.btime.v = time.Unix(0, fileInfo.CreationTime.Nanoseconds())
	return t, nil
}
const (
	fileBasicInfoClass fileInformationClass = iota
)
type fileInformationClass int
func getFileInformationByHandleEx(handle syscall.Handle, data *fileBasicInfo) (err error) {
	if findProcErr != nil {
		return findProcErr
	}

	r1, _, e1 := syscall.Syscall6(procGetFileInformationByHandleEx.Addr(), 4, uintptr(handle), uintptr(fileBasicInfoClass), uintptr(unsafe.Pointer(data)), unsafe.Sizeof(*data), 0, 0)
	if r1 == 0 {
		err = syscall.EINVAL
		if e1 != 0 {
			err = error(e1)
		}
	}
	return
}

// fileBasicInfo holds the C++ data for FileTimes.
//
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa364217(v=vs.85).aspx
type fileBasicInfo struct {
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
	ChangeTime     syscall.Filetime
	FileAttributes uint32
	_              uint32 // padding
}
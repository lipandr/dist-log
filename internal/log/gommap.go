package log

import (
	"reflect"
	"syscall"
	"unsafe"
)

type ProtFlags uint

const (
	ProtRead  ProtFlags = 0x1
	ProtWrite ProtFlags = 0x2
)

type MapFlags uint

const (
	MapShared MapFlags = 0x1
)

type SyncFlags uint

const (
	MsSync SyncFlags = 0x4
)

type MMap []byte

func Map(fd uintptr, prot ProtFlags, flags MapFlags) (MMap, error) {
	mMap, err := MapAt(0, fd, 0, -1, prot, flags)
	return mMap, err
}

func MapAt(addr uintptr, fd uintptr, offset, length int64, prot ProtFlags, flags MapFlags) (MMap, error) {
	if length == -1 {
		var stat syscall.Stat_t
		if err := syscall.Fstat(int(fd), &stat); err != nil {
			return nil, err
		}
		length = stat.Size
	}
	addr, err := mMapSyscall(addr, uintptr(length), uintptr(prot), uintptr(flags), fd, offset)
	if err != syscall.Errno(0) {
		return nil, err
	}
	mMap := MMap{}

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&mMap))
	dh.Data = addr
	dh.Len = int(length)
	dh.Cap = dh.Len
	return mMap, nil
}

func (mMap MMap) Sync(flags SyncFlags) error {
	rh := *(*reflect.SliceHeader)(unsafe.Pointer(&mMap))
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC, rh.Data, uintptr(rh.Len), uintptr(flags))
	if err != 0 {
		return err
	}
	return nil
}

func mMapSyscall(addr, length, prot, flags, fd uintptr, offset int64) (uintptr, error) {
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, length, prot, flags, fd, uintptr(offset))
	return addr, err
}

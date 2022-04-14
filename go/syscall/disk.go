package syscall

import (
	"syscall"

	filesystem_ "github.com/kaydxh/golang/go/filesystem"
)

type DiskUsage struct {
	stat *syscall.Statfs_t
}

func NewDiskUsage(path string) (*DiskUsage, error) {
	var stat syscall.Statfs_t
	mount, err := filesystem_.FindMount(path)
	if err != nil {
		return nil, err
	}

	err = syscall.Statfs(mount.Path, &stat)
	if err != nil {
		return nil, err
	}

	return &DiskUsage{&stat}, nil
}

// Free returns total free bytes on file system
func (du *DiskUsage) Free() uint64 {
	return du.stat.Bfree * uint64(du.stat.Bsize)
}

//  Avail returns total avail bytes on file system
func (du *DiskUsage) Avail() uint64 {
	return du.stat.Bavail * uint64(du.stat.Bsize)
}

// Size returns total size of the file system
func (du *DiskUsage) Size() uint64 {
	return uint64(du.stat.Blocks) * uint64(du.stat.Bsize)
}

// Size returns total used bytes on the file system
func (du *DiskUsage) Used() uint64 {
	return du.Size() - du.Free()
}

// Usage returns percentage of used on the file system
/*
https://github.com/coreutils/coreutils/blob/master/src/df.c#:~:text=pct%20%3D%20u100%20/%20nonroot_total%20%2B%20(u100%20%25%20nonroot_total%20!%3D%200)%3B
By default, ext2/3/4 filesystems reserve 5% of the space to be useable only by root. This is to avoid a normal user completely filling the disk which would
then cause system components to fail whenever they next needed to write to the disk
*/
func (du *DiskUsage) Usage() float32 {
	var deta float32
	u100 := du.Used() * 100
	nonrootTotal := du.Used() + du.Avail()
	if u100%nonrootTotal != 0 {
		deta = 1.0
	}
	return float32(u100)/float32(nonrootTotal) + deta
}

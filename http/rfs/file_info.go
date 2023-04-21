package rfs

import (
    "os"
    "time"
)

// FileInfo is a custom implementation of the os.FileInfo interface that allows setting the size of the file.
type FileInfo struct {
    name    string
    size    int64
    mode    os.FileMode
    modTime time.Time
    isDir   bool
    sys     interface{}
}

func (fi *FileInfo) Name() string {
    return fi.name
}

func (fi *FileInfo) Size() int64 {
    return fi.size
}

func (fi *FileInfo) Mode() os.FileMode {
    return fi.mode
}

func (fi *FileInfo) ModTime() time.Time {
    return fi.modTime
}

func (fi *FileInfo) IsDir() bool {
    return fi.isDir
}

func (fi *FileInfo) Sys() interface{} {
    return fi.sys
}

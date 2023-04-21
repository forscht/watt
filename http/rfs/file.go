package rfs

import (
    "io"
    "net/http"
    "os"
)

type File struct {
    http.File
    content  []byte
    size     int64
    offset   int64
    fileInfo os.FileInfo
}

func (f *File) Stat() (os.FileInfo, error) {
    return f.fileInfo, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
    switch whence {
    case io.SeekStart:
        f.offset = offset
    case io.SeekCurrent:
        f.offset += offset
    case io.SeekEnd:
        f.offset = f.size + offset
    }
    return f.offset, nil
}

func (f *File) Read(p []byte) (int, error) {
    if f.offset >= f.size {
        return 0, io.EOF
    }
    n := copy(p, f.content[f.offset:f.size])
    f.offset += int64(n)
    return n, nil
}

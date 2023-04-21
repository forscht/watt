// Package http provides a static file server using embedded files in Go.
// This code is inspired by a pull request in the gin-contrib/static library.
// Reference: https://github.com/gin-contrib/static/pull/20/files
package http

import (
    "embed"
    "io/fs"
    "net/http"

    "github.com/gin-gonic/contrib/static"

    "github.com/forscht/watt/http/rfs"
)

//go:embed static
// Define the embedded file system using the `embed` package's `FS` type.
// The `static` directory is embedded in the binary using the `//go:embed` directive,
// which allows it to be accessed as a file system using the `embed` package.
var server embed.FS

// Embed is a function that takes a target path as an argument and returns
// an implementation of the static.ServeFileSystem interface, which can be used
// with a Gin router to serve static files from the embedded file system.
// The embedded file system is defined by the "static" directory in the embedded
// files (using the `//go:embed` directive).
func Embed(targetPath string, placeholderValues map[string]string, regex bool) static.ServeFileSystem {
    fsys, err := fs.Sub(server, targetPath)
    if err != nil {
        panic(err)
    }

    return embedFileSystem{
        FileSystem: rfs.FileSystem{
            FileSystem:        http.FS(fsys),
            PlaceHolderValues: placeholderValues,
            Regex:             regex,
        },
    }
}

// Define a custom file system that wraps an `http.FileSystem` and implements
// the `Exists` method to check if a file exists in the embedded file system.
type embedFileSystem struct {
    http.FileSystem
}

// Exists is a method that takes a prefix and a path as arguments and returns
// a boolean value indicating whether the file with the specified path exists
// in the embedded file system. It uses the `Open` method of the embedded file
// system to check if the file can be opened without error. If there is an error,
// it means the file does not exist.
func (e embedFileSystem) Exists(prefix string, path string) bool {
    _, err := e.Open(path)
    if err != nil {
        return false
    }
    return true
}

// Package rfs provides a custom http.FileSystem implementation for serving files with placeholders replaced by their corresponding values.
package rfs

import (
    "bytes"
    "io"
    "net/http"
    "regexp"
)

// FileSystem is a custom struct that embeds the http.FileSystem interface and adds a PlaceHolderValues map for placeholder values
type FileSystem struct {
    http.FileSystem
    PlaceHolderValues map[string]string // map for placeholder values
    Regex             bool
}

// New creates a new instance of the FileSystem struct with the specified embedded http.FileSystem instance and placeholder map
// and returns a pointer to it.
func New(fs http.FileSystem, placeHolder map[string]string, regex bool) http.FileSystem {
    return &FileSystem{
        FileSystem:        fs,
        PlaceHolderValues: placeHolder, // initialize the map of placeholder values
        Regex:             regex,
    }
}

// Open is a method of the FileSystem struct that implements the http.FileSystem interface's Open method.
// It returns a http.File instance representing the requested file with any placeholders in its contents replaced by their
// corresponding values.
func (f FileSystem) Open(name string) (http.File, error) {
    // Open the file from the embedded http.FileSystem
    file, err := f.FileSystem.Open(name)
    if err != nil {
        return nil, err
    }

    // Get file info
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }

    // If the file is a directory, return it as is
    if fileInfo.IsDir() {
        return file, nil
    }
    fileInfo.Size()
    // Read the file contents
    fileContent, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    // Apply template to file content
    fileContent = f.ApplyTemplate(fileContent)

    // Create a new http.File with the modified content
    newFile := &File{
        File:    file,
        content: fileContent,
        size:    int64(len(fileContent)),
        offset:  0,
        fileInfo: &FileInfo{
            name:    fileInfo.Name(),
            size:    int64(len(fileContent)),
            mode:    fileInfo.Mode(),
            modTime: fileInfo.ModTime(),
            isDir:   fileInfo.IsDir(),
            sys:     fileInfo.Sys(),
        },
    }

    return newFile, nil // Return the newFile instance as a http.File instance
}

// ApplyTemplate is a method of the FileSystem struct that replaces placeholders in the input data with their corresponding values
// and returns the modified data
// If FileSystem.Regex == false replace values by {{key}} = value
func (f *FileSystem) ApplyTemplate(data []byte) []byte {

    // Loop over placeholder values
    for key, value := range f.PlaceHolderValues {
        // Try to compile the key as a regex pattern
        if f.Regex {
            re, err := regexp.Compile(key)
            if err == nil {
                // If the key can be compiled as a regex pattern, replace all occurrences of the pattern with the corresponding value
                data = re.ReplaceAll(data, []byte(value))
            }
        } else {
            data = bytes.ReplaceAll(data, []byte("{{"+key+"}}"), []byte(value))
        }
    }

    return data
}

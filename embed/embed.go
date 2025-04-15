package embed

import (
	"io/fs"
	"os"
	"time"
)

type embedFS struct {
	fs.FS
	modTime time.Time
}

type embedFile struct {
	fs.File
	modTime time.Time
}

type embedFileInfo struct {
	os.FileInfo
	modTime time.Time
}

func (f *embedFS) Open(name string) (fs.File, error) {
	file, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return &embedFile{File: file, modTime: f.modTime}, nil
}

func (f *embedFile) Stat() (os.FileInfo, error) {
	fileInfo, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return &embedFileInfo{FileInfo: fileInfo, modTime: f.modTime}, nil
}

func (f *embedFileInfo) ModTime() time.Time {
	return f.modTime
}

func NewFS(embedfs fs.FS, modTime time.Time) fs.FS {
	return &embedFS{FS: embedfs, modTime: modTime}
}

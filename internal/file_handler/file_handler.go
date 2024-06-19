package file_handler

import (
	"io/fs"
	"os"
	"path"
	"strconv"
)

var (
	fileWalCounter int = 0 // First aproach: only use one file for wal logs
)

type Options struct {
	DirName         string
	DirPerms        fs.FileMode
	FilePerms       fs.FileMode
	createFileFlags int
}

var DefaultOptions = &Options{
	DirName:         "WalFolder",
	DirPerms:        0755,
	FilePerms:       0640,
	createFileFlags: os.O_CREATE | os.O_RDWR,
}

func CreateWalFolder(opts Options) error {
	err := os.MkdirAll(opts.DirName, opts.DirPerms)
	if err != nil {
		return err
	}

	return nil
}

func CreateWalNewFile(opts Options) (*os.File, error) {
	filePath := path.Join(opts.DirName, "wal_"+strconv.Itoa(fileWalCounter)+".log")

	file, err := os.OpenFile(filePath, opts.createFileFlags, opts.FilePerms)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func OpenWal(opts *Options) (*os.File, error) {
	err := CreateWalFolder(*opts)
	if err != nil {
		return nil, err
	}

	file, err := CreateWalNewFile(*opts)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
}

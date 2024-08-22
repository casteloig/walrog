package file_handler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
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
	DirName:         "/tmp/WalFolder",
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
	fileName := fmt.Sprintf("wal_%03d.log", fileWalCounter)
	fileWalCounter++
	filePath := path.Join(opts.DirName, fileName)

	file, err := os.OpenFile(filePath, opts.createFileFlags, opts.FilePerms)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CreateCheckpointFile(opts Options) (*os.File, error) {
	filePath := path.Join(opts.DirName, "checkpoint")

	file, err := os.OpenFile(filePath, opts.createFileFlags, opts.FilePerms)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func OpenWal(opts *Options) (*os.File, *os.File, error) {
	err := CreateWalFolder(*opts)
	if err != nil {
		return nil, nil, err
	}

	walFile, err := CreateWalNewFile(*opts)
	if err != nil {
		return nil, nil, err
	}

	checkpointFile, err := CreateCheckpointFile(*opts)
	if err != nil {
		return nil, nil, err
	}

	return walFile, checkpointFile, nil
}

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

// Options defines the configuration options for managing WAL files and directories.
// Fields:
//   - DirName: The name of the directory where WAL files will be stored.
//   - DirPerms: The permissions to set for the WAL directory.
//   - FilePerms: The permissions to set for the WAL files.
//   - createFileFlags: Flags used when creating WAL files (e.g., read/write, create).
type Options struct {
	DirName         string
	DirPerms        fs.FileMode
	FilePerms       fs.FileMode
	createFileFlags int
}

// DefaultOptions provides a default configuration for managing WAL files and directories.
// Fields:
//   - DirName: Default directory for WAL files (/tmp/WalFolder).
//   - DirPerms: Default permissions for the WAL directory (0755).
//   - FilePerms: Default permissions for WAL files (0640).
//   - createFileFlags: Default flags for creating WAL files (os.O_CREATE | os.O_RDWR).
var DefaultOptions = &Options{
	DirName:         "/tmp/WalFolder",
	DirPerms:        0755,
	FilePerms:       0640,
	createFileFlags: os.O_CREATE | os.O_RDWR,
}

// CreateWalFolder() creates the directory for storing WAL files if it does not already exist.
// It uses the directory name and permissions specified in the provided Options.
//
// Parameters:
//   - opts: An Options struct containing the directory name and permissions.
//
// Returns:
//   - An error if the directory cannot be created, or nil if successful.
func CreateWalFolder(opts Options) error {
	err := os.MkdirAll(opts.DirName, opts.DirPerms)
	if err != nil {
		return fmt.Errorf("failed to create WAL folder: %w", err)
	}
	return nil
}

// CreateWalNewFile() creates a new WAL file in the specified directory with the given options.
// The file name is automatically generated in the format "wal_XXX.log", where XXX is a counter.
//
// Parameters:
//   - opts: An Options struct containing the directory name, file permissions, and creation flags.
//
// Returns:
//   - A pointer to the newly created os.File object.
//   - An error if the file cannot be created.
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

// CreateCheckpointFile() creates a new checkpoint file in the specified directory.
// Checkpoint files are used to store the state of the system at a specific point in time.
//
// Parameters:
//   - opts: An Options struct containing the directory name, file permissions, and creation flags.
//
// Returns:
//   - A pointer to the newly created os.File object.
//   - An error if the file cannot be created.
func CreateCheckpointFile(opts Options) (*os.File, error) {
	filePath := path.Join(opts.DirName, "checkpoint")

	file, err := os.OpenFile(filePath, opts.createFileFlags, opts.FilePerms)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// OpenWal() opens an existing WAL file for reading or writing.
// This function is used to access WAL files that have already been created.
//
// Parameters:
//   - filePath: The full path to the WAL file to be opened.
//   - opts: An Options struct containing the file permissions and creation flags.
//
// Returns:
//   - A pointer to the opened os.File object.
//   - An error if the file cannot be opened.
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

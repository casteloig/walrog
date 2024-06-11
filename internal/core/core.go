package core

import (
	"os"

	fh "github.com/casteloig/walrog/internal/file_handler"
)

type Wal struct {
	FileHandlerOpts fh.Options
	ActualFile      *os.File
}

func InitWal(options *fh.Options) (*Wal, error) {
	opts := options
	if options == nil {
		opts = fh.DefaultOptions
	}

	file, err := fh.OpenWal(opts)
	if err != nil {
		return nil, err
	}

	myWal := &Wal{
		FileHandlerOpts: *opts,
		ActualFile:      file,
	}

	return myWal, nil
}

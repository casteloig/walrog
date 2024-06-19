package core

import (
	"hash"
	"hash/crc32"
	"os"

	fh "github.com/casteloig/walrog/internal/file_handler"
)


type WalOptions struct {
	BufferSize      uint32
	SegmentSize     uint32 // must be multiple of BufferSize
	FileHandlerOpts *fh.Options
}

var DefaultWalOptions = &WalOptions{
	BufferSize:      2097152,  // 2Mb
	SegmentSize:     67108864, // 64Mb
	FileHandlerOpts: fh.DefaultOptions,
}

type BufferEntry struct {
	lastSeqNo uint64
	data      []byte
	crc       hash.Hash32
}

type Wal struct {
	Options  	*WalOptions
	HotFile  	*os.File // File that's being used
	Buffer      []*BufferEntry
}

// InitWal creates a new Wal
func InitWal(options *WalOptions) (*Wal, error) {
	// Get default Wal options
	opts := options
	if options == nil {
		opts = DefaultWalOptions
	}

	// Open Wal folder and first file
	file, err := fh.OpenWal(opts.FileHandlerOpts)
	if err != nil {
		return nil, err
	}

	// Create Wal and return
	w := &Wal{
		Options:  opts,
		HotFile:  file,
		Buffer:	nil,
	}

	return w, nil
}

func (w *Wal) writeBuffer(data []byte) error {

}

func (w *Wal) flushBuffer() error {

}


// encodeCRC32 returns the checksum of the data provided by the argument using IEEE poly
func calculate (data []byte) (uint32) {
	return crc32.ChecksumIEEE(data)
}

func validateCRC32(data uint32) {
	crc32.
}

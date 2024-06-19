package core

import (
	"bufio"
	"log"
	"os"

	fh "github.com/casteloig/walrog/internal/file_handler"
	"github.com/casteloig/walrog/internal/utils"
)

type WalOptions struct {
	BufferSize      uint32
	SegmentSize     uint32 // must be multiple of BufferSize
	FileHandlerOpts *fh.Options
}

// type BufferEntry struct {
// 	lastSeqNo uint64      // 8 bytes
// 	data      []byte      // data.Lenght()
// 	crc       hash.Hash32 // 4 bytes
// }

type Wal struct {
	Options *WalOptions
	HotFile *os.File // File that's being used
	Buffer  *bufio.Writer
}

var DefaultWalOptions = &WalOptions{
	BufferSize:      4194304,  // 4Mb
	SegmentSize:     67108864, // 64Mb
	FileHandlerOpts: fh.DefaultOptions,
}

var lastSeqNo uint32 = 0

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

	// Create buffer to write to the hot file
	// options.BufferSize is intended to be uint32 to persist across architectures,
	// we cast it into int because the function forces to do so
	writerBuffer := bufio.NewWriterSize(file, int(options.BufferSize))

	// Create Wal and return
	w := &Wal{
		Options: opts,
		HotFile: file,
		Buffer:  writerBuffer,
	}

	return w, nil
}

func (w *Wal) writeBuffer(data []byte) error {
	var bufferNewBytes []byte

	// First we add the LSN (4 bytes) to the buffer
	// in []byte
	lsnBytes := utils.Uint32ToBytes(uint32(lastSeqNo))
	bufferNewBytes = append(bufferNewBytes, lsnBytes...)

	// Then add data length (4 bytes)
	dataLength, err := utils.IntToUint32(len(data))
	if err != nil {
		log.Fatal(err)
	}
	dataLengthBytes := utils.Uint32ToBytes(dataLength)
	bufferNewBytes = append(bufferNewBytes, dataLengthBytes...)

	// Then add data (get length from content of dataLengthBytes)
	bufferNewBytes = append(bufferNewBytes, data...)

	// Calculate CRC appending data calculating checksum
	checksumRaw := append(lsnBytes, dataLengthBytes...)
	checksumRaw = append(checksumRaw, data...)
	crcBytes := utils.Uint32ToBytes(utils.CalculateCRC(checksumRaw))
	bufferNewBytes = append(bufferNewBytes, crcBytes...)

	fitsNewEntry := checkBufferNotLong(w, len(bufferNewBytes))
	if fitsNewEntry == true {
		w.Buffer.Write(bufferNewBytes)
	} else {
		log.Println("Flush to HotFile if fits or create a new one if does not fit")
		err = w.Buffer.Flush()
		if err != nil {
			log.Fatalf("Error flushing to file")
		}
	}

	return nil
}

func checkBufferNotLong(w *Wal, newEntryLength int) bool {
	if (w.Buffer.Size() + newEntryLength) > int(w.Options.BufferSize) {
		return false
	}

	return true
}

func (w *Wal) flushBuffer() error {
	// TO DO
	// check if buffer fits into HotFile
	//		if it does, flush
	// 		if it does not, close file, create a new one, edit HotFile and flush
	return nil
}

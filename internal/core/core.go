package core

import (
	"bufio"
	"fmt"
	"io"
	"os"

	fh "github.com/casteloig/walrog/internal/file_handler"
	utils "github.com/casteloig/walrog/internal/utils"
)

type WalOptions struct {
	BufferSize      uint32 // Size of the buffer
	SegmentSize     uint32 // Max size of the file. Must be multiple of BufferSize
	FileHandlerOpts *fh.Options
}

var DefaultWalOptions = &WalOptions{
	BufferSize:      4194304,  // 4Mb
	SegmentSize:     67108864, // 64Mb
	FileHandlerOpts: fh.DefaultOptions,
}

type Wal struct {
	Options        *WalOptions
	HotFile        *os.File // File that's being used
	CheckpointFile *os.File
	segmentUsed    int
	Buffer         *bufio.Writer
	lsn            uint32
}

type RecoveredEntry struct {
	lsn  uint32
	data []byte
}

// InitWal() creates a new Wal.
// If nil argument is passed, will use default options.
// Always use InitWal() after calling Recover() and make sure everything is recovered.
// InitWal() will delete all content in the Wal folder.
func InitWal(options *WalOptions) (*Wal, error) {
	// Get default options if no arg passed to function
	if options == nil {
		options = DefaultWalOptions
	}

	walFile, checkpointFile, err := fh.OpenWal(options.FileHandlerOpts)
	if err != nil {
		return nil, err
	}

	// Create buffer to write to the hot file
	writerBuffer := bufio.NewWriterSize(walFile, int(options.BufferSize))

	// Create Wal and return
	w := &Wal{
		Options:        options,
		HotFile:        walFile,
		CheckpointFile: checkpointFile,
		Buffer:         writerBuffer,
		segmentUsed:    0,
		lsn:            0,
	}

	return w, nil
}

// WriteBuffer() allows to write a slice of bytes on the Wal.
// It first writes to a buffer which will be dumped into a file when reaching
// WalOptions.BufferSize
func (w *Wal) WriteBuffer(data []byte) error {

	// create temp buffer before flushing any data
	tmpBuffer, err := w.createTmpBuff(data)
	if err != nil {
		return err
	}

	// Checks either buffer can be written, must be flushed or the hot file must be rotated first
	err = w.manageWriteFlow(tmpBuffer)
	if err != nil {
		return err
	}

	return nil
}

// recoverFile() reads entries from a given file and validates their integrity using CRC.
// If the CRC matches, the entry is stored in a slice of RecoveredEntry. Returns the recovered entries and an error if any issues occur.
func recoverFile(file *os.File) ([]RecoveredEntry, error) {
	var records []RecoveredEntry

	reader := bufio.NewReader(file)
	lsnBytes := make([]byte, 4)
	lengthBytes := make([]byte, 4)
	crcBytes := make([]byte, 4)

	for {
		// Read 4 bytes of LSN
		_, err := reader.Read(lsnBytes)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached end of file")
				break
			}
			return nil, fmt.Errorf("error reading LSN: %w", err)
		}

		// Read lengthData
		_, err = reader.Read(lengthBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading data length: %w", err)
		}
		dataLength := utils.BytesToUint32(lengthBytes)

		// Read data
		dataBytes := make([]byte, dataLength)
		_, err = reader.Read(dataBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading data: %w", err)
		}

		// Read 4 bytes of CRC
		_, err = reader.Read(crcBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading CRC: %w", err)
		}
		crcData := utils.BytesToUint32(crcBytes)

		// Calculate CRC of LSN, lengthData and data
		dataWithoutCRC := append(lsnBytes, lengthBytes...)
		dataWithoutCRC = append(dataWithoutCRC, dataBytes...)
		calculatedCRC := utils.CalculateCRC(dataWithoutCRC)

		// Compare CRC bytes
		if crcData != calculatedCRC {
			return nil, fmt.Errorf("CRC mismatch: read %v, calculated %v", crcData, calculatedCRC)
		}

		// Store data in the slice
		newRecord := RecoveredEntry{
			lsn:  utils.BytesToUint32(lsnBytes),
			data: dataBytes,
		}

		records = append(records, newRecord)
	}

	return records, nil
}

// FlushBuffer() forces a flush of the buffer to the Segment/WalFile
func (w *Wal) FlushBuffer() error {
	fmt.Println("Flushing buffer to file")
	w.segmentUsed += w.Buffer.Buffered()
	err := w.Buffer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing to file")
	}
	return nil
}

// checkBufferOverflow() checks if the new data fits within the buffer
// Return true if it leads to Overflow. Returns false if it actually fits
func (w *Wal) checkBufferOverflow(newEntryLength int) bool {
	// If buffer is full it will flush automatically
	return (w.Buffer.Buffered() + newEntryLength) > int(w.Options.BufferSize)
}

// changeHotFile() creates a new File/Segment and changes the Wal's pointer to the new file
func (w *Wal) changeHotFile(newFile *os.File) error {
	w.HotFile = newFile
	return nil
}

// createTmpBuff() creates a tmpBuff before writing to Bufio with all the data
func (w *Wal) createTmpBuff(data []byte) ([]byte, error) {
	var tmpBuffer []byte

	// First we add the LSN (4 bytes) to the buffer
	// in []byte
	newBytes := utils.Uint32ToBytes(w.lsn)
	tmpBuffer = utils.AppendBytesToSlice(tmpBuffer, newBytes)

	// Then add data length (4 bytes)
	dataLength, err := utils.IntToUint32(len(data))
	if err != nil {
		return nil, fmt.Errorf("failed to convert data length to uint32: %w", err)
	}
	newBytes = utils.Uint32ToBytes(dataLength)
	tmpBuffer = utils.AppendBytesToSlice(tmpBuffer, newBytes)

	// Then add data (get length from content of dataLengthBytes)
	tmpBuffer = utils.AppendBytesToSlice(tmpBuffer, data)

	// Calculate CRC and add it to the tmpBuffer
	crc := utils.CalculateCRC(tmpBuffer)
	newBytes = utils.Uint32ToBytes(crc)
	tmpBuffer = utils.AppendBytesToSlice(tmpBuffer, newBytes)

	return tmpBuffer, nil
}

// manageWriteFlow() manages all algorythms to write the data into the buffer and flush it to the hot file
// if needed.
func (w *Wal) manageWriteFlow(tmpBuffer []byte) error {
	// Check if tmpBuffer is bigger than Buffer max size
	if len(tmpBuffer) > int(w.Options.BufferSize) {
		return fmt.Errorf("data is bigger than buffer, data cannot be handled")
	}

	// Check if tmpBuffer fits real Buffer
	// If it fits (no overflow), enter the condition and Write
	if !w.checkBufferOverflow(len(tmpBuffer)) {
		w.Buffer.Write(tmpBuffer)
		return nil
	}

	// If buffer can be flushed into file
	// Flush and write
	if w.Buffer.Buffered() < (int(w.Options.SegmentSize) - w.segmentUsed) {
		err := w.FlushBuffer()
		if err != nil {
			return err
		}
		w.Buffer.Write(tmpBuffer)
		return nil
	}

	// If buffer cannot be flushed into file, we have to Rotate the new file
	// Rotate, flush and write
	newFile, err := fh.CreateWalNewFile(*w.Options.FileHandlerOpts)
	if err != nil {
		return err
	}
	err = w.changeHotFile(newFile)
	if err != nil {
		return err
	}
	err = w.FlushBuffer()
	if err != nil {
		return err
	}
	w.Buffer.Write(tmpBuffer)

	return nil
}

// TODO
func (w *Wal) Truncate(lsnFirst uint32, lsnLast uint32) error {
	return nil
}

// TODO
// 1. Close file properly
// 2. New func to recover file from LSN

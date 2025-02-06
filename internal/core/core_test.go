package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	fh "github.com/casteloig/walrog/internal/file_handler"
)

// InitWal creates a new Wal
func TestInitWal(t *testing.T) {
	// Create an instance of Wal with default options
	w, err := InitWal(nil)
	if err != nil {
		t.Fatalf("InitWal() failed: %v", err)
	}

	// Verify initialization
	if w.Options == nil {
		t.Error("Expected Wal options, got nil")
	}
	if w.HotFile == nil {
		t.Error("Wal file should not be nil after initialization")
	}
	if w.Buffer == nil {
		t.Error("Wal buffer should not be nil after initialization")
	}
}

// Test function for checkBufferOverflow method
func TestCheckBufferOverflow(t *testing.T) {
	bufferSize := 10
	testCases := []struct {
		name             string
		bufferContent    []byte
		newEntryLength   int
		expectedOverflow bool
	}{
		{
			name:             "Empty buffer NO overflow",
			bufferContent:    []byte{},
			newEntryLength:   5,
			expectedOverflow: false,
		},
		{
			name:             "Partial buffer NO overflow",
			bufferContent:    []byte{1, 2, 3},
			newEntryLength:   5,
			expectedOverflow: false,
		},
		{
			name:             "Partial buffer WITH overflow",
			bufferContent:    []byte{1, 2, 3},
			newEntryLength:   8,
			expectedOverflow: true,
		},
		{
			name:             "Empty buffer WITH overflow",
			bufferContent:    []byte{},
			newEntryLength:   12,
			expectedOverflow: true,
		},
		{
			name:             "Exact buffer size",
			bufferContent:    []byte{1, 2, 3, 4, 5},
			newEntryLength:   5,
			expectedOverflow: false,
		},
	}

	mockBuf := bufio.NewWriterSize(nil, bufferSize)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockBuf.Reset(nil)
			wal := &Wal{
				Options: &WalOptions{
					BufferSize: uint32(bufferSize),
				},
				Buffer: mockBuf,
			}

			// Simulate the current size of the buffer
			wal.Buffer.Write(tc.bufferContent)

			overflow := wal.checkBufferOverflow(tc.newEntryLength)
			if overflow != tc.expectedOverflow {
				t.Errorf("checkBufferOverflow() = %v; want %v", overflow, tc.expectedOverflow)
			}
		})
	}
}

// Test function for WriteBuffer method
func TestCreateTmpBuff(t *testing.T) {
	testCases := []struct {
		name           string
		lsn            uint32
		data           []byte
		expectedBuffer []byte
	}{
		{
			name:           "Case 1",
			lsn:            1,
			data:           []byte{1, 2, 3},
			expectedBuffer: []byte{1, 0, 0, 0, 3, 0, 0, 0, 1, 2, 3, 190, 45, 28, 49},
		},
		{
			name:           "Case 2",
			lsn:            11,
			data:           []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedBuffer: []byte{11, 0, 0, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 181, 89, 81, 142},
		},
		{
			name:           "Empty data",
			lsn:            11,
			data:           []byte{},
			expectedBuffer: []byte{11, 0, 0, 0, 0, 0, 0, 0, 63, 195, 72, 56},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wal := &Wal{lsn: tc.lsn}

			tmpBuffer, err := wal.createTmpBuff(tc.data)

			if tc.expectedBuffer == nil {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if !bytes.Equal(tmpBuffer, tc.expectedBuffer) {
					t.Errorf("expected buffer %v, got %v", tc.expectedBuffer, tmpBuffer)
				}
			}
		})
	}
}

// Test function for TestManageWriteFlow method
func TestManageWriteFlow(t *testing.T) {
	bufferSize := 10
	testCases := []struct {
		name             string
		segmentSize      int
		bufferContent    []byte
		tmpBuffer        []byte
		expectedError    error
		expectedBuffer   []byte
		expectedFlush    bool
		expectedRotation bool
	}{
		{
			name:             "Empty buffer NO overflow",
			segmentSize:      20,
			bufferContent:    []byte{},
			tmpBuffer:        []byte{1, 2, 3, 4, 5},
			expectedError:    nil,
			expectedBuffer:   []byte{1, 2, 3, 4, 5},
			expectedFlush:    false,
			expectedRotation: false,
		},
		{
			name:             "Buffer overflow, but can be flushed",
			segmentSize:      20,
			bufferContent:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
			tmpBuffer:        []byte{9, 10, 11},
			expectedError:    nil,
			expectedBuffer:   []byte{9, 10, 11},
			expectedFlush:    true,
			expectedRotation: false,
		},
		{
			name:             "Buffer needs rotation",
			segmentSize:      12,
			bufferContent:    []byte{1, 2, 3, 4, 5, 6, 7},
			tmpBuffer:        []byte{8, 9, 10, 11, 12},
			expectedError:    nil,
			expectedBuffer:   []byte{8, 9, 10, 11, 12},
			expectedFlush:    true,
			expectedRotation: true,
		},
		{
			name:             "Data bigger than buffer",
			segmentSize:      20,
			bufferContent:    []byte{},
			tmpBuffer:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expectedError:    fmt.Errorf("data is bigger than buffer, data cannot be handled"),
			expectedBuffer:   nil,
			expectedFlush:    false,
			expectedRotation: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flushedBuffer := &bytes.Buffer{}
			mockBuf := bufio.NewWriterSize(flushedBuffer, bufferSize)
			wal := &Wal{
				Options: &WalOptions{
					BufferSize:      uint32(bufferSize),
					SegmentSize:     uint32(tc.segmentSize),
					FileHandlerOpts: fh.DefaultOptions,
				},
				Buffer: mockBuf,
			}

			// Simulate the current size of the buffer
			wal.Buffer.Write(tc.bufferContent)
			wal.segmentUsed = len(tc.bufferContent)

			err := wal.manageWriteFlow(tc.tmpBuffer)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if mockBuf.Buffered() != len(tc.expectedBuffer) {
					t.Errorf("expected buffer %v, got %v", len(tc.expectedBuffer), mockBuf.Buffered())
				}
			}
		})
	}
}

func TestRecoverFile(t *testing.T) {
	testCases := []struct {
		name           string
		dataFile       []byte
		expectedResult []RecoveredEntry
	}{
		{
			name:     "Case 1: LSN=1, Hello World!",
			dataFile: []byte{1, 0, 0, 0, 12, 0, 0, 0, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 207, 169, 108, 170},
			expectedResult: []RecoveredEntry{
				{
					lsn:  1,
					data: []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33},
				},
			},
		},
		{
			name:     "Case 1: LSN=12, Hello World! + LSN=13, Bye World!",
			dataFile: []byte{12, 0, 0, 0, 12, 0, 0, 0, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 98, 175, 61, 28, 13, 0, 0, 0, 10, 0, 0, 0, 66, 121, 101, 32, 87, 111, 114, 108, 100, 33, 16, 211, 148, 16},
			expectedResult: []RecoveredEntry{
				{
					lsn:  12,
					data: []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33},
				},
				{
					lsn:  13,
					data: []byte{66, 121, 101, 32, 87, 111, 114, 108, 100, 33},
				},
			},
		},
	}

	// Writing existing data to file
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crate temp file
			file, err := os.CreateTemp("/tmp", "wal_test")
			if err != nil {
				t.Fatalf("Error creating temp file: %v", err)
			}
			defer os.Remove(file.Name()) // Make sure file is removed

			// Write test data into de file
			_, err = file.Write(tc.dataFile)
			if err != nil {
				t.Fatalf("Error writing data to file: %v", err)
			}

			// Close file after writing
			file.Close()

			// Reopen file to recover data
			file, err = os.Open(file.Name())
			if err != nil {
				t.Fatalf("Error reopening file for reading: %v", err)
			}
			defer file.Close() // Make sure file is closed

			// Recover data
			result, err := recoverFile(file)
			if err != nil {
				t.Fatalf("Error recovering file: %v", err)
			}

			// Check for expected results
			for i, entry := range tc.expectedResult {
				if i >= len(result) {
					t.Fatalf("Expected more entries in result")
				}

				if result[i].lsn != entry.lsn {
					t.Fatalf("Expected LSN %d, got %d", entry.lsn, result[i].lsn)
				}

				if !bytes.Equal(result[i].data, entry.data) {
					t.Fatalf("Expected data %v, got %v", entry.data, result[i].data)
				}
			}

			// Check if there is more entries (shall not be the case)
			if len(result) != len(tc.expectedResult) {
				t.Fatalf("Expected %d entries, got %d", len(tc.expectedResult), len(result))
			}
		})
	}
}

// TODO
// 1. Test using custom options

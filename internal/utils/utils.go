package utils

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const UINT32_MAX_NUMBER = 4294967295

// CalculateCRC returns the checksum of the data provided by the argument using IEEE poly
func CalculateCRC(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// IntToUint32 takes a integer and returns its conversion to uint32.
// If negative returns an error. If overflow caused by int bigger than uint32 returns an error.
func IntToUint32(i int) (uint32, error) {
	if (i < 0) || (i > UINT32_MAX_NUMBER) {
		return 0, errors.New("overflow transforming int to uint32")
	}
	return uint32(i), nil
}

// Uint32ToBytes takes an uint32 and return its conversion to bytes.
func Uint32ToBytes(i uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}

// BytesToUint32 takes a slice of bytes and return its conversion to uint32.
func BytesToUint32(i []byte) uint32 {
	return binary.LittleEndian.Uint32(i)
}

// appendBytesToSlice takes a slice of bytes as the first argument and another slice of bytes as the second argument.
// It returns a single slice with the second argument nested within the first.
func appendBytesToSlice(buf []byte, newElement []byte) []byte {
	return append(buf, newElement...)
}

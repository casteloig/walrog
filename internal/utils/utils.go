package utils

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const UINT32_MAX_NUMBER = 4294967295

// CalculateCRC returns the checksum of the data provided by the argument using the IEEE polynomial.
//
// Parameters:
//  - data: A slice of bytes for which the CRC checksum will be calculated.
//
// Returns:
//  - A uint32 value representing the CRC checksum.
func CalculateCRC(data []byte) uint32 {
    return crc32.ChecksumIEEE(data)
}

// IntToUint32 takes an integer and returns its conversion to uint32.
//
// Parameters:
//  - i: An integer value to be converted to uint32.
//
// Returns:
//  - A uint32 value if the conversion is successful.
//  - An error if the integer is negative or exceeds the maximum value of uint32.
func IntToUint32(i int) (uint32, error) {
    if (i < 0) || (i > UINT32_MAX_NUMBER) {
        return 0, errors.New("overflow transforming int to uint32")
    }
    return uint32(i), nil
}

// Uint32ToBytes takes a uint32 value and returns its conversion to a byte slice.
//
// Parameters:
//  - i: A uint32 value to be converted to a byte slice.
//
// Returns:
//  - A slice of 4 bytes representing the uint32 value in little-endian format.
func Uint32ToBytes(i uint32) []byte {
    buf := make([]byte, 4)
    binary.LittleEndian.PutUint32(buf, i)
    return buf
}

// BytesToUint32 takes a slice of bytes and returns its conversion to uint32.
//
// Parameters:
//  - i: A slice of 4 bytes to be converted to a uint32 value.
//
// Returns:
//  - A uint32 value represented by the byte slice in little-endian format.
func BytesToUint32(i []byte) uint32 {
    return binary.LittleEndian.Uint32(i)
}

// AppendBytesToSlice takes a slice of bytes as the first argument and another slice of bytes as the second argument.
// It returns a single slice with the second argument appended to the first.
//
// Parameters:
//  - buf: The original slice of bytes.
//  - newElement: The slice of bytes to append to the original slice.
//
// Returns:
//  - A new slice of bytes containing the original slice followed by the appended slice.
func AppendBytesToSlice(buf []byte, newElement []byte) []byte {
    return append(buf, newElement...)
}
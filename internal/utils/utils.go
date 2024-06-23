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

// Uint32ToBytes takes a uint32 and returns its conversion to []byte
func Uint32ToBytes(number uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, number)
	return buf
}

// IntToUint32 takes a integer and returns its conversion to uint32.
// If negative returns an error. If overflow caused by int bigger than uint32 returns an error.
func IntToUint32(i int) (uint32, error) {
	if (i < 0) || (i > UINT32_MAX_NUMBER) {
		return 0, errors.New("cannot convert negative int to uint32")
	}
	return uint32(i), nil
}

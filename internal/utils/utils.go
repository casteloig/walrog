package utils

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

// CalculateCRC returns the checksum of the data provided by the argument using IEEE poly
func CalculateCRC(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// func validateCRC32(data uint32) {
// 	crc32.
// }

func Uint32ToBytes(number uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, number)
	return buf
}

// Enhancement: avoid/alert overflow caused by int using 64 bits
func IntToUint32(i int) (uint32, error) {
	if i < 0 {
		return 0, errors.New("cannot convert negative int to uint32")
	}
	return uint32(i), nil
}

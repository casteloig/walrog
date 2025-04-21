package utils

import (
	"bytes"
	"testing"
)

func TestIntToUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    uint32
		wantErr bool
	}{
		{"positive int", 12345, 12345, false},
		{"zero", 0, 0, false},
		{"negative int", -1, 0, true},
		{"overflow int", 4294967296, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntToUint32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntToUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IntToUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32ToBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    uint32
		expected []byte
	}{
		{
			name:     "Zero",
			input:    0,
			expected: []byte{0, 0, 0, 0},
		},
		{
			name:     "One",
			input:    1,
			expected: []byte{1, 0, 0, 0},
		},
		{
			name:     "MaxUint32",
			input:    4294967295,
			expected: []byte{255, 255, 255, 255},
		},
		{
			name:     "Arbitrary number",
			input:    305419896, // 0x12345678
			expected: []byte{120, 86, 52, 18},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Uint32ToBytes(tc.input)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Uint32ToBytes(%d) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestBytesToUint32(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint32
	}{
		{
			name:     "Zero",
			input:    []byte{0, 0, 0, 0},
			expected: 0,
		},
		{
			name:     "One",
			input:    []byte{1, 0, 0, 0},
			expected: 1,
		},
		{
			name:     "MaxUint32",
			input:    []byte{255, 255, 255, 255},
			expected: 4294967295,
		},
		{
			name:     "Arbitrary number",
			input:    []byte{120, 86, 52, 18},
			expected: 305419896, // 0x12345678
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BytesToUint32(tc.input)
			if result != tc.expected {
				t.Errorf("Uint32ToBytes(%d) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestCalculateCRConBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "first example",
			input:    []byte{1, 0, 0, 0, 3, 0, 0, 0, 1, 2, 3},
			expected: []byte{190, 45, 28, 49},
		},
		{
			name:     "second example",
			input:    []byte{11, 0, 0, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: []byte{181, 89, 81, 142},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultCRC := CalculateCRC(tc.input)
			resultCRConBytes := Uint32ToBytes(resultCRC)
			if !bytes.Equal(resultCRConBytes, tc.expected) {
				t.Errorf("Uint32ToBytes(%d) = %v; want %v", tc.input, resultCRConBytes, tc.expected)
			}
		})
	}
}

func TestAppendBytesToSlice(t *testing.T) {
	testCases := []struct {
		name       string
		buf        []byte
		newElement []byte
		expected   []byte
	}{
		{
			name:       "Append to empty buffer",
			buf:        []byte{},
			newElement: []byte{1, 2, 3},
			expected:   []byte{1, 2, 3},
		},
		{
			name:       "Append to non-empty buffer",
			buf:        []byte{1, 2, 3},
			newElement: []byte{4, 5, 6},
			expected:   []byte{1, 2, 3, 4, 5, 6},
		},
		{
			name:       "Append empty slice to buffer",
			buf:        []byte{1, 2, 3},
			newElement: []byte{},
			expected:   []byte{1, 2, 3},
		},
		{
			name:       "Append slice to itself",
			buf:        []byte{1, 2, 3},
			newElement: []byte{1, 2, 3},
			expected:   []byte{1, 2, 3, 1, 2, 3},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendBytesToSlice(tt.buf, tt.newElement)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

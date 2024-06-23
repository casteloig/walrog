package utils

import (
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

// Helper function to compare byte slices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package file_handler

import (
	"os"
	"testing"
)

func TestWalOperations(t *testing.T) {
	// Common setup for all subtests
	os.RemoveAll("WalFolder")
	options := DefaultOptions

	err := CreateWalFolder(*options)
	if err != nil {
		t.Fatalf("CreateWalFolder() failed: %v", err)
	}

	t.Run("TestCreateWalFolder", func(t *testing.T) {
		// Verifies if the folder has been created
		_, err := os.Stat("WalFolder")
		if err != nil {
			t.Fatalf("Expected WalFolder to exist")
		}
	})

	t.Run("TestCreateWalFile", func(t *testing.T) {
		// Ensure the specific file does not exist
		os.Remove("WalFolder/wal_000.log")
		os.Remove("WalFolder/wal_001.log")

		// Calls test function
		file, err := CreateWalNewFile(*options)
		if err != nil {
			t.Fatalf("CreateWalNewFile() failed: %v", err)
		}

		// Verifies if the file has been created
		_, err = os.Stat("WalFolder/wal_000.log")
		if err != nil {
			t.Fatalf("Expected wal_000.log to exist")
		}

		// Verifies if the file pointer has been returned correctly
		if file == nil {
			t.Fatalf("Expected WalNewFile pointer to exist")
		}

		// Calls test function again, for second wal file
		file, err = CreateWalNewFile(*options)
		if err != nil {
			t.Fatalf("CreateWalNewFile() failed: %v", err)
		}

		// Verifies if the second file has been created
		_, err = os.Stat("WalFolder/wal_001.log")
		if err != nil {
			t.Fatalf("Expected wal_001.log to exist")
		}

		// Verifies if the second file pointer has been returned correctly
		if file == nil {
			t.Fatalf("Expected WalNewFile pointer to exist")
		}

	})

	// Clean up after all subtests
	os.RemoveAll("WalFolder")
}

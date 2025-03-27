package file_handler

import (
	"os"
	"path/filepath"
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
		_, err := os.Stat("/tmp/WalFolder")
		if err != nil {
			t.Fatalf("Expected WalFolder to exist")
		}
	})

	t.Run("TestCreateWalFile", func(t *testing.T) {
		// Ensure the specific file does not exist
		os.Remove("/tmp/WalFolder/wal_000.log")
		os.Remove("/tmp/WalFolder/wal_001.log")

		// Calls test function
		file, err := CreateWalNewFile(*options)
		if err != nil {
			t.Fatalf("CreateWalNewFile() failed: %v", err)
		}

		// Verifies if the file has been created
		_, err = os.Stat("/tmp/WalFolder/wal_000.log")
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
		_, err = os.Stat("/tmp/WalFolder/wal_001.log")
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

func TestCreateCheckpointFile(t *testing.T) {
    // Setup temporary directory for testing
    tempDir := t.TempDir()
    opts := Options{
        DirName:         tempDir,
        DirPerms:        0755,
        FilePerms:       0644,
        createFileFlags: os.O_CREATE | os.O_RDWR,
    }

    // Call CreateCheckpointFile
    file, err := CreateCheckpointFile(opts)
    if err != nil {
        t.Fatalf("CreateCheckpointFile failed: %v", err)
    }
    defer file.Close()

    // Verify the file exists
    checkpointPath := filepath.Join(tempDir, "checkpoint")
    if _, err := os.Stat(checkpointPath); os.IsNotExist(err) {
        t.Errorf("Checkpoint file was not created at %s", checkpointPath)
    }
}

func TestOpenWal(t *testing.T) {
    // Setup temporary directory for testing
    tempDir := t.TempDir()
    opts := &Options{
        DirName:         tempDir,
        DirPerms:        0755,
        FilePerms:       0644,
        createFileFlags: os.O_CREATE | os.O_RDWR,
    }

    // Call OpenWal
    walFile, checkpointFile, err := OpenWal(opts)
    if err != nil {
        t.Fatalf("OpenWal failed: %v", err)
    }
    defer walFile.Close()
    defer checkpointFile.Close()

    // Verify the WAL file exists
    walFilePath := filepath.Join(tempDir, "wal_000.log")
    if _, err := os.Stat(walFilePath); os.IsNotExist(err) {
        t.Errorf("WAL file was not created at %s", walFilePath)
    }

    // Verify the checkpoint file exists
    checkpointPath := filepath.Join(tempDir, "checkpoint")
    if _, err := os.Stat(checkpointPath); os.IsNotExist(err) {
        t.Errorf("Checkpoint file was not created at %s", checkpointPath)
    }
}
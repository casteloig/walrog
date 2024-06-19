package file_handler

import (
	"os"
	"testing"
)

// TestCreateWalFolder tests func createWalFolder()
func TestCreateWalFolder(t *testing.T) {
	// Limpia cualquier rastro anterior de pruebas
	os.RemoveAll("WalFolder")

	// Llama a la función que estamos probando
	err := CreateWalFolder()
	if err != nil {
		t.Fatalf("CreateWalFolder() failed: %v", err)
	}

	// Verifica que el directorio fue creado
	if _, err := os.Stat("WalFolder"); os.IsNotExist(err) {
		t.Fatalf("Expected WalFolder to exist")
	}

	// Verifica que el archivo de ejemplo fue creado
	if _, err := os.Stat("WalFolder/wal_0001.log"); os.IsNotExist(err) {
		t.Fatalf("Expected WalFolder/wal_0001.log to exist")
	}

	// // Limpia después de la prueba
	os.RemoveAll("WalFolder")
}

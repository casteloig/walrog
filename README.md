# Walrog

[![Build](https://github.com/casteloig/walrog/actions/workflows/main.yaml/badge.svg?branch=main)](https://github.com/casteloig/walrog/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/casteloig/walrog)](https://goreportcard.com/report/github.com/casteloig/walrog)
![Go Version](https://img.shields.io/badge/go-1.18+-blue)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)


Walrog is a Go library that implements a simple yet functional Write-Ahead Logging (WAL) system. It was designed as a personal project to understand and demonstrate the fundamentals of WAL, offering reliable data persistence and crash recovery mechanisms.

<div align="center">
  <img src="assets/images/walrog.png" width="260" alt="Walrog Logo" />
</div>

## ✨ Features

- Sequential logging with LSN, length, and CRC validation.
- Buffered writes to disk with automatic WAL file rotation.
- Recovery of valid records from existing WAL files.
- Configurable segmentation and initial checkpoint system.
- CRC-based data integrity checks.
- Unit tests covering the main functional use cases.

## 🚀 Basic Usage

### Installation

Clone the repository:

```bash
git clone https://github.com/casteloig/walrog.git
cd walrog
```

Make sure you have Go 1.18 or higher installed.

### Example

```golang
package main

import (
	"fmt"
	"github.com/casteloig/walrog/internal/core"
)

func main() {
	// Initialize with default options
	wal, err := core.InitWal(nil)
	if err != nil {
		panic(err)
	}

	// Write a message
	err = wal.WriteBuffer([]byte("Hello World!"))
	if err != nil {
		panic(err)
	}

	// Flush to disk
	err = wal.FlushBuffer()
	if err != nil {
		panic(err)
	}

	fmt.Println("Entry successfully written.")
}
```

## 📄 License

MIT License. See the [LICENSE](LICENSE) file for more details.
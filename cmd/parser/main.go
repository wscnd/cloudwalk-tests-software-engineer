package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/wscnd/cloudwalk-tests-software-engineer/internal/logparser"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

var (
	ErrLogInputProvided = errors.New("please provide a file path as an argument")
	ErrOpeningFile      = errors.New("failed to open file")
	ErrClosingFile      = errors.New("failed to close file")
)

func run() error {
	// -------------------------------------------------------------------------
	// Arguments processing
	if len(os.Args) < 2 {
		return ErrLogInputProvided
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrOpeningFile, err)
	}
	defer file.Close()
	stats, _ := file.Stat()
	fmt.Printf("Successfully opened and processed file: %s\nSize: %d bytes\n\n", filePath, stats.Size())

	// -------------------------------------------------------------------------
	// Parsing file

	logparser.Run(file)

	return nil
}

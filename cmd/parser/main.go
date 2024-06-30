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
	fmt.Println("Parsing completed. JSON file generated.")
}

var (
	ErrLogInputProvided = errors.New("please provide a file path as an argument")
	ErrOpeningFile      = errors.New("failed to open file")
	ErrClosingFile      = errors.New("failed to close file")
	ErrParsingLogFile   = errors.New("failed to parse log file")
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

	// -------------------------------------------------------------------------
	// Parsing file

	if err := logparser.Run(file); err != nil {
		return fmt.Errorf("%w: %s", ErrParsingLogFile, err)
	}

	return nil
}

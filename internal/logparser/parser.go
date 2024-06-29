package logparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type LogParser struct {
	logfile *os.File
}

func NewLogParser(logFile *os.File) *LogParser {
	return &LogParser{
		logfile: logFile,
	}
}

type Match struct {
	Matches int
}

func (lp *LogParser) detectMatches() int {
	scanner := bufio.NewScanner(lp.logfile)

	var count Match
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "InitGame:") {
			count.Matches++
		}
	}
	return count.Matches
}

func Run(file *os.File) {
	parser := NewLogParser(file)
	matchesDetected := parser.detectMatches()

	fmt.Printf("%+v matches found!\n", matchesDetected)
}

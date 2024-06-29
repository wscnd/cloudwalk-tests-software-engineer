package logparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type LogParser struct {
	logfile io.Reader
	matches *Match
}

func NewLogParser(logFile io.Reader, match *Match) *LogParser {
	return &LogParser{
		logfile: logFile,
		matches: match,
	}
}

type Match struct {
	Matches [][]string
}

func (lp *LogParser) detectMatches() error {
	scanner := bufio.NewScanner(lp.logfile)

	var lines []string
	var inMatch bool = false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "InitGame:") {
			// we are in a game
			if inMatch {
				lp.matches.Matches = append(lp.matches.Matches, lines)
				lines = nil
				inMatch = false
			} else {
				inMatch = true
			}
		} else {
			// lines with "---" are ignored
			if !strings.Contains(line, "---") {
				inMatch = true
				lines = append(lines, line)
			}
		}
	}
	if len(lines) != 0 {
		lp.matches.Matches = append(lp.matches.Matches, lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Run(file *os.File) {
	m := &Match{
		Matches: make([][]string, 0, 21),
	}

	parser := NewLogParser(file, m)
	matchesDetected := parser.detectMatches()

	fmt.Printf("%+v matches found!\n", matchesDetected)
}

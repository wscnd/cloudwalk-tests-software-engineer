package logparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	ErrScanningLogFile  = errors.New("error scanning log file")
	ErrDetectingMatches = errors.New("error detecting matches")
)

type LogParser struct {
	logfile       io.Reader
	matchesLog    [][]string
	matchesFound  chan []string
	matchesParsed chan *Match
}

func NewLogParser(logFile io.Reader, log [][]string) *LogParser {
	lp := &LogParser{
		matchesFound:  make(chan []string),
		matchesParsed: make(chan *Match),
	}

	if logFile != nil {
		lp.logfile = logFile
	}

	if log != nil {
		lp.matchesLog = log
	}
	return lp
}

func (lp *LogParser) parseMatchEvents(loglines []string) *Match {
	match := NewMatch(loglines)

	for _, line := range loglines {
		switch {
		// Process Kill Event
		case strings.Contains(line, "Kill"):
			lp.parseKills(line, match)

		// Parse ClientUserinfoChanged Event
		case strings.Contains(line, "ClientUserinfoChanged"):
			logs := strings.Split(line, "ClientUserinfoChanged: ")
			match.updatePlayerInfo(logs[1])
		}
	}
	return match
}

func (*LogParser) parseKills(line string, match *Match) {
	eventData := strings.Fields(line)

	match.TotalKills++
	match.updateKillCauses(eventData)

	killerID := eventData[2]
	victimID := eventData[3]

	switch {
	// case world killed
	case killerID == "1022":
		match.updatePlayerDeaths(victimID)

	// case player killed another player
	case killerID != victimID:
		match.updatePlayerKill(killerID)
		match.updatePlayerDeaths(victimID)

	// case player killed itself
	case killerID == victimID:
		match.updatePlayerDeaths(victimID)
	}
}

func (lp *LogParser) detectMatches(logfile io.Reader) error {
	scanner := bufio.NewScanner(logfile)
	defer close(lp.matchesFound)

	var matchLines []string
	var inMatch bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "InitGame:") {
			// we are in a game
			if inMatch {
				lp.matchesFound <- matchLines
				matchLines = nil
				inMatch = false
			} else {
				inMatch = true
			}
		} else {
			// lines with "---" are ignored
			if !strings.Contains(line, "---") {
				inMatch = true
				matchLines = append(matchLines, line)
			}
		}
	}
	// Edge case of the last InitGame processed
	if len(matchLines) != 0 {
		lp.matchesFound <- matchLines
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%w: %s", ErrScanningLogFile, err)
	}

	return nil
}

func (lp *LogParser) parseMatches() {
	defer close(lp.matchesParsed)

	for lines := range lp.matchesFound {
		lp.matchesParsed <- lp.parseMatchEvents(lines)
	}
}

func (lp *LogParser) sequentialDetectMatches() error {
	scanner := bufio.NewScanner(lp.logfile)

	var lines []string
	var inMatch bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "InitGame:") {
			// we are in a game
			if inMatch {
				lp.matchesLog = append(lp.matchesLog, lines)
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

	// Edge case of the last InitGame processed
	if len(lines) != 0 {
		lp.matchesLog = append(lp.matchesLog, lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Run(file *os.File) error {
	parser := NewLogParser(nil, nil)

	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := parser.detectMatches(file); err != nil {
			errChan <- fmt.Errorf("%w: %s", ErrDetectingMatches, err)
		}
	}()

	go func() {
		defer close(errChan)
		wg.Wait()
	}()

	go func() {
		parser.parseMatches()
	}()

	select {
	case err := <-errChan:
		return err
	default:
	}

	var matchData Matches = make(map[string]*Match)
	for match := range parser.matchesParsed {
		matchData[fmt.Sprintf("game-%d", len(matchData)+1)] = match
	}

	if err := matchData.toJSON(); err != nil {
		return err
	}

	return nil
}

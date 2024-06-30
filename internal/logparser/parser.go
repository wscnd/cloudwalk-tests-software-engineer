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
	matchesLog    [][]string    // Accumulated log lines grouped by matches
	matchesFound  chan []string // Channel for correctly identified match lines
	matchesParsed chan *Match   // Channel for parsed matches.
}

// NewLogParser creates a new LogParser instance.
// The arguments logFile and log are only used for testing.
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

// parseMatchEvents parses the given log lines from a indivual match into a Match object.
func (lp *LogParser) parseMatchEvents(loglines []string) *Match {
	match := NewMatch(loglines)

	for _, line := range loglines {
		switch {
		// Handle Kill events
		case strings.Contains(line, "Kill"):
			lp.parseKills(line, match)

		// Handle ClientUserinfoChanged events
		case strings.Contains(line, "ClientUserinfoChanged"):
			logs := strings.Split(line, "ClientUserinfoChanged: ")
			match.updatePlayerInfo(logs[1])
		}
	}
	return match
}

// parseKills updates the match object with kill event data.
func (*LogParser) parseKills(line string, match *Match) {
	eventData := strings.Fields(line)

	match.TotalKills++
	match.updateKillCauses(eventData)

	killerID := eventData[2]
	victimID := eventData[3]

	switch {
	// Handle the case where world is the killer
	case killerID == "1022":
		match.updatePlayerDeaths(victimID)

	// Handle the case where a player kills another player
	case killerID != victimID:
		match.updatePlayerKill(killerID)
		match.updatePlayerDeaths(victimID)

	// Caveat 4: Handle the case where a player kills themselves
	case killerID == victimID:
		match.updatePlayerDeaths(victimID)
	}
}

// detectMatches scans the log file, detect match boundaries
// and sends found lines to matchesFound channel
func (lp *LogParser) detectMatches(logfile io.Reader) error {
	scanner := bufio.NewScanner(logfile)
	defer close(lp.matchesFound)

	var matchLines []string
	var inMatch bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "InitGame:") {
			// We are in a game
			if inMatch {
				lp.matchesFound <- matchLines
				matchLines = nil
				inMatch = false
			} else {
				inMatch = true
			}
		} else {
			// Append lines
			if !strings.Contains(line, "---") {
				isKill := strings.Contains(line, "Kill")
				isPlayerUpdate := strings.Contains(line, "ClientUserinfoChanged")

				inMatch = true
				if isKill || isPlayerUpdate {
					matchLines = append(matchLines, line)
				}
			}
		}
	}
	// Edge case of the last InitGame processed
	// We send the remaining lines
	if len(matchLines) != 0 {
		lp.matchesFound <- matchLines
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%w: %s", ErrScanningLogFile, err)
	}

	return nil
}

// parseMatches processes the matches found and sent to matchesFound channel
// parses them and sends to matchesParsed channel.
func (lp *LogParser) parseMatches() {
	defer close(lp.matchesParsed)

	for lines := range lp.matchesFound {
		lp.matchesParsed <- lp.parseMatchEvents(lines)
	}
}

// Run initializes and runs the LogParser instance on the provided file.
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

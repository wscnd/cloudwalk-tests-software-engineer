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
	wg            *sync.WaitGroup
}

// NewLogParser creates a new LogParser instance.
// The arguments logFile and log are only used for testing.
func NewLogParser(logFile io.Reader, log [][]string) *LogParser {
	lp := &LogParser{
		matchesFound:  make(chan []string, 5),
		matchesParsed: make(chan *Match, 5),
		wg:            &sync.WaitGroup{},
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
func (lp *LogParser) detectMatches(logfile io.Reader, errChan chan<- error) {
	defer lp.wg.Done()
	defer close(lp.matchesFound)
	defer close(errChan)

	scanner := bufio.NewScanner(logfile)

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

				// errChan <- fmt.Errorf("%w: if we break here we should halt everything", ErrScanningLogFile)
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
		errChan <- fmt.Errorf("%w: %s", ErrScanningLogFile, err)
	}
}

// parseMatches processes the matches found and sent to matchesFound channel
// parses them and sends to matchesParsed channel.
func (lp *LogParser) parseMatches() {
	defer lp.wg.Done()
	defer close(lp.matchesParsed)

	for lines := range lp.matchesFound {
		lp.matchesParsed <- lp.parseMatchEvents(lines)
	}
}

// Run initializes and runs the LogParser instance on the provided file.
func Run(file *os.File) error {
	parser := NewLogParser(nil, nil)

	errChan := make(chan error, 1)

	parser.wg.Add(2)
	go parser.detectMatches(file, errChan)
	go parser.parseMatches()

	matchData := make(Matches)
	for {
		select {
		case err, ok := <-errChan:
			if ok && err != nil {
				return err
			}
			// The errChan is only being used in detectMatches, once we're done
			// detecting, we close it and this case is disabled.
			errChan = nil
		case match, moreData := <-parser.matchesParsed:
			if !moreData {
				parser.matchesParsed = nil
			} else {
				matchData[fmt.Sprintf("game-%d", len(matchData)+1)] = match
			}
		}

		if errChan == nil && parser.matchesParsed == nil {
			break
		}
	}

	if err := matchData.toJSON(); err != nil {
		return err
	}

	return nil
}

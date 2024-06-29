package logparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type LogParser struct {
	logfile    io.Reader
	matchesLog [][]string
}

func NewLogParser(logFile io.Reader, log [][]string) *LogParser {
	return &LogParser{
		logfile:    logFile,
		matchesLog: log,
	}
}

func (lp *LogParser) processMatches() *Matches {
	var matchData Matches
	for _, lines := range lp.matchesLog {
		matchData = append(matchData, lp.parseMatchEvents(lines))
	}

	for matchID, matchData := range matchData {
		fmt.Printf("Match ID: %d\n", matchID)
		fmt.Printf("Total Kills: %d\n", matchData.TotalKills)
		fmt.Println("Players:")
		for _, player := range matchData.Players {
			fmt.Printf("\tName: %s, Kills: %d, Deaths: %d\n", player.Name, player.Kills, player.Deaths)
		}
	}
	return &matchData
}

func (lp *LogParser) parseMatchEvents(lines []string) *Match {
	match := &Match{
		Players:  make(map[string]*PlayerData),
		MatchLog: lines,
	}

	for _, line := range lines {
		switch {
		// Process Kill Event
		case strings.Contains(line, "Kill"):
			eventData := strings.Fields(line)

			match.TotalKills++

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

		// Parse ClientUserinfoChanged Event
		case strings.Contains(line, "ClientUserinfoChanged"):
			logs := strings.Split(line, "ClientUserinfoChanged: ")
			match.updatePlayerInfo(logs[1])
		}
	}
	return match
}

func (lp *LogParser) detectMatches() error {
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

	// Finished parsing the last match
	if len(lines) != 0 {
		lp.matchesLog = append(lp.matchesLog, lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Run(file *os.File) error {
	log := make([][]string, 0, 21)

	parser := NewLogParser(file, log)
	if err := parser.detectMatches(); err != nil {
		return err
	}

	md := parser.processMatches()
	if err := md.toJSON(); err != nil {
		return err
	}

	return nil
}

package logparser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
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

type (
	PlayerData struct {
		Name   string `json:"name"`
		Kills  int    `json:"kills"`
		Deaths int    `json:"deaths"`
	}

	Match struct {
		TotalKills int
		Players    map[string]*PlayerData
		MatchLog   []string
	}

	Matches []*Match
)

func (m *Match) MarshalJSON() ([]byte, error) {
	type matchDataJSON struct {
		TotalKills int            `json:"total_kills"`
		Players    []string       `json:"players"`
		Kills      map[string]int `json:"kills"`
	}

	data := &matchDataJSON{
		TotalKills: m.TotalKills,
		Players:    make([]string, 0, len(m.Players)),
		Kills:      make(map[string]int, m.TotalKills),
	}

	for _, player := range m.Players {
		data.Players = append(data.Players, player.Name)
		data.Kills[player.Name] = player.Kills
	}

	return json.Marshal(data)
}

func (m *Matches) toJSON() error {
	output := make(map[string]*Match)

	for id, matchData := range *m {
		matchID := "game-" + strconv.Itoa(id+1)
		output[matchID] = matchData
	}
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("error generating json output: %w", err)
	}

	err = os.WriteFile("match_data.json", jsonOutput, 0o644)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
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
				if p, ok := match.Players[victimID]; !ok {
					match.Players[victimID] = &PlayerData{
						Deaths: 1,
					}
				} else {
					p.Deaths++
				}

			// case player killed another player
			case killerID != victimID:
				if p, ok := match.Players[killerID]; !ok {
					match.Players[killerID] = &PlayerData{
						Kills: 1,
					}
				} else {
					p.Kills++
				}
				if p, ok := match.Players[victimID]; !ok {
					match.Players[victimID] = &PlayerData{
						Deaths: 1,
					}
				} else {
					p.Deaths++
				}

			// case player killed itself
			case killerID == victimID:
				if p, ok := match.Players[victimID]; !ok {
					match.Players[victimID] = &PlayerData{
						Deaths: 1,
					}
				} else {
					p.Deaths++
				}
			}

		// Parse ClientUserinfoChanged Event
		case strings.Contains(line, "ClientUserinfoChanged"):
			logs := strings.Split(line, "ClientUserinfoChanged: ")

			eventData := logs[1]
			playerID := strings.Split(eventData, " ")[0]
			startIndex := strings.Index(eventData, "n\\") + len("n\\")
			endIndex := strings.Index(eventData, "\\t")
			playerNickName := eventData[startIndex:endIndex]

			if _, ok := match.Players[playerID]; !ok {
				match.Players[playerID] = &PlayerData{
					Name: playerNickName,
				}
			} else {
				match.Players[playerID].Name = playerNickName
			}
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
	if len(lines) != 0 {
		lp.matchesLog = append(lp.matchesLog, lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Run(file *os.File) {
	log := make([][]string, 0, 21)
	parser := NewLogParser(file, log)
	parser.detectMatches()
	md := parser.processMatches()
	md.toJSON()
}

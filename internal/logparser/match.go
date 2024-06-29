package logparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Match struct {
	TotalKills   int
	Players      map[string]*PlayerData
	MatchLog     []string
	KillsByMeans map[meansOfDeath]int
}

func NewMatch(loglines []string) *Match {
	return &Match{
		Players:      make(map[string]*PlayerData),
		MatchLog:     loglines,
		KillsByMeans: make(map[meansOfDeath]int),
	}
}

func (m *Match) updatePlayerKill(killerID string) {
	if p, ok := m.Players[killerID]; !ok {
		m.Players[killerID] = &PlayerData{
			Kills: 1,
		}
	} else {
		p.Kills++
	}
}

func (m *Match) updatePlayerDeaths(victimID string) {
	if p, ok := m.Players[victimID]; !ok {
		m.Players[victimID] = &PlayerData{
			Deaths: 1,
		}
	} else {
		p.Deaths++
	}
}

func (m *Match) updatePlayerInfo(log string) {
	playerID := strings.Split(log, " ")[0]
	startIndex := strings.Index(log, "n\\") + len("n\\")
	endIndex := strings.Index(log, "\\t")
	playerNickName := log[startIndex:endIndex]

	if _, ok := m.Players[playerID]; !ok {
		m.Players[playerID] = &PlayerData{
			Name: playerNickName,
		}
	} else {
		m.Players[playerID].Name = playerNickName
	}
}

func (m *Match) updateKillCauses(log []string) {
	lastPosition := len(log) - 1
	deathCause := log[lastPosition]
	cause := meansOfDeath(deathCause)

	m.KillsByMeans[cause]++
}

func (m *Match) MarshalJSON() ([]byte, error) {
	type matchDataJSON struct {
		TotalKills   int            `json:"total_kills"`
		Players      []string       `json:"players"`
		Kills        map[string]int `json:"kills"`
		KillsByMeans map[string]int `json:"kills_by_means"`
	}

	data := &matchDataJSON{
		TotalKills:   m.TotalKills,
		Players:      make([]string, 0, len(m.Players)),
		Kills:        make(map[string]int, m.TotalKills),
		KillsByMeans: make(map[string]int, len(m.KillsByMeans)),
	}

	for _, player := range m.Players {
		data.Players = append(data.Players, player.Name)
		data.Kills[player.Name] = player.Kills
	}

	for mod, count := range m.KillsByMeans {
		data.KillsByMeans[string(mod)] = count
	}

	return json.Marshal(data)
}

type Matches []*Match

func (ms *Matches) toJSON() error {
	output := make(map[string]*Match)

	for id, matchData := range *ms {
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

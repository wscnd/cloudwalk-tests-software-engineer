package logparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type (
	Match struct {
		TotalKills int
		Players    map[string]*PlayerData
		MatchLog   []string
	}

	Matches []*Match
)

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

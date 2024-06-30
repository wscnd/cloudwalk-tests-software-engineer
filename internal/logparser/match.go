package logparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	ErrParsingJSON     = errors.New("error generating json output: %w")
	ErrWritingJSONFile = errors.New("error writing json file: %w")
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
		KillsByMeans map[string]int `json:"kills_by_means,omitempty"`
	}

	data := matchDataJSON{
		TotalKills:   m.TotalKills,
		Players:      make([]string, 0, len(m.Players)),
		Kills:        make(map[string]int, len(m.Players)),
		KillsByMeans: make(map[string]int, len(m.KillsByMeans)),
	}

	for _, player := range m.Players {
		data.Players = append(data.Players, player.Name)
		data.Kills[player.Name] = player.Kills
	}

	sort.Slice(data.Players, func(i, j int) bool {
		return data.Kills[data.Players[i]] > data.Kills[data.Players[j]]
	})

	for mod, count := range m.KillsByMeans {
		data.KillsByMeans[string(mod)] = count
	}

	return json.Marshal(data)
}

type Matches map[string]*Match

func (ms *Matches) toJSON() error {
	output := make(map[string]*Match)

	for matchID, matchData := range *ms {
		output[matchID] = matchData
	}
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrParsingJSON, err)
	}

	err = os.WriteFile("match_data.json", jsonOutput, 0o644)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrWritingJSONFile, err)
	}

	return nil
}

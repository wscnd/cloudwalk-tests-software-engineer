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
	Score    int
	Players  map[string]*Player
	MatchLog []string
}

func NewMatch(loglines []string) *Match {
	return &Match{
		Players:  make(map[string]*Player),
		MatchLog: loglines,
	}
}

func (m *Match) findOrCreatePlayer(playerID string) *Player {
	p, ok := m.Players[playerID]
	if !ok {
		p = NewPlayer(playerID)
		m.Players[p.ID] = p
	}
	return p
}

func (m *Match) OnKill(log []string) {
	// Extract DeathCause
	lastPosition := len(log) - 1
	deathCause := log[lastPosition]
	cause := meansOfDeath(deathCause)

	// Update Match Score
	m.Score++

	// Extract Victim ID
	victimID := log[3]

	// There's always a victim
	// Caveat 4: Handle the case where a player kills themselves
	v := m.findOrCreatePlayer(victimID)
	v.RecordDeath(cause)

	// Extract Killer ID
	// Handle the case where world is the killer
	killerID := log[2]
	if killerID == "1022" {
		return
	}

	// Handle the case where a player kills another player
	k := m.findOrCreatePlayer(killerID)
	k.RecordKill(v)
}

func (m *Match) OnPlayerInfoChange(log string) {
	playerID := strings.Split(log, " ")[0]
	startIndex := strings.Index(log, "n\\") + len("n\\")
	endIndex := strings.Index(log, "\\t")
	name := log[startIndex:endIndex]

	p := m.findOrCreatePlayer(playerID)
	p.UpdateNickname(name)
}

func (m *Match) MarshalJSON() ([]byte, error) {
	type matchDataJSON struct {
		Score        int            `json:"total_kills"`
		Players      []string       `json:"players"`
		Kills        map[string]int `json:"kills"`
		KillsByMeans map[string]int `json:"kills_by_means,omitempty"`
	}

	data := matchDataJSON{
		Score:        m.Score,
		Players:      make([]string, 0, len(m.Players)),
		Kills:        make(map[string]int, len(m.Players)),
		KillsByMeans: make(map[string]int, len(m.Players)),
	}

	for _, player := range m.Players {
		data.Players = append(data.Players, player.Name)
		data.Kills[player.Name] = player.Kills

		for _, dc := range player.DeathCauses {
			_, ok := data.KillsByMeans[string(dc)]
			if !ok {
				data.KillsByMeans[string(dc)] = 1
			} else {
				data.KillsByMeans[string(dc)] += +1
			}

		}
	}

	sort.Slice(data.Players, func(i, j int) bool {
		return data.Kills[data.Players[i]] > data.Kills[data.Players[j]]
	})

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

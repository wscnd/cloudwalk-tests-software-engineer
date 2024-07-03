package logparser

type Player struct {
	ID          string         `json:"-"`
	Name        string         `json:"name"`
	Kills       int            `json:"kills"`
	DeathCauses []meansOfDeath `json:"deathCauses"`
}

func (p *Player) RecordKill(victim *Player) {
	// Handle the case where a player kills another player
	p.Kills++
}

func (p *Player) RecordDeath(mean meansOfDeath) {
	p.DeathCauses = append(p.DeathCauses, mean)
}

func (p *Player) UpdateNickname(name string) {
	p.Name = name
}

func NewPlayer(id string) *Player {
	p := new(Player)
	p.ID = id
	return p
}

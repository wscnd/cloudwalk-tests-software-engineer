package logparser

type PlayerData struct {
	Name   string `json:"name"`
	Kills  int    `json:"kills"`
	Deaths int    `json:"deaths"`
}

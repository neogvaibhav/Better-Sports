package helpers

type Team struct {
	Name    string   `json:"name,omitempty"`
	Players []string `json:"players,omitempty"`
}

type UpdateTeamPlayers struct {
	ID        string   `json:"id,omitempty"`
	PlayerIds []string `json:"player_ids,omitempty"`
}

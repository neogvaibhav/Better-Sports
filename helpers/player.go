package helpers

type Player struct {
	Name     string `json:"name,omitempty"`
	Grade    string `json:"grade,omitempty"`
	Position string `json:"position,omitempty"`
	Goals    *int   `json:"goals"`
	Assists  *int   `json:"assists"`
	Fouls    *int   `json:"fouls"`
}

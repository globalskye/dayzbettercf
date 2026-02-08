package player

import "time"

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Members   []Member  `json:"members,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Member struct {
	GroupID   int64     `json:"group_id"`
	PlayerID  int64     `json:"player_id"`
	CftoolsID string    `json:"cftools_id"`
	Alias     string    `json:"alias"`
	Player    *Player   `json:"player,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

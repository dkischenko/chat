package models

type Message struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	UFrom     string `json:"UFrom"`
	CreatedAt int    `json:"created_at"`
}

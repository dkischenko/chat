package models

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"userName" bson:"username"`
	PasswordHash string `json:"passwordHash" bson:"password"`
	Key          string `json:"key" bson:"key"`
	IsOnline     bool   `json:"isOnline" bson:"isOnline"`
	LastOnline   *int   `json:"lastOnline"`
}

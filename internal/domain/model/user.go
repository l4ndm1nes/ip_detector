package model

type User struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IP           string `json:"ip"`
	Country      string `json:"country,omitempty"`
	PasswordHash string `json:"-"`
}

package models

type Request struct {
	ID                int    `json:"id"`
	Date1             string `json:"date1"`
	Date2             string `json:"date2"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password"`
	Token             string `json:"token,omitempty"`
	Role              string `json:"role,omitempty"`
	Line              int    `json:"line"`
}
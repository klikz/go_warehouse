package models

type User struct {
	ID                int    `json:"-"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password"`
	Token             string `json:"token,omitempty"`
	Role              string `json:"role,omitempty"`
}

type ParsedToken struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

package models

type Model struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Comment string `json:"comment"`
}

package models

type Akt struct {
	UserName      string `json:"username"`
	Component_id  int    `json:"component_id"`
	Comment       string `json:"comment"`
	UserID        int
	Quantity      float64
	Checkpoint_id int
	Photo         string `json:"photo"`
}

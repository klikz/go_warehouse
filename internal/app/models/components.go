package models

type Component struct {
	Available     float64 `json:"available"`
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Checkpoint    string  `json:"checkpoint"`
	Checkpoint_id int     `json:"checkpoint_id"`
	Unit          string  `json:"unit"`
	Specs         string  `json:"specs"`
	Photo         string  `json:"photo"`
	Time          string  `json:"time"`
	Type          string  `json:"type"`
	Type_id       int     `json:"type_id"`
	Weight        float64 `json:"weight"`
	Status        int     `jspn:"status,omitempty"`
	Token         string  `json:"token,omitempty"`
	Quantity      float64 `json:"quantity,omitempty"`
}

type Last struct {
	Serial        string `json:"serial"`
	Model_id      int    `json:"model_id"`
	Model         string `json:"model"`
	Checkpoint_id int    `json:"checkpoint_id"`
	Line          string `json:"line"`
	Product_id    int    `json:"product_id"`
	Time          string `json:"time"`
}

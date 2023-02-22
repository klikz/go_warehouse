package models

type ImportModel struct {
	Token       string  `json:"token"`
	Name        string  `json:"name"`
	LotID       int     `json:"lot_id"`
	Comment     string  `json:"comment"`
	BatchID     int     `json:"batch_id"`
	ContainerID int     `json:"container_id"`
	Quantity    float64 `json:"quantity"`
	R_Quantity  float64 `json:"r_quantity"`
	ComponentID int     `json:"component_id"`
	IncomeID    int     `json:"income_id"`
	File64      string  `json:"file64"`
	Unit        string  `json:"unit"`
}

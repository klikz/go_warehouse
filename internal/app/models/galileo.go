package models

type Galileo struct {
	Barcode        string `json:"Barcode"`
	Name           string `json:"name"`
	OpCode         string `json:"OpCode"`
	TypeFreon      string `json:"TypeFreon"`
	Result         string `json:"Result"`
	ProgQuantity   string `json:"ProgQuantity"`
	Quantity       string `json:"Quantity"`
	CycleTotalTime string `json:"CycleTotalTime"`
	Time           string `json:"Time"`
}

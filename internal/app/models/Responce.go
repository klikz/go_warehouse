package models

type Responce struct {
	Result string      `json:"result"`
	Err    interface{} `json:"error"`
	Data   interface{} `json:"data"`
}

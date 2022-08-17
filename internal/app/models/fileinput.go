package models

type FileInput struct {
	Code          string  `title:"code"` // sheet可选，不声明则选择首个sheet页读写
	Name          string  `title:"name"`
	Unit          string  `title:"unit"`
	Quantity      float64 `title:"quantity"`
	Available     float64
	Checkpoint_id int
	ID            int
}

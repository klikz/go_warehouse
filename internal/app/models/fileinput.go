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

type FileInput2 struct {
	Code     string  `title:"code"`
	Unit     string  `title:"unit"`
	Quantity float64 `title:"quantity"`
	Comment  string  `title:"comment"`
}

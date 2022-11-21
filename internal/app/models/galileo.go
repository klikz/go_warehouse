package models

type Galileo struct {
	Serial          string  `json:"serial"`
	OpCode          string  `json:"opcode"`
	TypeFreon       string  `json:"type_freon"`
	PreVacuum       float32 `json:"pre_vacuum"`
	ConturPressure  float32 `json:"contur_pressure"`
	Vacuum          float32 `json:"vacuum"`
	PoiskUtechek    float32 `json:"poisk_utechek"`
	RealQuantity    float32 `json:"real_quantity"`
	ProgramQuantity float32 `json:"program_quantity"`
	RefPressure     float32 `json:"ref_pressure"`
	RefTemp         float32 `json:"ref_temp"`
	Time            string  `json:"time"`
}

package models

type Laboratory struct {
	StartTime  string `json:"start_time"`
	Serial     string `json:"serial"`
	Line       int    `json:"line"`
	Point      int    `json:"point"`
	EndTime    string `json:"end_time"`
	Duration   int    `json:"duration"`
	Model      string `json:"model"`
	Compressor string `json:"compressor"`
	Result     string `json:"result"`
}

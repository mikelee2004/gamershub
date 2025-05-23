package models

type Statistics struct {
	Rank      string  `json:"rank"`
	PeakRank  string  `json:"peak_rank"`
	KD        float32 `json:"kd"`
	WinRate   float32 `json:"win_rate"`
	HSPercent float32 `json:"hs_percent"`
}

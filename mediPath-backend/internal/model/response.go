package model

type PredictResponse struct {
	Disease     string  `json:"disease"`
	Confidence  float64 `json:"confidence"`
	Explanation string  `json:"explanation"`
	Routine     string  `json:"routine"`
}

package model

type PredictRequest struct {
	Drugs []string `json:"drugs"`
}
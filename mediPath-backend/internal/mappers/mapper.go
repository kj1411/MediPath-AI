package mappers

import "mediPath-backend/internal/model"

var ATCCodeToDisease = map[string]string{
	"A": "Digestive system disorder",
	"B": "Blood disorder",
	"C": "Cardiovascular disease",
	"D": "Skin disease",
	"G": "Genitourinary disease",
	"H": "Hormonal disorder",
	"J": "Infectious disease",
	"L": "Cancer or immune disorder",
	"M": "Musculoskeletal disorder",
	"N": "Neurological disorder",
	"R": "Respiratory disease",
	"V": "General therapeutic or miscellaneous condition",
}

func MapATCToDisease(code string) string {
	if disease, ok := ATCCodeToDisease[code]; ok {
		return disease
	}
	return "Unknown Disease"
}

func MapPredictionToResponse(
	disease string,
	confidence float64,
	explanation string,
	routine string,
) model.PredictResponse {

	return model.PredictResponse{
		Disease:     disease,
		Confidence:  confidence,
		Explanation: explanation,
		Routine:     routine,
	}
}

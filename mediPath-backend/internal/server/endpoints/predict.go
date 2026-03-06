package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"mediPath-backend/internal/mappers"
	"mediPath-backend/internal/model"
	"mediPath-backend/internal/services"
)

type PredictHandler struct {
	ML    *services.MLService
	Agent *services.AgentService
	Drug  *services.DrugService
}

func (h *PredictHandler) Handle(w http.ResponseWriter, r *http.Request) {

	var req model.PredictRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Drugs) == 0 {
		http.Error(w, "Drugs list cannot be empty", http.StatusBadRequest)
		return
	}

	// Normalize drugs
	drugs := h.Drug.NormalizeDrugs(req.Drugs)

	log.Println("Received drugs:", drugs)

	// ML prediction
	mlResult, err := h.ML.Predict(drugs)
	if err != nil {
		log.Println("ML Error:", err)
		http.Error(w, "ML prediction failed", http.StatusInternalServerError)
		return
	}

	log.Println("ML Result:", mlResult)

	atcCode := mlResult["disease"].(string)
	disease := mappers.MapATCToDisease(atcCode)
	if disease == "Unknown" {
		resp := mappers.MapPredictionToResponse(
			"Unknown",
			0,
			"Unable to determine condition from provided medicines.",
			"Please consult a healthcare professional.",
		)

		json.NewEncoder(w).Encode(resp)
		return
	}

	confidence, ok := mlResult["confidence"].(float64)
	if !ok {
		http.Error(w, "Invalid ML confidence output", 500)
		return
	}

	// AI explanation
	explanation, routine, err :=
		h.Agent.GeneratePatientGuidance(disease)

	if err != nil {
		log.Println("AI Error:", err)
	}

	resp := mappers.MapPredictionToResponse(
		disease,
		confidence,
		explanation,
		routine,
	)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Response error:", err)
	}
}

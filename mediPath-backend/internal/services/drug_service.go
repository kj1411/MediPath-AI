package services

import "strings"

type DrugService struct{}

func NewDrugService() *DrugService {
	return &DrugService{}
}

func (d *DrugService) NormalizeDrugs(drugs []string) []string {

	var cleaned []string

	for _, drug := range drugs {

		drug = strings.TrimSpace(drug)
		drug = strings.ToLower(drug)

		if drug != "" {
			cleaned = append(cleaned, drug)
		}
	}

	return cleaned
}
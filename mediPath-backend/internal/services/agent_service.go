package services

type AgentService struct {
	AI *AIService
}

func NewAgentService(ai *AIService) *AgentService {
	return &AgentService{AI: ai}
}

func (a *AgentService) GeneratePatientGuidance(
	disease string,
) (string, string, error) {

	explanation, routine, err := a.AI.Explain(disease)

	if err != nil {
		return "", "", err
	}

	return explanation, routine, nil
}

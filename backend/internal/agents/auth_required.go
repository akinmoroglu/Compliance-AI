package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"compliance-checker/internal/models"
)

type authOutput struct {
	Violations []models.AgentViolation `json:"violations"`
}

func RunAuthRequiredAgent(ctx context.Context, copy string, category string) ([]models.AgentViolation, int32, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer cleanup()

	// model.ResponseSchema removed to bypass genai.TypeString loop bug

	systemPrompt, err := os.ReadFile("prompts/auth_required_agent.txt")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read auth required prompt: %w", err)
	}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	input := fmt.Sprintf("Category: %s\n\nAd copy:\n%s", category, copy)
	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		return nil, 0, err
	}

	tokens := int32(0)
	if resp.UsageMetadata != nil {
		tokens = resp.UsageMetadata.TotalTokenCount
	}

	rawStr := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	
	rawStr = strings.TrimPrefix(strings.TrimSpace(rawStr), "```json")
	rawStr = strings.TrimPrefix(strings.TrimSpace(rawStr), "```")
	rawStr = strings.TrimSuffix(strings.TrimSpace(rawStr), "```")
	rawStr = strings.TrimSpace(rawStr)

	var out authOutput
	if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
		return nil, 0, fmt.Errorf("failed to parse auth required output: %w", err)
	}
	return out.Violations, tokens, nil
}

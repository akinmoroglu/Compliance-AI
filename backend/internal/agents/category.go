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

func RunCategoryClassifier(ctx context.Context, copy string) (*models.CategoryOutput, int32, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer cleanup()

	// enforce output schema
	// model.ResponseSchema removed to bypass genai.TypeString loop bug

	systemPrompt, err := os.ReadFile("prompts/category_classifier.txt")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read category prompt: %w", err)
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf("Ad copy to classify:\n\n%s", copy)))
	if err != nil {
		return nil, 0, err
	}

	tokens := int32(0)
	if resp.UsageMetadata != nil {
		tokens = resp.UsageMetadata.TotalTokenCount
	}

	rawStr := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	
	// Strip markdown blocks if present and trim whitespace
	rawStr = strings.TrimPrefix(strings.TrimSpace(rawStr), "```json")
	rawStr = strings.TrimPrefix(strings.TrimSpace(rawStr), "```")
	rawStr = strings.TrimSuffix(strings.TrimSpace(rawStr), "```")
	rawStr = strings.TrimSpace(rawStr)
	
	var out models.CategoryOutput
	if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
		return nil, 0, fmt.Errorf("failed to parse category output: %w", err)
	}
	return &out, tokens, nil
}

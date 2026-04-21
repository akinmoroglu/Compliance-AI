package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/google/generative-ai-go/genai"
    "compliance-checker/internal/models"
)

// RunAlignmentAgent checks for contradictions between ad copy, visual creative,
// and landing page key claims. Only fires when both a creative AND a landing page exist.
//
// Critical design rule: receives raw image bytes — NOT a text description of the image.
// Sending a text summary would lose visual detail ("game of telephone" problem).
// For video input, pass the first frame as imageBytes.
func RunAlignmentAgent(
    ctx context.Context,
    imageBytes []byte,
    fullCopy string,
    lpClaims *models.LPKeyClaims,
    existingViolationCodes []string,
) ([]models.AlignmentViolation, int32, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // model.ResponseSchema removed to bypass genai.TypeString loop bug

    systemPrompt, err := os.ReadFile("prompts/alignment_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read alignment prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    // Build multipart input: image + structured text
    claimsJSON, _ := json.Marshal(lpClaims)

    textInput := fmt.Sprintf(
        "Ad copy:\n%s\n\nLanding page key claims (JSON):\n%s\n\nExisting violation codes already found (do not re-flag these):\n%s",
        fullCopy,
        string(claimsJSON),
        strings.Join(existingViolationCodes, ", "),
    )

    var parts []genai.Part
    if len(imageBytes) > 0 {
        mimeType := http.DetectContentType(imageBytes)
        parts = append(parts, genai.Blob{MIMEType: mimeType, Data: imageBytes})
    }
    parts = append(parts, genai.Text(textInput))

    resp, err := model.GenerateContent(ctx, parts...)
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

    type alignmentOutput struct {
        AlignmentViolations []models.AlignmentViolation `json:"alignment_violations"`
    }
    var out alignmentOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse alignment output: %w", err)
    }

    return out.AlignmentViolations, tokens, nil
}

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

type personalAttributeOutput struct {
    ViolationFound bool   `json:"violation_found"`
    Explanation    string `json:"explanation"`
}

// RunPersonalAttributeAgent checks for P_PERSONAL_ATTRIBUTE violations.
// imageExtractedText and audioTranscript are empty in Phase 1, populated in Phase 1b.
func RunPersonalAttributeAgent(ctx context.Context, copy string, imageExtractedText string, audioTranscript string) ([]models.AgentViolation, int32, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // Removed ResponseSchema for this specific agent. 
    // Gemini 2.5 Flash has a known bug where `genai.TypeString` in schemas 
    // forces it into infinite deterministic repetition loops for this specific system prompt.
    // The ResponseMIMEType="application/json" in client.go plus prompt formatting is sufficient.
    systemPrompt, err := os.ReadFile("prompts/personal_attribute_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read pa prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    // Build input
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("Ad copy:\n\n%s", copy))
    if imageExtractedText != "" {
        sb.WriteString(fmt.Sprintf("\n\nIMAGE EXTRACTED TEXT:\n%s", imageExtractedText))
    }
    if audioTranscript != "" {
        sb.WriteString(fmt.Sprintf("\n\nAUDIO TRANSCRIPT:\n%s", audioTranscript))
    }

    resp, err := model.GenerateContent(ctx, genai.Text(sb.String()))
    if err != nil {
        return nil, 0, err
    }

    tokens := int32(0)
    if resp.UsageMetadata != nil {
        tokens = resp.UsageMetadata.TotalTokenCount
    }

    rawStr := strings.TrimSpace(string(resp.Candidates[0].Content.Parts[0].(genai.Text)))
    rawStr = strings.TrimPrefix(rawStr, "```json")
    rawStr = strings.TrimPrefix(rawStr, "```")
    rawStr = strings.TrimSuffix(rawStr, "```")
    rawStr = strings.TrimSpace(rawStr)

    if rawStr == "" {
        return nil, tokens, nil
    }

    var out personalAttributeOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse pa output (raw: %q): %w", rawStr, err)
    }

    if out.ViolationFound {
        source := "copy"
        if strings.Contains(out.Explanation, "[Found in image text]") {
            source = "image_text"
        } else if strings.Contains(out.Explanation, "[Found in voiceover]") {
            source = "voiceover"
        }
        return []models.AgentViolation{{Code: "P_PERSONAL_ATTRIBUTE", Explanation: out.Explanation, Source: source}}, tokens, nil
    }
    return nil, tokens, nil
}

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

type violationsOutput struct {
    Violations []models.AgentViolation `json:"violations"`
}

// RunViolationsAgent checks for P_, R_, and M_ violations in ad copy.
// In Phase 1 (first call): imageExtractedText and audioTranscript are empty strings.
// In Phase 1b (second call): they contain OCR output and voiceover transcript.
// The agent prompt handles both cases — it only checks media text if provided.
func RunViolationsAgent(ctx context.Context, copy string, category string, imageExtractedText string, audioTranscript string) ([]models.AgentViolation, int32, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // model.ResponseSchema removed to bypass genai.TypeString loop bug

    systemPrompt, err := os.ReadFile("prompts/violations_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read violations prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    // Build input — conditionally include media text
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("Category: %s\n\nAd copy:\n%s", category, copy))
    if imageExtractedText != "" {
        sb.WriteString(fmt.Sprintf("\n\nIMAGE EXTRACTED TEXT (text found in the image or video frames via OCR):\n%s", imageExtractedText))
    }
    if audioTranscript != "" {
        sb.WriteString(fmt.Sprintf("\n\nAUDIO TRANSCRIPT (voiceover from the video):\n%s", audioTranscript))
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

    var out violationsOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse violations output (raw: %q): %w", rawStr, err)
    }

    // Detect source from explanation prefix set by the agent prompt
    for i := range out.Violations {
        switch {
        case strings.HasPrefix(out.Violations[i].Explanation, "[Found in image text]"):
            out.Violations[i].Source = "image_text"
        case strings.HasPrefix(out.Violations[i].Explanation, "[Found in voiceover]"):
            out.Violations[i].Source = "voiceover"
        default:
            out.Violations[i].Source = "copy"
        }
    }

    return out.Violations, tokens, nil
}

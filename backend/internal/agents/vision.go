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

// RunVisionAgent analyzes ad creative for OCR text extraction and visual policy violations.
// Accepts either raw image bytes OR a slice of video frame bytes (JPEG).
// For video: all frames are sent in a single Gemini call as sequential image parts.
// The agent prompt instructs Gemini to note the frame number when flagging video violations.
func RunVisionAgent(ctx context.Context, imageBytes []byte, videoFrames [][]byte) (*models.VisionAgentOutput, int32, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // model.ResponseSchema removed to bypass genai.TypeString loop bug

    systemPrompt, err := os.ReadFile("prompts/vision_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read vision prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    // Build multipart content parts
    var parts []genai.Part
    parts = append(parts, genai.Text("Analyze the following ad creative:"))

    if len(imageBytes) > 0 {
        mimeType := http.DetectContentType(imageBytes)
        parts = append(parts, genai.Blob{MIMEType: mimeType, Data: imageBytes})
    }

    for i, frame := range videoFrames {
        parts = append(parts, genai.Text(fmt.Sprintf("[Video frame %d of %d]", i+1, len(videoFrames))))
        parts = append(parts, genai.Blob{MIMEType: "image/jpeg", Data: frame})
    }

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

    var out models.VisionAgentOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse vision output: %w", err)
    }

    // Set source field on all visual violations
    for i := range out.VisualViolations {
        out.VisualViolations[i].Source = "image_text"
    }

    return &out, tokens, nil
}

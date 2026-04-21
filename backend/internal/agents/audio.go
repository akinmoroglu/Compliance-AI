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

// RunAudioAgent transcribes the audio track of a video ad.
// Transcription only — no violation checking. Violation checking on the
// transcript happens in Phase 1b via RunViolationsAgent and RunPersonalAttributeAgent.
func RunAudioAgent(ctx context.Context, audioBytes []byte) (*models.AudioAgentOutput, int32, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // model.ResponseSchema removed to bypass genai.TypeString loop bug

    systemPrompt, err := os.ReadFile("prompts/audio_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read audio prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    parts := []genai.Part{
        genai.Blob{MIMEType: "audio/mpeg", Data: audioBytes},
        genai.Text("Transcribe this audio track as instructed in the system prompt."),
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

    var out models.AudioAgentOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse audio output: %w", err)
    }

    return &out, tokens, nil
}

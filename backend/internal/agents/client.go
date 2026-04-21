package agents

import (
    "context"
    "os"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

// newModel creates a Gemini generative model with temperature=0 and JSON output enforced.
// The model name is read from GEMINI_MODEL env var, defaulting to "gemini-2.0-flash".
// All agents use this helper. To upgrade a single agent to a more powerful model,
// temporarily override GEMINI_MODEL — no code change required.
func newModel(ctx context.Context) (*genai.GenerativeModel, func(), error) {
    client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
    if err != nil {
        return nil, nil, err
    }
    modelName := os.Getenv("GEMINI_MODEL")
    if modelName == "" {
        modelName = "gemini-2.0-flash"
    }
    model := client.GenerativeModel(modelName)
    // Raise temperature to 0.4 to thoroughly break deterministic repetition loops.
    // Cap max tokens to 8192 (Gemini maximum) to allow long JSON generations but still prevent infinite loops.
    temp := float32(0.4)
    model.SetTemperature(temp)
    model.SetMaxOutputTokens(8192)
    model.ResponseMIMEType = "application/json"
    cleanup := func() { client.Close() }
    return model, cleanup, nil
}

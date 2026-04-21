package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"
    "strings"

    firecrawl "github.com/mendableai/firecrawl-go"
    "github.com/google/generative-ai-go/genai"
    "compliance-checker/internal/models"
)

// RunLPAgent fetches the landing page via Firecrawl and runs Gemini analysis
// for LP-level violations and structured key claim extraction.
//
// Flexibility note: Firecrawl can be replaced with a Playwright-based fetcher
// by swapping fetchLP() without changing the agent's analysis logic. The
// interface contract is: fetchLP() returns (markdownContent string, jsRendered bool, err error).
func RunLPAgent(ctx context.Context, url string) (*models.LPAgentOutput, int32, error) {
    // Step 1: Fetch landing page via Firecrawl with 8-second timeout
    fetchCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
    defer cancel()

    lpContent, jsRendered, fetchErr := fetchLandingPage(fetchCtx, url)
    if fetchErr != nil {
        return &models.LPAgentOutput{LPAnalysisStatus: "fetch_failed"}, 0, nil
    }

    lpStatus := "analyzed"
    if jsRendered {
        lpStatus = "js_rendered"
    }

    // Truncate: first 4000 chars (body) + last 1000 chars (footer disclaimers)
    lpExcerpt := buildLPExcerpt(lpContent)

    // Step 2: Gemini analysis
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, 0, err
    }
    defer cleanup()

    // model.ResponseSchema removed to bypass genai.TypeString loop bug

    systemPrompt, err := os.ReadFile("prompts/lp_agent.txt")
    if err != nil {
        return nil, 0, fmt.Errorf("failed to read lp prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    input := fmt.Sprintf("Landing page URL: %s\n\nLanding page content:\n%s", url, lpExcerpt)
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

    var out models.LPAgentOutput
    if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
        return nil, 0, fmt.Errorf("failed to parse lp output: %w", err)
    }

    out.LPAnalysisStatus = lpStatus
    out.LPTextExcerpt = lpExcerpt

    // Tag all LP violations with source = "landing_page"
    for i := range out.Violations {
        out.Violations[i].Source = "landing_page"
    }

    return &out, tokens, nil
}

// fetchLandingPage uses Firecrawl to retrieve and render the landing page.
// Returns: markdown content, whether JS rendering was required, error.
func fetchLandingPage(ctx context.Context, url string) (string, bool, error) {
    apiKey := os.Getenv("FIRECRAWL_API_KEY")
    if apiKey == "" {
        return "", false, fmt.Errorf("FIRECRAWL_API_KEY not set")
    }

    app, err := firecrawl.NewFirecrawlApp(apiKey, "")
    if err != nil {
        return "", false, err
    }

    params := &firecrawl.ScrapeParams{
        Formats: []string{"markdown"},
    }

    // Run in goroutine to respect ctx cancellation
    type result struct {
        content    string
        jsRendered bool
        err        error
    }
    ch := make(chan result, 1)

    go func() {
        scrapeResult, err := app.ScrapeURL(url, params)
        if err != nil {
            ch <- result{err: err}
            return
        }
        content := ""
        if scrapeResult.Markdown != "" {
            content = scrapeResult.Markdown
        }
        // Firecrawl uses headless Chromium — treat as JS-rendered if content
        // is shorter than expected for the page (heuristic: under 500 chars is suspicious)
        jsRendered := len(content) < 500
        ch <- result{content: content, jsRendered: jsRendered}
    }()

    select {
    case <-ctx.Done():
        return "", false, fmt.Errorf("lp fetch timeout: %w", ctx.Err())
    case r := <-ch:
        return r.content, r.jsRendered, r.err
    }
}

// buildLPExcerpt extracts first 4000 chars (body) + last 1000 chars (footer).
// The footer section typically contains disclaimers critical to compliance checking.
func buildLPExcerpt(content string) string {
    const maxBody = 4000
    const footerLen = 1000

    if len(content) <= maxBody+footerLen {
        return content
    }

    body := content[:maxBody]
    footer := content[len(content)-footerLen:]
    return body + "\n\n[...content truncated...]\n\n" + footer
}

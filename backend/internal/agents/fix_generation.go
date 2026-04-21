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

type fixBatchOutput struct {
	Fixes []models.FixOutput `json:"fixes"`
}

var fixFrameworks = map[string]string{
	"P_HATE_SPEECH":         "No fix possible. Explain the prohibition.",
	"P_ILLEGAL_DRUGS":       "No fix possible. Explain the prohibition.",
	"P_IMPERSONATION":       "No fix possible. Remove unauthorized use of brand or celebrity identity.",
	"P_PERSONAL_ATTRIBUTE":  "Rewrite from 'you have this condition' framing to 'this product supports this function' framing. Example: 'Are you struggling with debt?' → 'Designed to help you reach financial freedom.'",
	"P_BULLYING_HARASSMENT": "No fix possible. Remove targeting of specific individuals or groups.",
	"A_GAMBLING":            "This category requires Meta written permission. No creative fix available — the advertiser must obtain authorization before running this ad.",
	"A_CRYPTO_EXCHANGE":     "This category requires Meta written permission and regulatory licensing. No creative fix available — the advertiser must obtain authorization.",
	"R_WEIGHT_LOSS_CLAIM":   "Remove specific numeric claims and timeframes. Replace with aspirational non-specific language. Examples: 'Support your weight management goals' / 'Designed to complement a healthy lifestyle.' Cannot include specific amounts, timeframes, or guarantees.",
	"R_BODY_SHAMING":        "Reframe from problem-focused to aspiration-focused. Remove language implying the reader's current body is shameful. Replace with: 'Feel confident' / 'Support your goals' / 'Designed for you.'",
	"R_HEALTH_CLAIM":        "Remove the disease or condition name from the claim. Replace with general wellness language. 'Treats arthritis' → 'Supports healthy joints.' Add disclaimer: 'These statements have not been evaluated by the FDA.'",
	"R_GUARANTEED_OUTCOME":  "Remove guarantee language. Replace with directional benefit language. 'Guaranteed 30% returns' → 'Designed for consistent growth.' Add risk disclaimer for financial products.",
	"R_SCAREWARE":           "Remove fake system UI elements. Replace with branded product imagery and direct benefit messaging.",
	"M_URGENCY":             "Add the actual offer end date if there is one. If the urgency is fabricated, remove it or make it verifiable.",
	"M_NARRATIVE":           "Remove fabricated social proof screenshots. Replace with verifiable testimonials with proper attribution.",
	"M_HYGIENE":             "Correct typos and grammar. Replace low-resolution or watermarked images. Remove all-caps and symbol spam.",
	"M_SENSATIONALISM":      "Replace sensationalist headlines with direct, factual product claims. Remove fake play button overlays.",
	"M_DISCLOSURE":          "Add required disclaimer for the ad category. Supplements: 'These statements have not been evaluated by the FDA. This product is not intended to diagnose, treat, cure, or prevent any disease.' Health/fitness: 'Individual results may vary.' Financial: include APR and link to full terms.",
	"_suggested_copy":       "Rewrite the entire primary text of the ad to remove all violations listed. Preserve the core product message and call-to-action. Output only the rewritten copy — no explanations.",
}

func RunFixGenerationAgent(ctx context.Context, originalCopy string, category string, violations []models.AgentViolation) (map[string]string, int32, error) {
	if len(violations) == 0 {
		return map[string]string{}, 0, nil
	}

	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer cleanup()

	// model.ResponseSchema removed to bypass genai.TypeString loop bug

	systemPrompt, err := os.ReadFile("prompts/fix_generation_agent.txt")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read fix generation prompt: %w", err)
	}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	// Build the batched input
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Category: %s\n", category))
	sb.WriteString(fmt.Sprintf("Original copy: %s\n\n", originalCopy))
	sb.WriteString("All violations found: ")
	codes := make([]string, len(violations))
	for i, v := range violations {
		codes[i] = v.Code
	}
	sb.WriteString(strings.Join(codes, ", "))
	sb.WriteString("\n\nViolations to fix:\n")
	for _, v := range violations {
		framework := fixFrameworks[v.Code]
		if framework == "" {
			framework = "Rewrite the copy to remove the violation."
		}
		sb.WriteString(fmt.Sprintf("- Code: %s\n  Fix framework: %s\n\n", v.Code, framework))
	}
	sb.WriteString("\n\nAlso output a '_suggested_copy' entry: a complete rewrite of the original primary text that removes ALL violations above. Keep the core message and CTA. Code: '_suggested_copy'.")

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
		return map[string]string{}, tokens, nil
	}

	var out fixBatchOutput
	if err := json.Unmarshal([]byte(rawStr), &out); err != nil {
		return nil, 0, fmt.Errorf("failed to parse fix output (raw: %q): %w", rawStr, err)
	}

	result := make(map[string]string)
	for _, f := range out.Fixes {
		result[f.Code] = f.SuggestedFix
	}
	return result, tokens, nil
}

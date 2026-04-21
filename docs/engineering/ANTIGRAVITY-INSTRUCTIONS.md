# Antigravity Implementation Instructions
# Ad Compliance Checker — Local MVP

---

## What You Are Building

You are building a Go backend API for an AI-powered Meta Ad Compliance Checker. The Vue 3 frontend already exists and is functional. You will:

1. Build the Go backend service from scratch
2. Make minimal, targeted changes to the existing Vue frontend to wire it to the real API
3. Do NOT redesign or restructure the frontend — only change the data layer

---

## What Already Exists (Do Not Touch Unless Instructed)

The frontend project is at the root of this directory. Key files:

- `src/components/ComplianceDashboard.vue` — the main UI component (4-step accordion)
- `src/App.vue` — root layout with header
- `src/main.ts`, `tailwind.config.js`, `vite.config.ts` — config files, do not modify

**What the frontend does today:**
- Step 1: User selects platform (Meta), region, age range (min/max)
- Step 2: User enters ad copy (primaryText, headline, description) and optional landing page URL
- Step 3: Shows a loading spinner (currently uses `setTimeout` to simulate)
- Step 4: Shows a hardcoded mock compliance result

**Your job in the frontend:** Replace the `setTimeout` simulation with a real API call, and replace the hardcoded Step 4 result with dynamic data from the backend.

---

## Backend: Project Structure to Create

Create a `backend/` directory at the root of the project with this exact structure:

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handler.go
│   ├── agents/
│   │   ├── category.go
│   │   ├── violations.go
│   │   ├── personal_attribute.go
│   │   ├── auth_required.go
│   │   └── fix_generation.go
│   ├── scoring/
│   │   └── engine.go
│   ├── db/
│   │   ├── postgres.go
│   │   └── checks.go
│   └── models/
│       └── types.go
├── prompts/
│   ├── category_classifier.txt
│   ├── violations_agent.txt
│   ├── personal_attribute_agent.txt
│   ├── auth_required_agent.txt
│   └── fix_generation_agent.txt
├── docker-compose.yml
├── .env.example
├── .gitignore
└── go.mod
```

---

## Environment Setup

**`backend/docker-compose.yml`**
```yaml
version: '3.9'
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: compliance_checker
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data

volumes:
  pg_data:
```

**`backend/.env.example`**
```
GEMINI_API_KEY=your_google_ai_studio_key_here
DATABASE_URL=postgres://dev:dev@localhost:5432/compliance_checker?sslmode=disable
PORT=8080
```

**`backend/.gitignore`**
```
.env
```

**`backend/go.mod`** — run these commands to initialize:
```bash
cd backend
go mod init compliance-checker
go get github.com/google/generative-ai-go/genai@latest
go get google.golang.org/api/option
go get github.com/lib/pq
go get github.com/google/uuid
```

---

## Database Schema

**`backend/internal/db/postgres.go`**

On startup, run this migration automatically:

```go
package db

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
	migrate()
	log.Println("database connected")
}

func migrate() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS checks (
			id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			status      TEXT NOT NULL DEFAULT 'complete',
			input_copy  TEXT NOT NULL,
			result      JSONB,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
```

---

## Shared Types

**`backend/internal/models/types.go`**

```go
package models

type CheckRequest struct {
	Platform      string   `json:"platform"`
	Region        string   `json:"region"`
	AgeMin        int      `json:"age_min"`
	AdFormat      string   `json:"ad_format"`
	PrimaryText   string   `json:"primary_text"`
	Headline      string   `json:"headline"`
	Description   string   `json:"description"`
	LandingPageURL string  `json:"landing_page_url"`
}

type Violation struct {
	Code          string `json:"code"`
	Title         string `json:"title"`
	Severity      string `json:"severity"`
	Explanation   string `json:"explanation"`
	SuggestedFix  string `json:"suggested_fix"`
}

type CheckResult struct {
	CheckID         string      `json:"check_id"`
	Action          string      `json:"action"`           // "not_compliant" | "compliant"
	ComplianceScore int         `json:"compliance_score"` // 0–100
	RiskCategory    string      `json:"risk_category"`    // "High Risk" | "Medium Risk" | "Low Risk"
	Category        string      `json:"category"`
	CategoryState   string      `json:"category_state"`
	Violations      []Violation `json:"violations"`
}

// Internal agent output types
type AgentViolation struct {
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
}

type CategoryOutput struct {
	Category      string `json:"category"`
	CategoryState string `json:"category_state"`
	Reasoning     string `json:"reasoning"`
}

type PersonalAttributeOutput struct {
	ViolationFound bool   `json:"violation_found"`
	Explanation    string `json:"explanation"`
}

type FixOutput struct {
	Code         string `json:"code"`
	SuggestedFix string `json:"suggested_fix"`
}
```

---

## Scoring Engine

**`backend/internal/scoring/engine.go`**

This is pure Go logic — no LLM involved. Implement exactly as specified:

```go
package scoring

import "strings"

type Result struct {
	ComplianceScore int
	RiskCategory    string
	Action          string
}

func Score(violations []string, categoryState string, authorized bool) Result {
	var p, a, r, m []string
	for _, v := range violations {
		switch {
		case strings.HasPrefix(v, "P_"):
			p = append(p, v)
		case strings.HasPrefix(v, "A_"):
			a = append(a, v)
		case strings.HasPrefix(v, "R_"):
			r = append(r, v)
		case strings.HasPrefix(v, "M_"):
			m = append(m, v)
		}
	}

	// Rule 1: Any prohibited violation
	if len(p) > 0 {
		return Result{0, "High Risk", "not_compliant"}
	}
	// Rule 2: Unauthorized authorization-required category
	if len(a) > 0 && !authorized {
		return Result{0, "High Risk", "not_compliant"}
	}
	// Rule 3: Restricted violations
	if len(r) > 0 {
		score := 50 - (10 * (len(r) - 1)) - (5 * len(m))
		if score < 1 {
			score = 1
		}
		risk := "Medium Risk"
		if score < 30 {
			risk = "High Risk"
		}
		return Result{score, risk, "compliant"}
	}
	// Rule 4: Quality/moderate violations only
	if len(m) > 0 {
		score := 75 - (5 * (len(m) - 1))
		if score < 45 {
			score = 45
		}
		risk := "Low Risk"
		if score < 65 {
			risk = "Medium Risk"
		}
		return Result{score, risk, "compliant"}
	}
	// Rule 5: Clean — standard category
	if categoryState == "ALLOWED" {
		return Result{95, "Low Risk", "compliant"}
	}
	// Rule 6: Clean — sensitive category
	return Result{78, "Low Risk", "compliant"}
}
```

Write unit tests for this function in `backend/internal/scoring/engine_test.go`. Test all 6 rules.

---

## Gemini Agent Setup

**Shared Gemini client helper — create in `backend/internal/agents/client.go`:**

```go
package agents

import (
	"context"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func newModel(ctx context.Context) (*genai.GenerativeModel, func(), error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, nil, err
	}
	model := client.GenerativeModel("gemini-1.5-flash")
	temp := float32(0)
	model.SetTemperature(temp)
	model.ResponseMIMEType = "application/json"
	cleanup := func() { client.Close() }
	return model, cleanup, nil
}
```

---

## Agent 1 — Category Classifier

**`backend/internal/agents/category.go`**

Reads system prompt from `prompts/category_classifier.txt`.  
Returns: category name, category_state (PROHIBITED / AUTH_REQUIRED / RESTRICTED / ALLOWED), reasoning.

```go
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"compliance-checker/internal/models"
)

func RunCategoryClassifier(ctx context.Context, copy string) (*models.CategoryOutput, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	// enforce output schema
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"category":       {Type: genai.TypeString},
			"category_state": {Type: genai.TypeString},
			"reasoning":      {Type: genai.TypeString},
		},
		Required: []string{"category", "category_state", "reasoning"},
	}

	systemPrompt, err := os.ReadFile("prompts/category_classifier.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read category prompt: %w", err)
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf("Ad copy to classify:\n\n%s", copy)))
	if err != nil {
		return nil, err
	}

	raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
	var out models.CategoryOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse category output: %w", err)
	}
	return &out, nil
}
```

**`backend/prompts/category_classifier.txt`** — write this content exactly:

```
You are a Meta ad category classifier. Your only job is to assign exactly one category to an ad based on its copy.

OUTPUT FORMAT: JSON with fields: category (string), category_state (string), reasoning (string).

category_state must be exactly one of: PROHIBITED, AUTH_REQUIRED, RESTRICTED, ALLOWED

PROHIBITED CATEGORIES (set category_state = PROHIBITED):
- Weapons & Explosives Sale
- Illegal Drugs Sale or Promotion (NOT CBD — that is AUTH_REQUIRED)
- Counterfeit Goods
- Tobacco & Nicotine Products
- Deepfake & AI Identity Tools
- Covert Surveillance / Stalkerware
- Adult Explicit Content / Escort Services
- Human Exploitation / Trafficking
- Payday Loans / Short-term Loans 90 days or less
- Binary Options / CFD Trading
- ICO (Initial Coin Offerings)

AUTH_REQUIRED CATEGORIES (set category_state = AUTH_REQUIRED):
- Online Gambling & Betting
- Dating Platforms
- Cryptocurrency Exchange / Trading / Lending
- CBD Products
- Prescription Drugs
- Drug & Alcohol Addiction Treatment

RESTRICTED CATEGORIES (set category_state = RESTRICTED):
- Health — Weight Loss (Mild)
- Health — Medical Aesthetics & Cosmetic Procedures
- Health — General Wellness & Supplements
- Health — Mental Health & Therapy
- Financial — Investment & Trading (non-crypto)
- Financial — Insurance
- Financial — Long-term Loans (over 90 days)
- Financial — Credit Cards
- Health — OTC Drugs
- Reproductive Health & Sexual Wellness
- Software — Security & Antivirus
- E-commerce — General (if ad makes health or financial claims)
- Alcohol (beer, wine, spirits — allowed with 18+ targeting)

ALLOWED CATEGORIES (set category_state = ALLOWED):
- E-commerce — Apparel & Fashion
- E-commerce — Beauty & Skincare (no medical claims)
- E-commerce — Luxury Goods (authorized resellers)
- E-commerce — Electronics & Tech
- E-commerce — Home & Garden
- Software — Games & Entertainment
- Software — Productivity & Utilities
- Food & Beverage — General (no health claims)
- Education & Online Courses
- Travel & Hospitality
- Events & Entertainment
- Non-profit & Charity
- B2B Services
- Real Estate
- Automotive
- Sports & Fitness (no specific health claims)
- Any category not listed above that does not involve health, finance, or regulated goods

INSTRUCTIONS:
1. Read the ad copy carefully.
2. Assign the single best-matching category.
3. Set category_state based on the category tier above.
4. Write a one-sentence reasoning explaining your choice.
5. When in doubt between two categories, pick the more restrictive one.
```

---

## Agent 2 — Violations Agent

**`backend/internal/agents/violations.go`**

Checks for P_, R_, and M_ violations. Returns array of violations found.

```go
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"compliance-checker/internal/models"
)

type violationsOutput struct {
	Violations []models.AgentViolation `json:"violations"`
}

func RunViolationsAgent(ctx context.Context, copy string, category string) ([]models.AgentViolation, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"violations": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"code":        {Type: genai.TypeString},
						"explanation": {Type: genai.TypeString},
					},
					Required: []string{"code", "explanation"},
				},
			},
		},
		Required: []string{"violations"},
	}

	systemPrompt, err := os.ReadFile("prompts/violations_agent.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read violations prompt: %w", err)
	}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	input := fmt.Sprintf("Category: %s\n\nAd copy:\n%s", category, copy)
	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		return nil, err
	}

	raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
	var out violationsOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse violations output: %w", err)
	}
	return out.Violations, nil
}
```

**`backend/prompts/violations_agent.txt`** — write this content exactly:

```
You are a Meta advertising policy compliance agent. Your job is to detect policy violations in ad copy.

You will receive: the ad category and the ad copy.
You output: a JSON object with a "violations" array. If no violations, return {"violations": []}.

VIOLATION CODES YOU MUST CHECK (output the exact code string):

P_ PROHIBITED VIOLATIONS (any of these = immediate serious violation):
- P_HATE_SPEECH: Ad attacks people based on race, ethnicity, religion, gender, sexual orientation, disability, or national origin. Slurs, dehumanizing language, calls for discrimination.
- P_ILLEGAL_DRUGS: Ad promotes sale, purchase, or use of controlled substances (cocaine, heroin, MDMA, etc.). NOT CBD (that is an A_ code checked separately).
- P_IMPERSONATION: Ad falsely uses a real brand's name, logo, or identity to deceive. "Official Nike sale" from a non-Nike seller. Celebrity endorsement not authorized by that celebrity.
- P_PERSONAL_ATTRIBUTE: Ad copy directly asserts or implies that the reader has a personal attribute. Protected attributes: race, religion, age, sexual orientation, disability, named health conditions, financial vulnerability, criminal record. THREE violation patterns: (1) Direct assertion: "As a diabetic..." (2) Rhetorical question implying condition: "Are you struggling with debt?" / "Tired of dealing with anxiety?" / "Is your credit score holding you back?" (3) Audience address implying attribute: "For people dealing with depression..." / "For anyone struggling with addiction..." NOT a violation: "For fitness enthusiasts" / "For small business owners" / "Feel your best" / "Support healthy joints."
- P_BULLYING_HARASSMENT: Ad targets and mocks a specific real person or group in a harassing or demeaning way.

R_ RESTRICTED VIOLATIONS (fixable — provide specific explanation):
- R_WEIGHT_LOSS_CLAIM: Ad makes a specific numeric weight loss claim with amount and timeframe ("Lose 20kg in 30 days") or uses absolute guarantee language for weight loss. NOT a violation: "Support your weight goals" / aspirational language without specific numbers.
- R_BODY_SHAMING: Ad uses language implying the reader's current body or appearance is shameful, wrong, or a problem. "Tired of your belly fat?" / "Embarrassed by your skin?" / "Do you hate how you look?" NOT a violation: "Feel your best" / positive aspiration.
- R_HEALTH_CLAIM: Non-prescription product claims to treat, cure, heal, prevent, or diagnose a specific named medical condition or disease. "Cures arthritis" / "Prevents Alzheimer's" / "Reverses hair loss." NOT a violation: "Supports joint health" / "May help with energy levels."
- R_GUARANTEED_OUTCOME: Ad guarantees a specific financial return or outcome. "Guaranteed 30% returns" / "Make $5000/month guaranteed" / "Double your investment." NOT a violation: general benefit language.
- R_SCAREWARE: Ad uses fake system UI, fake OS alerts, or fake virus warnings to create fear. "Your phone has 3 viruses!" shown as a fake notification.

M_ MODERATE VIOLATIONS (reduce score — provide explanation):
- M_URGENCY: Ad uses countdown timers, "today only," "only X left in stock," or unverifiable limited-time pressure with no specified end date. Standard tactic — flag but score leniently.
- M_NARRATIVE: Ad uses fabricated text message screenshots, fake DM screenshots, fake comment screenshots, or "store closing after X years" narratives.
- M_HYGIENE: Ad has 3 or more typos, uses low-resolution or watermarked images (describe if text mentions visual), excessive all-caps or emoji spam, or generally unprofessional presentation signaling scam risk.
- M_SENSATIONALISM: Ad uses "DOCTORS HATE HIM," "You won't believe," "This one weird trick," fake play button on a static image, or headlines that don't match the actual content.
- M_DISCLOSURE: Ad is missing required disclaimers for its category. Supplement ad with no FDA disclaimer. Health/fitness results with no "results may vary." Financial offer with no APR or terms.

IMPORTANT RULES:
1. Only flag violations you are confident about. Do not flag borderline cases as violations.
2. For P_PERSONAL_ATTRIBUTE, apply the three-pattern test strictly. Aspirational language ("Feel your best", "Support healthy joints") is NEVER a violation.
3. Output the exact violation code string — no variations.
4. For each violation, write a specific explanation quoting the exact problematic copy.
5. If no violations found, return {"violations": []}.
```

---

## Agent 3 — Personal Attribute Agent

**`backend/internal/agents/personal_attribute.go`**

Always runs on every ad. Focused single-purpose prompt.

```go
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"compliance-checker/internal/models"
)

func RunPersonalAttributeAgent(ctx context.Context, copy string) ([]models.AgentViolation, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"violation_found": {Type: genai.TypeBoolean},
			"explanation":     {Type: genai.TypeString},
		},
		Required: []string{"violation_found", "explanation"},
	}

	systemPrompt, err := os.ReadFile("prompts/personal_attribute_agent.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read personal attribute prompt: %w", err)
	}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf("Ad copy:\n\n%s", copy)))
	if err != nil {
		return nil, err
	}

	raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
	var out models.PersonalAttributeOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse personal attribute output: %w", err)
	}

	if out.ViolationFound {
		return []models.AgentViolation{{Code: "P_PERSONAL_ATTRIBUTE", Explanation: out.Explanation}}, nil
	}
	return nil, nil
}
```

**`backend/prompts/personal_attribute_agent.txt`**:

```
You are a specialist in Meta's Personal Attribute Assertion policy (Section 2.12).

Meta prohibits ad copy that asserts or implies that the person viewing the ad has a specific personal attribute.

Protected attributes: race, ethnicity, national origin, religion, age, sexual orientation, gender identity, disability, named physical or mental health conditions, vulnerable financial status, criminal record.

THREE VIOLATION PATTERNS:
1. DIRECT ASSERTION: "As a diabetic, you know..." / "As a Christian..." / "As someone with ADHD..."
2. RHETORICAL QUESTION IMPLYING CONDITION: "Are you struggling with debt?" / "Tired of dealing with anxiety?" / "Is your credit score holding you back?" / "Are you over 50 and worried about joints?"
3. AUDIENCE ADDRESS IMPLYING ATTRIBUTE: "For people dealing with diabetes..." / "For anyone struggling with addiction..." / "For those battling depression..."

THE TEST: Does this copy ASSUME the reader has a specific personal attribute?
- If YES → violation
- If the copy describes a product benefit or uses aspirational language WITHOUT implying the reader currently has a problem → NOT a violation

NOT VIOLATIONS:
- "For fitness enthusiasts" (interest, not a condition)
- "For small business owners" (professional role)
- "Feel your best every day" (aspirational)
- "Supports healthy blood sugar" (product function, no assumption about reader)
- "Support your energy levels" (general wellness)
- "For people interested in healthy living" (interest-based)

Output JSON: {"violation_found": true/false, "explanation": "quote the specific copy that violates, or empty string if no violation"}
```

---

## Agent 4 — Auth Required Agent

**`backend/internal/agents/auth_required.go`**

```go
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"compliance-checker/internal/models"
)

type authOutput struct {
	Violations []models.AgentViolation `json:"violations"`
}

func RunAuthRequiredAgent(ctx context.Context, copy string, category string) ([]models.AgentViolation, error) {
	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"violations": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"code":        {Type: genai.TypeString},
						"explanation": {Type: genai.TypeString},
					},
					Required: []string{"code", "explanation"},
				},
			},
		},
		Required: []string{"violations"},
	}

	systemPrompt, err := os.ReadFile("prompts/auth_required_agent.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read auth required prompt: %w", err)
	}
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(string(systemPrompt))},
	}

	input := fmt.Sprintf("Category: %s\n\nAd copy:\n%s", category, copy)
	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		return nil, err
	}

	raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
	var out authOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse auth required output: %w", err)
	}
	return out.Violations, nil
}
```

**`backend/prompts/auth_required_agent.txt`**:

```
You are a Meta authorization-required category specialist. You check if an ad promotes services that require explicit Meta written permission.

You will receive the ad category and copy. You output a JSON violations array. If no violations, return {"violations": []}.

AUTHORIZATION-REQUIRED VIOLATION CODES:

A_GAMBLING: Ad promotes online gambling, sports betting, casino games, or poker platforms without indicating they hold Meta written permission. Keywords: "bet," "odds," "casino," "poker," "sports betting," "wager."

A_CRYPTO_EXCHANGE: Ad promotes a cryptocurrency exchange, trading platform, crypto lending, or crypto investment service without indicating they hold Meta written permission and appropriate regulatory licenses. Keywords: "trade crypto," "buy Bitcoin," "crypto exchange," "DeFi," "earn crypto." NOTE: General crypto education content (no buy/sell/trade CTA) is NOT a violation.

If the category passed to you is not related to gambling or crypto, return {"violations": []}.

For each violation, explain specifically what in the copy triggers the flag.
```

---

## Agent 5 — Fix Generation Agent

**`backend/internal/agents/fix_generation.go`**

Single batched call — all violations at once. Do NOT make one call per violation.

```go
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
}

func RunFixGenerationAgent(ctx context.Context, originalCopy string, category string, violations []models.AgentViolation) (map[string]string, error) {
	if len(violations) == 0 {
		return map[string]string{}, nil
	}

	model, cleanup, err := newModel(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"fixes": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"code":          {Type: genai.TypeString},
						"suggested_fix": {Type: genai.TypeString},
					},
					Required: []string{"code", "suggested_fix"},
				},
			},
		},
		Required: []string{"fixes"},
	}

	systemPrompt, err := os.ReadFile("prompts/fix_generation_agent.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read fix generation prompt: %w", err)
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

	resp, err := model.GenerateContent(ctx, genai.Text(sb.String()))
	if err != nil {
		return nil, err
	}

	raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
	var out fixBatchOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("failed to parse fix output: %w", err)
	}

	result := make(map[string]string)
	for _, f := range out.Fixes {
		result[f.Code] = f.SuggestedFix
	}
	return result, nil
}
```

**`backend/prompts/fix_generation_agent.txt`**:

```
You are a Meta ad copy rewriter. You receive an original ad copy, its category, all violations found, and fix frameworks for each violation.

Your job: generate one specific, actionable rewrite suggestion per violation code.

Rules:
1. Each fix must directly address the specific violation in the original copy.
2. The fix must not introduce new violations — you can see all violations found, so avoid rewriting one thing in a way that creates another problem.
3. For P_ violations (Prohibited): explain the prohibition clearly — do not suggest a creative rewrite.
4. For A_ violations (Authorization Required): explain the authorization requirement — do not suggest a creative rewrite.
5. For R_ and M_ violations: provide a specific rewrite of the problematic copy. Quote the original and suggest the replacement.
6. Keep fixes concise and actionable — one or two sentences maximum.
7. Do not explain Meta policy in depth — just give the fix.

Output JSON: {"fixes": [{"code": "VIOLATION_CODE", "suggested_fix": "your specific rewrite suggestion here"}]}
```

---

## Violation Metadata Map

Create this in **`backend/internal/models/metadata.go`** — used to enrich agent output with human-readable titles and severity labels:

```go
package models

type ViolationMeta struct {
	Title    string
	Severity string
}

var ViolationMetadata = map[string]ViolationMeta{
	"P_HATE_SPEECH":         {"Hate Speech", "PROHIBITED"},
	"P_ILLEGAL_DRUGS":       {"Illegal Drug Promotion", "PROHIBITED"},
	"P_IMPERSONATION":       {"Brand or Identity Impersonation", "PROHIBITED"},
	"P_PERSONAL_ATTRIBUTE":  {"Personal Attribute Assertion", "PROHIBITED"},
	"P_BULLYING_HARASSMENT": {"Bullying or Harassment", "PROHIBITED"},
	"A_GAMBLING":            {"Unauthorized Gambling Promotion", "AUTHORIZATION_REQUIRED"},
	"A_CRYPTO_EXCHANGE":     {"Unauthorized Crypto Exchange Promotion", "AUTHORIZATION_REQUIRED"},
	"R_WEIGHT_LOSS_CLAIM":   {"Specific Weight Loss Claim", "RESTRICTED"},
	"R_BODY_SHAMING":        {"Negative Body Image or Body Shaming", "RESTRICTED"},
	"R_HEALTH_CLAIM":        {"Unproven Health or Medical Claim", "RESTRICTED"},
	"R_GUARANTEED_OUTCOME":  {"Guaranteed Financial Outcome", "RESTRICTED"},
	"R_SCAREWARE":           {"Scareware or Fake System Alert", "RESTRICTED"},
	"M_URGENCY":             {"Urgency Tactic", "MODERATE"},
	"M_NARRATIVE":           {"Fabricated Social Proof", "MODERATE"},
	"M_HYGIENE":             {"Low Production Quality", "MODERATE"},
	"M_SENSATIONALISM":      {"Clickbait or Sensationalism", "MODERATE"},
	"M_DISCLOSURE":          {"Missing Required Disclaimer", "MODERATE"},
}
```

---

## Main Handler

**`backend/internal/api/handler.go`**

This is the orchestrator. It runs all agents, collects results, scores, and returns the full result.

```go
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"compliance-checker/internal/agents"
	"compliance-checker/internal/models"
	"compliance-checker/internal/scoring"
)

func HandleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req models.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Build combined copy string from all copy fields
	fullCopy := req.PrimaryText
	if req.Headline != "" {
		fullCopy += "\nHeadline: " + req.Headline
	}
	if req.Description != "" {
		fullCopy += "\nDescription: " + req.Description
	}

	ctx := r.Context()

	// Phase 1: Category Classification
	categoryResult, err := agents.RunCategoryClassifier(ctx, fullCopy)
	if err != nil {
		log.Printf("category classifier error: %v", err)
		http.Error(w, `{"error":"classification failed"}`, http.StatusInternalServerError)
		return
	}

	// Fast exit for PROHIBITED categories
	var allAgentViolations []models.AgentViolation
	if categoryResult.CategoryState == "PROHIBITED" {
		allAgentViolations = append(allAgentViolations, models.AgentViolation{
			Code:        "P_PROHIBITED_CATEGORY",
			Explanation: "This ad category is prohibited on Meta: " + categoryResult.Category,
		})
	} else {
		// Phase 2: Parallel agents — fault-tolerant (partial result > total failure)
		var mu sync.Mutex
		var wg sync.WaitGroup

		runAgent := func(fn func() ([]models.AgentViolation, error), name string) {
			defer wg.Done()
			v, err := fn()
			if err != nil {
				log.Printf("agent %s failed: %v", name, err)
				return
			}
			mu.Lock()
			allAgentViolations = append(allAgentViolations, v...)
			mu.Unlock()
		}

		wg.Add(3)
		go runAgent(func() ([]models.AgentViolation, error) {
			return agents.RunViolationsAgent(ctx, fullCopy, categoryResult.Category)
		}, "violations")
		go runAgent(func() ([]models.AgentViolation, error) {
			return agents.RunPersonalAttributeAgent(ctx, fullCopy)
		}, "personal_attribute")
		go runAgent(func() ([]models.AgentViolation, error) {
			return agents.RunAuthRequiredAgent(ctx, fullCopy, categoryResult.Category)
		}, "auth_required")
		wg.Wait()
	}

	// Deduplicate violations by code
	seen := map[string]bool{}
	var uniqueViolations []models.AgentViolation
	for _, v := range allAgentViolations {
		if !seen[v.Code] {
			seen[v.Code] = true
			uniqueViolations = append(uniqueViolations, v)
		}
	}

	// Phase 3: Scoring Engine
	violationCodes := make([]string, len(uniqueViolations))
	for i, v := range uniqueViolations {
		violationCodes[i] = v.Code
	}
	scoreResult := scoring.Score(violationCodes, categoryResult.CategoryState, false)

	// Fix Generation — single batched call
	fixes, err := agents.RunFixGenerationAgent(ctx, fullCopy, categoryResult.Category, uniqueViolations)
	if err != nil {
		log.Printf("fix generation error: %v", err)
		// Non-fatal — continue without fixes
	}

	// Enrich violations with metadata and fixes
	var enrichedViolations []models.Violation
	for _, v := range uniqueViolations {
		meta, ok := models.ViolationMetadata[v.Code]
		if !ok {
			meta = models.ViolationMeta{Title: v.Code, Severity: "UNKNOWN"}
		}
		fix := ""
		if fixes != nil {
			fix = fixes[v.Code]
		}
		enrichedViolations = append(enrichedViolations, models.Violation{
			Code:         v.Code,
			Title:        meta.Title,
			Severity:     meta.Severity,
			Explanation:  v.Explanation,
			SuggestedFix: fix,
		})
	}

	if enrichedViolations == nil {
		enrichedViolations = []models.Violation{}
	}

	result := models.CheckResult{
		Action:          scoreResult.Action,
		ComplianceScore: scoreResult.ComplianceScore,
		RiskCategory:    scoreResult.RiskCategory,
		Category:        categoryResult.Category,
		CategoryState:   categoryResult.CategoryState,
		Violations:      enrichedViolations,
	}

	json.NewEncoder(w).Encode(result)
}
```

---

## Entry Point

**`backend/cmd/server/main.go`**

```go
package main

import (
	"log"
	"net/http"
	"os"
	"compliance-checker/internal/api"
	"compliance-checker/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/checks", api.HandleCheck)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
```

Add godotenv: `go get github.com/joho/godotenv`

---

## Frontend Changes — Targeted Only

Make these **three changes only** in `src/components/ComplianceDashboard.vue`. Do not change the UI, layout, styles, or any other logic.

### Change 1: Add API result state variables

In the `<script setup>` section, add these new reactive variables after the existing ones:

```typescript
// API result state
const checkResult = ref<{
  action: string
  compliance_score: number
  risk_category: string
  category: string
  violations: Array<{
    code: string
    title: string
    severity: string
    explanation: string
    suggested_fix: string
  }>
} | null>(null)
const checkError = ref<string | null>(null)
```

### Change 2: Replace the nextStep function

Replace the existing `nextStep()` function with this version that calls the real API:

```typescript
async function nextStep() {
  if (currentStep.value === 2) {
    isChecking.value = true
    currentStep.value = 3
    checkError.value = null

    try {
      const response = await fetch('http://localhost:8080/checks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          platform: selectedPlatform.value,
          region: selectedRegion.value,
          age_min: parseInt(selectedMinAge.value),
          ad_format: adFormat.value,
          primary_text: adCopy.value.primaryText,
          headline: adCopy.value.headline,
          description: adCopy.value.description,
          landing_page_url: landingPageUrl.value
        })
      })

      if (!response.ok) {
        throw new Error('API request failed')
      }

      checkResult.value = await response.json()
      isChecking.value = false
      currentStep.value = 4
    } catch (err) {
      isChecking.value = false
      checkError.value = 'Analysis failed. Please check the backend is running and try again.'
      currentStep.value = 2
    }
  } else {
    currentStep.value++
  }
}
```

### Change 3: Replace Step 4 hardcoded results with dynamic data

In the template, find the **Step 4 Results Dashboard** section (the `v-if="currentStep === 4"` div). Replace only the **right column** (the `col-span-3` div containing the hardcoded issues) with this dynamic version:

```html
<!-- Right: Results Checklist (3 cols) -->
<div class="col-span-3 space-y-4">

  <!-- Score Summary Card -->
  <div class="border border-border rounded-xl bg-card p-5">
    <div class="flex items-center justify-between mb-3">
      <div>
        <div class="text-3xl font-bold" :class="checkResult?.compliance_score === 0 ? 'text-destructive' : checkResult?.compliance_score >= 70 ? 'text-green-600' : 'text-amber-500'">
          {{ checkResult?.compliance_score ?? 0 }}<span class="text-lg font-normal text-muted-foreground">/100</span>
        </div>
        <div class="text-sm font-semibold mt-1" :class="checkResult?.risk_category === 'High Risk' ? 'text-destructive' : checkResult?.risk_category === 'Medium Risk' ? 'text-amber-500' : 'text-green-600'">
          {{ checkResult?.risk_category }}
        </div>
      </div>
      <div class="text-right">
        <div class="text-xs text-muted-foreground">Category</div>
        <div class="text-sm font-semibold">{{ checkResult?.category }}</div>
        <div class="text-xs mt-1 px-2 py-0.5 rounded-full inline-block"
          :class="checkResult?.action === 'not_compliant' ? 'bg-destructive/10 text-destructive' : 'bg-green-100 text-green-700'">
          {{ checkResult?.action === 'not_compliant' ? 'Not Compliant' : 'Compliant' }}
        </div>
      </div>
    </div>
    <!-- Score bar -->
    <div class="w-full bg-muted rounded-full h-2">
      <div class="h-2 rounded-full transition-all duration-500"
        :class="checkResult?.compliance_score === 0 ? 'bg-destructive' : checkResult?.compliance_score >= 70 ? 'bg-green-500' : 'bg-amber-400'"
        :style="`width: ${checkResult?.compliance_score ?? 0}%`">
      </div>
    </div>
  </div>

  <!-- Violations List -->
  <div v-if="checkResult?.violations && checkResult.violations.length > 0"
       class="border border-border rounded-xl bg-card overflow-hidden">
    <div class="bg-destructive/10 border-b border-destructive/20 p-4 flex items-center justify-between">
      <div class="flex items-center gap-2 text-destructive font-bold">
        <IconAlertCircle :size="20" />
        {{ checkResult.violations.length }} Issue{{ checkResult.violations.length > 1 ? 's' : '' }} Found
      </div>
      <div class="bg-destructive text-destructive-foreground text-xs font-bold px-2 py-0.5 rounded-full">
        {{ checkResult.risk_category }}
      </div>
    </div>
    <div class="p-6 space-y-6">
      <div v-for="violation in checkResult.violations" :key="violation.code">
        <div class="flex items-start gap-2">
          <IconAlertCircle class="text-destructive mt-0.5 shrink-0" :size="16" />
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-1">
              <h4 class="font-bold text-sm text-foreground">{{ violation.title }}</h4>
              <span class="text-[10px] font-bold px-1.5 py-0.5 rounded"
                :class="violation.severity === 'PROHIBITED' ? 'bg-destructive/10 text-destructive' : violation.severity === 'AUTHORIZATION_REQUIRED' ? 'bg-orange-100 text-orange-700' : violation.severity === 'RESTRICTED' ? 'bg-amber-100 text-amber-700' : 'bg-blue-100 text-blue-700'">
                {{ violation.severity }}
              </span>
            </div>
            <p class="text-sm text-muted-foreground">{{ violation.explanation }}</p>
            <div v-if="violation.suggested_fix" class="mt-2 bg-green-50 border border-green-200 rounded-lg p-3">
              <p class="text-xs font-bold text-green-700 mb-1">Suggested Fix</p>
              <p class="text-sm text-green-800">{{ violation.suggested_fix }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- No violations -->
  <div v-else class="bg-green-50 border border-green-200 p-4 rounded-xl flex items-start gap-3">
    <IconCheck class="text-green-600 shrink-0 mt-0.5" />
    <p class="text-sm text-green-800 font-medium">No policy violations detected. This ad appears compliant with Meta's advertising standards.</p>
  </div>

</div>
```

---

## How to Run Locally

```bash
# 1. Start the database
cd backend
docker compose up -d

# 2. Copy env file and add your Gemini API key
cp .env.example .env
# Edit .env and set GEMINI_API_KEY=your_key

# 3. Start the Go server
go run cmd/server/main.go

# 4. In a new terminal, start the Vue frontend
cd ..  # back to project root
npm install
npm run dev
# Frontend runs at http://localhost:5173
# Backend runs at http://localhost:8080
```

---

## Validation — 10 Test Cases

After the full pipeline is running, test these inputs in the UI and verify the expected outputs:

| # | Primary Text | Expected violation | Expected action | Expected score |
|---|---|---|---|---|
| 1 | "Lose 20kg in 30 days with our supplement!" | R_WEIGHT_LOSS_CLAIM | compliant | ~50 |
| 2 | "Are you struggling with debt? We can help." | P_PERSONAL_ATTRIBUTE | not_compliant | 0 |
| 3 | "Bet on Premier League today — big wins await!" | A_GAMBLING | not_compliant | 0 |
| 4 | "Get fit this summer. Shop our activewear." | none | compliant | 90+ |
| 5 | "DOCTORS HATE HIM! This one trick cures arthritis!" | R_HEALTH_CLAIM + M_SENSATIONALISM | compliant | ~35 |
| 6 | "Buy Rolex watches — genuine Swiss luxury deals" | P_IMPERSONATION or P_COUNTERFEIT | not_compliant | 0 |
| 7 | "Only 3 spots left! Sign up today or miss out!" | M_URGENCY | compliant | 75 |
| 8 | "Trade Bitcoin on our platform — high returns guaranteed" | A_CRYPTO_EXCHANGE | not_compliant | 0 |
| 9 | "Tired of your belly fat? Our pill melts it in 7 days!" | R_WEIGHT_LOSS_CLAIM + R_BODY_SHAMING | compliant | ~40 |
| 10 | "For Christians looking to invest ethically." | P_PERSONAL_ATTRIBUTE | not_compliant | 0 |

**Pass threshold:** 8/10 correct. Below 7 means prompts need calibration.

---

## Constraints

- Do not use any framework (Gin, Echo, Fiber). Use Go's standard `net/http` only.
- Do not add authentication or middleware for the MVP.
- Do not change the Vue frontend design, layout, or styles — only the three changes specified above.
- All prompts must be read from `.txt` files, not hardcoded in Go.
- Temperature must be 0 for all Gemini calls.
- All Gemini calls must use ResponseSchema enforcement — never parse free-form text.
- The scoring engine must remain pure Go with no LLM involvement.

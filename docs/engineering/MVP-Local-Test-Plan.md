# Ad Compliance Checker — Local MVP Plan

**Stack:** Go · PostgreSQL · Gemini · Vue 3  
**Scope:** Text-only · Synchronous · 12 violation codes · Single ad submission  
**Goal:** Validate category classification accuracy, scoring logic, and end-to-end latency before writing production infrastructure

---

## Decisions Locked

| Decision | Choice | Reason |
|---|---|---|
| LLM provider | Google Gemini (gemini-1.5-flash) | Aligns with Google tech stack |
| Agent architecture | Separate LLM call per agent | Separation of concerns, accurate prompts, Go goroutines make parallelism cheap |
| MVP scope | Text-only (no image) | Covers ~75% of violations, validates core logic first |
| Auth-required path | Always assume unauthorized | Simplest default, correct for most users |
| Violation codes in MVP | 12 (see Section 5) | Covers all 4 tiers, most common violations |
| API key for local | Google AI Studio key | Free tier, no GCP setup needed for local |
| Gemini output mode | JSON schema enforcement | Required for determinism — no free-form text |
| Fix Generation calls | Single batched call (all violations at once) | Free tier = 15 RPM. 1 call/ad for fixes instead of N calls keeps us safely under limit |
| Phase 2 error handling | Fault-tolerant goroutines (not errgroup) | Partial result > total failure if one agent has a network blip |

---

## Section 1 — Prerequisites

Install these before starting:

```bash
# Required
go 1.22+
docker + docker compose
node 20+ + npm
git

# Verify
go version
docker --version
node --version
```

Get a Gemini API key:  
→ https://aistudio.google.com/app/apikey  
(Free tier is enough for local testing)

---

## Section 2 — Project Structure

```
compliance-checker/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, wire everything together
├── internal/
│   ├── api/
│   │   └── handler.go           # HTTP handlers: POST /checks, GET /checks/{id}
│   ├── agents/
│   │   ├── category.go          # Phase 1: Category Classifier
│   │   ├── violations.go        # Phase 2: Prohibited + Restricted + Quality agents
│   │   ├── personal_attribute.go # Phase 2: Personal Attribute Agent (runs on all ads)
│   │   └── fix_generation.go    # Fix Generation: one call per violation
│   ├── scoring/
│   │   └── engine.go            # Scoring Engine — pure Go, no LLM
│   ├── db/
│   │   ├── postgres.go          # Connection + migrations
│   │   └── checks.go            # DB queries for checks table
│   └── models/
│       └── types.go             # Shared types: ViolationCode, CheckResult, etc.
├── frontend/                    # Vue 3 app (see Section 7)
├── prompts/                     # System prompts as .txt files (easier to edit)
│   ├── category_classifier.txt
│   ├── violations_agent.txt
│   ├── personal_attribute_agent.txt
│   └── fix_generation_agent.txt
├── docker-compose.yml
├── .env.example
└── go.mod
```

> **Why `prompts/` as separate files?**  
> Prompts will change constantly as you calibrate accuracy. Keeping them as files (not hardcoded strings) means you can edit and re-test without recompiling. In production, these move to a config store.

---

## Section 3 — Local Environment Setup

**docker-compose.yml**
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

**.env.example** (copy to `.env`, fill in your key)
```env
GEMINI_API_KEY=your_key_here
DATABASE_URL=postgres://dev:dev@localhost:5432/compliance_checker?sslmode=disable
PORT=8080
```

**Start the database:**
```bash
docker compose up -d
go run cmd/server/main.go
# Frontend: cd frontend && npm install && npm run dev
```

---

## Section 4 — Database Schema

One table for the MVP. Keep it simple.

```sql
CREATE TABLE checks (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  status      TEXT NOT NULL DEFAULT 'pending',  -- pending | processing | complete | error
  input_copy  TEXT NOT NULL,
  result      JSONB,                             -- full CheckResult stored as JSON
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Run this once on startup (or use a migration file).

---

## Section 5 — MVP Violation Codes (12 of 34)

These 12 cover all four tiers and the highest-frequency violations.  
Full 34-code coverage comes after MVP validation.

| Code | Tier | Why include in MVP |
|---|---|---|
| `P_HATE_SPEECH` | P_ | Common, clear signal, good test case |
| `P_ILLEGAL_DRUGS` | P_ | Common in DTC health space |
| `P_PERSONAL_ATTRIBUTE` | P_ | Unique Meta rule, ADSENTRY missed it — must validate |
| `P_IMPERSONATION` | P_ | High frequency in performance marketing |
| `A_CRYPTO_EXCHANGE` | A_ | Common in Rockads customer base |
| `A_GAMBLING` | A_ | Common, clear authorization requirement |
| `R_WEIGHT_LOSS_CLAIM` | R_ | Most enforced Meta rule in health vertical |
| `R_BODY_SHAMING` | R_ | Directly related to weight loss |
| `R_HEALTH_CLAIM` | R_ | High frequency in supplements/DTC |
| `R_GUARANTEED_OUTCOME` | R_ | Common in finance and health |
| `M_URGENCY` | M_ | Nearly universal in performance ads |
| `M_DISCLOSURE` | M_ | Easy to check, high signal |

---

## Section 6 — Agent Architecture

### How it runs (Go goroutines)

```
POST /checks received
        │
        ▼
[1] Category Classifier Agent  (1 LLM call — blocks until done)
        │
        ├── If PROHIBITED category → skip to scoring, violations = [P_ code]
        │
        ▼
[2] Phase 2 — 4 agents run in PARALLEL (goroutines + errgroup)
        ├── Violations Agent       (1 LLM call — checks P_, R_, M_ codes)
        ├── Personal Attribute Agent (1 LLM call — always runs)
        ├── Auth-Required Agent    (1 LLM call — checks A_ codes)
        └── [Image Agent]          (DEFERRED — not in MVP)
        │
        ▼  (collect all violations, deduplicate)
[3] Scoring Engine             (pure Go — no LLM)
        │
        ▼
[4] Fix Generation             (1 batched LLM call — all violations at once)
        │
        ▼
[5] Write result to PostgreSQL, return response
```

**Key Go pattern for Phase 2 — fault-tolerant goroutines:**

> **Why not errgroup?** errgroup cancels all agents if one fails. For MVP, a network blip on the Auth agent shouldn't discard results from the Violations agent. Fault-tolerant: collect what succeeds, log what fails.

```go
var mu sync.Mutex
var allViolations []models.Violation
var wg sync.WaitGroup

runAgent := func(fn func() ([]models.Violation, error), name string) {
    defer wg.Done()
    v, err := fn()
    if err != nil {
        log.Printf("agent %s failed: %v", name, err)
        return // log and continue — don't fail the whole check
    }
    mu.Lock()
    allViolations = append(allViolations, v...)
    mu.Unlock()
}

wg.Add(3)
go runAgent(func() ([]models.Violation, error) {
    return agents.RunViolationsAgent(ctx, copy, category)
}, "violations")
go runAgent(func() ([]models.Violation, error) {
    return agents.RunPersonalAttributeAgent(ctx, copy)
}, "personal_attribute")
go runAgent(func() ([]models.Violation, error) {
    return agents.RunAuthRequiredAgent(ctx, copy, category)
}, "auth_required")

wg.Wait()
```

---

### Agent 1 — Category Classifier

**Gemini call:**
```go
// Model: gemini-1.5-flash
// Temperature: 0
// Response schema: enforce JSON
```

**Output schema (enforce via Gemini response_schema):**
```json
{
  "category": "Health — Weight Loss",
  "category_state": "RESTRICTED",
  "reasoning": "Ad promotes a weight loss supplement with specific outcome claims."
}
```

**Prompt (prompts/category_classifier.txt):**  
Give it Section 3 of Meta-Policy-Spec-Compliance-Checker.md — the category taxonomy table only. Keep the prompt under 1000 tokens. One task: classify.

---

### Agent 2 — Violations Agent

**Covers:** P_ codes + R_ codes + M_ codes from the 12 MVP codes  
(excludes A_ codes — that's the Auth-Required Agent)

**Output schema:**
```json
{
  "violations": [
    {
      "code": "R_WEIGHT_LOSS_CLAIM",
      "explanation": "Copy states 'Lose 20kg in 30 days' — specific numeric claim with timeframe."
    }
  ]
}
```

If no violations: `{ "violations": [] }`

**Prompt:** Section 4 (P_ codes) + Section 6 (R_ codes) + Section 7 (M_ codes) for the 10 non-A_ MVP codes only. Keep it focused — don't dump the whole spec.

---

### Agent 3 — Personal Attribute Agent

**Always runs regardless of category.**

**Output schema:**
```json
{
  "violation_found": true,
  "explanation": "'Are you struggling with debt?' implies financial vulnerability."
}
```

**Prompt:** Section 8 of the spec only. This agent does one job — the prompt can be very short and very precise.

---

### Agent 4 — Auth-Required Agent

**Covers:** A_CRYPTO_EXCHANGE, A_GAMBLING (MVP scope)

**Output schema:**
```json
{
  "violations": [
    {
      "code": "A_GAMBLING",
      "explanation": "Ad promotes online betting without indicating Meta written permission."
    }
  ]
}
```

---

### Scoring Engine (pure Go — no LLM)

```go
// internal/scoring/engine.go

type Result struct {
    ComplianceScore int    
    RiskCategory    string 
    Action          string 
}

func Score(violations []string, categoryState string, authorized bool) Result {
    var p, a, r, m []string
    for _, v := range violations {
        switch {
        case strings.HasPrefix(v, "P_"): p = append(p, v)
        case strings.HasPrefix(v, "A_"): a = append(a, v)
        case strings.HasPrefix(v, "R_"): r = append(r, v)
        case strings.HasPrefix(v, "M_"): m = append(m, v)
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
        if score < 1 { score = 1 }
        risk := "Medium Risk"
        if score < 30 { risk = "High Risk" }
        return Result{score, risk, "compliant"}
    }
    // Rule 4: Quality/moderate violations only
    if len(m) > 0 {
        score := 75 - (5 * (len(m) - 1))
        if score < 45 { score = 45 }
        risk := "Low Risk"
        if score < 65 { risk = "Medium Risk" }
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

Unit test this before wiring LLM calls. It's pure logic — 100% testable.

---

### Fix Generation Agent

**Single batched call — all violations at once.**

> **Why batched?** Free tier is 15 RPM. N parallel fix calls per ad burns the quota fast. One batched call costs 1 request regardless of violation count, keeps you safely under the limit, and also lets the agent see all violations together — preventing a fix for one violation from accidentally introducing another.

**Input (all violations in one payload):**
```json
{
  "category": "Health — Weight Loss",
  "original_copy": "Lose 20kg in 30 days with our new supplement! Limited time offer.",
  "all_violations_found": ["R_WEIGHT_LOSS_CLAIM", "M_URGENCY"],
  "violations_to_fix": [
    {
      "code": "R_WEIGHT_LOSS_CLAIM",
      "fix_framework": "Remove specific numeric claim and timeframe. Replace with aspirational language. Cannot include specific amounts, timeframes, or guarantees."
    },
    {
      "code": "M_URGENCY",
      "fix_framework": "No fix required if used in isolation. If stacked with other violations, consider removing countdown timers or disclosing the actual sale end date."
    }
  ]
}
```

> Including `category` and `all_violations_found` gives the agent full context — it won't rewrite away one violation only to introduce a new one.

**Output schema:**
```json
{
  "fixes": [
    {
      "code": "R_WEIGHT_LOSS_CLAIM",
      "suggested_fix": "Support your weight management goals with our new supplement."
    },
    {
      "code": "M_URGENCY",
      "suggested_fix": "Shop our supplement — offer available while stocks last."
    }
  ]
}
```

**Prompt (prompts/fix_generation_agent.txt):**  
Short. You receive the full context above. Generate one rewrite per violation. Do not explain — just the rewritten copy. Rewrites must not introduce new violations.

---

## Section 7 — API Contract (MVP)

**POST /checks**
```json
// Request
{ "copy": "Lose 20kg in 30 days! Limited time offer." }

// Response (immediate)
{ "check_id": "uuid", "status": "processing" }
```

**GET /checks/{id}**
```json
{
  "check_id": "uuid",
  "status": "complete",
  "action": "compliant",
  "compliance_score": 35,
  "risk_category": "Medium Risk",
  "category": "Health — Weight Loss",
  "violations": [
    {
      "code": "R_WEIGHT_LOSS_CLAIM",
      "title": "Specific Weight Loss Claim",
      "severity": "RESTRICTED",
      "policy_reference": "https://transparency.meta.com/policies/ad-standards/",
      "explanation": "Copy states 'Lose 20kg in 30 days' — specific numeric claim with timeframe.",
      "suggested_fix": "Support your weight management goals with our supplement."
    },
    {
      "code": "M_URGENCY",
      "title": "Urgency Tactic",
      "severity": "MODERATE",
      "explanation": "'Limited time offer' with no end date specified.",
      "suggested_fix": "Add the actual offer end date, or remove the urgency framing."
    }
  ]
}
```

> For MVP, process synchronously — no polling needed. Return the full result directly from POST /checks. Add async processing only when latency becomes an issue.

---

## Section 8 — Vue Frontend (Minimal)

Three views only:

**View 1 — Submit Form**
```
┌─────────────────────────────────────────┐
│  Ad Compliance Checker                  │
│                                         │
│  Ad Copy                                │
│  ┌───────────────────────────────────┐  │
│  │ Paste your ad copy here...        │  │
│  │                                   │  │
│  └───────────────────────────────────┘  │
│                                         │
│         [ Check Compliance ]            │
└─────────────────────────────────────────┘
```

**View 2 — Loading**
```
Analysing your ad...  ⏳
```

**View 3 — Results**
```
┌──────────────────────────────────────────┐
│  Compliance Score        35 / 100        │
│  ████░░░░░░░░░░░░  Medium Risk           │
│  Status: compliant (2 issues found)      │
├──────────────────────────────────────────┤
│  ⚠ R_  Specific Weight Loss Claim        │
│  "Lose 20kg in 30 days..." implies a     │
│  specific numeric claim.                 │
│  Fix → "Support your weight management   │
│  goals with our supplement."             │
├──────────────────────────────────────────┤
│  ⚠ M_  Urgency Tactic                   │
│  'Limited time offer' with no end date.  │
│  Fix → Add actual offer end date.        │
└──────────────────────────────────────────┘
```

No routing, no Pinia, no auth. One page. Get the data rendering correctly first.

---

## Section 9 — Go Dependencies

```bash
go mod init compliance-checker
go get github.com/google/generative-ai-go/genai@latest
go get google.golang.org/api/option
go get github.com/lib/pq                     # PostgreSQL driver
go get github.com/google/uuid
go get golang.org/x/sync                     # errgroup for parallel agents
```

**Gemini client setup:**
```go
import (
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
model := client.GenerativeModel("gemini-1.5-flash")
model.SetTemperature(0)  // determinism
// Set ResponseMIMEType + ResponseSchema for structured output
```

---

## Section 10 — Build Order

Do this in sequence. Each step is testable independently.

| Step | What to build | Test signal |
|---|---|---|
| 1 | DB schema + Go server skeleton + POST/GET handlers (no LLM) | `curl localhost:8080/checks` returns dummy JSON |
| 2 | Scoring Engine in pure Go | Unit tests pass for all 6 scoring rules |
| 3 | Category Classifier Agent | 10 manual test inputs classified correctly |
| 4 | Violations Agent + Personal Attribute Agent (parallel) | Known violating copies return correct codes |
| 5 | Auth-Required Agent | Crypto/gambling copy returns A_ codes |
| 6 | Fix Generation Agent | Fixes are coherent rewrites, not hallucinations |
| 7 | Wire all agents → scoring engine → DB write | End-to-end test with 20 test cases |
| 8 | Vue form + results view | Can submit and read results in browser |

**Start with Step 2, not Step 1.** The scoring engine has zero LLM dependency — you can write and test it in an hour. Getting it right before touching the LLM means your LLM output → score mapping is trustworthy from the start.

---

## Section 11 — Test Cases for Validation

Run these 10 inputs once the full pipeline is wired. Expected outputs are pre-determined.

| # | Input copy | Expected top violation | Expected action | Expected score |
|---|---|---|---|---|
| 1 | "Lose 20kg in 30 days with our supplement!" | R_WEIGHT_LOSS_CLAIM | compliant | ~50 |
| 2 | "Are you struggling with debt? We can help." | P_PERSONAL_ATTRIBUTE | not_compliant | 0 |
| 3 | "Bet on Premier League today — big wins await!" | A_GAMBLING | not_compliant | 0 |
| 4 | "Get fit this summer. Shop our activewear." | none | compliant | 95 |
| 5 | "DOCTORS HATE HIM! This one trick cures arthritis!" | R_HEALTH_CLAIM + M_SENSATIONALISM | compliant | ~35 |
| 6 | "Buy Rolex watches — genuine Swiss luxury" | P_IMPERSONATION or P_COUNTERFEIT | not_compliant | 0 |
| 7 | "Only 3 spots left! Sign up today or miss out!" | M_URGENCY | compliant | 75 |
| 8 | "Trade Bitcoin on our platform — high returns guaranteed" | A_CRYPTO_EXCHANGE | not_compliant | 0 |
| 9 | "Tired of your belly fat? Our pill melts it away in 7 days!" | R_WEIGHT_LOSS_CLAIM + R_BODY_SHAMING | compliant | ~40 |
| 10 | "For Christians looking to invest ethically." | P_PERSONAL_ATTRIBUTE | not_compliant | 0 |

**Pass criteria:** 8/10 correct. Anything below 7 means the prompts need calibration before moving forward.

---

## Out of Scope for MVP

- Image / video analysis
- Landing page URL fetch
- Async processing (email notification)
- Authorization path (user indicates they hold Meta auth)
- Full 34-code coverage (22 codes deferred)
- Age targeting check (R_AGE_TARGETING_BREACH)
- Authentication / user accounts
- Rate limiting

---

## Estimated Effort

| Component | Effort |
|---|---|
| DB + Go skeleton + API handlers | 0.5 day |
| Scoring engine + unit tests | 0.5 day |
| Gemini integration + all 4 agents | 1.5 days |
| Prompt calibration (iterate until 8/10 pass) | 1 day |
| Vue frontend | 1 day |
| **Total** | **~4.5 days** |

Prompt calibration is the wildcard. Budget more time here if the category classifier misclassifies edge cases.

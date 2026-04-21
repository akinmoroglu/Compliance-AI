# Antigravity Implementation Instructions
# Ad Compliance Checker — v2 (Multimodal)

---

## BEFORE YOU START — Read This First

**Read this entire document before writing a single line of code.** Later blocks correct and extend earlier ones. Implementing block by block without reading ahead will produce inconsistencies.

**This document replaces the v1 Antigravity instructions entirely.** Do not merge, do not cross-reference v1. Implement solely from this document.

**Recommended implementation order:**

1. `go.mod` — add new dependencies
2. `internal/models/types.go` — all types; everything else depends on this
3. `internal/models/metadata.go` — add new violation codes
4. `internal/db/postgres.go` — schema migration
5. `internal/agents/client.go` — model config flag
6. `internal/cache/cache.go` — new file
7. `internal/preprocessor/preprocessor.go` — new file
8. `internal/agents/violations.go` — updated signature
9. `internal/agents/personal_attribute.go` — updated signature
10. `internal/agents/vision.go` — new file
11. `internal/agents/audio.go` — new file
12. `internal/agents/lp.go` — new file
13. `internal/agents/alignment.go` — new file
14. `internal/agents/fix_generation.go` — minor update (_suggested_copy)
15. `internal/scoring/engine.go` — updated action strings
16. `internal/advisory/advisory.go` — new file
17. `internal/api/handler.go` — full rewrite; depends on all of the above
18. `cmd/server/main.go` — minor update
19. `src/components/AdvisoryPanel.vue` — new Vue component
20. `src/components/ComplianceDashboard.vue` — full frontend rewrite

**Do not skip steps or reorder.** The handler (step 17) references all packages above it.

---

## BLOCK 1 — Architecture, File Structure, Types, API Contract

---

## What You Are Building

You are upgrading the existing Go backend from v1 (text-only) to v2 (multimodal). v2 adds image analysis, video analysis, audio transcription, landing page fetching, and cross-source alignment checking. The Vue 3 frontend will also be updated with new inputs and result displays.

**This document is a complete replacement for the v1 Antigravity instructions. Do not merge with the v1 document — implement solely from this document.**

Your tasks:
1. Modify the existing Go backend at `backend/` to implement the v2 pipeline
2. Add new agent files, a pre-processor, a cache layer, and an advisory system
3. Update the existing handler, scoring engine, types, and client helper
4. Update the Vue frontend with new input fields and result components

---

## What Already Exists (Do Not Reinvent)

The following files exist from v1 and must be **modified** (not deleted and recreated):

- `backend/internal/agents/client.go` — Gemini client helper
- `backend/internal/agents/violations.go` — Violations agent
- `backend/internal/agents/personal_attribute.go` — Personal attribute agent
- `backend/internal/agents/auth_required.go` — Auth required agent
- `backend/internal/agents/fix_generation.go` — Fix generation agent
- `backend/internal/agents/category.go` — Category classifier (no changes needed)
- `backend/internal/api/handler.go` — Orchestrator/handler (full rewrite)
- `backend/internal/scoring/engine.go` — Scoring engine (update)
- `backend/internal/models/types.go` — Shared types (update)
- `backend/internal/models/metadata.go` — Violation metadata (update)
- `backend/internal/db/postgres.go` — DB init (minor update)

The following prompt files **already exist** and must NOT be overwritten:

- `backend/prompts/category_classifier.txt` — unchanged from v1
- `backend/prompts/violations_agent.txt` — already updated for v2
- `backend/prompts/personal_attribute_agent.txt` — unchanged from v1
- `backend/prompts/auth_required_agent.txt` — unchanged from v1
- `backend/prompts/fix_generation_agent.txt` — unchanged from v1
- `backend/prompts/vision_agent.txt` — new, already written
- `backend/prompts/audio_agent.txt` — new, already written
- `backend/prompts/lp_agent.txt` — new, already written
- `backend/prompts/alignment_agent.txt` — new, already written

---

## Backend: Complete v2 File Structure

The full structure after v2 implementation:

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    (modify — new env vars)
├── internal/
│   ├── api/
│   │   └── handler.go                 (FULL REWRITE — v2 orchestrator)
│   ├── agents/
│   │   ├── client.go                  (modify — model config flag, multipart support)
│   │   ├── category.go                (no change)
│   │   ├── violations.go              (modify — new parameters)
│   │   ├── personal_attribute.go      (modify — new parameters)
│   │   ├── auth_required.go           (no change)
│   │   ├── fix_generation.go          (modify — new input fields)
│   │   ├── vision.go                  (CREATE NEW)
│   │   ├── audio.go                   (CREATE NEW)
│   │   ├── lp.go                      (CREATE NEW)
│   │   └── alignment.go               (CREATE NEW)
│   ├── preprocessor/
│   │   └── preprocessor.go            (CREATE NEW)
│   ├── cache/
│   │   └── cache.go                   (CREATE NEW)
│   ├── advisory/
│   │   └── advisory.go                (CREATE NEW)
│   ├── scoring/
│   │   └── engine.go                  (modify — new action strings, ALIGN_ tier)
│   ├── db/
│   │   ├── postgres.go                (modify — schema update)
│   │   └── checks.go                  (no change)
│   └── models/
│       ├── types.go                   (FULL REWRITE — all new types)
│       └── metadata.go                (modify — add new violation codes)
├── prompts/
│   ├── category_classifier.txt        (no change)
│   ├── violations_agent.txt           (already updated)
│   ├── personal_attribute_agent.txt   (no change)
│   ├── auth_required_agent.txt        (no change)
│   ├── fix_generation_agent.txt       (no change)
│   ├── vision_agent.txt               (already written)
│   ├── audio_agent.txt                (already written)
│   ├── lp_agent.txt                   (already written)
│   └── alignment_agent.txt            (already written)
├── docker-compose.yml                 (modify — add ffmpeg note)
├── .env.example                       (modify — add new keys)
├── .gitignore
└── go.mod                             (modify — add new dependencies)
```

Frontend files to create or modify:

```
src/
├── components/
│   ├── ComplianceDashboard.vue        (FULL REWRITE — v2 inputs and results)
│   └── AdvisoryPanel.vue              (CREATE NEW)
└── App.vue                            (no change)
```

---

## Environment Setup

**`backend/.env.example`** — replace existing file with:

```
GEMINI_API_KEY=your_google_ai_studio_key_here
FIRECRAWL_API_KEY=your_firecrawl_api_key_here
DATABASE_URL=postgres://dev:dev@localhost:5432/compliance_checker?sslmode=disable
PORT=8080
GEMINI_MODEL=gemini-2.0-flash
MAX_VIDEO_DURATION_SECONDS=120
```

`GEMINI_MODEL`: model string used by all agents. Default `gemini-2.0-flash`. Can be changed to `gemini-1.5-pro` per agent if needed — see the agent client section.

`MAX_VIDEO_DURATION_SECONDS`: configurable cap for accepted video length. Default 120. Engineering decision — do not hardcode this value in Go; always read from env with fallback to 120.

**New Go dependencies** — run from `backend/`:

```bash
go get github.com/mendableai/firecrawl-go@latest
go get github.com/joho/godotenv
```

ffmpeg is a system binary, not a Go package. For local development, install via:
- macOS: `brew install ffmpeg`
- Ubuntu/Debian: `apt-get install -y ffmpeg`

For Docker deployment, add to the Go server Dockerfile:
```dockerfile
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*
```

**`backend/docker-compose.yml`** — no change needed for local development (postgres only). The Go server runs locally with ffmpeg installed on the host machine.

---

## Database Schema Update

**`backend/internal/db/postgres.go`** — update the `migrate()` function to add new columns:

```go
func migrate() {
    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS checks (
            id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            status          TEXT NOT NULL DEFAULT 'complete',
            input_copy      TEXT NOT NULL,
            content_hash    TEXT,
            result          JSONB,
            created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
        );

        ALTER TABLE checks ADD COLUMN IF NOT EXISTS content_hash TEXT;
    `)
    if err != nil {
        log.Fatalf("migration failed: %v", err)
    }
}
```

`content_hash` stores the SHA-256 cache key for the request. Used by the cache layer to look up prior results. See the Cache section in Block 2.

---

## Shared Types

**`backend/internal/models/types.go`** — full replacement:

```go
package models

// ─────────────────────────────────────────────
// REQUEST TYPES
// ─────────────────────────────────────────────

// CheckRequest is populated by the handler from the multipart form (or JSON fallback).
// When a file is present the handler uses multipart/form-data.
// When no file is present the handler accepts application/json for backward compatibility.
type CheckRequest struct {
    Platform       string   `json:"platform"`
    Countries      []string `json:"countries"`       // empty = global / all countries
    AgeMin         int      `json:"age_min"`
    AgeMax         int      `json:"age_max"`
    AdFormat       string   `json:"ad_format"`       // "single" | "carousel"
    PrimaryText    string   `json:"primary_text"`
    Headline       string   `json:"headline"`
    Description    string   `json:"description"`
    LandingPageURL string   `json:"landing_page_url"`
    ForceRefresh   bool     `json:"force_refresh"`
    // File fields are NOT JSON-serializable — handled as []byte in the preprocessor
}

// ─────────────────────────────────────────────
// PIPELINE CONTEXT — passed between phases
// ─────────────────────────────────────────────

// RequestContext holds all pre-processed media and extracted content.
// Populated by the pre-processor before Phase 1 begins.
// Passed by pointer to all agent goroutines and Phase 1b.
// Each field has exactly one writer — no mutex needed on the struct itself.
type RequestContext struct {
    // From pre-processor (populated before Phase 1)
    ImageBytes       []byte   // nil if no image uploaded
    VideoFrames      [][]byte // JPEG frame bytes extracted by ffmpeg (nil if no video)
    AudioBytes       []byte   // MP3 audio track bytes (nil if no video or no audio track)
    HasImage         bool
    HasVideo         bool

    // From Phase 1 agents (written by individual goroutines; non-overlapping fields)
    AudioTranscript       string // g6: Audio Agent
    AudioAnalysisStatus   string // g6: "transcribed"|"no_speech_detected"|"low_quality_audio"|"music_only"
    ImageExtractedText    string // g4: Vision Agent OCR (Part 1 of vision prompt)
    VisualCheckConfidence string // g4: "high"|"partial"|"low"
    DetectedLanguage      string // g4 or g6: first writer wins
    LPKeyClaims           *LPKeyClaims // g5: LP Agent structured claims — used by Alignment Agent
    LPAnalysisStatus      string       // g5: "analyzed"|"fetch_failed"|"js_rendered"
    Category              string       // Sequential: set before Phase 1
    CategoryState         string       // Sequential: set before Phase 1
}

// ─────────────────────────────────────────────
// AGENT OUTPUT TYPES
// ─────────────────────────────────────────────

type AgentViolation struct {
    Code        string `json:"code"`
    Explanation string `json:"explanation"`
    Source      string `json:"source"` // "copy" | "image_text" | "voiceover" | "landing_page"
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

// VisionAgentOutput is the structured output from the Vision Agent.
// Part 1: OCR of all visible text. Part 2: visual violations.
type VisionAgentOutput struct {
    ExtractedText         string          `json:"extracted_text"`
    VisualViolations      []AgentViolation `json:"visual_violations"`
    VisualCheckConfidence string          `json:"visual_check_confidence"` // "high" | "partial" | "low"
    DetectedLanguage      string          `json:"detected_language"`
}

// AudioAgentOutput is the structured output from the Audio Agent.
// Transcription only — no violation checking.
type AudioAgentOutput struct {
    Transcript          string `json:"transcript"`
    DetectedLanguage    string `json:"detected_language"`
    AudioAnalysisStatus string `json:"audio_analysis_status"` // "transcribed" | "no_speech_detected" | "low_quality_audio" | "music_only"
    HasVoiceover        bool   `json:"has_voiceover"`
    DurationNote        string `json:"duration_note"`
}

// LPKeyClaims holds structured claims extracted from the landing page.
type LPKeyClaims struct {
    ProductName     string   `json:"product_name"`
    PriceClaims     []string `json:"price_claims"`
    HealthClaims    []string `json:"health_claims"`
    GuaranteeClaims []string `json:"guarantee_claims"`
    HasDisclaimer   bool     `json:"has_disclaimer"`
    DisclaimerText  string   `json:"disclaimer_text"`
    LPSummary       string   `json:"lp_summary"`
    Domain          string   `json:"domain"`
}

// LPAgentOutput combines LP violations and extracted key claims.
type LPAgentOutput struct {
    Violations       []AgentViolation `json:"violations"`
    KeyClaims        LPKeyClaims      `json:"key_claims"`
    LPAnalysisStatus string           `json:"lp_analysis_status"` // "analyzed"|"not_provided"|"fetch_failed"|"js_rendered"
    LPTextExcerpt    string           `json:"-"` // not sent to client; stored internally for future use
}

// AlignmentViolation is a cross-source contradiction finding.
type AlignmentViolation struct {
    Code              string `json:"code"`
    Severity          string `json:"severity"` // always "MODERATE" for ALIGN_ codes
    Explanation       string `json:"explanation"`
    SourceA           string `json:"source_a"`
    SourceB           string `json:"source_b"`
    ClaimInSourceA    string `json:"claim_in_source_a"`
    ClaimInSourceB    string `json:"claim_in_source_b"`
}

// Advisory is a non-violation recommendation.
// Advisories do NOT affect the compliance score.
type Advisory struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// ─────────────────────────────────────────────
// RESPONSE TYPES
// ─────────────────────────────────────────────

// Violation is an enriched policy violation returned in the API response.
type Violation struct {
    Code         string `json:"code"`
    Title        string `json:"title"`
    Severity     string `json:"severity"`
    Explanation  string `json:"explanation"`
    SuggestedFix string `json:"suggested_fix"`
    Source       string `json:"source"` // "copy" | "image_text" | "voiceover" | "landing_page"
}

// CheckResult is the complete API response body.
type CheckResult struct {
    CheckID               string               `json:"check_id"`
    Action                string               `json:"action"`
    // Action values: "ALLOWED" | "PROHIBITED" | "RESTRICTED" |
    //                "AUTHORIZATION_REQUIRED" | "ALLOWED_WITH_RESTRICTIONS"
    ComplianceScore       int                  `json:"compliance_score"`
    RiskCategory          string               `json:"risk_category"` // "High Risk" | "Medium Risk" | "Low Risk"
    Category              string               `json:"category"`
    CategoryState         string               `json:"category_state"`
    DetectedLanguage      string               `json:"detected_language"`
    VisualCheckConfidence string               `json:"visual_check_confidence"` // "high"|"partial"|"low"|""
    LPAnalysisStatus      string               `json:"lp_analysis_status"`      // "analyzed"|"not_provided"|"fetch_failed"|"js_rendered"
    Violations            []Violation          `json:"violations"`
    AlignmentViolations   []AlignmentViolation `json:"alignment_violations"`
    Advisories            []Advisory           `json:"advisories"`
    SuggestedCopy         string               `json:"suggested_copy"`
}
```

---

## Violation Metadata Map

**`backend/internal/models/metadata.go`** — add new v2 violation codes to the existing map. The existing v1 codes remain unchanged. Add these entries:

```go
// Add to ViolationMetadata map (existing v1 entries remain):

// Vision Agent codes (P_ and R_ tier):
"P_ADULT_EXPLICIT":            {"Explicit Adult Content", "PROHIBITED"},
"P_VIOLENCE_GRAPHIC":          {"Graphic Violence", "PROHIBITED"},
"P_WEAPONS_VISUAL":            {"Weapons Display", "PROHIBITED"},
"P_PROFANITY_VISUAL":          {"Profane Visual Content", "PROHIBITED"},
"P_CHILD_SAFETY_VISUAL":       {"Child Safety Violation", "PROHIBITED"},
"R_BEFORE_AFTER_WEIGHT_VISUAL":   {"Before/After Weight Loss Image", "RESTRICTED"},
"R_BEFORE_AFTER_COSMETIC_VISUAL": {"Before/After Cosmetic Image", "RESTRICTED"},
"R_BODY_SHAMING_VISUAL":          {"Body Shaming Visual", "RESTRICTED"},
"R_SKIN_WHITENING_VISUAL":        {"Skin Whitening Visual", "RESTRICTED"},
"R_SEXUAL_WELLNESS_VISUAL":       {"Sexual Wellness Visual", "RESTRICTED"},
"M_SUGGESTIVE_VISUAL":            {"Suggestive Imagery", "MODERATE"},
"M_SHOCK_IMAGERY":                {"Shock Imagery", "MODERATE"},
"M_SCAREWARE_VISUAL":             {"Scareware Visual", "MODERATE"},

// LP Agent codes:
"P_COUNTERFEIT":               {"Counterfeit Goods", "PROHIBITED"},
"P_MISINFORMATION":            {"Dangerous Misinformation", "PROHIBITED"},
"P_HUMAN_EXPLOITATION":        {"Human Exploitation Content", "PROHIBITED"},
"R_HEALTH_CLAIM":              {"Unproven Health or Medical Claim", "RESTRICTED"},
"R_GUARANTEED_OUTCOME":        {"Guaranteed Financial Outcome", "RESTRICTED"},
"R_BEFORE_AFTER_WEIGHT":       {"Before/After Weight Loss Claim", "RESTRICTED"},
"R_FINANCIAL_DISCLOSURE":      {"Missing Financial Disclosure", "RESTRICTED"},

// Alignment Agent codes:
"ALIGN_PRODUCT_MISMATCH":  {"Ad/LP Product Mismatch", "MODERATE"},
"ALIGN_CLAIM_MISMATCH":    {"Ad/LP Claim Contradiction", "MODERATE"},
"ALIGN_PRICE_MISMATCH":    {"Ad/LP Price Contradiction", "MODERATE"},
"ALIGN_VISUAL_COPY_MISMATCH": {"Visual/Copy Contradiction", "MODERATE"},
```

Note: some LP codes (`R_HEALTH_CLAIM`, `R_GUARANTEED_OUTCOME`) already exist in v1 metadata for copy-based violations. These are the same codes — the source field on the `Violation` struct distinguishes where the violation was found. Do not add duplicate entries.

---

## API Contract

### Endpoint

```
POST /checks
```

### Request — two accepted content types

**Option A: multipart/form-data** (when image or video is attached)

The frontend sends `FormData`. All text fields are form values. Files are form file fields.

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `platform` | string | yes | "Meta" |
| `countries` | string (JSON array) | no | e.g. `'["US","GB","TR"]'`. Empty or omitted = global. |
| `age_min` | string (int) | yes | e.g. `"18"` |
| `age_max` | string (int) | yes | e.g. `"65"` |
| `ad_format` | string | yes | `"single"` \| `"carousel"` |
| `primary_text` | string | yes | |
| `headline` | string | no | |
| `description` | string | no | |
| `landing_page_url` | string | no | |
| `force_refresh` | string | no | `"true"` \| `"false"`. Bypasses cache. |
| `image` | file | no | Mutually exclusive with `video`. |
| `video` | file | no | Mutually exclusive with `image`. |

**Option B: application/json** (text-only, no file — backward-compatible)

Same fields as above, countries as actual JSON array, age_min/age_max as integers.

The handler detects the content type and parses accordingly. JSON path used when `Content-Type` is `application/json`. Multipart path used otherwise.

### Response

```json
{
  "check_id": "550e8400-e29b-41d4-a716-446655440000",
  "action": "PROHIBITED",
  "compliance_score": 0,
  "risk_category": "High Risk",
  "category": "Health — Weight Loss (Mild)",
  "category_state": "RESTRICTED",
  "detected_language": "en",
  "visual_check_confidence": "high",
  "lp_analysis_status": "analyzed",
  "violations": [
    {
      "code": "R_WEIGHT_LOSS_CLAIM",
      "title": "Specific Weight Loss Claim",
      "severity": "RESTRICTED",
      "explanation": "Ad copy states 'Lose 20kg in 30 days' — specific numeric claim with timeframe.",
      "suggested_fix": "Replace with: 'Support your weight management goals with our supplement.'",
      "source": "copy"
    },
    {
      "code": "P_ADULT_EXPLICIT",
      "title": "Explicit Adult Content",
      "severity": "PROHIBITED",
      "explanation": "[Found in image] Uploaded image contains explicit adult imagery.",
      "suggested_fix": "Remove explicit imagery. Use product-focused visuals.",
      "source": "image_text"
    }
  ],
  "alignment_violations": [
    {
      "code": "ALIGN_CLAIM_MISMATCH",
      "severity": "MODERATE",
      "explanation": "Ad copy claims FDA approval but landing page includes an FDA disclaimer.",
      "source_a": "ad_copy",
      "source_b": "landing_page",
      "claim_in_source_a": "FDA approved formula",
      "claim_in_source_b": "These statements have not been evaluated by the Food and Drug Administration."
    }
  ],
  "advisories": [
    {
      "code": "ADV_MISSING_GEO_TARGETING",
      "message": "This ad category (Health — Weight Loss) has geo-specific restrictions. No target countries were specified — country-level compliance rules could not be verified."
    }
  ],
  "suggested_copy": "Support your weight management goals with our supplement. Individual results may vary."
}
```

**Action value mapping:**

| Action | Meaning |
|--------|---------|
| `PROHIBITED` | Contains P_ violations — ad cannot run |
| `AUTHORIZATION_REQUIRED` | Contains A_ violations — requires Meta pre-approval |
| `RESTRICTED` | Contains R_ violations — ad can run with fixes |
| `ALLOWED_WITH_RESTRICTIONS` | Only M_ or ALIGN_ violations — ad can run but should be improved |
| `ALLOWED` | No violations — compliant |

Note: v1 used `"not_compliant"` and `"compliant"`. v2 uses the five values above. Update all frontend references accordingly.

**Advisory codes and their human-readable messages:**

| Code | Message |
|------|---------|
| `ADV_MISSING_AGE_TARGETING` | "This ad category requires age-gating (18+ or 21+ depending on region). No minimum age was set above the threshold — targeting may not be compliant." |
| `ADV_MISSING_GEO_TARGETING` | "This ad category has geo-specific restrictions. No target countries were specified — country-level compliance rules could not be verified." |
| `ADV_LP_NOT_PROVIDED` | "No landing page URL was provided. Destination-level violations could not be checked." |
| `ADV_LP_FETCH_FAILED` | "The landing page could not be fetched. Destination-level violations were not checked." |
| `ADV_LP_JS_RENDERED` | "The landing page required JavaScript rendering. Content may have been partially captured." |
| `ADV_VISUAL_PARTIAL` | "Visual analysis confidence was partial or low due to image resolution or composition. Some visual violations may have been missed." |
| `ADV_AUDIO_NOT_CHECKED` | "A video was provided but the audio track could not be analyzed. Voiceover violations were not checked." |
| `ADV_VIDEO_LONG_CLIP` | "Video exceeds the recommended duration for compliance checking. Frame sampling may have missed content in longer segments." |

---

*END OF BLOCK 1*

---

## BLOCK 2 — Pre-processor, Cache, Gemini Client Update, Orchestrator

---

## Updated Gemini Client Helper

**`backend/internal/agents/client.go`** — replace existing file entirely:

```go
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
    temp := float32(0)
    model.SetTemperature(temp)
    model.ResponseMIMEType = "application/json"
    cleanup := func() { client.Close() }
    return model, cleanup, nil
}
```

No other changes to the client helper. Agents that send multipart content (Vision, Audio, Alignment) build their own `[]genai.Part` slices using `genai.ImageData()` and `genai.Blob{}` — they do not need a separate helper function.

---

## Pre-processor

**`backend/internal/preprocessor/preprocessor.go`** — create new file:

```go
package preprocessor

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"

    "compliance-checker/internal/models"
)

const maxVideoSizeMB = 200

// Output holds everything the pre-processor extracted from the HTTP request.
// This is the input to all Phase 1 goroutines.
type Output struct {
    Req         *models.CheckRequest
    ImageBytes  []byte   // nil if no image
    ImageSize   int64
    VideoFrames [][]byte // nil if no video; JPEG bytes per frame
    AudioBytes  []byte   // nil if no video or no audio track found
    VideoSize   int64
    HasImage    bool
    HasVideo    bool
    VideoDurationSeconds float64 // 0 if not a video
}

// Parse reads the HTTP request and extracts all inputs.
// Handles both multipart/form-data (when a file is present) and
// application/json (text-only, backward-compatible).
func Parse(r *http.Request) (*Output, error) {
    out := &Output{Req: &models.CheckRequest{}}

    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "application/json") {
        // JSON path — text-only, no file
        if err := parseJSON(r, out.Req); err != nil {
            return nil, fmt.Errorf("json parse: %w", err)
        }
        return out, nil
    }

    // Multipart path
    if err := r.ParseMultipartForm(maxVideoSizeMB * 1024 * 1024); err != nil {
        return nil, fmt.Errorf("multipart parse: %w", err)
    }
    if err := parseFormFields(r, out.Req); err != nil {
        return nil, fmt.Errorf("form fields: %w", err)
    }

    // Check for image file
    imageFile, imageHeader, err := r.FormFile("image")
    if err == nil {
        defer imageFile.Close()
        out.ImageBytes, err = io.ReadAll(imageFile)
        if err != nil {
            return nil, fmt.Errorf("read image: %w", err)
        }
        out.ImageSize = imageHeader.Size
        out.HasImage = true
    }

    // Check for video file
    videoFile, videoHeader, err := r.FormFile("video")
    if err == nil {
        defer videoFile.Close()

        if videoHeader.Size > maxVideoSizeMB*1024*1024 {
            return nil, fmt.Errorf("video exceeds %dMB limit", maxVideoSizeMB)
        }

        videoBytes, err := io.ReadAll(videoFile)
        if err != nil {
            return nil, fmt.Errorf("read video: %w", err)
        }

        // Detect file extension from filename for ffmpeg
        ext := filepath.Ext(videoHeader.Filename)
        if ext == "" {
            ext = ".mp4" // safe default
        }

        frames, audioBytes, duration, err := extractVideoContent(videoBytes, ext)
        if err != nil {
            return nil, fmt.Errorf("video extraction: %w", err)
        }

        out.VideoFrames = frames
        out.AudioBytes = audioBytes
        out.VideoSize = videoHeader.Size
        out.VideoDurationSeconds = duration
        out.HasVideo = true
    }

    return out, nil
}

// parseFormFields reads all text fields from the multipart form into CheckRequest.
func parseFormFields(r *http.Request, req *models.CheckRequest) error {
    req.Platform = r.FormValue("platform")
    req.PrimaryText = r.FormValue("primary_text")
    req.Headline = r.FormValue("headline")
    req.Description = r.FormValue("description")
    req.LandingPageURL = r.FormValue("landing_page_url")
    req.ForceRefresh = r.FormValue("force_refresh") == "true"
    req.AdFormat = r.FormValue("ad_format")

    ageMin, _ := strconv.Atoi(r.FormValue("age_min"))
    ageMax, _ := strconv.Atoi(r.FormValue("age_max"))
    req.AgeMin = ageMin
    req.AgeMax = ageMax

    // Countries come as a JSON array string: '["US","GB","TR"]'
    // Parse it; on failure treat as global (empty slice)
    countriesRaw := r.FormValue("countries")
    if countriesRaw != "" && countriesRaw != "[]" {
        // Strip brackets and split by comma — simple parse, avoids json import cycle
        trimmed := strings.Trim(countriesRaw, "[]")
        for _, c := range strings.Split(trimmed, ",") {
            c = strings.Trim(strings.TrimSpace(c), `"`)
            if c != "" {
                req.Countries = append(req.Countries, c)
            }
        }
    }

    return nil
}

// parseJSON reads a JSON body into CheckRequest.
// encoding/json is imported at the package level above.
func parseJSON(r *http.Request, req *models.CheckRequest) error {
    return json.NewDecoder(r.Body).Decode(req)
}

// extractVideoContent runs ffmpeg to extract frames and audio from a video file.
// Returns: JPEG frame bytes (capped at 30), MP3 audio bytes (nil if no audio), duration in seconds.
//
// Frame extraction strategy:
//   Step 1 — scene-change detection (threshold 0.4). Captures frames at visual transitions.
//   Step 2 — if fewer than 3 frames detected, fall back to uniform 1fps sampling.
//   Cap at 30 frames in either case.
//
// Flexibility note: if ffmpeg is unavailable or performance is inadequate, this function
// can be replaced with Gemini API native video processing (pass raw video bytes directly
// to the Vision Agent as a genai.Blob with the video MIME type). The agent interface
// (VideoFrames [][]byte) stays the same — simply pass one entry containing the full
// video bytes and update the Vision Agent to handle it. No orchestrator changes needed.
func extractVideoContent(videoBytes []byte, ext string) (frames [][]byte, audioBytes []byte, duration float64, err error) {
    tmpDir, err := os.MkdirTemp("", "compliance-video-*")
    if err != nil {
        return nil, nil, 0, err
    }
    defer os.RemoveAll(tmpDir)

    // Write video bytes to temp file
    videoPath := filepath.Join(tmpDir, "input"+ext)
    if err := os.WriteFile(videoPath, videoBytes, 0644); err != nil {
        return nil, nil, 0, err
    }

    // Get video duration via ffprobe
    duration = getVideoDuration(videoPath)

    // Step 1: Scene-change frame extraction
    framesDir := filepath.Join(tmpDir, "frames")
    os.MkdirAll(framesDir, 0755)

    cmd := exec.Command("ffmpeg",
        "-i", videoPath,
        "-vf", "select='gt(scene,0.4)'",
        "-vsync", "vfr",
        "-q:v", "2",
        "-frames:v", "30", // hard cap
        filepath.Join(framesDir, "frame_%03d.jpg"),
    )
    cmd.Run() // error is non-fatal — check frame count below

    frameFiles, _ := filepath.Glob(filepath.Join(framesDir, "*.jpg"))

    // Step 2: Fallback to uniform 1fps if fewer than 3 frames
    if len(frameFiles) < 3 {
        os.RemoveAll(framesDir)
        os.MkdirAll(framesDir, 0755)

        maxFrames := int(duration) // 1fps for duration gives N frames
        if maxFrames > 30 {
            maxFrames = 30
        }
        if maxFrames < 1 {
            maxFrames = 10 // minimum attempt
        }

        cmd = exec.Command("ffmpeg",
            "-i", videoPath,
            "-r", "1",
            "-q:v", "2",
            "-frames:v", strconv.Itoa(maxFrames),
            filepath.Join(framesDir, "frame_%03d.jpg"),
        )
        cmd.Run()
        frameFiles, _ = filepath.Glob(filepath.Join(framesDir, "*.jpg"))
    }

    // Cap at 30 frames
    if len(frameFiles) > 30 {
        frameFiles = frameFiles[:30]
    }

    for _, f := range frameFiles {
        b, readErr := os.ReadFile(f)
        if readErr != nil {
            continue
        }
        frames = append(frames, b)
    }

    // Extract audio track as MP3
    audioPath := filepath.Join(tmpDir, "audio.mp3")
    audioCmd := exec.Command("ffmpeg",
        "-i", videoPath,
        "-q:a", "0",
        "-map", "a",
        audioPath,
    )
    if audioErr := audioCmd.Run(); audioErr == nil {
        audioBytes, _ = os.ReadFile(audioPath)
    }
    // Audio extraction failure is non-fatal (video may have no audio track)

    return frames, audioBytes, duration, nil
}

// getVideoDuration returns the duration in seconds using ffprobe.
// Returns 0 on failure — callers treat 0 as unknown.
func getVideoDuration(videoPath string) float64 {
    cmd := exec.Command("ffprobe",
        "-v", "error",
        "-show_entries", "format=duration",
        "-of", "csv=p=0",
        videoPath,
    )
    out, err := cmd.Output()
    if err != nil {
        return 0
    }
    d, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
    return d
}
```

---

## Cache Layer

**`backend/internal/cache/cache.go`** — create new file:

```go
package cache

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "sort"
    "sync"
    "time"

    "compliance-checker/internal/models"
)

const TTL = 24 * time.Hour

type entry struct {
    result    models.CheckResult
    expiresAt time.Time
}

var (
    mu    sync.RWMutex
    store = make(map[string]*entry)
)

// BuildKey creates a deterministic SHA-256 cache key from request inputs.
// File identity is represented by byte length (imageSize, videoSize) — fast and
// sufficient for MVP. force_refresh is NOT part of the key; callers skip cache
// lookup when force_refresh is true.
func BuildKey(req *models.CheckRequest, imageSize int64, videoSize int64) string {
    countries := make([]string, len(req.Countries))
    copy(countries, req.Countries)
    sort.Strings(countries) // normalise order

    raw := fmt.Sprintf("%s|%d|%d|%s|%s|%s|%s|%v|%d|%d",
        req.PrimaryText,
        req.AgeMin,
        req.AgeMax,
        req.LandingPageURL,
        req.Headline,
        req.Description,
        countries,
        req.Platform,
        imageSize,
        videoSize,
    )
    h := sha256.Sum256([]byte(raw))
    return hex.EncodeToString(h[:])
}

// Get returns a cached result if it exists and has not expired.
func Get(key string) (*models.CheckResult, bool) {
    mu.RLock()
    defer mu.RUnlock()
    e, ok := store[key]
    if !ok || time.Now().After(e.expiresAt) {
        return nil, false
    }
    result := e.result
    return &result, true
}

// Set stores a result under the given key with a 24-hour TTL.
func Set(key string, result models.CheckResult) {
    mu.Lock()
    defer mu.Unlock()
    store[key] = &entry{result: result, expiresAt: time.Now().Add(TTL)}
}
```

---

## Orchestrator — Full Handler Rewrite

**`backend/internal/api/handler.go`** — full replacement:

```go
package api

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"

    "compliance-checker/internal/advisory"
    "compliance-checker/internal/agents"
    "compliance-checker/internal/cache"
    "compliance-checker/internal/db"
    "compliance-checker/internal/models"
    "compliance-checker/internal/preprocessor"
    "compliance-checker/internal/scoring"
)

func HandleCheck(w http.ResponseWriter, r *http.Request) {
    // CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    w.Header().Set("Content-Type", "application/json")

    // ─────────────────────────────────────────
    // STEP 1: Parse request via pre-processor
    // ─────────────────────────────────────────
    ppOut, err := preprocessor.Parse(r)
    if err != nil {
        log.Printf("preprocessor error: %v", err)
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }
    req := ppOut.Req

    // Build full copy string from all copy fields (used by text-based agents)
    fullCopy := req.PrimaryText
    if req.Headline != "" {
        fullCopy += "\nHeadline: " + req.Headline
    }
    if req.Description != "" {
        fullCopy += "\nDescription: " + req.Description
    }

    ctx := r.Context()

    // ─────────────────────────────────────────
    // STEP 2: Cache check
    // ─────────────────────────────────────────
    cacheKey := cache.BuildKey(req, ppOut.ImageSize, ppOut.VideoSize)
    if !req.ForceRefresh {
        if cached, ok := cache.Get(cacheKey); ok {
            log.Printf("cache hit for key %s", cacheKey[:8])
            json.NewEncoder(w).Encode(cached)
            return
        }
    }

    // ─────────────────────────────────────────
    // STEP 3: Pipeline context
    // ─────────────────────────────────────────
    // RequestContext is the shared mutable store passed between phases.
    // Phase 1 goroutines write to non-overlapping fields — no mutex needed on ctx fields.
    // allAgentViolations is shared across goroutines — protected by mu.
    rctx := &models.RequestContext{
        ImageBytes:  ppOut.ImageBytes,
        VideoFrames: ppOut.VideoFrames,
        AudioBytes:  ppOut.AudioBytes,
        HasImage:    ppOut.HasImage,
        HasVideo:    ppOut.HasVideo,
    }

    var (
        mu                 sync.Mutex
        allAgentViolations []models.AgentViolation
    )

    // Helper: run an agent goroutine, collect violations, log errors without crashing.
    // Each call to runAgent decrements the provided WaitGroup on completion.
    makeRunner := func(wg *sync.WaitGroup) func(fn func() ([]models.AgentViolation, error), name string) {
        return func(fn func() ([]models.AgentViolation, error), name string) {
            defer wg.Done()
            v, err := fn()
            if err != nil {
                log.Printf("agent [%s] failed: %v", name, err)
                return
            }
            mu.Lock()
            allAgentViolations = append(allAgentViolations, v...)
            mu.Unlock()
        }
    }

    // ─────────────────────────────────────────
    // STEP 4: Sequential — Category Classifier
    // ─────────────────────────────────────────
    categoryResult, err := agents.RunCategoryClassifier(ctx, fullCopy)
    if err != nil {
        log.Printf("category classifier error: %v", err)
        http.Error(w, `{"error":"classification failed"}`, http.StatusInternalServerError)
        return
    }
    rctx.Category = categoryResult.Category
    rctx.CategoryState = categoryResult.CategoryState

    // Fast exit for PROHIBITED category — no parallel agents needed
    if categoryResult.CategoryState == "PROHIBITED" {
        allAgentViolations = append(allAgentViolations, models.AgentViolation{
            Code:        "P_PROHIBITED_CATEGORY",
            Explanation: "This ad category is entirely prohibited on Meta: " + categoryResult.Category,
            Source:      "copy",
        })
        finaliseAndRespond(w, req, rctx, allAgentViolations, nil, cacheKey, ppOut, categoryResult)
        return
    }

    // ─────────────────────────────────────────
    // STEP 5: Phase 1 — Parallel agents
    // ─────────────────────────────────────────
    // Goroutines fire conditionally based on what inputs are available.
    // All goroutines write to rctx (non-overlapping fields) or allAgentViolations (mutex-protected).

    var wg1 sync.WaitGroup
    run1 := makeRunner(&wg1)

    // g1 — Violations Agent: copy-based policy check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, error) {
        return agents.RunViolationsAgent(ctx, fullCopy, rctx.Category, "", "")
    }, "violations")

    // g2 — Personal Attribute Agent: copy-based PA check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, error) {
        return agents.RunPersonalAttributeAgent(ctx, fullCopy, "", "")
    }, "personal_attribute")

    // g3 — Auth Required Agent: copy-based authorization check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, error) {
        return agents.RunAuthRequiredAgent(ctx, fullCopy, rctx.Category)
    }, "auth_required")

    // g4 — Vision Agent: image or video frame analysis + OCR
    // Only fires when image or video frames are present.
    if rctx.HasImage || rctx.HasVideo {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            visionOut, err := agents.RunVisionAgent(ctx, rctx.ImageBytes, rctx.VideoFrames)
            if err != nil {
                log.Printf("agent [vision] failed: %v", err)
                return
            }
            // Write to rctx — Vision Agent is the only writer to these fields
            rctx.ImageExtractedText = visionOut.ExtractedText
            rctx.VisualCheckConfidence = visionOut.VisualCheckConfidence
            if rctx.DetectedLanguage == "" {
                rctx.DetectedLanguage = visionOut.DetectedLanguage
            }
            // Collect visual violations
            mu.Lock()
            allAgentViolations = append(allAgentViolations, visionOut.VisualViolations...)
            mu.Unlock()
        }()
    }

    // g5 — LP Agent: landing page fetch (via Firecrawl) + analysis
    // Only fires when a landing page URL is provided.
    if req.LandingPageURL != "" {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            lpOut, err := agents.RunLPAgent(ctx, req.LandingPageURL)
            if err != nil {
                log.Printf("agent [lp] failed: %v", err)
                rctx.LPAnalysisStatus = "fetch_failed"
                return
            }
            // Write to rctx — LP Agent is the only writer to these fields
            // Note: LPTextExcerpt stays on lpOut only; rctx carries structured claims
            rctx.LPKeyClaims = &lpOut.KeyClaims
            rctx.LPAnalysisStatus = lpOut.LPAnalysisStatus
            // Collect LP violations
            mu.Lock()
            allAgentViolations = append(allAgentViolations, lpOut.Violations...)
            mu.Unlock()
        }()
    }

    // g6 — Audio Agent: transcription of video audio track
    // Only fires when a video was uploaded AND an audio track was extracted.
    if rctx.HasVideo && len(rctx.AudioBytes) > 0 {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            audioOut, err := agents.RunAudioAgent(ctx, rctx.AudioBytes)
            if err != nil {
                log.Printf("agent [audio] failed: %v", err)
                return
            }
            // Write to rctx — Audio Agent is the only writer to AudioTranscript
            rctx.AudioTranscript = audioOut.Transcript
            if rctx.DetectedLanguage == "" {
                rctx.DetectedLanguage = audioOut.DetectedLanguage
            }
            rctx.AudioAnalysisStatus = audioOut.AudioAnalysisStatus
        }()
    }

    wg1.Wait() // Block until all Phase 1 goroutines complete

    // ─────────────────────────────────────────
    // STEP 6: Phase 1b — Supplementary violations pass
    // ─────────────────────────────────────────
    // After Phase 1 completes, rctx now has ImageExtractedText and AudioTranscript.
    // Run Violations Agent and Personal Attribute Agent again with media-sourced text.
    // Near-zero latency penalty — LP Agent (~8s) was the Phase 1 bottleneck.
    if rctx.ImageExtractedText != "" || rctx.AudioTranscript != "" {
        var wg1b sync.WaitGroup
        run1b := makeRunner(&wg1b)

        wg1b.Add(1)
        go run1b(func() ([]models.AgentViolation, error) {
            return agents.RunViolationsAgent(ctx, fullCopy, rctx.Category, rctx.ImageExtractedText, rctx.AudioTranscript)
        }, "violations_1b")

        wg1b.Add(1)
        go run1b(func() ([]models.AgentViolation, error) {
            return agents.RunPersonalAttributeAgent(ctx, fullCopy, rctx.ImageExtractedText, rctx.AudioTranscript)
        }, "personal_attribute_1b")

        wg1b.Wait()
    }

    // ─────────────────────────────────────────
    // STEP 7: Phase 2 — Alignment Agent (conditional)
    // ─────────────────────────────────────────
    // Fires only when both a visual creative AND a landing page were provided.
    // Receives raw image bytes (NOT a text description) to prevent lossy summarisation.
    var alignmentViolations []models.AlignmentViolation
    if (rctx.HasImage || rctx.HasVideo) && rctx.LPKeyClaims != nil {
        existingCodes := violationCodes(allAgentViolations)

        // Select visual reference for Alignment Agent.
        // Image ads: use the original image bytes.
        // Video ads: use the first extracted frame as the representative visual.
        var alignImageBytes []byte
        if rctx.HasImage {
            alignImageBytes = rctx.ImageBytes
        } else if rctx.HasVideo && len(rctx.VideoFrames) > 0 {
            alignImageBytes = rctx.VideoFrames[0]
        }

        alignmentViolations, err = agents.RunAlignmentAgent(
            ctx,
            alignImageBytes,
            fullCopy,
            rctx.LPKeyClaims,
            existingCodes,
        )
        if err != nil {
            log.Printf("agent [alignment] failed: %v", err)
            // Non-fatal — continue without alignment results
        }
    }

    finaliseAndRespond(w, req, rctx, allAgentViolations, alignmentViolations, cacheKey, ppOut, categoryResult)
}

// finaliseAndRespond handles deduplication, scoring, fix generation,
// advisory generation, enrichment, caching, DB save, and response encoding.
// Extracted to a shared function so the PROHIBITED fast-exit path and the
// normal path both produce the same response structure.
func finaliseAndRespond(
    w http.ResponseWriter,
    req *models.CheckRequest,
    rctx *models.RequestContext,
    rawViolations []models.AgentViolation,
    alignmentViolations []models.AlignmentViolation,
    cacheKey string,
    ppOut *preprocessor.Output,
    categoryResult *models.CategoryOutput,
) {
    ctx := context.Background() // use background ctx for post-pipeline work

    // ─── Deduplication ───
    // Keep first occurrence of each violation code.
    // Exception: same code from different sources (copy vs image_text vs voiceover)
    // are kept as separate entries — they represent distinct problems.
    seen := map[string]bool{}
    var uniqueViolations []models.AgentViolation
    for _, v := range rawViolations {
        dedupeKey := v.Code + "|" + v.Source
        if !seen[dedupeKey] {
            seen[dedupeKey] = true
            uniqueViolations = append(uniqueViolations, v)
        }
    }

    // ─── Scoring Engine ───
    codes := violationCodes(uniqueViolations)
    alignCodes := make([]string, len(alignmentViolations))
    for i, a := range alignmentViolations {
        alignCodes[i] = a.Code
    }
    allCodesForScoring := append(codes, alignCodes...)
    scoreResult := scoring.Score(allCodesForScoring, categoryResult.CategoryState, false)

    // ─── Fix Generation ───
    fixes, err := agents.RunFixGenerationAgent(ctx, req.PrimaryText, categoryResult.Category, uniqueViolations)
    if err != nil {
        log.Printf("fix generation error: %v", err)
        // Non-fatal
    }

    // ─── Suggested copy: top-level rewrite if violations exist ───
    suggestedCopy := ""
    if fixes != nil {
        suggestedCopy = fixes["_suggested_copy"] // Fix Generation Agent outputs this key
    }

    // ─── Enrich violations with metadata and fixes ───
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
            Source:       v.Source,
        })
    }
    if enrichedViolations == nil {
        enrichedViolations = []models.Violation{}
    }
    if alignmentViolations == nil {
        alignmentViolations = []models.AlignmentViolation{}
    }

    // ─── Advisory Generation ───
    // Pure deterministic Go logic — no LLM. See advisory package.
    advisories := advisory.Generate(req, rctx, ppOut.VideoDurationSeconds)

    // ─── Build result ───
    maxVideoDuration, _ := strconv.ParseFloat(os.Getenv("MAX_VIDEO_DURATION_SECONDS"), 64)
    if maxVideoDuration == 0 {
        maxVideoDuration = 120
    }

    result := models.CheckResult{
        Action:                scoreResult.Action,
        ComplianceScore:       scoreResult.ComplianceScore,
        RiskCategory:          scoreResult.RiskCategory,
        Category:              categoryResult.Category,
        CategoryState:         categoryResult.CategoryState,
        DetectedLanguage:      rctx.DetectedLanguage,
        VisualCheckConfidence: rctx.VisualCheckConfidence,
        LPAnalysisStatus:      lpAnalysisStatus(req, rctx),
        Violations:            enrichedViolations,
        AlignmentViolations:   alignmentViolations,
        Advisories:            advisories,
        SuggestedCopy:         suggestedCopy,
    }

    // ─── Cache and persist ───
    cache.Set(cacheKey, result)
    saveCheckToDB(cacheKey, req.PrimaryText, result)

    json.NewEncoder(w).Encode(result)
}

// violationCodes extracts just the code strings from an AgentViolation slice.
func violationCodes(violations []models.AgentViolation) []string {
    codes := make([]string, len(violations))
    for i, v := range violations {
        codes[i] = v.Code
    }
    return codes
}

// lpAnalysisStatus returns the correct status string based on what happened with LP fetching.
func lpAnalysisStatus(req *models.CheckRequest, rctx *models.RequestContext) string {
    if req.LandingPageURL == "" {
        return "not_provided"
    }
    if rctx.LPAnalysisStatus != "" {
        return rctx.LPAnalysisStatus
    }
    return "analyzed"
}

// saveCheckToDB persists the result to PostgreSQL.
// Non-fatal on error — compliance checking proceeds regardless of DB availability.
func saveCheckToDB(contentHash string, inputCopy string, result models.CheckResult) {
    resultJSON, err := json.Marshal(result)
    if err != nil {
        log.Printf("db marshal error: %v", err)
        return
    }
    _, err = db.DB.Exec(
        `INSERT INTO checks (input_copy, content_hash, result) VALUES ($1, $2, $3)`,
        inputCopy, contentHash, resultJSON,
    )
    if err != nil {
        log.Printf("db insert error: %v", err)
    }
}
```

**Important implementation notes for the orchestrator:**

1. `rctx.DetectedLanguage` may be written by both Vision Agent and Audio Agent. Use the "first writer wins" pattern shown above — check `if rctx.DetectedLanguage == ""` before writing. Both agents run concurrently but the language detection for the same ad should agree.

2. The `finaliseAndRespond` function uses `context.Background()` for post-pipeline work (fix generation, DB save). This is intentional — the HTTP request context may have already timed out by this point. Fix generation and DB writes must not be abandoned.

3. `fixes["_suggested_copy"]` — the Fix Generation Agent will be updated in Block 3 to also output a top-level `_suggested_copy` key with a complete rewritten primary text. Per-violation fixes remain in `fixes[violationCode]`.

4. `rctx.VisualCheckConfidence` — populated by Vision Agent during Phase 1. If no image/video was provided, this remains empty string, which is correct and maps to the empty string in the API response.

5. The Alignment Agent receives `rctx.ImageBytes` (the original image) — NOT the video frames. For video input, pass the first frame (`rctx.VideoFrames[0]`) as the image. The Alignment Agent's job is copy/visual vs landing page cross-check, not full video analysis.

---

*END OF BLOCK 2*

---

## BLOCK 3 — Agent Updates, New Agents, Scoring Engine, Advisory System

---

## Note on Types

The `RequestContext` and `LPAgentOutput` definitions in Block 1 are authoritative and complete — they include all fields referenced in the handler and agents. Implement from Block 1 directly; no corrections needed.

---

## Updated Agent: violations.go

**`backend/internal/agents/violations.go`** — update the function signature and input builder:

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

type violationsOutput struct {
    Violations []models.AgentViolation `json:"violations"`
}

// RunViolationsAgent checks for P_, R_, and M_ violations in ad copy.
// In Phase 1 (first call): imageExtractedText and audioTranscript are empty strings.
// In Phase 1b (second call): they contain OCR output and voiceover transcript.
// The agent prompt handles both cases — it only checks media text if provided.
func RunViolationsAgent(ctx context.Context, copy string, category string, imageExtractedText string, audioTranscript string) ([]models.AgentViolation, error) {
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
        return nil, err
    }

    raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
    var out violationsOutput
    if err := json.Unmarshal([]byte(raw), &out); err != nil {
        return nil, fmt.Errorf("failed to parse violations output: %w", err)
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

    return out.Violations, nil
}
```

---

## Updated Agent: personal_attribute.go

**`backend/internal/agents/personal_attribute.go`** — same pattern as violations update:

```go
// RunPersonalAttributeAgent checks for P_PERSONAL_ATTRIBUTE violations.
// imageExtractedText and audioTranscript are empty in Phase 1, populated in Phase 1b.
func RunPersonalAttributeAgent(ctx context.Context, copy string, imageExtractedText string, audioTranscript string) ([]models.AgentViolation, error) {
    // ... setup identical to v1 ...

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
    // ... parse PersonalAttributeOutput same as v1 ...

    if out.ViolationFound {
        source := "copy"
        if strings.Contains(out.Explanation, "[Found in image text]") {
            source = "image_text"
        } else if strings.Contains(out.Explanation, "[Found in voiceover]") {
            source = "voiceover"
        }
        return []models.AgentViolation{{Code: "P_PERSONAL_ATTRIBUTE", Explanation: out.Explanation, Source: source}}, nil
    }
    return nil, nil
}
```

---

## New Agent: vision.go

**`backend/internal/agents/vision.go`** — create new file:

```go
package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/google/generative-ai-go/genai"
    "compliance-checker/internal/models"
)

// RunVisionAgent analyzes ad creative for OCR text extraction and visual policy violations.
// Accepts either raw image bytes OR a slice of video frame bytes (JPEG).
// For video: all frames are sent in a single Gemini call as sequential image parts.
// The agent prompt instructs Gemini to note the frame number when flagging video violations.
func RunVisionAgent(ctx context.Context, imageBytes []byte, videoFrames [][]byte) (*models.VisionAgentOutput, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, err
    }
    defer cleanup()

    model.ResponseSchema = &genai.Schema{
        Type: genai.TypeObject,
        Properties: map[string]*genai.Schema{
            "extracted_text": {Type: genai.TypeString},
            "visual_violations": {
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
            "visual_check_confidence": {Type: genai.TypeString},
            "detected_language":       {Type: genai.TypeString},
        },
        Required: []string{"extracted_text", "visual_violations", "visual_check_confidence", "detected_language"},
    }

    systemPrompt, err := os.ReadFile("prompts/vision_agent.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to read vision prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    // Build multipart content parts
    var parts []genai.Part
    parts = append(parts, genai.Text("Analyze the following ad creative:"))

    if len(imageBytes) > 0 {
        mimeType := http.DetectContentType(imageBytes)
        parts = append(parts, genai.ImageData(mimeType, imageBytes))
    }

    for i, frame := range videoFrames {
        parts = append(parts, genai.Text(fmt.Sprintf("[Video frame %d of %d]", i+1, len(videoFrames))))
        parts = append(parts, genai.ImageData("image/jpeg", frame))
    }

    resp, err := model.GenerateContent(ctx, parts...)
    if err != nil {
        return nil, err
    }

    raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
    var out models.VisionAgentOutput
    if err := json.Unmarshal([]byte(raw), &out); err != nil {
        return nil, fmt.Errorf("failed to parse vision output: %w", err)
    }

    // Set source field on all visual violations
    for i := range out.VisualViolations {
        out.VisualViolations[i].Source = "image_text"
    }

    return &out, nil
}
```

---

## New Agent: audio.go

**`backend/internal/agents/audio.go`** — create new file:

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

// RunAudioAgent transcribes the audio track of a video ad.
// Transcription only — no violation checking. Violation checking on the
// transcript happens in Phase 1b via RunViolationsAgent and RunPersonalAttributeAgent.
func RunAudioAgent(ctx context.Context, audioBytes []byte) (*models.AudioAgentOutput, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, err
    }
    defer cleanup()

    model.ResponseSchema = &genai.Schema{
        Type: genai.TypeObject,
        Properties: map[string]*genai.Schema{
            "transcript":            {Type: genai.TypeString},
            "detected_language":     {Type: genai.TypeString},
            "audio_analysis_status": {Type: genai.TypeString},
            "has_voiceover":         {Type: genai.TypeBoolean},
            "duration_note":         {Type: genai.TypeString},
        },
        Required: []string{"transcript", "detected_language", "audio_analysis_status", "has_voiceover"},
    }

    systemPrompt, err := os.ReadFile("prompts/audio_agent.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to read audio prompt: %w", err)
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
        return nil, err
    }

    raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
    var out models.AudioAgentOutput
    if err := json.Unmarshal([]byte(raw), &out); err != nil {
        return nil, fmt.Errorf("failed to parse audio output: %w", err)
    }

    return &out, nil
}
```

---

## New Agent: lp.go

**`backend/internal/agents/lp.go`** — create new file:

```go
package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

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
func RunLPAgent(ctx context.Context, url string) (*models.LPAgentOutput, error) {
    // Step 1: Fetch landing page via Firecrawl with 8-second timeout
    fetchCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
    defer cancel()

    lpContent, jsRendered, fetchErr := fetchLandingPage(fetchCtx, url)
    if fetchErr != nil {
        return &models.LPAgentOutput{LPAnalysisStatus: "fetch_failed"}, nil
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
            "key_claims": {
                Type: genai.TypeObject,
                Properties: map[string]*genai.Schema{
                    "product_name":      {Type: genai.TypeString},
                    "price_claims":      {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
                    "health_claims":     {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
                    "guarantee_claims":  {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
                    "has_disclaimer":    {Type: genai.TypeBoolean},
                    "disclaimer_text":   {Type: genai.TypeString},
                    "lp_summary":        {Type: genai.TypeString},
                    "domain":            {Type: genai.TypeString},
                },
                Required: []string{"product_name", "has_disclaimer", "lp_summary", "domain"},
            },
        },
        Required: []string{"violations", "key_claims"},
    }

    systemPrompt, err := os.ReadFile("prompts/lp_agent.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to read lp prompt: %w", err)
    }
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(string(systemPrompt))},
    }

    input := fmt.Sprintf("Landing page URL: %s\n\nLanding page content:\n%s", url, lpExcerpt)
    resp, err := model.GenerateContent(ctx, genai.Text(input))
    if err != nil {
        return nil, err
    }

    raw := resp.Candidates[0].Content.Parts[0].(genai.Text)
    var out models.LPAgentOutput
    if err := json.Unmarshal([]byte(raw), &out); err != nil {
        return nil, fmt.Errorf("failed to parse lp output: %w", err)
    }

    out.LPAnalysisStatus = lpStatus
    out.LPTextExcerpt = lpExcerpt

    // Tag all LP violations with source = "landing_page"
    for i := range out.Violations {
        out.Violations[i].Source = "landing_page"
    }

    return &out, nil
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
        if scrapeResult.Markdown != nil {
            content = *scrapeResult.Markdown
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
```

**Add this field to `models.LPAgentOutput`** (Block 1 definition is missing it):

```go
type LPAgentOutput struct {
    Violations       []AgentViolation `json:"violations"`
    KeyClaims        LPKeyClaims      `json:"key_claims"`
    LPAnalysisStatus string           `json:"lp_analysis_status"`
    LPTextExcerpt    string           `json:"-"` // internal use only, not serialised
}
```

---

## New Agent: alignment.go

**`backend/internal/agents/alignment.go`** — create new file:

```go
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
) ([]models.AlignmentViolation, error) {
    model, cleanup, err := newModel(ctx)
    if err != nil {
        return nil, err
    }
    defer cleanup()

    model.ResponseSchema = &genai.Schema{
        Type: genai.TypeObject,
        Properties: map[string]*genai.Schema{
            "alignment_violations": {
                Type: genai.TypeArray,
                Items: &genai.Schema{
                    Type: genai.TypeObject,
                    Properties: map[string]*genai.Schema{
                        "code":               {Type: genai.TypeString},
                        "severity":           {Type: genai.TypeString},
                        "explanation":        {Type: genai.TypeString},
                        "source_a":           {Type: genai.TypeString},
                        "source_b":           {Type: genai.TypeString},
                        "claim_in_source_a":  {Type: genai.TypeString},
                        "claim_in_source_b":  {Type: genai.TypeString},
                    },
                    Required: []string{"code", "severity", "explanation", "source_a", "source_b", "claim_in_source_a", "claim_in_source_b"},
                },
            },
        },
        Required: []string{"alignment_violations"},
    }

    systemPrompt, err := os.ReadFile("prompts/alignment_agent.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to read alignment prompt: %w", err)
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
        parts = append(parts, genai.ImageData(mimeType, imageBytes))
    }
    parts = append(parts, genai.Text(textInput))

    resp, err := model.GenerateContent(ctx, parts...)
    if err != nil {
        return nil, err
    }

    raw := resp.Candidates[0].Content.Parts[0].(genai.Text)

    type alignmentOutput struct {
        AlignmentViolations []models.AlignmentViolation `json:"alignment_violations"`
    }
    var out alignmentOutput
    if err := json.Unmarshal([]byte(raw), &out); err != nil {
        return nil, fmt.Errorf("failed to parse alignment output: %w", err)
    }

    return out.AlignmentViolations, nil
}
```

---

## Updated Agent: fix_generation.go

**`backend/internal/agents/fix_generation.go`** — two changes only:

**Change 1:** Add `_suggested_copy` to the `fixFrameworks` map:

```go
"_suggested_copy": "Rewrite the entire primary text of the ad to remove all violations listed. Preserve the core product message and call-to-action. Output only the rewritten copy — no explanations.",
```

**Change 2:** Add `_suggested_copy` to the batched input so the agent generates it. In the `sb` builder section, after listing all violations, append:

```go
sb.WriteString("\n\nAlso output a '_suggested_copy' entry: a complete rewrite of the original primary text that removes ALL violations above. Keep the core message and CTA. Code: '_suggested_copy'.")
```

No other changes to fix_generation.go. The existing `result[f.Code] = f.SuggestedFix` loop will naturally store the `_suggested_copy` entry in the map, which the handler reads as `fixes["_suggested_copy"]`.

---

## Updated Scoring Engine

**`backend/internal/scoring/engine.go`** — full replacement:

```go
package scoring

import "strings"

type Result struct {
    ComplianceScore int
    RiskCategory    string
    Action          string
}

// Score applies the compliance scoring waterfall.
// Rules are evaluated in strict priority order — first match wins.
// ALIGN_ codes are treated as M_ tier (Rule 4): moderate severity, fixable.
// Advisories are NOT passed to this function — they never affect the score.
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
        case strings.HasPrefix(v, "M_") || strings.HasPrefix(v, "ALIGN_"):
            m = append(m, v)
        }
    }

    // Rule 1: Any prohibited violation → zero score, cannot run
    if len(p) > 0 {
        return Result{0, "High Risk", "PROHIBITED"}
    }

    // Rule 2: Unauthorized auth-required category → zero score, needs Meta approval
    if len(a) > 0 && !authorized {
        return Result{0, "High Risk", "AUTHORIZATION_REQUIRED"}
    }

    // Rule 3: Restricted violations → fixable, score 1–50
    if len(r) > 0 {
        score := 50 - (10 * (len(r) - 1)) - (5 * len(m))
        if score < 1 {
            score = 1
        }
        risk := "Medium Risk"
        if score < 30 {
            risk = "High Risk"
        }
        return Result{score, risk, "RESTRICTED"}
    }

    // Rule 4: Moderate or alignment violations only → can run, should improve
    if len(m) > 0 {
        score := 75 - (5 * (len(m) - 1))
        if score < 45 {
            score = 45
        }
        risk := "Low Risk"
        if score < 65 {
            risk = "Medium Risk"
        }
        return Result{score, risk, "ALLOWED_WITH_RESTRICTIONS"}
    }

    // Rule 5: Clean — standard allowed category
    if categoryState == "ALLOWED" {
        return Result{95, "Low Risk", "ALLOWED"}
    }

    // Rule 6: Clean — sensitive or authorization-required category with no violations
    return Result{78, "Low Risk", "ALLOWED"}
}
```

**Update `backend/internal/scoring/engine_test.go`** — add test cases for:
- R_ + ALIGN_ combination (should be RESTRICTED, not ALLOWED_WITH_RESTRICTIONS)
- ALIGN_ alone (should be ALLOWED_WITH_RESTRICTIONS)
- New action string values (all 5 states must be tested)
- Preserve existing 6 test cases from v1

---

## Advisory System

**`backend/internal/advisory/advisory.go`** — create new file:

```go
package advisory

import (
    "fmt"
    "os"
    "strconv"

    "compliance-checker/internal/models"
    "compliance-checker/internal/preprocessor"
)

// Age thresholds: minimum age_min required to avoid ADV_MISSING_AGE_TARGETING.
// IMPORTANT: Keys must exactly match the category strings output by category_classifier.txt.
// If the user sets age_min below the threshold for a listed category, the advisory fires.
var ageGatingThresholds = map[string]int{
    "Alcohol — Beer, Wine, Spirits":              18,
    "Online Gambling & Betting":                  18,
    "Dating Platforms":                           18,
    "Health — Sexual Wellness":                   18,
    "Health — Reproductive Health":               18,
    "CBD Products":                               18,
    "Drug & Alcohol Addiction Treatment":         18,
    "Weapons Accessories & Safety Equipment":     18,
    "Health — Cannabis / Hemp":                   21, // stricter — US/CA/MX markets
}

// Categories where geo-targeting advisory fires when no countries are specified.
// IMPORTANT: Keys must exactly match the category strings output by category_classifier.txt.
var geoAdvisoryCategories = map[string]bool{
    "Online Gambling & Betting":                    true,
    "CBD Products":                                 true, // US only per Meta policy
    "Cryptocurrency Exchange / Trading / Lending":  true,
    "Drug & Alcohol Addiction Treatment":           true, // US only per Meta policy
    "Health — Cannabis / Hemp":                     true, // US/CA/MX only per Meta policy
}

// Generate produces the advisory list for a given request context.
// This is pure Go — no LLM calls. All logic is deterministic.
// Advisories do NOT affect the compliance score.
func Generate(req *models.CheckRequest, rctx *models.RequestContext, videoDurationSeconds float64) []models.Advisory {
    var advisories []models.Advisory

    category := rctx.Category

    // ADV_MISSING_AGE_TARGETING
    // Fires when the ad category requires age-gating and age_min is below the threshold.
    if threshold, ok := ageGatingThresholds[category]; ok {
        if req.AgeMin < threshold {
            advisories = append(advisories, models.Advisory{
                Code: "ADV_MISSING_AGE_TARGETING",
                Message: fmt.Sprintf(
                    "This ad category (%s) requires a minimum age target of %d+. Your current age_min is %d.",
                    category, threshold, req.AgeMin,
                ),
            })
        }
    }

    // ADV_MISSING_GEO_TARGETING
    // Fires when the ad category has country-specific restrictions but no countries were specified.
    if geoAdvisoryCategories[category] && len(req.Countries) == 0 {
        advisories = append(advisories, models.Advisory{
            Code: "ADV_MISSING_GEO_TARGETING",
            Message: fmt.Sprintf(
                "This ad category (%s) has geo-specific compliance requirements. No target countries were specified — country-level rules could not be verified.",
                category,
            ),
        })
    }

    // ADV_LP_NOT_PROVIDED
    if req.LandingPageURL == "" {
        advisories = append(advisories, models.Advisory{
            Code:    "ADV_LP_NOT_PROVIDED",
            Message: "No landing page URL was provided. Destination-level violations could not be checked.",
        })
    }

    // ADV_LP_FETCH_FAILED
    if rctx.LPAnalysisStatus == "fetch_failed" {
        advisories = append(advisories, models.Advisory{
            Code:    "ADV_LP_FETCH_FAILED",
            Message: "The landing page could not be fetched. Destination-level violations were not checked.",
        })
    }

    // ADV_LP_JS_RENDERED
    // Informational only — Firecrawl handled it, but content may be partial.
    if rctx.LPAnalysisStatus == "js_rendered" {
        advisories = append(advisories, models.Advisory{
            Code:    "ADV_LP_JS_RENDERED",
            Message: "The landing page required JavaScript rendering. Content was captured but may be incomplete.",
        })
    }

    // ADV_VISUAL_PARTIAL
    if rctx.VisualCheckConfidence == "partial" || rctx.VisualCheckConfidence == "low" {
        advisories = append(advisories, models.Advisory{
            Code:    "ADV_VISUAL_PARTIAL",
            Message: "Visual analysis confidence was partial or low. Some image-based violations may have been missed.",
        })
    }

    // ADV_AUDIO_NOT_CHECKED
    // Fires when a video was uploaded but no audio transcript was produced.
    if rctx.HasVideo && rctx.AudioTranscript == "" {
        advisories = append(advisories, models.Advisory{
            Code:    "ADV_AUDIO_NOT_CHECKED",
            Message: "A video was provided but the audio track could not be analyzed. Voiceover violations were not checked.",
        })
    }

    // ADV_VIDEO_LONG_CLIP
    // Fires when video duration exceeds MAX_VIDEO_DURATION_SECONDS.
    maxDuration := 120.0
    if envVal := os.Getenv("MAX_VIDEO_DURATION_SECONDS"); envVal != "" {
        if v, err := strconv.ParseFloat(envVal, 64); err == nil {
            maxDuration = v
        }
    }
    if videoDurationSeconds > maxDuration {
        advisories = append(advisories, models.Advisory{
            Code: "ADV_VIDEO_LONG_CLIP",
            Message: fmt.Sprintf(
                "Video duration (%.0fs) exceeds the %.0fs analysis limit. Frame sampling may have missed content in later segments.",
                videoDurationSeconds, maxDuration,
            ),
        })
    }

    return advisories
}
```

---

*END OF BLOCK 3*

---

## BLOCK 4 — Frontend Changes, Test Cases, Constraints, How to Run

---

## Frontend — TypeScript Types

Update the `checkResult` reactive ref in `ComplianceDashboard.vue` to the full v2 type:

```typescript
const checkResult = ref<{
  check_id: string
  action: 'ALLOWED' | 'PROHIBITED' | 'RESTRICTED' | 'AUTHORIZATION_REQUIRED' | 'ALLOWED_WITH_RESTRICTIONS'
  compliance_score: number
  risk_category: string
  category: string
  category_state: string
  detected_language: string
  visual_check_confidence: 'high' | 'partial' | 'low' | ''
  lp_analysis_status: 'analyzed' | 'not_provided' | 'fetch_failed' | 'js_rendered' | ''
  violations: Array<{
    code: string
    title: string
    severity: string
    explanation: string
    suggested_fix: string
    source: 'copy' | 'image_text' | 'voiceover' | 'landing_page'
  }>
  alignment_violations: Array<{
    code: string
    severity: string
    explanation: string
    source_a: string
    source_b: string
    claim_in_source_a: string
    claim_in_source_b: string
  }>
  advisories: Array<{
    code: string
    message: string
  }>
  suggested_copy: string
} | null>(null)
```

Also add these new reactive state variables in `<script setup>`:

```typescript
// File upload state
const uploadedImageFile = ref<File | null>(null)
const uploadedVideoFile = ref<File | null>(null)
const uploadedFilePreviewUrl = ref<string | null>(null)
const uploadedFileType = ref<'image' | 'video' | null>(null)

// New targeting inputs
const selectedCountries = ref<string[]>([])  // empty = global
const countrySearchQuery = ref('')
const selectedMaxAge = ref('65')  // now actually sent to API

// Advanced options
const forceRefresh = ref(false)
const showAdvancedOptions = ref(false)
```

---

## Frontend — Step 1: Targeting Changes

### Country Selector (Change 1)

Replace the existing `<select v-model="selectedRegion">` with a searchable multi-select tag input. Build this as inline Vue (no external library).

**Behavior:**
- A text input where the user types a country name
- A dropdown appears below showing matching countries from a hardcoded list of ~60 major countries
- Selecting a country adds it as a removable chip/tag above the input
- An "X" button on each chip removes that country
- If no countries are selected, the label reads "All Countries (Global)" — this is the default state
- The list must include at minimum: United States, United Kingdom, Germany, France, Turkey, Brazil, Australia, Canada, India, UAE, Saudi Arabia, Nigeria, South Africa, Japan, South Korea, Singapore, Netherlands, Spain, Italy, Poland, Mexico, Indonesia, Thailand, Malaysia, Philippines, Egypt, Pakistan, Bangladesh, Vietnam, Argentina, Colombia, Chile, Sweden, Norway, Denmark, Finland, Switzerland, Austria, Belgium, Portugal, Czech Republic, Romania, Hungary, Greece, Israel, New Zealand, Ireland, Hong Kong, Taiwan, Ukraine, and all EU countries

**Data sent to API:** `countries` as a JSON-encoded array string in FormData (e.g. `'["US","GB","TR"]'`), or an empty array for global.

### Age Range (Change 2)

- Min age options: 13, 18, 21, 25, 35, 45, 55, 65
- Max age options: 17, 24, 34, 44, 54, 64, 65 ("No max / 65+")
- **Both values must be sent to the API**: `age_min` and `age_max`
- Default: min=18, max=65

---

## Frontend — Step 2: Creative & Copy Changes

### Real File Upload (Change 3)

Replace `triggerUploadSingle()` mock with a real hidden `<input type="file">` element:

```html
<input
  ref="fileInputRef"
  type="file"
  accept="image/jpeg,image/png,image/gif,image/webp,video/mp4,video/quicktime,video/x-msvideo"
  class="hidden"
  @change="handleFileSelected"
/>
```

**`handleFileSelected` logic:**

```typescript
function handleFileSelected(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const isVideo = file.type.startsWith('video/')
  const isImage = file.type.startsWith('image/')

  if (isVideo) {
    uploadedVideoFile.value = file
    uploadedImageFile.value = null
    uploadedFileType.value = 'video'
  } else if (isImage) {
    uploadedImageFile.value = file
    uploadedVideoFile.value = null
    uploadedFileType.value = 'image'
  }

  // Create object URL for preview
  if (uploadedFilePreviewUrl.value) {
    URL.revokeObjectURL(uploadedFilePreviewUrl.value)
  }
  uploadedFilePreviewUrl.value = URL.createObjectURL(file)
}
```

**Preview rendering:**
- If `uploadedFileType === 'image'`: render `<img :src="uploadedFilePreviewUrl">` in the upload zone
- If `uploadedFileType === 'video'`: render `<video :src="uploadedFilePreviewUrl" controls class="w-full rounded-lg">` in the upload zone
- When a video is detected, show an informational note below the upload zone:
  `"Video detected — copy, visuals, and audio track will all be analyzed. This takes 10–15 seconds."`

### Multipart Form Submission (Change 4)

Replace the `fetch()` JSON body in `nextStep()` with a `FormData` builder:

```typescript
async function nextStep() {
  if (currentStep.value === 2) {
    isChecking.value = true
    currentStep.value = 3
    checkError.value = null

    try {
      const formData = new FormData()
      formData.append('platform', selectedPlatform.value)
      formData.append('age_min', selectedMinAge.value)
      formData.append('age_max', selectedMaxAge.value)
      formData.append('ad_format', adFormat.value)
      formData.append('primary_text', adCopy.value.primaryText)
      formData.append('headline', adCopy.value.headline)
      formData.append('description', adCopy.value.description)
      formData.append('landing_page_url', landingPageUrl.value)
      formData.append('force_refresh', forceRefresh.value ? 'true' : 'false')
      formData.append('countries', JSON.stringify(selectedCountries.value))

      if (uploadedImageFile.value) {
        formData.append('image', uploadedImageFile.value)
      }
      if (uploadedVideoFile.value) {
        formData.append('video', uploadedVideoFile.value)
      }

      // Do NOT set Content-Type header manually — browser sets it with correct boundary
      const response = await fetch('http://localhost:8080/checks', {
        method: 'POST',
        body: formData
      })

      if (!response.ok) throw new Error('API request failed')

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

### Force Refresh Toggle (Change 5)

Below the Landing Page URL input, add a subtle toggle:

```html
<div class="mt-2">
  <button
    type="button"
    @click="showAdvancedOptions = !showAdvancedOptions"
    class="text-xs text-muted-foreground hover:text-foreground"
  >
    {{ showAdvancedOptions ? '▲ Hide' : '▼ Advanced options' }}
  </button>
  <div v-if="showAdvancedOptions" class="mt-2 flex items-center gap-2">
    <input type="checkbox" id="forceRefresh" v-model="forceRefresh" class="rounded" />
    <label for="forceRefresh" class="text-xs text-muted-foreground">
      Re-fetch this page (bypass 24-hour cache)
    </label>
  </div>
</div>
```

### Carousel: Image-Only (Change 6)

In the carousel card upload handler, restrict the file input to images only:

```html
<input type="file" accept="image/jpeg,image/png,image/gif,image/webp" ... />
```

---

## Frontend — Step 3: Loading State Changes

### Dynamic Loading Steps (Change 7 + 8)

Replace the two hardcoded loading items with a computed list based on what was submitted:

```typescript
const loadingSteps = computed(() => {
  const steps = [
    { label: 'Classifying ad category', done: true },
    { label: 'Analyzing copy & policy rules', done: false },
  ]
  if (uploadedImageFile.value || uploadedVideoFile.value) {
    steps.push({ label: 'Analyzing visual creative', done: false })
  }
  if (uploadedVideoFile.value) {
    steps.push({ label: 'Transcribing audio track', done: false })
  }
  if (landingPageUrl.value) {
    steps.push({ label: 'Fetching and analyzing landing page', done: false })
  }
  if ((uploadedImageFile.value || uploadedVideoFile.value) && landingPageUrl.value) {
    steps.push({ label: 'Cross-checking ad against landing page', done: false })
  }
  return steps
})
```

Update the loading subtext: `"Running full multimodal analysis. This typically takes 10–15 seconds."`

Animate the steps sequentially: use a `setInterval` that marks one additional step as `done: true` every 2 seconds. Reset interval when `isChecking` becomes false.

---

## Frontend — Step 4: Results Changes

### Score Card: 6-State Action Badge (Change 10)

Replace the binary `not_compliant` check with a full action map:

```typescript
const actionConfig = computed(() => {
  const map: Record<string, { label: string; bgClass: string; textClass: string }> = {
    PROHIBITED:               { label: 'Prohibited',               bgClass: 'bg-destructive/10',  textClass: 'text-destructive' },
    AUTHORIZATION_REQUIRED:   { label: 'Authorization Required',   bgClass: 'bg-orange-100',      textClass: 'text-orange-700' },
    RESTRICTED:               { label: 'Restricted',               bgClass: 'bg-amber-100',       textClass: 'text-amber-700' },
    ALLOWED_WITH_RESTRICTIONS:{ label: 'Allowed With Restrictions', bgClass: 'bg-yellow-100',     textClass: 'text-yellow-700' },
    ALLOWED:                  { label: 'Compliant',                bgClass: 'bg-green-100',       textClass: 'text-green-700' },
  }
  return map[checkResult.value?.action ?? 'ALLOWED'] ?? map['ALLOWED']
})
```

### Violation Source Badge (Change 13)

For each violation in the list, render a source badge beside the severity badge:

```typescript
const sourceBadgeConfig: Record<string, { label: string; bgClass: string; textClass: string }> = {
  copy:         { label: 'Ad copy',      bgClass: 'bg-muted/50',    textClass: 'text-muted-foreground' },
  image_text:   { label: 'Image text',   bgClass: 'bg-amber-100',   textClass: 'text-amber-700' },
  voiceover:    { label: 'Voiceover',    bgClass: 'bg-blue-100',    textClass: 'text-blue-700' },
  landing_page: { label: 'Landing page', bgClass: 'bg-purple-100',  textClass: 'text-purple-700' },
}
```

Show this badge only when source is not `'copy'` — it is the default and does not need a label.

### Visual Confidence Disclaimer (Change 14)

If `checkResult.visual_check_confidence` is `'partial'` or `'low'`, render this note directly above the violations list:

```html
<div v-if="checkResult?.visual_check_confidence === 'partial' || checkResult?.visual_check_confidence === 'low'"
     class="bg-amber-50 border border-amber-200 rounded-lg p-3 text-xs text-amber-800 flex items-center gap-2">
  <IconAlertCircle :size="14" class="shrink-0" />
  Visual analysis confidence was partial. Some image-based violations may have been missed.
</div>
```

### Detected Language Badge (Change 15)

In the results header, if `detected_language` is not `'en'` and not empty, show:

```html
<span v-if="checkResult?.detected_language && checkResult.detected_language !== 'en'"
      class="text-xs bg-muted/50 border border-border px-2 py-0.5 rounded-full text-muted-foreground">
  Analyzed in: {{ languageLabel(checkResult.detected_language) }}
</span>
```

Implement `languageLabel()` with a simple map for common languages (tr=Turkish, de=German, fr=French, es=Spanish, pt=Portuguese, ar=Arabic, it=Italian, nl=Dutch, pl=Polish, ru=Russian). Return the ISO code itself as fallback.

### Alignment Violations Section (Change 12)

After the main violations list, add a separate accordion section for alignment violations:

```html
<div v-if="checkResult?.alignment_violations?.length > 0"
     class="border border-border rounded-xl bg-card overflow-hidden">
  <div class="bg-amber-500/10 border-b border-amber-500/20 p-4 flex items-center gap-2">
    <IconAlertCircle :size="20" class="text-amber-600" />
    <span class="font-bold text-amber-700">
      {{ checkResult.alignment_violations.length }} Consistency Issue{{ checkResult.alignment_violations.length > 1 ? 's' : '' }}
    </span>
    <span class="text-xs text-amber-600 ml-2">Ad vs. Landing Page</span>
  </div>
  <div class="p-6 space-y-6">
    <div v-for="av in checkResult.alignment_violations" :key="av.code">
      <div class="text-sm font-bold text-foreground mb-2">{{ av.explanation }}</div>
      <div class="grid grid-cols-2 gap-3">
        <div class="bg-muted/30 rounded-lg p-3 border border-border">
          <div class="text-[10px] font-bold text-muted-foreground uppercase mb-1">{{ av.source_a }}</div>
          <div class="text-sm text-foreground italic">"{{ av.claim_in_source_a }}"</div>
        </div>
        <div class="bg-muted/30 rounded-lg p-3 border border-border">
          <div class="text-[10px] font-bold text-muted-foreground uppercase mb-1">{{ av.source_b }}</div>
          <div class="text-sm text-foreground italic">"{{ av.claim_in_source_b }}"</div>
        </div>
      </div>
    </div>
  </div>
</div>
```

### Suggested Copy Card (Change 16)

After the violations and alignment sections, if `suggested_copy` is non-empty:

```html
<div v-if="checkResult?.suggested_copy" class="border border-border rounded-xl bg-card overflow-hidden">
  <div class="bg-green-50 border-b border-green-200 p-4 flex items-center justify-between">
    <span class="font-bold text-green-800">Suggested Revised Copy</span>
    <button @click="copySuggestedCopy" class="text-xs font-bold text-green-700 hover:underline">
      Copy to clipboard
    </button>
  </div>
  <div class="p-4">
    <p class="text-sm text-foreground bg-muted/30 p-4 rounded-lg">{{ checkResult.suggested_copy }}</p>
    <div class="mt-3 flex justify-end">
      <button
        disabled
        class="bg-muted text-muted-foreground px-4 py-2 rounded-lg text-sm font-bold cursor-not-allowed opacity-60"
        title="Coming in v3"
      >
        Apply to Meta Draft (v3)
      </button>
    </div>
  </div>
</div>
```

---

## Frontend — Advisory Panel Component

**`src/components/AdvisoryPanel.vue`** — create new file:

```vue
<script setup lang="ts">
import { IconInfoCircle } from '@tabler/icons-vue'

defineProps<{
  advisories: Array<{ code: string; message: string }>
}>()
</script>

<template>
  <div v-if="advisories.length > 0" class="border border-amber-200 rounded-xl bg-amber-50 overflow-hidden">
    <div class="px-4 py-3 border-b border-amber-200 flex items-center gap-2">
      <IconInfoCircle :size="18" class="text-amber-600 shrink-0" />
      <span class="font-bold text-amber-800 text-sm">
        {{ advisories.length }} Recommendation{{ advisories.length > 1 ? 's' : '' }}
      </span>
      <span class="text-xs text-amber-600 ml-1">— these do not affect your compliance score</span>
    </div>
    <div class="p-4 space-y-3">
      <div
        v-for="advisory in advisories"
        :key="advisory.code"
        class="flex items-start gap-2"
      >
        <div class="h-1.5 w-1.5 rounded-full bg-amber-500 mt-2 shrink-0"></div>
        <p class="text-sm text-amber-800">{{ advisory.message }}</p>
      </div>
    </div>
  </div>
</template>
```

Import and render `AdvisoryPanel` in `ComplianceDashboard.vue` inside the Step 4 results section, directly below the score card and above the violations list:

```html
<AdvisoryPanel :advisories="checkResult?.advisories ?? []" />
```

---

## How to Run Locally (v2)

```bash
# Prerequisites
# Install ffmpeg (required for video processing)
# macOS:   brew install ffmpeg
# Ubuntu:  sudo apt-get install -y ffmpeg

# 1. Start the database
cd backend
docker compose up -d

# 2. Set up environment
cp .env.example .env
# Edit .env:
#   GEMINI_API_KEY=your_google_ai_studio_key
#   FIRECRAWL_API_KEY=your_firecrawl_key   (get free key at firecrawl.dev)
#   GEMINI_MODEL=gemini-2.0-flash          (default, leave as-is)
#   MAX_VIDEO_DURATION_SECONDS=120          (default, leave as-is)

# 3. Install Go dependencies
go mod tidy

# 4. Start the Go backend
go run cmd/server/main.go

# 5. In a new terminal, start the Vue frontend
cd ..
npm install
npm run dev
# Frontend: http://localhost:5173
# Backend:  http://localhost:8080

# For text-only testing (no ffmpeg or Firecrawl needed):
# Submit the form without uploading a file and without a landing page URL.
# The pipeline falls back cleanly to text-only analysis.
```

---

## Validation — 15 Test Cases

Test these after the full pipeline is running. Cases 1–6 are text-only (no file required). Cases 7–15 require files or URLs.

| # | Input | Expected violations | Expected action | Expected score |
|---|-------|--------------------|-----------------|-|
| 1 | Copy: "Lose 20kg in 30 days with our supplement!" | R_WEIGHT_LOSS_CLAIM (copy) | RESTRICTED | ~50 |
| 2 | Copy: "Are you struggling with debt? We can help." | P_PERSONAL_ATTRIBUTE (copy) | PROHIBITED | 0 |
| 3 | Copy: "Bet on Premier League today — big wins await!" | A_GAMBLING | AUTHORIZATION_REQUIRED | 0 |
| 4 | Copy: "Get fit this summer. Shop our activewear." | none | ALLOWED | 95 |
| 5 | Copy: "DOCTORS HATE HIM! This one trick cures arthritis!" | R_HEALTH_CLAIM + M_SENSATIONALISM | RESTRICTED | ~35 |
| 6 | Copy: "Trade Bitcoin on our platform — guaranteed returns" | A_CRYPTO_EXCHANGE + R_GUARANTEED_OUTCOME | AUTHORIZATION_REQUIRED | 0 |
| 7 | Image with clear before/after weight loss comparison | R_BEFORE_AFTER_WEIGHT_VISUAL (image_text) | RESTRICTED | ~50 |
| 8 | Copy: "FDA approved formula" + LP with FDA disclaimer text | ALIGN_CLAIM_MISMATCH | ALLOWED_WITH_RESTRICTIONS | ~70 |
| 9 | Video with voiceover: "Are you struggling with anxiety?" | P_PERSONAL_ATTRIBUTE (voiceover) | PROHIBITED | 0 |
| 10 | Image with text overlay: "Guaranteed results in 7 days" | R_GUARANTEED_OUTCOME (image_text) | RESTRICTED | ~50 |
| 11 | LP URL that returns 404 or times out | ADV_LP_FETCH_FAILED advisory; no LP violations | Depends on copy | — |
| 12 | Alcohol category ad, age_min=13 | ADV_MISSING_AGE_TARGETING advisory | Depends on copy | — |
| 13 | Online Gambling category, no countries selected | ADV_MISSING_GEO_TARGETING advisory | AUTHORIZATION_REQUIRED | 0 |
| 14 | Clean image + clean copy + compliant LP | none; no advisories except possibly ADV_LP_JS_RENDERED | ALLOWED | 95 |
| 15 | Video submitted, no LP provided | ADV_LP_NOT_PROVIDED + ADV_AUDIO_NOT_CHECKED (if audio fails) | Depends on copy | — |

**Pass threshold:** 12/15 correct. For cases 11–15, the advisory code must be present in the `advisories[]` array; the compliance score being unaffected by advisories must be verified.

---

## Constraints

**Keep all v1 constraints:**
- Do not use any Go HTTP framework (Gin, Echo, Fiber). Use standard `net/http` only.
- Do not add authentication or middleware.
- All prompts must be read from `.txt` files, not hardcoded in Go.
- Temperature must be 0 for all Gemini calls.
- All Gemini calls must use `ResponseSchema` enforcement — never parse free-form JSON text.
- The scoring engine must remain pure Go with no LLM involvement.

**New v2 constraints:**

1. **Fix invalid Go syntax in preprocessor.go**: The `parseJSON` function has an `import "encoding/json"` statement inside the function body. Move `encoding/json` to the package-level `import` block at the top of the file. This is the only known syntax error in the Block 2 specification.

2. **Do not hardcode model name**: Always read from `os.Getenv("GEMINI_MODEL")` with a fallback to `"gemini-2.0-flash"`. The model string must never appear as a string literal in Go code.

3. **Do not hardcode video duration limit**: Always read from `os.Getenv("MAX_VIDEO_DURATION_SECONDS")` with a fallback to `120`. Must not appear as a numeric literal in Go code.

4. **ffmpeg is a system binary, not a Go package**: Do not attempt to import or install ffmpeg via `go get`. It must be installed separately (see How to Run Locally). Wrap all `exec.Command("ffmpeg", ...)` calls with error handling that returns a descriptive error if the binary is not found (`exec.LookPath("ffmpeg")`).

5. **Firecrawl API key failure is non-fatal**: If `FIRECRAWL_API_KEY` is empty or Firecrawl returns an error, the LP Agent must return `LPAnalysisStatus: "fetch_failed"` and the handler must continue — it must not return an HTTP 500.

6. **Advisory system is pure Go, no LLM**: The `advisory.Generate()` function must contain no Gemini calls, no HTTP calls, no goroutines. It is deterministic and synchronous.

7. **Advisories are never passed to `scoring.Score()`**: The scoring engine receives only `AgentViolation` codes. Advisory codes (prefixed `ADV_`) must never appear in the violations slice passed to scoring.

8. **Do not mix AI vendors**: All LLM calls use Gemini via the `google/generative-ai-go` SDK. Do not add OpenAI, Anthropic, or any other AI vendor dependency.

9. **Frontend: do not set `Content-Type` header manually on multipart fetch**: The browser sets it automatically with the correct multipart boundary when using `FormData`. Manually setting it will break the boundary string and cause the backend to fail parsing.

10. **Frontend: revoke object URLs on file change**: Call `URL.revokeObjectURL(uploadedFilePreviewUrl.value)` before creating a new preview URL to prevent memory leaks on repeated file selection.

11. **Carousel format remains image-only**: The carousel card upload must only accept `image/*` MIME types. Video upload is available only in single-ad format.

---

*END OF BLOCK 4 — ANTIGRAVITY-INSTRUCTIONS-v2.md IS COMPLETE*

package models

// ─────────────────────────────────────────────
// REQUEST TYPES
// ─────────────────────────────────────────────

// CheckRequest is populated by the handler from the multipart form (or JSON fallback).
// When a file is present the handler uses multipart/form-data.
// When no file is present the handler accepts application/json for backward compatibility.
type CheckRequest struct {
    Platform       string   `json:"platform"`
    Region         string   `json:"region"`          // "global" | "us" | "uk" | "eu"
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

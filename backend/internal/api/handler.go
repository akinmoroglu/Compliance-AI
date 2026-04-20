package api

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"

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
    reqStart := time.Now()
    log.Println("📥 [API] Handling incoming ad compliance check...")
    log.Println("⏳ [Pipeline] STEP 1: Pre-processing media inputs (Extracting frames/audio)...")
    pStart := time.Now()
    ppOut, err := preprocessor.Parse(r)
    log.Printf("   -> Pre-processing finished in %v", time.Since(pStart))
    
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

    // Prepend targeting context so all agents are aware of region and age targeting.
    regionLabel := map[string]string{
        "us": "United States",
        "uk": "United Kingdom",
        "eu": "European Union",
    }[req.Region]
    if regionLabel == "" {
        regionLabel = "Global (all regions)"
    }
    ageMax := req.AgeMax
    if ageMax == 0 {
        ageMax = 65
    }
    fullCopy = fmt.Sprintf("[TARGETING CONTEXT]\nRegion: %s\nAge range: %d–%d\n\n%s",
        regionLabel, req.AgeMin, ageMax, fullCopy)

    ctx := r.Context()

    // ─────────────────────────────────────────
    // STEP 2: Cache check
    // ─────────────────────────────────────────
    log.Println("⏳ [Pipeline] STEP 2: Checking SHA-256 fingerprint cache...")
    cacheKey := cache.BuildKey(req, ppOut.ImageSize, ppOut.VideoSize)
    if !req.ForceRefresh {
        if cached, ok := cache.Get(cacheKey); ok {
            log.Printf("🎯 [Cache] Hit for key %s. Returning instantly.", cacheKey[:8])
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
        sessionTotalTokens int32
        tokenMu            sync.Mutex
    )

    trackTokens := func(name string, count int32) {
        if count > 0 {
            log.Printf("🪙 [Tokens] %s used %d tokens", name, count)
            tokenMu.Lock()
            sessionTotalTokens += count
            tokenMu.Unlock()
        }
    }

    // Helper: run an agent goroutine, collect violations, log errors without crashing.
    makeRunner := func(wg *sync.WaitGroup) func(fn func() ([]models.AgentViolation, int32, error), name string) {
        return func(fn func() ([]models.AgentViolation, int32, error), name string) {
            agentStart := time.Now()
            defer wg.Done()
            v, tokens, err := fn()
            trackTokens(name, tokens)
            if err != nil {
                log.Printf("❌ [Agent] %s failed after %v: %v", name, time.Since(agentStart), err)
                return
            }
            log.Printf("✅ [Agent] %s complete in %v", name, time.Since(agentStart))
            mu.Lock()
            allAgentViolations = append(allAgentViolations, v...)
            mu.Unlock()
        }
    }

    // ─────────────────────────────────────────
    // STEP 4: Sequential — Category Classifier
    // ─────────────────────────────────────────
    log.Println("⏳ [Pipeline] STEP 4: Running Category Classifier Agent...")
    cStart := time.Now()
    categoryResult, cTokens, err := agents.RunCategoryClassifier(ctx, fullCopy)
    trackTokens("category", cTokens)
    log.Printf("   -> Category Classifier finished in %v", time.Since(cStart))
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
        finaliseAndRespond(w, req, rctx, allAgentViolations, nil, cacheKey, ppOut, categoryResult, reqStart, sessionTotalTokens)
        return
    }

    // ─────────────────────────────────────────
    // STEP 5: Phase 1 — Parallel agents
    // ─────────────────────────────────────────
    log.Println("🚀 [Pipeline] STEP 5: Launching Phase 1 Parallel Core Agents...")
    // Goroutines fire conditionally based on what inputs are available.
    // All goroutines write to rctx (non-overlapping fields) or allAgentViolations (mutex-protected).

    var wg1 sync.WaitGroup
    run1 := makeRunner(&wg1)

    // g1 — Violations Agent: copy-based policy check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, int32, error) {
        return agents.RunViolationsAgent(ctx, fullCopy, rctx.Category, "", "")
    }, "violations")

    // g2 — Personal Attribute Agent: copy-based PA check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, int32, error) {
        return agents.RunPersonalAttributeAgent(ctx, fullCopy, "", "")
    }, "personal_attribute")

    // g3 — Auth Required Agent: copy-based authorization check
    wg1.Add(1)
    go run1(func() ([]models.AgentViolation, int32, error) {
        return agents.RunAuthRequiredAgent(ctx, fullCopy, rctx.Category)
    }, "auth_required")

    // g4 — Vision Agent: image or video frame analysis + OCR
    // Only fires when image or video frames are present.
    if rctx.HasImage || rctx.HasVideo {
        wg1.Add(1)
        go func() {
            agentStart := time.Now()
            defer wg1.Done()
            visionOut, tokens, err := agents.RunVisionAgent(ctx, rctx.ImageBytes, rctx.VideoFrames)
            trackTokens("vision", tokens)
            if err != nil {
                log.Printf("❌ [Agent] vision failed after %v: %v", time.Since(agentStart), err)
                return
            }
            log.Printf("✅ [Agent] vision complete in %v", time.Since(agentStart))
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
            agentStart := time.Now()
            defer wg1.Done()
            lpOut, tokens, err := agents.RunLPAgent(ctx, req.LandingPageURL)
            trackTokens("lp", tokens)
            if err != nil {
                log.Printf("❌ [Agent] lp failed after %v: %v", time.Since(agentStart), err)
                rctx.LPAnalysisStatus = "fetch_failed"
                return
            }
            log.Printf("✅ [Agent] lp complete in %v", time.Since(agentStart))
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
            agentStart := time.Now()
            defer wg1.Done()
            audioOut, tokens, err := agents.RunAudioAgent(ctx, rctx.AudioBytes)
            trackTokens("audio", tokens)
            if err != nil {
                log.Printf("❌ [Agent] audio failed after %v: %v", time.Since(agentStart), err)
                return
            }
            log.Printf("✅ [Agent] audio complete in %v", time.Since(agentStart))
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
    log.Println("🚦 [Pipeline] STEP 6: Waiting on media transcripts for Phase 1b (if applicable)...")
    // After Phase 1 completes, rctx now has ImageExtractedText and AudioTranscript.
    // Run Violations Agent and Personal Attribute Agent again with media-sourced text.
    // Near-zero latency penalty — LP Agent (~8s) was the Phase 1 bottleneck.
    if rctx.ImageExtractedText != "" || rctx.AudioTranscript != "" {
        var wg1b sync.WaitGroup
        run1b := makeRunner(&wg1b)

        wg1b.Add(1)
        go run1b(func() ([]models.AgentViolation, int32, error) {
            return agents.RunViolationsAgent(ctx, fullCopy, rctx.Category, rctx.ImageExtractedText, rctx.AudioTranscript)
        }, "violations_1b")

        wg1b.Add(1)
        go run1b(func() ([]models.AgentViolation, int32, error) {
            return agents.RunPersonalAttributeAgent(ctx, fullCopy, rctx.ImageExtractedText, rctx.AudioTranscript)
        }, "personal_attribute_1b")

        wg1b.Wait()
    }

    // ─────────────────────────────────────────
    // STEP 7: Phase 2 — Alignment Agent
    // ─────────────────────────────────────────
    log.Println("🔍 [Pipeline] STEP 7: Running Cross-Source Alignment Agent...")
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

        alignStart := time.Now()
        var alignTokens int32
        alignmentViolations, alignTokens, err = agents.RunAlignmentAgent(
            ctx,
            alignImageBytes,
            fullCopy,
            rctx.LPKeyClaims,
            existingCodes,
        )
        trackTokens("alignment", alignTokens)
        if err != nil {
            log.Printf("❌ [Agent] alignment failed after %v: %v", time.Since(alignStart), err)
            // Non-fatal — continue without alignment results
        } else {
            log.Printf("✅ [Agent] alignment complete in %v", time.Since(alignStart))
        }
    }

    finaliseAndRespond(w, req, rctx, allAgentViolations, alignmentViolations, cacheKey, ppOut, categoryResult, reqStart, sessionTotalTokens)
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
    reqStart time.Time,
    sessionTotalTokens int32,
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
    log.Println("🛠️ [Pipeline] STEP 8: Generating suggested fixes...")
    fStart := time.Now()
    fixes, fixTokens, err := agents.RunFixGenerationAgent(ctx, req.PrimaryText, categoryResult.Category, uniqueViolations)
    if fixTokens > 0 {
        log.Printf("🪙 [Tokens] fix_generation used %d tokens", fixTokens)
        sessionTotalTokens += fixTokens
    }
    log.Printf("   -> Fix generation finished in %v", time.Since(fStart))
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
    log.Printf("🪙 [Tokens] Total tokens used for session: %d", sessionTotalTokens)
    log.Printf("✅ [Pipeline] SUCCESS! Encoding and responding. Total latency: %v", time.Since(reqStart))
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

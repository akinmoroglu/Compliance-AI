package advisory

import (
    "fmt"
    "os"
    "strconv"

    "compliance-checker/internal/models"
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
    // Fires when the ad category has country-specific restrictions and region is global/unspecified.
    if geoAdvisoryCategories[category] && (req.Region == "" || req.Region == "global") {
        advisories = append(advisories, models.Advisory{
            Code: "ADV_MISSING_GEO_TARGETING",
            Message: fmt.Sprintf(
                "This ad category (%s) has geo-specific compliance requirements. No specific region was selected — country-level rules could not be verified.",
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

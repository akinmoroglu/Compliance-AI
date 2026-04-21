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

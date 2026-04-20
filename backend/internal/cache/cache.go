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

    raw := fmt.Sprintf("%s|%d|%d|%s|%s|%s|%s|%s|%v|%d|%d",
        req.PrimaryText,
        req.AgeMin,
        req.AgeMax,
        req.LandingPageURL,
        req.Headline,
        req.Description,
        req.Region,
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

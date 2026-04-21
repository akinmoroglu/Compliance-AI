> **Source:** https://pozitim.atlassian.net/wiki/spaces/ROC/pages/3403579458
> **Last updated:** 2026-04-13 (v2) — "Add landing_page_url field, result fields, error code INVALID_LANDING_PAGE_URL, and API-OQ-7"

# API Contract: Ad Compliance Checker (Phase 1)

| Field | Value |
| --- | --- |
| Product | Rockads Marketing Panel |
| Feature | Ad Compliance Checker |
| Phase | 1 of 3 |
| Status | Draft |
| Author | Akin Moroglu |
| Last Updated | April 2026 |
| Parent doc | PRD: Ad Compliance AI: Phase 1 (Meta) |

---

## Overview

This document defines the API contract between the frontend (Vue 3) and backend (Golang + Encore.dev) for the Ad Compliance Checker feature. Both teams should treat this as the source of truth for request/response shapes. Changes require sign-off from both FE and BE.

**Base URL:** `/compliance/v1`
**Auth:** All endpoints require the existing Rockads panel session token (same auth model as the rest of the panel; no new auth mechanism).
**Content-Type:** `application/json` except where noted (file upload uses `multipart/form-data`).

> **Note for BE:** Encore.dev generates a typed TypeScript client from service definitions. Once endpoints are implemented, share the generated client with FE to avoid manual fetch wrappers.

---

## Design Decisions

**Why all checks are async at the API level:**
The API always returns a `check_id` immediately on submission; it never blocks until evaluation completes. The FE handles the two UX modes (synchronous-feeling vs. true async) by polling behavior:

* For images and videos up to 30 seconds: FE polls every 2 seconds for up to 35 seconds. For most requests this returns quickly and feels synchronous to the user.
* For videos over 30 seconds: FE detects video duration client-side before upload, shows the async holding screen immediately after submission, and stops polling. BE sends the email notification on completion.

This keeps the API simple (one mode, not two) and lets the FE own the UX distinction.

**Why file upload is a separate step:**
Large video files should not be bundled into the check submission body. A dedicated upload endpoint returns a `file_id`; the check submission references that ID. This also makes retry logic cleaner if the check fails.

**Credit handling:**
Credits are reserved on submission and deducted only on successful completion. A failed or timed-out check releases the reservation; no credits are lost.

**Landing page evaluation:**
The `landing_page_url` field is optional. If provided, the backend fetches the page at submission time and includes its content in the AI evaluation. If the page is unreachable, the check proceeds without it and the result includes a `landing_page_fetch_warning`. This is a non-blocking behavior — a failed landing page fetch does not cause the check to fail.

---

## Endpoints

### 1. Upload a Creative File

`POST /compliance/v1/uploads`
**Content-Type:** `multipart/form-data`

Uploads an image or video file and returns a `file_id` to reference in the check submission. Uploaded files expire after 24 hours if not used in a check.

**Request:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `file` | binary | Yes | Image or video file |

**Response** `200`:

```json
{
  "file_id": "file_abc123",
  "filename": "my-ad.mp4",
  "size_bytes": 8421376,
  "media_type": "video",
  "duration_seconds": 45,
  "expires_at": "2026-04-10T12:00:00Z"
}
```

**Field notes:**

* `media_type`: `"image"` or `"video"`
* `duration_seconds`: present only for video files; `null` for images
* `expires_at`: FE should not cache `file_id` beyond this timestamp

**Errors:**

| Code | HTTP | Meaning |
| --- | --- | --- |
| `FILE_TOO_LARGE` | 413 | File exceeds maximum allowed size (limit TBD; see OQ-4) |
| `UNSUPPORTED_FORMAT` | 415 | File type is not accepted (accepted types TBD; see OQ-4) |
| `UPLOAD_FAILED` | 500 | Server-side upload error; safe to retry |

---

### 2. Submit a Compliance Check

`POST /compliance/v1/checks`

Submits an ad for compliance evaluation. Always returns immediately with a `check_id`. Credits are reserved at this point.

**Request (single format):**

```json
{
  "platform": "meta",
  "region": "global",
  "age_min": 25,
  "age_max": 34,
  "ad_format": "single",
  "primary_text": "Lose 20kg in 30 days!",
  "headline": "100% Natural",
  "description": "No exercise needed.",
  "landing_page_url": "https://example.com/weight-loss-product",
  "creatives": [
    {
      "file_id": "file_abc123"
    }
  ]
}
```

**Request (carousel format):**

```json
{
  "platform": "meta",
  "region": "us",
  "age_min": 18,
  "age_max": 44,
  "ad_format": "carousel",
  "primary_text": "Discover our full range.",
  "landing_page_url": "https://example.com/products",
  "creatives": [
    {
      "file_id": "file_card1",
      "headline": "Product One",
      "description": "Best for beginners."
    },
    {
      "file_id": "file_card2",
      "headline": "Product Two",
      "description": "Advanced formula."
    }
  ]
}
```

**Request fields:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `platform` | string | Yes | `"meta"` only in Phase 1 |
| `region` | string | Yes | `"global"` or ISO country/region code (e.g., `"us"`, `"eu"`) |
| `age_min` | integer | Yes | Minimum of target age range |
| `age_max` | integer | Yes | Maximum of target age range |
| `ad_format` | string | Yes | `"single"` or `"carousel"` |
| `primary_text` | string | No | Applies to single and carousel (global copy) |
| `headline` | string | No | Single format only |
| `description` | string | No | Single format only |
| `landing_page_url` | string | No | HTTPS URL of the ad's destination page. If provided, BE fetches and evaluates the page as part of the check. If omitted or unreachable, check runs on creative and copy only. |
| `creatives` | array | Yes | At least one item required |
| `creatives[].file_id` | string | Yes | ID from the upload endpoint |
| `creatives[].headline` | string | No | Carousel cards only |
| `creatives[].description` | string | No | Carousel cards only |

**Response** `200`:

```json
{
  "check_id": "chk_xyz789",
  "status": "pending",
  "is_async": false,
  "credits_reserved": 10,
  "created_at": "2026-04-09T10:00:00Z"
}
```

**Field notes:**

* `is_async`: `true` if the uploaded video duration exceeds 30 seconds. FE uses this to decide whether to poll or show the async holding screen.
* `credits_reserved`: number of credits held; deducted on completion, released on failure.

**Errors:**

| Code | HTTP | Meaning |
| --- | --- | --- |
| `INSUFFICIENT_CREDITS` | 402 | User does not have enough AI credits; FE should prompt top-up |
| `FILE_NOT_FOUND` | 404 | One or more `file_id` values are invalid or expired |
| `INVALID_REQUEST` | 400 | Missing required fields or invalid values |
| `INVALID_LANDING_PAGE_URL` | 400 | URL provided is not a valid HTTPS URL. FE should validate client-side first (US-14); this is the server-side safety net. |

---

### 3. Get Check Status and Results

`GET /compliance/v1/checks/{check_id}`

Returns the current status of a check. FE polls this endpoint after submission.

**Response** `200`:

```json
{
  "check_id": "chk_xyz789",
  "status": "completed",
  "is_async": false,
  "platform": "meta",
  "ad_format": "single",
  "created_at": "2026-04-09T10:00:00Z",
  "completed_at": "2026-04-09T10:00:18Z",
  "credits_consumed": 10,
  "result": {
    "passed": false,
    "landing_page_evaluated": true,
    "landing_page_fetch_warning": null,
    "violations": [
      {
        "id": "viol_001",
        "title": "Misleading Health Claim",
        "severity": "high_risk",
        "policy_reference": "Meta Advertising Policies: Health and Wellness",
        "explanation": "The primary text makes a guaranteed weight loss claim ('Lose 20kg in 30 days!'). Meta prohibits guaranteed outcomes for health or wellness results without scientific substantiation.",
        "suggested_fix": "Remove the guarantee framing. For example: 'Thousands of customers have seen results with our natural formula. See what it can do for you.'"
      }
    ]
  }
}
```

**Status values:**

| Status | Meaning | FE action |
| --- | --- | --- |
| `pending` | Check queued, not yet started | Keep polling |
| `processing` | AI evaluation in progress | Keep polling |
| `completed` | Results ready | Show results |
| `failed` | Evaluation failed; credits not deducted | Show error + retry option |

**Violation object:**

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | Unique per violation within a check |
| `title` | string | Short label |
| `severity` | string | `"high_risk"`, `"medium_risk"`, or `"low_risk"` |
| `policy_reference` | string | Named Meta policy (not just a description) |
| `explanation` | string | Plain-language reason for violation |
| `suggested_fix` | string | Concrete rewrite or modification |

**Landing page fields in result:**

| Field | Type | Notes |
| --- | --- | --- |
| `landing_page_evaluated` | boolean | `true` if a URL was provided and successfully fetched |
| `landing_page_fetch_warning` | string or null | Human-readable warning if the URL was provided but could not be fetched. `null` if fetch succeeded or no URL was provided. |

**When** `passed` is `true`: `violations` is an empty array `[]`.

**Errors:**

| Code | HTTP | Meaning |
| --- | --- | --- |
| `CHECK_NOT_FOUND` | 404 | `check_id` does not exist or does not belong to this user |

---

### 4. Get Credit Balance

`GET /credits/v1/balance`

Returns the user's current AI credit balance. FE calls this before enabling the "Check Compliance" button.

> This endpoint likely already exists in the panel. BE to confirm or share the existing contract. If the path differs, FE should use the existing endpoint rather than a new one.

**Response** `200`:

```json
{
  "balance": 150,
  "unit": "ai_credits"
}
```

---

## Polling Behavior (FE Responsibility)

```
FE submits check
      |
      v
is_async = false?
   YES: poll GET /checks/{check_id} every 2s
         |
         |-- status = "completed" --> show results
         |-- status = "failed"    --> show error + retry
         |-- 35 seconds elapsed with no result --> treat as async (show email state)

   NO (is_async = true): show async holding screen, stop polling
         |
         BE sends email on completion
         User clicks deep-link --> FE calls GET /checks/{check_id} --> show results
```

---

## Error Response Format

All errors follow this structure:

```json
{
  "code": "INSUFFICIENT_CREDITS",
  "message": "Your AI credit balance is too low to run a compliance check. Please top up to continue.",
  "details": {}
}
```

`message` is safe to display directly in the UI. `details` is optional and may contain additional context for debugging.

---

## Full Error Code Reference

| Code | HTTP | Endpoint(s) | Notes |
| --- | --- | --- | --- |
| `INSUFFICIENT_CREDITS` | 402 | POST /checks | FE blocks before API call (NFR-6); this is a safety net |
| `FILE_TOO_LARGE` | 413 | POST /uploads | Show inline upload error |
| `UNSUPPORTED_FORMAT` | 415 | POST /uploads | Show inline upload error |
| `UPLOAD_FAILED` | 500 | POST /uploads | Safe to retry |
| `FILE_NOT_FOUND` | 404 | POST /checks | `file_id` expired or invalid |
| `INVALID_REQUEST` | 400 | POST /checks | Missing or malformed fields |
| `INVALID_LANDING_PAGE_URL` | 400 | POST /checks | URL provided is not a valid HTTPS URL. FE should validate client-side first (US-14); this is the server-side safety net. |
| `CHECK_NOT_FOUND` | 404 | GET /checks/{id} | Wrong ID or unauthorized |
| `EVALUATION_FAILED` | (in body) | GET /checks/{id} | `status: "failed"`; no HTTP error, credits released |

---

## Open Items

| # | Item | Owner | Needed by |
| --- | --- | --- | --- |
| API-OQ-1 | Maximum file size for uploads (image and video separately) | Engineering | Before FE can implement upload validation (US-05) |
| API-OQ-2 | Accepted file formats (e.g., JPG, PNG, MP4, MOV) | Engineering | Before FE can implement upload validation (US-05) |
| API-OQ-3 | Credit cost per check (maps to PRD OQ-1) | Product + Finance | Before `credits_reserved` value can be shown in UI |
| API-OQ-4 | Does a credit balance endpoint already exist? If so, what is the path and response shape? | Engineering | Before FE can implement credit guard (NFR-6) |
| API-OQ-5 | How does the deep-link in the async email authenticate the user? Session-based (user must be logged in) or token-based (one-time link)? | Engineering | Required for US-07 async email flow |
| API-OQ-6 | What is the maximum number of carousel cards the backend can process in one check? (FE currently limits to 10; confirm BE agrees) | Engineering | Alignment on carousel card limit (FR-2) |
| API-OQ-7 | How does the backend fetch landing pages? Plain HTTP GET with a standard user-agent, or a headless browser (e.g., Chromium via Playwright)? Many advertiser landing pages are JS-rendered, bot-protected (Cloudflare), or geo-restricted. Plain HTTP will return empty or blocked content for a significant portion of URLs. What is the fetch timeout? | Engineering | Determines reliability of landing page analysis. If plain HTTP is used, the feature will silently fail on a large share of real-world URLs. Headless browser adds infrastructure cost and complexity. |

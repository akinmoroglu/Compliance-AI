> **HOW TO USE THIS FILE**
>
> **Step 1 — Create the page structure in Confluence:**
> - Go to your Rockads space → click **+ Create**
> - Title: `Product Specs` → Save as a top-level page
> - Inside it, create child page: `Ad Compliance AI`
> - Inside that, create child page: `PRD — Ad Compliance Checker (Phase 1)`
>
> **Step 2 — Paste the content:**
> - Open the PRD page in edit mode
> - In the editor, type `/` → search **Markdown** → select "Insert Markdown"
> - Paste everything below the divider line
>
> **Step 3 — After pasting, do these manually in the Confluence editor:**
> - Select the metadata table at the top → click the table → apply **"Page Properties"** macro (wrap the table in it). This makes the metadata searchable across Confluence.
> - Select the JTBD quote block → format it as a **"Quote"** block (highlight → toolbar → Quote)
> - Add page labels: `prd`, `compliance`, `meta`, `phase-1`, `draft`
>
> **Delete this instruction block before sharing with your team.**

---

# PRD: Ad Compliance AI — Phase 1 (Meta)

| Field | Value |
|---|---|
| Product | Rockads Marketing Panel |
| Feature | Ad Compliance Checker |
| Phase | 1 of 3 |
| Status | 🟡 Draft |
| Author | Akin Moroglu |
| Last Updated | April 2026 |
| Reviewers | [Add names] |

---

## 1. Problem & Context

### Background

Rockads' core business is renting agency-tier ad accounts to agencies and DTC brands, allowing customers to run ads without the overhead of managing account health, credit limits, and platform restrictions.

In practice, the panel's active user base is segmented into two distinct groups with very different behaviors:

- **Resellers** (e.g., ROAS, Elite and similar): The primary active panel users. They engage regularly with the top-up flow to fund accounts on behalf of their clients. They use Rockads as a financial and account management layer, not as a marketing tool.
- **Direct advertisers (small players)**: Agencies and DTC brands who signed up for Rockads but largely do not use the panel beyond initial onboarding. They represent the untapped engagement opportunity.

Rockads is building SaaS features to increase panel activation among signed-up users and to position the product as a complete marketing platform. The Ad Compliance Checker is the second SaaS feature, following the existing ad automation feature. The goal is to give users — particularly direct advertisers — meaningful reasons to return to and engage with the panel regularly.

### The Problem

Ad account bans and rejection loops are the most damaging pain point for Rockads customers. A single non-compliant ad can trigger a policy strike; repeated violations can get the rented account flagged or suspended — which directly hurts Rockads' core product.

Customers currently have no systematic way to check their ad creative and copy for policy violations before submission. They rely on trial and error: submit, get rejected, guess what's wrong, resubmit. This loop wastes budget, delays campaigns, and erodes trust in the ad account provider.

### Job to be Done

> *"When I'm about to launch a Meta ad, I want to know if my creative or copy will get rejected — before I spend a cent — so I can fix it without losing time or burning my account's trust score."*

### Why Now

Meta's enforcement of ad policies has become increasingly automated and aggressive, particularly in high-risk verticals (health, finance, supplements, crypto) which are common among Rockads' customer base. The risk of non-compliance is higher than ever, making a compliance checker a high-value, timely addition.

---

## 2. Goals & Success Metrics

### Strategic Goal

Increase panel activation among signed-up users who are not regularly engaging with the product. The compliance checker, alongside the existing ad automation feature, should give direct advertisers a concrete reason to open the panel outside of account funding.

This is not about converting resellers into heavier SaaS users — their use case is different. The target is the broader base of direct advertiser accounts that signed up but have low or zero recurring panel engagement.

### Product Goals

| Goal | Definition | Target |
|---|---|---|
| G1 — Activated User Reach | % of inactive signed-up users (no activity in 30 days) who run at least one check within 60 days of launch | 20% |
| G2 — Repeat Engagement | % of first-time checkers who return for a second check within 30 days | 40% |
| G3 — SaaS Feature Compounding | Retention rate of users using both compliance checker AND ad automation vs. either alone | Hypothesis to validate post-launch |
| G4 — Revenue Signal | AI credit consumption per activated user | Track as leading indicator |

**Note:** G1 requires an outbound activation push (email, in-panel prompt). The feature alone will not reach dormant users passively. This is a go-to-market dependency — see OQ-9.

### What We Are Not Optimizing For

- Reseller top-up behavior (different user segment, different product surface)
- Raw number of checks run (easily inflated by power users)
- Zero ad rejections (we reduce risk; we cannot guarantee outcomes)
- Speed of compliance check (accuracy and trust take priority over latency in Phase 1)

### Key Signal to Watch

Do users who run a compliance check return to the panel within the following 7 days at a higher rate than those who don't? Targets (20%, 40%) are directional and should be revisited at the 2-week mark with real data.

---

## 3. Scope

### Phase 1 — In Scope

| Area | Detail |
|---|---|
| Platform | Meta (Facebook + Instagram) only |
| Compliance type | Platform policy compliance only |
| Ad formats | Single image, single video, carousel (multi-card) |
| Input | Visual upload + ad copy (primary text, headline, description) + optional landing page URL |
| Visual analysis | Includes OCR — text overlaid on images is extracted and evaluated |
| Output | Violations list with severity, policy reference, explanation, suggested fix |
| Session model | One-shot checks only — no history, no saved reports |
| Monetization | AI credit consumption per check (cost TBD) |
| Async video | Videos >30 seconds trigger async flow + email notification |

### Phase 1 — Explicitly Out of Scope

- Google Ads and TikTok compliance (Phase 2+)
- Brand compliance and legal compliance checking
- Check history and saved reports (Phase 2)
- Bulk upload / batch processing
- API access for compliance checks
- "Publish Ad" button on results screen (infrastructure exists; deferred to Phase 2)
- Automated fix application (suggestions shown; user applies manually)
- Public-facing compliance tool as acquisition funnel (PLG initiative — Phase 3)

### Phase Roadmap

| Phase | Scope |
|---|---|
| Phase 1 | Meta platform compliance · single & carousel ads · one-shot checks · panel-embedded |
| Phase 2 | "Publish Ad" from results · Google Ads · check history · saved reports |
| Phase 3 | TikTok · bulk processing · public-facing PLG funnel (compliance as lead magnet → signup → KYC) |

---

## 4. Functional Requirements

### FR-1: Platform & Context Selection

The user selects the ad platform and provides targeting context before uploading any creative.

- User selects Meta as the target platform. Google Ads and TikTok are visible but disabled with a "Coming Soon" label.
- User selects a target region (default: All Regions & Countries) and an age range (default: 25–34). These inputs refine the compliance evaluation — some Meta policies are stricter for certain geographies and age brackets (e.g., financial products, alcohol, health claims near minors).
- The system passes platform, region, and age range as context to the evaluation engine.

### FR-2: Ad Creative & Copy Input

**Single Image or Video:**
- User uploads one image or video via drag-and-drop or file picker
- User enters: Primary Text, Headline, Description (optional)

**Carousel:**
- User enters a Global Primary Text applying to the entire carousel
- User adds 2–10 carousel cards; each card has its own media upload, headline, and description
- Cards are presented horizontally; the active card is highlighted

**Landing Page URL (both formats):**
- User provides the destination URL the ad will link to (optional but strongly recommended)
- Meta's review system explicitly evaluates the landing page alongside the creative: "The products and services promoted in an ad must match those promoted on the landing page." A landing page mismatch is a common and non-obvious rejection trigger.
- If provided, the backend fetches the landing page and evaluates it as part of the check (see FR-4)
- If omitted, the check runs against creative and copy only; the UI should make clear that landing page evaluation is skipped

**Input validation:**
- At least one creative must be uploaded before the check can be submitted
- Accepted file types and size limits to be defined by Engineering (see OQ-4)
- Landing page URL, if provided, must be a valid and reachable HTTPS URL. Invalid or unreachable URLs should be surfaced as an inline warning before submission, not a blocking error.

### FR-3: Compliance Evaluation Trigger

The "Check Compliance" button is disabled until required inputs are provided.

- **Images and videos ≤30 seconds:** Synchronous. User sees processing state with animated indicators before results appear.
- **Videos >30 seconds:** Asynchronous. User sees a holding screen explaining that results will be emailed. System sends email on completion with a deep-link back to results.

### FR-4: AI Evaluation Engine (Backend)

The backend evaluates content against a hardcoded Meta policy knowledge base embedded in the prompt.

- **Visual analysis:** Evaluates images for policy-violating elements (e.g., before/after imagery, sensationalist imagery, implied medical outcomes). Includes OCR — text overlaid on images is extracted and treated as ad copy.
- **Copy analysis:** Primary text, headline, and description evaluated for prohibited claims, restricted language, and policy violations.
- **Landing page analysis (when URL is provided):** The backend fetches the landing page at check submission time and evaluates its content against the same policy knowledge base. Key checks include: (1) consistency between ad claims and landing page content, (2) prohibited claims on the destination page, (3) required disclosures present for regulated categories (financial, health, crypto). If the landing page is unreachable at fetch time, the check proceeds without it and the violation list includes a warning noting the page could not be evaluated.

Each violation in the response must contain:

| Field | Description |
|---|---|
| Violation title | Short label (e.g., "Misleading Health Claim") |
| Severity | High Risk / Medium Risk / Low Risk |
| Policy reference | Specific Meta policy violated (named, not just described) |
| Explanation | Why this content violates the policy, in plain language |
| Suggested fix | Concrete rewrite or modification the user can apply |

### FR-5: Results Display

Results are presented in a split-panel layout.

- **Left panel:** Simulated preview of the ad in a Meta feed using the uploaded creative and copy
- **Right panel:** All violations as cards, ordered High → Medium → Low
- Passed state shown with positive confirmation if no violations found
- "Start New Check" resets the entire flow
- **⚠️ Session-scoped:** Results are not recoverable after navigating away (no history in Phase 1). The UI must display a visible warning before the user leaves or resets.

### FR-6: Credit Consumption

- Each check consumes a defined number of AI credits from the user's Rockads balance
- If balance is insufficient, the system blocks the check and prompts top-up
- Credit consumption is recorded per check for billing and analytics

---

## 5. Non-Functional Requirements

| ID | Requirement | Detail |
|---|---|---|
| NFR-1 | Response time | Synchronous checks must return results in ≤30 seconds for 95% of requests. Processing state UI must appear immediately on submission. |
| NFR-2 | Async video flow | Videos >30s trigger async without exception. Email sent within 5 min of completion, including a deep-link to results. |
| NFR-3 | Accuracy baseline | No excessive false positives. Before launch: run 20 test cases (10 violating, 10 compliant) to establish qualitative baseline. Engineering/QA responsibility. |
| NFR-4 | Graceful failure | Backend failure or timeout → clear error + retry option. Credits not consumed on failed checks. Upload failures communicated inline. |
| NFR-5 | Platform fit | Embedded in Rockads panel, follows existing auth model. FE: Vue 3, shadcn-vue, Tailwind v4, Pinia, Vue Router. BE: Golang + Encore.dev + PostgreSQL, deployed on GCP Kubernetes (Autopilot). Async jobs via Encore native pub/sub. |
| NFR-6 | Credit guard | Credit balance checked before any backend call. Insufficient balance blocked at frontend, not after API call. |

---

## 6. Open Questions

| # | Question | Owner | Blocks |
|---|---|---|---|
| OQ-1 | What is the AI credit cost per compliance check? | Product + Finance | Credit guard UI, pricing communication |
| OQ-3 | Which AI model will be used for evaluation? | Engineering | Cost-per-check calculation, latency estimates, OCR accuracy |
| OQ-4 | Accepted file types and max upload size? | Engineering | Input validation (FR-2) |
| OQ-5 | Who owns the Meta policy knowledge base? Who updates it when Meta changes policies? | Product + Engineering | Risk of silent inaccuracy post-launch |
| OQ-6 | Which email system is used for async video notifications? | Engineering | NFR-2 async flow |
| OQ-7 | Is there a rate limit per user beyond credit consumption? | Product | Abuse prevention, infra cost modeling |
| OQ-8 | Full rollout or segmented rollout on launch day? | Product + Growth | G1 metric baseline |
| OQ-9 | What is the activation channel for inactive users? (email, push, in-panel prompt) | Product + Growth | G1 is unreachable without outbound strategy |
| OQ-10 | Should the ad automation feature and compliance checker be cross-promoted within the panel? (e.g., after a compliance check passes, prompt user to set up automation) | Product | If yes, this changes the post-results UI and represents an opportunity to compound SaaS engagement |
| OQ-11 | How should the backend handle landing page fetch for slow or bot-protected pages? Some advertiser landing pages use Cloudflare protection, JavaScript-only rendering, or geo-blocks that will cause the fetch to fail or return empty content. Does the backend use a headless browser or a plain HTTP fetch? What is the timeout? | Engineering | Determines reliability of landing page analysis and fallback behavior |
| OQ-10 | Should compliance checker and ad automation cross-promote within the panel? | Product | Post-results UI, SaaS compounding opportunity |

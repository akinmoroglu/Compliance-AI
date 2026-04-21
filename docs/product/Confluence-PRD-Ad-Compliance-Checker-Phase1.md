> **Source:** https://pozitim.atlassian.net/wiki/spaces/ROC/pages/3403350080
> **Last updated:** 2026-04-13 (v4) — "Add landing page URL to FR-2, FR-4, scope, and OQ-11"

# PRD: Ad Compliance AI: Phase 1 (Meta)

| Field | Value |
| --- | --- |
| Product | Rockads Marketing Panel |
| Feature | Ad Compliance Checker |
| Phase | 1 of 3 |
| Status | Draft |
| Author | Akin Moroglu |
| Last Updated | April 2026 |

---

## 1. Problem & Context

### Background

Rockads' core business is renting agency-tier ad accounts to agencies and DTC brands, allowing customers to run ads without the overhead of managing account health, credit limits, and platform restrictions.

In practice, the panel's active user base is segmented into two distinct groups with very different behaviors:

* **Resellers** (e.g., ROAS, Elite and similar): These are the primary active panel users. They engage regularly with the top-up flow to fund accounts on behalf of their clients. They use Rockads as a financial and account management layer, not as a marketing tool.
* **Direct advertisers (small players)**: These are agencies and DTC brands who signed up for Rockads but largely do not use the panel beyond initial onboarding. They represent the untapped engagement opportunity.

Rockads is building SaaS features to increase panel activation among signed-up users and to position the product as a complete marketing platform. The Ad Compliance Checker is the second SaaS feature, following the existing ad automation feature. The goal is to give users (particularly direct advertisers) meaningful reasons to return to and engage with the panel regularly.

### The Problem

Ad account bans and rejection loops are the most damaging pain point for Rockads customers. A single non-compliant ad can trigger a policy strike; repeated violations can get the rented account flagged or suspended; this directly hurts Rockads' core product.

Customers currently have no systematic way to check their ad creative and copy for policy violations before submission. They rely on trial and error: submit, get rejected, guess what's wrong, resubmit. This loop wastes budget, delays campaigns, and erodes trust in the ad account provider.

### Job to be Done

> "When I'm about to launch a Meta ad, I want to know if my creative or copy will get rejected (before I spend a cent) so I can fix it without losing time or burning my account's trust score."

### Why Now

Meta's enforcement of ad policies has become increasingly automated and aggressive, particularly in high-risk verticals (health, finance, supplements, crypto) which are common among Rockads' customer base. The risk of non-compliance is higher than ever, making a compliance checker a high-value, timely addition.

---

## 2. Goals & Success Metrics

### Strategic Goal

Increase panel activation among signed-up users who are not regularly engaging with the product. The compliance checker, alongside the existing ad automation feature, should give direct advertisers a concrete reason to open the panel outside of account funding.

This is not about converting resellers into heavier SaaS users; their use case is different. The target is the broader base of direct advertiser accounts that signed up but have low or zero recurring panel engagement.

### Product Goals

**G1: Activated User Reach:** At least 20% of previously inactive signed-up users (defined as: registered but no panel activity in the last 30 days) run at least one compliance check within 60 days of launch. This requires an activation push (e.g., email, in-panel prompt); the feature alone will not reach dormant users passively.

**G2: Repeat Engagement:** Of users who complete their first check, at least 40% return to run a second check within 30 days. Repeat usage distinguishes genuine value from one-time curiosity.

**G3: SaaS Feature Compounding:** Track whether users who use both the compliance checker and the ad automation feature show higher retention than users of either feature alone. If confirmed, the case for building more SaaS features is strengthened.

**G4: Revenue Signal:** Credit consumption per activated user validates willingness-to-pay for AI-powered features. This informs credit pricing decisions (currently unresolved, see OQ-1).

### What We Are Not Optimizing For

* Reseller top-up behavior (different user segment, different product surface)
* Raw number of checks run (easily inflated by a small number of power users)
* Zero ad rejections (we reduce risk; we cannot guarantee outcomes)
* Speed of compliance check (accuracy and trust take priority over latency in Phase 1)

### Measurement Notes

The 20% and 40% targets are directional. They should be baselined against actual inactive user counts after launch and revisited at the 2-week mark. The most important early signal: do users who run a compliance check return to the panel within the following 7 days at a higher rate than those who don't?

---

## 3. Scope

### Phase 1 — In Scope

* **Platform:** Meta (Facebook + Instagram) only
* **Compliance type:** Platform policy compliance only (no brand guidelines, no legal review)
* **Ad formats:** Single image, single video, carousel (multi-card)
* **Input:** Visual upload (image or video file) + ad copy fields (primary text, headline, description) + optional landing page URL
* **Visual analysis:** Includes OCR; text overlaid on images is extracted and evaluated alongside copy
* **Output:** List of policy violations, each with severity rating, the specific Meta policy rule violated, and a suggested rewrite or fix
* **Session model:** One-shot checks only. No history, no saved reports
* **Monetization:** AI credit consumption per check (credit cost TBD)
* **Async video processing:** Videos longer than 30 seconds trigger an async flow with email notification upon completion

### Phase 1 — Explicitly Out of Scope

* Google Ads compliance (Phase 2)
* TikTok compliance (Phase 2)
* Brand compliance checking (tone of voice, brand guidelines)
* Legal compliance checking (regional advertising law, GDPR, FTC disclosures)
* Check history and saved reports (Phase 2)
* Bulk upload / batch processing of multiple ads
* API access for compliance checks
* Integration with Meta Ads Manager (direct submission after check)
* Automated fix application (suggestions are shown; user applies them manually)
* "Publish Ad" button on results screen (infrastructure exists; deferred to Phase 2 to keep Phase 1 lean)
* Public-facing compliance tool as a user acquisition funnel (PLG initiative, separate project, post Phase 1)

### Phase Roadmap (for context)

| Phase | Scope |
| --- | --- |
| Phase 1 | Meta platform compliance, single & carousel ads, one-shot checks, panel-embedded |
| Phase 2 | "Publish Ad" from results screen, Google Ads, check history, saved reports |
| Phase 3 | TikTok, bulk processing, public-facing PLG funnel (compliance as lead magnet → signup → KYC/account connect) |

---

## 4. Functional Requirements

### FR-1: Platform & Context Selection

The user selects the ad platform and provides targeting context before uploading any creative.

* User selects Meta as the target platform. Google Ads and TikTok are visible but disabled with a "Coming Soon" label.
* User selects a target region (default: All Regions & Countries) and an age range (default: 25-34). These inputs refine the compliance evaluation; some Meta policies are stricter for certain geographies and age brackets (e.g., financial products, alcohol, health claims near minors).
* The system stores the selected platform, region, and age range and passes them as context to the evaluation engine.

### FR-2: Ad Creative & Copy Input

The user provides the full ad content to be evaluated.

**Ad Format — Single Image or Video:**

* User uploads one image or video file via drag-and-drop or file picker.
* User enters: Primary Text, Headline, Description (optional).

**Ad Format — Carousel:**

* User enters a Global Primary Text that applies to the entire carousel.
* User adds between 2 and 10 carousel cards.
* Each card has its own: media upload (image or video), headline, and description.
* Cards are presented horizontally; the active card being edited is highlighted.

**Landing Page URL (both formats):**

* User provides the destination URL the ad will link to (optional but strongly recommended).
* Meta's review system explicitly evaluates the landing page alongside the creative: "The products and services promoted in an ad must match those promoted on the landing page." A landing page mismatch is a common and non-obvious rejection trigger.
* If provided, the backend fetches the landing page and evaluates it as part of the check (see FR-4).
* If omitted, the check runs against creative and copy only; the UI should make clear that landing page evaluation is skipped.

**Input validation:**

* At least one creative (image or video) must be uploaded before the check can be submitted.
* Accepted file types and size limits are defined by the engineering team based on what the AI evaluation backend can process.
* Landing page URL, if provided, must be a valid and reachable HTTPS URL. Invalid or unreachable URLs should be surfaced as an inline warning before submission, not a blocking error.

### FR-3: Compliance Evaluation Trigger

The user initiates the check via a "Check Compliance" button, which is disabled until required inputs are provided.

* For images and short videos (30 seconds or shorter): the check runs synchronously. The user sees a processing state with animated status indicators before results appear.
* For videos longer than 30 seconds: the check runs asynchronously. The user is immediately shown a holding screen explaining that the check is in progress and they will receive an email notification when results are ready. The system sends this email upon completion.

### FR-4: AI Evaluation Engine (Backend)

The backend evaluates the submitted creative and copy against a hardcoded Meta policy knowledge base embedded in the prompt.

* **Visual analysis:** The AI evaluates image content for policy-violating visual elements (e.g., before/after imagery, sensationalist imagery, implied medical outcomes). This includes OCR; text overlaid on images is extracted and treated as ad copy.
* **Copy analysis:** Primary text, headline, and description are evaluated for prohibited claims, restricted language, and policy violations.
* **Landing page analysis (when URL is provided):** The backend fetches the landing page at check submission time and evaluates its content against the same policy knowledge base. Key checks include: (1) consistency between ad claims and landing page content, (2) prohibited claims on the destination page, (3) required disclosures present for regulated categories (financial, health, crypto). If the landing page is unreachable at fetch time, the check proceeds without it and the violation list includes a warning noting the page could not be evaluated.
* The evaluation produces a structured response per violation, containing:

    * **Violation title:** Short label (e.g., "Misleading Health Claim")
    * **Severity:** One of three levels: High Risk, Medium Risk, Low Risk
    * **Policy reference:** The specific Meta policy being violated (named, not just described)
    * **Explanation:** Why this content violates the policy, in plain language
    * **Suggested fix:** A concrete rewrite or modification the user can apply

### FR-5: Results Display

Results are presented in a split-panel layout.

* **Left panel:** A simulated preview of how the ad would appear in a Meta feed (Facebook or Instagram), using the uploaded creative and entered copy.
* **Right panel:** A list of all detected violations, each displayed as a card with severity badge, policy reference, explanation, and suggested fix.
* Violations are ordered by severity (High → Medium → Low).
* If no violations are found, a "passed" state is shown with a positive confirmation message.
* A "Start New Check" action resets the entire flow.
* Results are session-scoped. If the user navigates away from the results page, results are not recoverable (no history in Phase 1). The UI should make this clear with a visible warning before the user leaves or resets.

### FR-6: Credit Consumption

* Each compliance check consumes a defined number of AI credits from the user's Rockads balance.
* The credit cost per check is to be defined (see Open Questions).
* If the user has insufficient credits, the system blocks the check and prompts the user to top up.
* Credit consumption is recorded per check for billing and analytics purposes.

---

## 5. Non-Functional Requirements

| # | Requirement | Detail |
| --- | --- | --- |
| NFR-1 | Response Time | Synchronous checks (images and videos 30s or shorter) must return results within 30 seconds for at least 95% of requests. Processing state UI must appear immediately upon submission. |
| NFR-2 | Async Video Flow | Videos over 30 seconds must trigger the async path without exception. Email notification must be sent within 5 minutes of result completion, including a direct deep-link back to the results view. |
| NFR-3 | Accuracy | No excessive false positives. Before launch, at least 20 test cases (10 violating, 10 compliant) must be run to establish a qualitative accuracy baseline. |
| NFR-4 | Graceful Failure | Backend failures or timeouts must show a clear error with retry option. Credits must not be consumed on failed checks. File upload errors must be communicated inline. |
| NFR-5 | Embedding & Platform Fit | Feature is embedded in the Rockads panel. Frontend: Vue 3 (Composition API), shadcn-vue, Tailwind CSS v4, Pinia, Vue Router. Backend: Golang on Encore.dev, PostgreSQL, GCP Kubernetes Autopilot. Async video uses Encore's native pub/sub. |
| NFR-6 | Credit Guard | AI credit balance must be checked before initiating any evaluation call. Insufficient balance must block the check at the frontend, not after the API call. |

---

## 6. Open Questions

| # | Question | Owner | Impact if Unresolved |
| --- | --- | --- | --- |
| OQ-1 | What is the AI credit cost per compliance check? | Product + Finance | Blocks credit guard implementation and any user-facing pricing communication |
| OQ-3 | Which AI model will be used for visual and copy evaluation? | Engineering | Affects cost-per-check calculation, latency estimates, and OCR accuracy |
| OQ-4 | What are the accepted file types and max file size for uploads? | Engineering | Required for input validation (FR-2) and user-facing upload constraints |
| OQ-5 | How will the Meta policy knowledge base be maintained? Who is responsible for updating it when Meta changes policies? | Product + Engineering | Stale policies = false confidence; high compliance risk for customers |
| OQ-6 | What email system is used for async video notifications? | Engineering | Required for NFR-2 async flow |
| OQ-7 | Is there a rate limit per user per day/month, beyond credit consumption? | Product | Affects abuse prevention and infrastructure cost modeling |
| OQ-8 | Will the compliance checker be surfaced in the panel navigation for all users immediately, or rolled out to a segment first? | Product + Growth | Affects G1 (adoption) metric baseline and rollout plan |
| OQ-9 | How will inactive users be notified that the compliance checker exists? A passive in-panel placement will not reach users who never open the panel. Is there a planned activation channel (email campaign, push notification, re-engagement flow)? | Product + Growth | G1 target is unreachable without an outbound activation strategy; this is a go-to-market dependency, not just a product decision |
| OQ-10 | Should the ad automation feature and compliance checker be cross-promoted within the panel? (e.g., after a compliance check passes, prompt user to set up automation) | Product | If yes, this changes the post-results UI and represents an opportunity to compound SaaS engagement |
| OQ-11 | How should the backend handle landing page fetch for slow or bot-protected pages? Some advertiser landing pages use Cloudflare protection, JavaScript-only rendering, or geo-blocks that will cause the fetch to fail or return empty content. Does the backend use a headless browser or a plain HTTP fetch? What is the timeout? | Engineering | Determines reliability of landing page analysis and fallback behavior |

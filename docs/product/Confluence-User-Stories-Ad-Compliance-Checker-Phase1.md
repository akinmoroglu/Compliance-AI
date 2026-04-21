> **Source:** https://pozitim.atlassian.net/wiki/spaces/ROC/pages/3404234791
> **Last updated:** 2026-04-13 (v2) — "Add US-14: landing page URL input"

# User Stories: Ad Compliance Checker (Phase 1)

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

## How to Read This Document

Each story follows the format:

> As a **[user]**, when **[context]**, I want to **[action]**, so that **[benefit]**.

Acceptance criteria follow the Given/When/Then pattern. A story is "done" only when all criteria pass.

**Priority (MoSCoW):**

* **M** = Must Have (Phase 1 is broken without it)
* **S** = Should Have (important but has a workaround)
* **C** = Could Have (nice to have; cut if time is tight)

---

## Epic 1: Platform and Context Selection

### US-01: Select Target Platform

**Priority:** M

As a panel user, when I open the compliance checker, I want to select Meta as my target platform so that I know my ad will be evaluated against Meta's specific policies.

**Acceptance Criteria:**

_Given_ I am on the compliance checker page,
_When_ I see the platform selection step,
_Then_ Meta is shown as an active, selectable option and Google Ads and TikTok are visible but disabled with a "Coming Soon" label.

_Given_ I have not selected a platform,
_When_ I try to proceed to the next step,
_Then_ the "Next Step" button is disabled and I cannot advance.

_Given_ I select Meta,
_When_ I proceed,
_Then_ my platform selection is retained throughout the session and passed to the evaluation engine.

---

### US-02: Set Target Region and Age Range

**Priority:** M

As a panel user, when I am setting up my compliance check, I want to specify a target region and age range so that the evaluation accounts for geography and demographic-specific policy restrictions (e.g., stricter rules on financial or health content near minors).

**Acceptance Criteria:**

_Given_ I am on the platform selection step,
_When_ the step loads,
_Then_ the region defaults to "All Regions & Countries" and the age range defaults to 25-34.

_Given_ I change the region or age range,
_When_ I proceed to the next step,
_Then_ my selections are retained and passed as context to the evaluation engine alongside the platform selection.

_Given_ I do not change the defaults,
_When_ the check runs,
_Then_ the evaluation uses "All Regions & Countries" and age range 25-34 as context.

---

## Epic 2: Ad Creative and Copy Input

### US-03: Upload a Single Image or Video and Enter Ad Copy

**Priority:** M

As a panel user, when I want to check a standard single-format ad, I want to upload one image or video and enter my ad copy so that the system has everything it needs to evaluate the full ad.

**Acceptance Criteria:**

_Given_ I am on the creative input step with "Single Image or Video" selected,
_When_ I upload a file via drag-and-drop or file picker,
_Then_ a preview of the uploaded file is shown in the upload area.

_Given_ I have uploaded a file,
_When_ I enter Primary Text, Headline, and optionally Description,
_Then_ all entered text is retained and associated with the uploaded creative.

_Given_ I have not uploaded any file,
_When_ I try to click "Check Compliance",
_Then_ the button is disabled and I cannot proceed.

_Given_ I have uploaded a file but left all copy fields empty,
_When_ I try to click "Check Compliance",
_Then_ the button is enabled (copy is optional; the visual alone can be evaluated).

_Given_ I am on the creative input step,
_When_ the step loads,
_Then_ a "Landing Page URL" field is visible (see US-14 for full acceptance criteria).

---

### US-04: Build and Edit a Carousel Ad

**Priority:** M

As a panel user, when I want to check a carousel-format ad, I want to add multiple cards with individual media and copy so that each card is evaluated as part of the full carousel.

**Acceptance Criteria:**

_Given_ I select "Carousel Format",
_When_ the format loads,
_Then_ one empty card is shown and a Global Primary Text field is available.

_Given_ I am in carousel mode,
_When_ I click the "+" button,
_Then_ a new empty card is added; I can add up to 10 cards total.

_Given_ I click on a card,
_When_ the card editor opens,
_Then_ I can upload media and enter a Headline and Description specific to that card.

_Given_ I have not uploaded media for the first carousel card,
_When_ I try to click "Check Compliance",
_Then_ the button is disabled.

_Given_ I have uploaded media for at least the first card,
_When_ I click "Check Compliance",
_Then_ the check proceeds with all cards submitted together.

---

### US-14: Provide a Landing Page URL for Full Evaluation

**Priority:** M

As a panel user, when I am setting up my compliance check, I want to optionally enter the landing page URL my ad links to so that the system can verify the destination page is also compliant (and that the ad's claims match what the page actually says).

**Acceptance Criteria:**

_Given_ I am on the creative input step,
_When_ the step loads,
_Then_ a "Landing Page URL" input field is visible below the ad copy fields with a label indicating it is optional but recommended.

_Given_ I enter a valid HTTPS URL,
_When_ I move to the next field or submit,
_Then_ no validation error is shown and the URL is included in the check submission.

_Given_ I enter a URL that is not a valid HTTPS URL (e.g., missing protocol, malformed),
_When_ I attempt to submit the check,
_Then_ an inline validation warning appears near the URL field explaining the issue; the check is not blocked but the user is informed the URL will not be evaluated.

_Given_ I leave the landing page URL field empty,
_When_ I submit the check,
_Then_ the check proceeds using creative and copy only; the results screen includes a visible notice that landing page analysis was skipped because no URL was provided.

_Given_ I provide a URL that is unreachable at the time of the check,
_When_ results are returned,
_Then_ the violation list includes a non-blocking notice ("Landing page could not be fetched — URL analysis skipped") and the rest of the results remain valid.

---

### US-05: See Accepted File Constraints

**Priority:** S

As a panel user, when I am uploading a creative, I want to see the accepted file types and maximum file size so that I do not waste time uploading an incompatible file.

**Acceptance Criteria:**

_Given_ I am on the creative input step,
_When_ the upload area is visible,
_Then_ accepted file formats and the maximum file size are displayed near the upload zone.

_Given_ I upload a file that exceeds the size limit or uses an unsupported format,
_When_ the upload is attempted,
_Then_ an inline error message explains the issue without clearing any other input I have entered.

---

## Epic 3: Compliance Check Execution

### US-06: See a Processing State During Evaluation

**Priority:** M

As a panel user, when I submit my ad for compliance checking, I want to see a clear processing state so that I know the system is working and I do not think it has crashed.

**Acceptance Criteria:**

_Given_ I click "Check Compliance" and my video is 30 seconds or shorter,
_When_ the evaluation begins,
_Then_ a processing screen appears immediately with animated status indicators; I am not shown a frozen screen at any point.

_Given_ the evaluation completes successfully,
_When_ results are ready,
_Then_ the processing screen transitions automatically to the results view.

---

### US-07: Receive Email When a Long Video Check Completes

**Priority:** M

As a panel user, when I submit a video longer than 30 seconds, I want to be informed that the check will take longer and receive an email when it is done so that I can close the browser and return when the results are ready.

**Acceptance Criteria:**

_Given_ I upload a video longer than 30 seconds and click "Check Compliance",
_When_ the check is triggered,
_Then_ I am immediately shown a holding screen that explains the check is processing and that I will receive an email with a link to the results.

_Given_ the async evaluation completes,
_When_ results are ready,
_Then_ an email is sent to my registered address within 5 minutes, containing a direct link that takes me back to the results view in the panel.

_Given_ I click the link in the email,
_When_ the page loads,
_Then_ my results are shown without requiring me to re-submit.

---

### US-08: Be Blocked If Credits Are Insufficient

**Priority:** M

As a panel user, when I try to run a compliance check without sufficient AI credits, I want to be clearly told and directed to top up so that I understand why the check did not run.

**Acceptance Criteria:**

_Given_ my AI credit balance is below the cost of one compliance check,
_When_ I click "Check Compliance",
_Then_ the check does not run; I see an inline message explaining I have insufficient credits and a prompt to top up.

_Given_ I have sufficient credits,
_When_ I click "Check Compliance",
_Then_ the check runs and the credit cost is deducted from my balance upon completion.

_Given_ the check fails due to a backend error (not insufficient credits),
_Then_ credits are not deducted and I am shown an error with a retry option.

---

## Epic 4: Results

### US-09: See Detailed Violations with Severity and Suggested Fixes

**Priority:** M

As a panel user, when the compliance check is complete, I want to see each policy violation with its severity level, the specific Meta policy it references, an explanation, and a suggested fix so that I know exactly what to change before resubmitting.

**Acceptance Criteria:**

_Given_ the evaluation finds one or more violations,
_When_ results are shown,
_Then_ each violation is displayed as a card containing: violation title, severity badge (High Risk / Medium Risk / Low Risk), the specific Meta policy violated (by name), a plain-language explanation, and a concrete suggested rewrite or fix.

_Given_ multiple violations are found,
_When_ results are displayed,
_Then_ violations are ordered from High Risk to Medium Risk to Low Risk.

_Given_ the evaluation finds no violations,
_When_ results are shown,
_Then_ a clear "passed" confirmation state is shown with a positive message; no violation cards are shown.

---

### US-10: See a Simulated Ad Preview Alongside Results

**Priority:** S

As a panel user, when I am reviewing my compliance results, I want to see a simulated Meta feed preview of my ad so that I can visually connect the violations to the actual ad content.

**Acceptance Criteria:**

_Given_ results are displayed,
_When_ I view the results page,
_Then_ the left panel shows a simulated Facebook or Instagram feed card using my uploaded creative and entered copy.

_Given_ I submitted a carousel,
_When_ results are shown,
_Then_ the preview shows the carousel format, not a single card.

---

### US-11: Be Warned Before Leaving the Results Page

**Priority:** M

As a panel user, when I am about to navigate away from my results, I want to see a warning so that I do not accidentally lose my compliance report (which cannot be recovered in Phase 1).

**Acceptance Criteria:**

_Given_ I am on the results page,
_When_ I click "Start New Check" or attempt to navigate away,
_Then_ a confirmation warning is shown explaining that results will not be saved and asking me to confirm before proceeding.

_Given_ I confirm I want to leave,
_When_ I proceed,
_Then_ the results are cleared and I return to the start of the flow.

_Given_ I dismiss the warning,
_When_ I return to the results page,
_Then_ my results are still visible.

---

## Epic 5: Error Handling

### US-12: See a Recoverable Error If the Evaluation Fails

**Priority:** M

As a panel user, when the AI evaluation backend fails or times out, I want to see a clear error message and be able to retry so that I am not left confused and my credits are not wasted.

**Acceptance Criteria:**

_Given_ the evaluation backend fails or times out,
_When_ the error occurs,
_Then_ I am shown a clear error message (not a blank screen or spinner) with a "Try Again" option.

_Given_ the failure occurred after I submitted,
_When_ the error is shown,
_Then_ no AI credits have been deducted from my balance.

_Given_ I click "Try Again",
_When_ the retry is triggered,
_Then_ my previously uploaded creative and copy are retained; I do not need to re-enter anything.

---

### US-13: See an Inline Error If My File Upload Fails

**Priority:** M

As a panel user, when my file upload fails, I want to see an inline error near the upload zone so that I can fix the issue without losing any copy or other inputs I have already entered.

**Acceptance Criteria:**

_Given_ I attempt to upload a file and the upload fails,
_When_ the error occurs,
_Then_ an inline error message appears near the upload zone explaining what went wrong.

_Given_ the upload error appears,
_When_ I review the form,
_Then_ all previously entered copy fields (Primary Text, Headline, Description) remain intact.

_Given_ I resolve the upload issue and re-upload,
_When_ the new upload succeeds,
_Then_ the error message is dismissed and I can proceed normally.

---

## Out of Scope for Phase 1

The following are explicitly deferred. Do not build or design for these in this phase:

* Saving or accessing past compliance check results (history)
* Bulk ad upload
* "Publish Ad" from results screen
* Brand or legal compliance checking
* Google Ads or TikTok evaluation

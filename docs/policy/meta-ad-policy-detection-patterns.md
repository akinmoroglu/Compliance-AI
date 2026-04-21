# Meta Ad Policy — Violation Detection Patterns
**Document:** 2 of 2 — Detection Patterns for AI Agents  
**Compiled:** April 2026  
**Status:** Interpretive — NOT official Meta policy  
**Depends on:** `meta-ad-policy.md` (Document 1) — all section references (§) point there

> ⚠️ **IMPORTANT:** This document contains inferred detection patterns derived from Meta's official policies. It is not sourced verbatim from Meta. Every pattern references the Document 1 section it is based on. Patterns should be treated as strong indicators, not guarantees. Meta's automated review may apply additional undocumented logic.

---

## Severity Levels

All violations in this document are tagged with one of three severity levels:

| Level | Meaning | Expected Outcome |
|---|---|---|
| **CRITICAL** | Hard policy violation with no authorization path. Ad will be rejected. | Block submission or show blocking error |
| **WARNING** | Policy violation that will likely cause rejection or require authorization. | Show prominent warning, require acknowledgment |
| **ADVISORY** | Borderline, context-dependent, or jurisdiction-specific risk. | Show informational flag; allow submission with user acknowledgment |

---

## Agent 1: Category Detection (`category.go`)

**Purpose:** Identify which restricted or sensitive policy categories apply to this ad. Output feeds into all downstream agents. A single ad can belong to multiple categories.

**Input:** Headline, primary text, description, landing page text (if provided), ad format, targeting settings.

**Output:** Array of matched categories, each with confidence signal.

---

### CAT-01 · Alcohol

**Doc 1 ref:** §4.1

**Signals (any match = category active):**
- Copy contains product names or types: beer, wine, spirits, whiskey, vodka, rum, gin, tequila, champagne, prosecco, cocktail, lager, ale, cider, mead, sake, liqueur, brandy, bourbon, scotch
- Copy references: "happy hour", "drinks on us", "open bar", "brewery", "winery", "distillery", "pour", "cheers", "toast"
- Copy references a venue where alcohol is the primary business: bar, pub, nightclub, lounge (when associated with drinking)
- Logo or brand name of a known alcohol brand appears in copy or image alt text

**Exclusions:** Cooking with alcohol (non-promotional context), alcohol-free products.

---

### CAT-02 · Tobacco and Nicotine

**Doc 1 ref:** §4.2

**Signals:**
- Copy contains: cigarette, cigar, tobacco, nicotine, vape, vaping, e-cigarette, e-cig, vaporizer, hookah, shisha, snus, snuff, dip, chewing tobacco, heated tobacco, IQOS, Juul
- Copy references: "smoking", "smoke break", "roll your own", "rolling papers", "pipe tobacco"
- Brand names associated with tobacco or vaping products

**Exclusions:** Cessation products (nicotine patch, nicotine gum, Nicorette) — these are allowed with proper framing. If copy mentions cessation or quitting, treat as ADVISORY rather than CRITICAL.

---

### CAT-03 · Online Gambling and Games

**Doc 1 ref:** §4.3

**Signals:**
- Copy contains: casino, poker, slots, blackjack, roulette, baccarat, craps, sports betting, bet, wager, odds, sportsbook, lottery, lotto, bingo, scratch card, jackpot, spin to win, fantasy sports (with entry fee)
- Copy references: "win real money", "cash prizes", "play for money", "deposit bonus", "free spins", "welcome bonus", "first bet"
- Brand names of known gambling or betting platforms

**Exclusions:** Free-to-play games with no monetary value; skill competitions with no entry fee.

---

### CAT-04 · Dating Services

**Doc 1 ref:** §4.4

**Signals:**
- Copy contains: dating app, dating site, dating service, matchmaking, find singles, meet singles, online dating, dating profile, relationship app, "find your match", "meet someone", "swipe right"
- Landing page URL contains: match, date, tinder, bumble, hinge, eharmony, pof, zoosk, or similar

**Exclusions:** Social networking not primarily for romantic/sexual connection; professional networking.

---

### CAT-05 · Weapons and Related

**Doc 1 ref:** §4.5

**Signals:**
- Copy contains: firearm, gun, rifle, pistol, handgun, shotgun, AR-15, AK-47, ammo, ammunition, bullets, magazine, suppressor, silencer, scope, holster, 3D-printed gun, ghost gun, knife (non-culinary context), taser, pepper spray, explosive, grenade, landmine
- Copy references: "shooting range", "gun store", "self-defense weapon", "weapon accessories", "weapon modifications"

**Exclusions:** Firearms safety equipment (safes, locks), hunting/outdoor apparel (non-weapon), general sporting goods stores where weapons are not the primary focus.

---

### CAT-06 · Health and Wellness

**Doc 1 ref:** §4.8

**Signals — Weight loss sub-category:**
- Copy contains: weight loss, lose weight, diet pill, fat burner, slimming, keto supplement, intermittent fasting supplement, appetite suppressant, metabolism booster, BMI, "shed pounds", "drop dress size"
- Copy references specific dietary supplements, meal replacements, or weight management products

**Signals — Cosmetic procedures sub-category:**
- Copy contains: Botox, botulinum toxin, dermal filler, rhinoplasty, liposuction, tummy tuck, breast augmentation, blepharoplasty, facelift, labiaplasty, hair transplant, micro-needling, chemical peel, laser resurfacing, micropigmentation, anti-aging treatment
- Copy references cosmetic surgery clinics, aesthetic medicine, med spa

**Signals — Reproductive health sub-category:**
- Copy contains: condom, contraception, IUD, birth control, IVF, fertility treatment, ovulation, pregnancy test, abortion, erectile dysfunction, ED treatment, Viagra, Cialis, HSDD, premature ejaculation, femtech

---

### CAT-07 · Drugs and Pharmaceuticals

**Doc 1 ref:** §4.9

**Signals — Prescription drugs:**
- Copy contains specific prescription drug brand names or generic equivalents
- Copy references "prescription", "your doctor", "consult a physician", "licensed pharmacy", "telehealth prescription"
- Landing page URL is a pharmacy or telehealth platform domain

**Signals — Cannabis and CBD:**
- Copy contains: cannabis, marijuana, weed, THC, CBD, cannabidiol, hemp oil (with CBD), edibles, gummies (in cannabis context), dispensary, cannabis oil
- Copy references: "recreational use", "medical marijuana", "420", cannabis slang

**Signals — Illicit drugs:**
- Copy contains references to cocaine, heroin, MDMA, ecstasy, meth, methamphetamine, fentanyl, LSD, ketamine, psychedelics (outside research context), drug paraphernalia (bong, rolling papers in non-tobacco context)

---

### CAT-08 · Financial Products

**Doc 1 ref:** §4.10

**Signals — Prohibited products (detect for CRITICAL violation):**
- Copy contains: payday loan, paycheck advance, same-day cash, "cash until payday", bail bond, binary option, CFD, "contract for difference", ICO, "initial coin offering", penny auction

**Signals — Restricted products (detect for WARNING if no auth indicators):**
- Copy contains: credit card, personal loan, mortgage, insurance policy, investment opportunity, investment returns, APR, interest rate, refinance, debt consolidation, student loan refinancing, annuity

---

### CAT-09 · Cryptocurrency

**Doc 1 ref:** §4.11

**Signals:**
- Copy contains: Bitcoin, Ethereum, Solana, BNB, crypto, cryptocurrency, DeFi, blockchain, NFT, Web3, staking, yield farming, trading platform, crypto exchange, crypto wallet (with buying/selling/swapping/staking features), crypto lending, crypto mining
- Copy references: "buy crypto", "trade crypto", "earn crypto", "crypto returns", "invest in Bitcoin"

**Exclusions:** News or educational content about crypto without any offer or investment solicitation; storage-only wallets; tax services for crypto companies.

---

### CAT-10 · Addiction Treatment

**Doc 1 ref:** §4.7

**Signals:**
- Copy contains: rehab, rehabilitation, detox, addiction treatment, drug treatment, alcohol treatment, recovery program, sobriety, substance abuse, 12-step, intervention
- Copy references: addiction counseling, addiction therapy, withdrawal, recovery center, treatment facility

---

### CAT-11 · Social Issues, Elections, and Politics (SIEP)

**Doc 1 ref:** §7

**Signals:**
- Copy contains: election, vote, ballot, candidate, political party, senator, congressman, president, referendum, civic, democracy, political campaign
- Copy contains social issue framing: immigration policy, gun control, abortion rights, climate policy, racial justice, LGBTQ+ rights in a political framing
- Copy mentions political figures by name or title in a persuasive framing
- Copy contains "Paid for by" (indicates SIEP compliance attempt — check for completeness)

---

### CAT-12 · Hazardous Goods

**Doc 1 ref:** §4.6

**Signals:**
- Copy contains: hazardous material, explosive chemical, toxic substance, corrosive, flammable liquid, chemical weapon, dangerous goods (in sales context)

---

### CAT-13 · Historical Artifacts

**Doc 1 ref:** §4.12

**Signals:**
- Copy contains: historical artifact, antiquity for sale, ancient relic, archaeological item, pre-Columbian art (in sales context)

---

### CAT-14 · Animals

**Doc 1 ref:** §4.13

**Signals:**
- Copy references selling or trading live animals in a peer-to-peer context
- Copy references products derived from endangered species: ivory, rhino horn, shark fin, tiger parts, whale products

---

### CAT-15 · Human Body Parts

**Doc 1 ref:** §4.14

**Signals:**
- Copy contains: sell kidney, buy blood, buy organs, sell plasma (in sales context, not donation)

---

### CAT-16 · Crisis Exploitation

**Doc 1 ref:** §4.15

**Signals:**
- Copy references a current or recent named crisis (pandemic, named disaster, named conflict) in combination with sales language, urgency language, or product promotion
- Copy contains: "due to COVID", "amid the pandemic", "in light of [crisis name]" + sales/urgency framing

---

## Agent 2: Violation Detection (`violations.go`)

**Purpose:** Given the active categories (from Agent 1) and the full ad copy, detect specific policy violations. For each violation, output: rule ID, severity, matched text, and plain-language explanation.

**Input:** Full ad copy (headline + primary text + description + OCR-extracted image text), active categories from Agent 1, ad targeting settings.

**Output:** Array of violations, each with rule ID, severity, matched snippet, and explanation.

---

### SECTION A: Copy-Level Violations (All Categories)

---

#### V-001 · Personal Attribute Assertion
**Severity:** CRITICAL  
**Doc 1 ref:** §5.2  
**Handled by:** Dedicated Agent 3 (see below). Do not duplicate here — forward all personal attribute matches to Agent 3.

---

#### V-002 · Profanity
**Severity:** CRITICAL  
**Doc 1 ref:** §3.9  
**Detection:** Match against a profanity word list. Flag any occurrence in copy or OCR-extracted image text.  
**Note:** Include creative spellings, character substitutions (e.g., "sh!t", "f*ck"), and common euphemisms used to evade detection.  
**Fix signal:** Remove or replace with non-profane language. (→ Agent 5)

---

#### V-003 · Vaccine Discouragement
**Severity:** CRITICAL  
**Doc 1 ref:** §3.6  
**Patterns to match:**
- "vaccines cause [harm]", "vaccines are dangerous", "vaccine injuries", "don't get vaccinated"
- "anti-vax", "vaccine-free", "unvaccinated and proud"
- Implying vaccination is harmful or ineffective without medical basis
- "natural immunity is better than vaccines", "refuse the jab"

**Fix signal:** Remove anti-vaccine claims entirely. (→ Agent 5)

---

#### V-004 · Misinformation Indicators
**Severity:** WARNING  
**Doc 1 ref:** §6  
**Patterns:**
- Claims that contradict established scientific consensus (e.g., "COVID was a hoax", "climate change is fake")
- References to content known to be debunked by fact checkers
- "Doctors don't want you to know", "suppressed research", "what they're hiding from you"
- Conspiracy framing applied to medical, scientific, or electoral topics

**Note:** AI cannot definitively confirm fact-check status. Flag for human review.  
**Fix signal:** Remove unsupported claims; add evidence citations. (→ Agent 5)

---

#### V-005 · Fake Urgency / Scarcity
**Severity:** WARNING  
**Doc 1 ref:** §3.4, §3.5  
**Patterns:**
- "Only [N] left", "Limited stock", "Going fast", "Almost sold out"
- "Offer expires [tonight/today/in X hours]", "Last chance", "Flash sale ends at midnight"
- "Act now or miss out", "Don't wait", "This deal won't last"
- Countdown timers with false urgency framing

**Context rule:** Urgency is NOT inherently prohibited — only if used deceptively (e.g., "only 3 left" repeated indefinitely, or false deadlines). Flag as WARNING for human review.  
**Fix signal:** Remove or substantiate urgency claims. (→ Agent 5)

---

#### V-006 · Unsubstantiated Guaranteed Outcome Claims
**Severity:** WARNING  
**Doc 1 ref:** §3.4, §4.8.1, §4.10  
**Patterns:**
- "Guaranteed to [lose X kg / clear debt / earn X% returns / cure Y]"
- "100% results", "Always works", "Never fails"
- "Clinically proven to [cure/eliminate/prevent]" — without citing the study
- Absolute language implying certainty: "You WILL lose weight", "You WILL earn money"
- "Money-back guarantee" combined with outcome claims (guarantee on refund policy is fine; guarantee on outcomes is not)

**Context rule:** Heightened CRITICAL severity when combined with CAT-06 (health), CAT-07 (drugs), or CAT-08 (financial).

---

#### V-007 · Celebrity Endorsement Without Substantiation
**Severity:** WARNING  
**Doc 1 ref:** §3.4  
**Patterns:**
- "As seen on [TV show / news program]" — especially with fake logos
- "Endorsed by [celebrity name]" without verifiable endorsement
- "Recommended by [doctor / expert / celebrity]" without citation
- Fake media logos (CNN, BBC, etc.) in imagery or copy

**Fix signal:** Remove unsubstantiated endorsement claims; if real, add citation/link. (→ Agent 5)

---

#### V-008 · Miracle Cure / "Secret" Language
**Severity:** WARNING  
**Doc 1 ref:** §3.4, §4.8, §4.9  
**Patterns:**
- "This one trick / secret / hack", "What doctors don't want you to know"
- "Ancient remedy", "Forbidden formula", "Secret ingredient"
- "Miracle cure", "Miracle pill", "The cure they don't want you to find"
- "I lost 20 kg with this one weird trick"

**Context rule:** Elevate to CRITICAL if combined with active CAT-07 (drugs/pharmaceuticals) or claims about medical conditions.  
**Fix signal:** Replace with feature-benefit language without conspiratorial framing. (→ Agent 5)

---

### SECTION B: Health and Wellness Violations

**Applies when:** CAT-06 is active.

---

#### V-010 · Weight Loss Before/After Side-by-Side
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.1  
**Detection:**
- Copy explicitly describes a side-by-side comparison: "before and after", "see the transformation", "compare my results", "side-by-side results"
- Image metadata or alt text suggests before/after layout

**Exception rule:** If ad is for a fitness class (Pilates, yoga, CrossFit, gym), before/after is allowed. Check CAT-06 sub-type.  
**Fix signal:** Remove comparison framing; show product in use without comparison. (→ Agent 5)

---

#### V-011 · Weight Loss Body Fat Close-Up
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.1  
**Detection:**
- Copy references close-up imagery of body fat being pinched or grabbed: "shows exactly where the fat is", "pinch test", "love handle close-up"
- Image description or OCR references this visual style

**Fix signal:** Use full-body or product-focused imagery instead. (→ Agent 5)

---

#### V-012 · Cosmetic Anti-Aging Before/After Side-by-Side
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.1  
**Detection:**
- Ad is for Botox, dermal filler, anti-aging product, or wrinkle treatment AND copy contains before/after language or imagery

**Exception rule:** Close-up images of treated areas ARE allowed for anti-aging. Only the side-by-side comparison is prohibited.  
**Fix signal:** Use single-phase close-up showing result; remove comparison. (→ Agent 5)

---

#### V-013 · Negative Body Image / Body-Shaming Language
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.1  
**Patterns:**
- "Hate your belly", "embarrassed by your weight", "ashamed of your body"
- "Tired of being fat", "Sick of your love handles", "Disgusted by your appearance"
- "Your body is wrong", "Your skin is a problem", "Fix your flaws"
- "No one wants to see that", "Don't let [body part] ruin your confidence"
- Language implying the user's current appearance is a failure or defect

**Fix signal:** Reframe to positive/aspirational: "feel your best", "reach your goals", "invest in yourself". (→ Agent 5)

---

#### V-014 · Medical Cure/Treatment Claim for Non-Rx Product
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8, §4.9  
**Patterns:**
- "Cures [disease]", "treats [medical condition]", "prevents [disease]", "reverses [condition]", "eliminates [medical symptom]"
- "Clinically proven to treat", "medically proven to cure"
- "FDA-approved cure" (unless actually FDA-approved and properly licensed)
- Disease-specific claims for supplements, OTC products, or non-prescription items: "cures diabetes", "treats cancer", "prevents Alzheimer's"

**Context rule:** CRITICAL only for non-prescription products. If CAT-07 (prescription drugs) is active and authorization has been claimed, downgrade to WARNING.  
**Fix signal:** Replace with structure-function claims: "may support", "helps maintain", "formulated to assist". Add disclaimer. (→ Agent 5)

---

#### V-015 · Skin Whitening / Permanent Color Change
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.1  
**Patterns:**
- "Permanently whiten skin", "permanent skin lightening", "bleach your skin", "achieve lighter skin permanently"
- Products explicitly marketed to permanently alter skin color

**Fix signal:** Remove permanent skin color change claims; if product is a brightening/tone-evening cosmetic, reframe accordingly. (→ Agent 5)

---

#### V-016 · Adult Sexual Product / Service
**Severity:** CRITICAL  
**Doc 1 ref:** §4.8.2  
**Patterns:**
- Explicit references to: adult arousal products, sex toys, adult entertainment venues (strip clubs, adult theaters), "orgasmic therapy", "sexual fantasies hotel", G-spot augmentation, penile enlargement (for pleasure, not medical)
- Sexual services or escort-adjacent services

**Fix signal:** Remove prohibited product/service references. Reframe to health focus if applicable. (→ Agent 5)

---

### SECTION C: Financial Product Violations

**Applies when:** CAT-08 is active.

---

#### V-020 · Prohibited Financial Product (Hard Ban)
**Severity:** CRITICAL  
**Doc 1 ref:** §4.10.1  
**Patterns:**
- "Payday loan", "payday advance", "paycheck loan", "cash til payday"
- "Bail bond", "bail bondsman"
- Loans explicitly described as 90 days or less / same-day repayment
- "Binary option", "binary trading", "fixed return investment"
- "Contract for difference", "CFD trading", "CFD broker"
- "Initial coin offering", "ICO", "token sale", "token launch" (in financial investment context)
- "Penny auction", "bidding auction", "bid credits"

**Fix signal:** This product category cannot be advertised on Meta. Remove ad or replace with allowed financial product. (→ Agent 5)

---

#### V-021 · Investment Promotion via Direct Message (US)
**Severity:** CRITICAL (US-targeted ads only)  
**Doc 1 ref:** §4.10.4  
**Patterns:**
- "Message us to invest", "DM for investment details", "Contact us via Messenger for returns"
- Call to action directing users to private messaging for investment purposes

**Context rule:** Only flag for ads targeting US audiences. Check targeting settings.  
**Fix signal:** Remove DM-based investment CTA; redirect to landing page with proper disclosures. (→ Agent 5)

---

#### V-022 · Misleading Investment Return Claims
**Severity:** WARNING  
**Doc 1 ref:** §3.4, §4.10  
**Patterns:**
- "Guaranteed [X]% returns", "earn X% daily/weekly/monthly guaranteed"
- "Risk-free investment", "no-risk returns", "safe guaranteed profits"
- "Double your money in [timeframe]", "10x your investment"
- "Passive income guaranteed", "financial freedom guaranteed"

**Fix signal:** Replace guaranteed language with hedged language; add required financial disclosures. (→ Agent 5)

---

#### V-023 · Missing Financial Disclosure
**Severity:** WARNING  
**Doc 1 ref:** §4.10.2  
**Detection:** Ad is for a financial product (loan, credit card, insurance) AND contains specific financial terms (APR, interest rate, loan amount, repayment terms) WITHOUT visible disclosure language.

**Patterns indicating disclosure is present (suppress warning if matched):** "Subject to terms and conditions", "Representative APR", "Terms apply", "Representative example", legal disclaimer text

**Fix signal:** Add legally required disclosure text per applicable jurisdiction. (→ Agent 5)

---

#### V-024 · Deceptive Student Loan Claims
**Severity:** WARNING  
**Doc 1 ref:** §4.10.1  
**Patterns:**
- "Student loan forgiveness" — without referencing official government programs
- "Guaranteed student loan forgiveness", "eliminate student debt immediately"
- "Student loan relief" combined with urgency or fee-based framing

**Fix signal:** Reference official forgiveness programs by name; add disclaimer. (→ Agent 5)

---

### SECTION D: Tobacco and Alcohol Violations

---

#### V-030 · Prohibited Tobacco / Nicotine Product
**Severity:** CRITICAL  
**Doc 1 ref:** §4.2  
**Detection:** CAT-02 is active AND ad is promoting the product (not cessation).

**Cessation exception:** If copy prominently features quitting language AND product is nicotine patch, nicotine gum, or known FDA/WHO-approved cessation product — downgrade to ADVISORY.  
**Fix signal:** Cannot be advertised on Meta. Remove ad or replace with cessation product framing if applicable. (→ Agent 5)

---

#### V-031 · Alcohol in Prohibited Country
**Severity:** CRITICAL  
**Doc 1 ref:** §4.1  
**Detection:** CAT-01 is active AND targeting includes any country from the prohibited list: Afghanistan, Brunei, Bangladesh, Egypt, Ghana, India, Iran, Kuwait, Libya, Lithuania, Nigeria, Norway, Pakistan, Russia, Saudi Arabia, Sri Lanka, Turkey, UAE, Yemen.  
**Fix signal:** Exclude prohibited countries from targeting. (→ Agent 5)

---

#### V-032 · Alcohol Targeting Under 18
**Severity:** CRITICAL  
**Doc 1 ref:** §4.1  
**Detection:** CAT-01 is active AND minimum age in targeting settings is below 18.  
**Fix signal:** Increase minimum age targeting to 18+. (→ Agent 5)

---

### SECTION E: Gambling Violations

---

#### V-040 · Gambling Without Authorization Indicator
**Severity:** WARNING → escalates to CRITICAL if confirmed unauthorized  
**Doc 1 ref:** §4.3  
**Detection:** CAT-03 is active AND no authorization indicator present.

**Authorization indicators (suppress warning if matched):** Ad account level authorization flag (if accessible); ad explicitly references licensed operation in copy ("licensed by [regulator]").  
**Fix signal:** Obtain Meta authorization before running. (→ Agent 5)

---

#### V-041 · Gambling Targeting Under 18
**Severity:** CRITICAL  
**Doc 1 ref:** §4.3  
**Detection:** CAT-03 is active AND minimum age in targeting settings is below 18.  
**Fix signal:** Increase minimum age targeting to 18+. (→ Agent 5)

---

#### V-042 · Gambling in Unsupported Market
**Severity:** CRITICAL  
**Doc 1 ref:** §4.3  
**Detection:** CAT-03 is active AND targeting includes any of the 100+ unsupported gambling markets listed in Doc 1 §4.3.  
**Fix signal:** Exclude all unsupported gambling markets from geo-targeting. (→ Agent 5)

---

#### V-043 · Deceptive Gambling Claims
**Severity:** CRITICAL  
**Doc 1 ref:** §4.3  
**Patterns:**
- "Guaranteed wins", "beat the house guaranteed", "rigged in your favor"
- "Winning system", "never lose strategy", "foolproof betting system"

**Fix signal:** Remove deceptive gambling claims. (→ Agent 5)

---

### SECTION F: Drug Violations

---

#### V-050 · Illicit Drug Promotion
**Severity:** CRITICAL  
**Doc 1 ref:** §4.9.1  
**Detection:** CAT-07 (illicit drugs signals) active AND copy is promotional (not educational/news/advocacy).  
**Fix signal:** Cannot be advertised on Meta. Remove. (→ Agent 5)

---

#### V-051 · Drug Paraphernalia
**Severity:** CRITICAL  
**Doc 1 ref:** §4.9.1, §4.2  
**Patterns:**
- "Buy bong", "rolling papers for sale", "vaporizer for weed", "dab rig", "grinder for herb"
- Products described for use with illicit drugs or tobacco/nicotine (outside cessation context)

**Fix signal:** Remove product listing or reframe for clearly non-drug use. (→ Agent 5)

---

#### V-052 · CBD Without Authorization Indicator (US)
**Severity:** WARNING  
**Doc 1 ref:** §4.9.4  
**Detection:** Copy contains CBD, cannabidiol, or hemp-derived CBD AND ad targets US audience AND no LegitScript/authorization indicator present.  
**Fix signal:** Obtain LegitScript certification + Meta written authorization before running. (→ Agent 5)

---

#### V-053 · CBD Health Claim
**Severity:** CRITICAL  
**Doc 1 ref:** §4.9.4  
**Patterns:**
- "CBD cures [condition]", "CBD treats anxiety/depression/pain", "CBD prevents seizures"
- Any disease treatment, cure, or prevention claim applied to CBD or hemp products

**Fix signal:** Remove health/disease claims. CBD can only be described by product features, not medical outcomes. (→ Agent 5)

---

#### V-054 · THC / Psychoactive Cannabis Product
**Severity:** CRITICAL  
**Doc 1 ref:** §4.9.4  
**Detection:** Copy contains THC, psychoactive, "gets you high", "recreational cannabis", "cannabis edibles" (with THC context), "marijuana flower" for sale.  
**Fix signal:** Cannot be advertised on Meta. Remove. (→ Agent 5)

---

### SECTION G: SIEP Violations

---

#### V-060 · SIEP Without Authorization Indicator
**Severity:** WARNING  
**Doc 1 ref:** §7  
**Detection:** CAT-11 is active AND no "Paid for by" disclaimer or Meta authorization indicator found in copy.  
**Fix signal:** Complete Meta SIEP authorization; add "Paid for by [name]" disclaimer to all ads. (→ Agent 5)

---

#### V-061 · False Voting Information
**Severity:** CRITICAL  
**Doc 1 ref:** §7  
**Patterns:**
- Incorrect date, location, or method of voting: "Vote on [wrong date]", "Voting has been extended to [wrong date]"
- "Voting is cancelled", "polls are closed", "the election has been moved"
- False census participation information

**Fix signal:** Remove false voting/census information. (→ Agent 5)

---

#### V-062 · Deepfake / Deceptive AI-Generated Political Content
**Severity:** CRITICAL  
**Doc 1 ref:** §7  
**Patterns:**
- Copy references AI-generated or fabricated content featuring real political figures without disclosure
- "Realistic video of [politician name] saying..." without clear AI label

**Note:** Visual detection requires image analysis beyond copy scan. Flag copy signals only.  
**Fix signal:** Label AI-generated content clearly; do not fabricate statements by real political figures. (→ Agent 5)

---

### SECTION H: Suicide, Self-Harm, and Eating Disorder Violations

---

#### V-070 · Suicide/Self-Harm Encouragement
**Severity:** CRITICAL  
**Doc 1 ref:** §5.5  
**Patterns:**
- Explicit encouragement of self-harm or suicide — any form
- "Death as an escape", "end your suffering permanently", "there's no way out"
- Meme-style copy with dark depression/hopelessness framing: "just end it", "numb the pain forever"
- References to suicide methods or self-harm methods

**Fix signal:** Remove all such content entirely. (→ Agent 5)

---

#### V-071 · Eating Disorder Content
**Severity:** CRITICAL  
**Doc 1 ref:** §5.5  
**Patterns:**
- "Extreme fasting", "starvation diet", "eating as little as possible to lose weight fast"
- "Thigh gap", "collar bone showing", "rib visibility" framed as goals or achievements
- "Before/after" showing extreme weight loss with visible bones or hollow appearance in recovery context
- "Admission of extreme weight loss behavior"

**Fix signal:** Remove content referencing extreme restriction behaviors or body ideals. (→ Agent 5)

---

### SECTION I: Adult Nudity and Sexual Content Violations

---

#### V-080 · Sexually Suggestive Language or Imagery Description
**Severity:** CRITICAL  
**Doc 1 ref:** §5.1  
**Patterns in copy:**
- Explicit or implicit sexual language, innuendo, or emoji used suggestively (🍆, 🍑, 💦 in sexual context)
- "Sexually satisfying", "you'll love what happens in bed", "spice up your sex life" (as a promotional hook)
- References to sexual acts or positions in promotional copy
- "Provocative", "seductive", "will turn you on" framing

**Fix signal:** Remove sexually suggestive language; reframe for health/wellness focus if applicable. (→ Agent 5)

---

#### V-081 · Excessive Focus on Body Parts
**Severity:** WARNING  
**Doc 1 ref:** §5.1  
**Patterns in copy or image description:**
- Copy references or imagery emphasizes abs, buttocks, chests, or other body parts in a sexualized way (distinct from fitness/health context)

**Context rule:** Fitness product showing abs in workout context = acceptable. Same imagery in a sexual framing = WARNING.  
**Fix signal:** Use full-body imagery or health-focused framing. (→ Agent 5)

---

### SECTION J: Violence and Graphic Content Violations

---

#### V-090 · Shocking / Graphic / Sensational Content Description
**Severity:** CRITICAL  
**Doc 1 ref:** §5.4  
**Patterns in copy or image description:**
- References to: exposed organs, open wounds, graphic medical imagery, visible bones or innards
- "Graphic images of [medical condition]", "shocking before/after medical procedure"
- Firearms, bladed weapons, or explosives depicted as pointed at the viewer (described in copy or alt text)

**Fix signal:** Replace with clinically neutral imagery; avoid graphic medical visuals. (→ Agent 5)

---

### SECTION K: Bullying and Harassment Violations

---

#### V-100 · Degrading or Shaming Content
**Severity:** CRITICAL  
**Doc 1 ref:** §5.3  
**Patterns:**
- Copy targeting a named public or private individual with degrading language
- "Loser", "disgusting", "pathetic", "shameful" applied to a specific person or group in an ad context
- Content clearly designed to shame a specific individual

**Fix signal:** Remove targeting of individuals; reframe to aspirational or factual language. (→ Agent 5)

---

## Agent 3: Personal Attribute Detection (`personal_attribute.go`)

**Purpose:** Detect copy that asserts or implies personal attributes of the ad's audience. This is the single most commonly violated rule in Meta advertising and requires dedicated detection. Output: flag status, matched text, attribute type, severity, and suggested rewrite.

**Doc 1 ref:** §5.2

**Severity:** CRITICAL for all direct assertions; WARNING for implied or indirect assertions.

---

### The Core Test

**A violation occurs when ad copy implies or directly states that the reader (user who will see the ad) HAS a specific personal attribute.** The copy does not need to name the reader — it is sufficient that the framing assumes the reader belongs to a sensitive group.

**Not a violation:** Describing who the product serves in general terms ("for people who want to improve their fitness") without implying the reader has a condition.

**Violation:** Assuming the reader has a condition ("Are you struggling with diabetes?") or attribute ("As a fellow [group]...").

---

### PA-001 · Health Condition Assertion

**Severity:** CRITICAL  
**Attribute type:** Physical or mental health, medical condition

**Patterns:**
- "Are you [condition]?", "Do you have [condition]?", "Living with [condition]?"
- "Fellow [condition] sufferers", "My fellow [condition] patients"
- "Managing your [condition]", "dealing with your [condition]", "coping with your [disease]"
- "For people with [condition] in [location]" — when the location makes it feel targeted to a specific person
- Naming conditions: diabetes, cancer, HIV/AIDS, depression, anxiety, bipolar, schizophrenia, epilepsy, Alzheimer's, Parkinson's, MS, ADHD, autism, chronic pain, fibromyalgia, IBS, Crohn's, PCOS, endometriosis, etc.
- Mental health framing: "struggling with mental health", "dealing with your anxiety", "for anxiety sufferers", "depressed people"

**Example violations:**
- "Fellow diabetics in Atlanta — manage your blood sugar with our app" ← CRITICAL
- "Are you struggling with depression? Try our therapy app." ← CRITICAL
- "Living with chronic pain? We can help." ← CRITICAL (direct assumption)

**Allowed alternatives:**
- "Our app supports people managing blood sugar" (no assumption about reader)
- "Mental health support for those who need it" (no specific assumption)

---

### PA-002 · Financial Vulnerability Assertion

**Severity:** CRITICAL  
**Attribute type:** Vulnerable financial status

**Patterns:**
- "Struggling with debt?", "Drowning in debt?", "Can't pay your bills?"
- "Bad credit?", "No credit history?", "Rejected for loans?"
- "Living paycheck to paycheck?", "Can't make ends meet?"
- "Facing bankruptcy?", "Behind on your mortgage?", "Missed a payment?"
- "For people in debt in [location]"

**Example violations:**
- "Struggling with debt in New York? Our service can help." ← CRITICAL
- "Can't pay your bills? We offer short-term loans." ← CRITICAL

**Allowed alternatives:**
- "Debt management solutions for households" (no assumption about reader)
- "Personal loans — apply in minutes" (neutral, no vulnerability implied)

---

### PA-003 · Religious Belief Assertion

**Severity:** CRITICAL  
**Attribute type:** Religion or beliefs

**Patterns:**
- "As a [religion]...", "For [religion] people...", "Fellow [religion] members..."
- "Are you [religion]?", "Muslims in [city]", "Christians in [location]"
- "For believers in [faith]", "Your faith community"
- Naming religions: Muslim, Christian, Jewish, Hindu, Buddhist, Sikh, etc. in a targeting context

**Example violations:**
- "Are you a Muslim living in New York City? Check out our latest deals!" ← CRITICAL (Meta's own example)
- "As a Christian, you know the importance of family values — that's why we created..." ← CRITICAL

**Allowed alternatives:**
- "Products inspired by faith traditions" (informational, not targeting-by-assertion)
- Description of product features without attributing religion to the reader

---

### PA-004 · Age Assertion

**Severity:** CRITICAL  
**Attribute type:** Age

**Patterns:**
- "Are you over [age]?", "For people [age]+", "Fellow seniors..."
- "As someone in their [decade]...", "At [age], you know that..."
- "Baby boomers dealing with [issue]", "Millennials struggling with [problem]"
- Age + health condition compound: "Over 50 and worried about your heart?"

**Example violations:**
- "Over 50 and worried about your health? Try our supplement." ← CRITICAL
- "As a senior, you know that joint pain is a daily struggle." ← CRITICAL

**Allowed alternatives:**
- "Joint support supplements — formulated for long-term wellness"
- Age requirements stated neutrally: "Must be 18+ to apply" (not asserting an attribute about the reader)

---

### PA-005 · Sexual Orientation or Gender Identity Assertion

**Severity:** CRITICAL  
**Attribute type:** Sexual orientation or practices, gender identity

**Patterns:**
- "As a gay man...", "For LGBTQ+ people dealing with...", "Fellow trans women..."
- "Are you gay?", "Are you questioning your gender?"
- "[City]'s gay community — this is for you"
- Sexual orientation + condition compound: "Gay men in [city] dealing with HIV"

**Example violations:**
- "Gay men in San Francisco — get tested with our app." ← CRITICAL
- "For trans women looking for affirming healthcare" ← WARNING (informational but implies identity)

**Allowed alternatives:**
- "Inclusive healthcare for everyone" (not asserting orientation/identity)
- "LGBTQ+-affirming services" in a general informational context (not directly asserting reader is LGBTQ+)

---

### PA-006 · Race or Ethnicity Assertion

**Severity:** CRITICAL  
**Attribute type:** Race, ethnicity, national origin

**Patterns:**
- "As a [race/ethnicity]...", "For [ethnic group] families...", "Fellow [ethnicity] community members..."
- "[Ethnic group] in [city] — this is for you"
- Race + health compound: "Black women at higher risk of [condition]" used to address reader directly

**Example violations:**
- "Latino families in Miami — save on groceries with our app" ← CRITICAL
- "For Black women managing hypertension" ← CRITICAL

**Allowed alternatives:**
- "Grocery savings for families everywhere"
- "Hypertension management tools — for everyone" (can mention health disparities in educational content without addressing reader)

---

### PA-007 · Voting Status or Political Affiliation

**Severity:** CRITICAL  
**Attribute type:** Voting status, political beliefs

**Patterns:**
- "Fellow conservatives...", "As a liberal...", "For Trump/Biden supporters..."
- "Are you registered to vote?", "Voter in [location]..."
- "For people who believe in [political position]..." directly addressing reader

**Allowed alternatives:**
- General civic messaging without assuming the reader's political affiliation or voting status

---

### PA-008 · Criminal Record Assertion

**Severity:** CRITICAL  
**Attribute type:** Criminal record

**Patterns:**
- "Have you been arrested?", "Do you have a criminal record?", "For people with a felony..."
- "Expunge your record", "Post-incarceration services for you"
- "Dealing with a DUI?" — directly addressing reader

**Allowed alternatives:**
- "Legal services for record expungement" (service description, not reader-attribute assertion)

---

### PA-009 · Trade Union Membership

**Severity:** WARNING  
**Attribute type:** Trade union membership

**Patterns:**
- "Fellow union members...", "Are you a union worker?", "For [union name] members..."

---

### PA-010 · Implied Personal Situation (General)

**Severity:** WARNING  
**Attribute type:** General personal situation — intrusive or inappropriate implication

**Patterns:**
- "We know you're going through a tough time"
- "You've been struggling recently"
- "You deserve better than your current situation"
- Any copy that implies knowledge of the reader's specific personal circumstances without basis

**Note:** These are softer violations but still prohibited if they make the user feel surveilled or that the ad knows their personal situation.

---

## Agent 4: Authorization Requirement Check (`auth_required.go`)

**Purpose:** Determine whether the active ad category requires Meta authorization and whether any authorization signals are present in the ad submission. Output: list of authorization requirements, evidence of authorization (or lack thereof), and recommended action.

**Input:** Active categories from Agent 1, ad account metadata (if accessible), copy and landing page.

**Output per required category:** auth_required (boolean), auth_evidence_found (boolean), action_required.

---

### AUTH-01 · Online Gambling Authorization

**Required:** Yes  
**Doc 1 ref:** §4.3  
**Authorization evidence signals in copy:**
- "Licensed by [named regulator]", "Authorized operator", "License number: [X]"
- Landing page includes licensing details

**Action if missing:** Block or warn. Advertiser must complete Meta authorization process (Authorizations and Verifications tab in Meta Business Suite) and submit regulatory license.

---

### AUTH-02 · Dating Services Authorization

**Required:** Yes  
**Doc 1 ref:** §4.4  
**Authorization evidence signals:**
- Meta written permission obtained (account-level flag, if accessible)
- No content-level signal reliably confirms this — flag all dating ads for authorization verification

**Action if missing:** Warn. Advertiser must request permission from Meta via the official form.

---

### AUTH-03 · Drug and Alcohol Addiction Treatment (US-targeted)

**Required:** Yes (US only)  
**Doc 1 ref:** §4.7  
**Authorization evidence signals:**
- Copy references "LegitScript certified", "licensed treatment facility"
- Landing page has LegitScript certification badge

**Context rule:** Only flag for ads targeting US audiences. Check targeting settings.  
**Action if missing:** Warn. Advertiser must obtain LegitScript certification and Meta written permission.

---

### AUTH-04 · Prescription Drugs (US, CA, NZ)

**Required:** Yes (in US, Canada, New Zealand only)  
**Doc 1 ref:** §4.9.2  
**Authorization evidence signals:**
- Copy contains "licensed pharmacy", "valid prescription required", "LegitScript certified"
- Disclaimer present: "consult a licensed health professional"

**Context rule:** Only flag for ads targeting US, Canada, or New Zealand audiences.  
**Action if missing:** Warn. Advertiser must obtain LegitScript certification or complete Meta's pharmaceutical manufacturer review.

---

### AUTH-05 · CBD Products (US)

**Required:** Yes (US only)  
**Doc 1 ref:** §4.9.4  
**Authorization evidence signals:**
- Copy or landing page references LegitScript certification for CBD
- No reliable copy-level signal — flag all CBD ads targeting US audiences

**Action if missing:** Warn. Advertiser must obtain LegitScript certification + Meta written authorization.

---

### AUTH-06 · Cryptocurrency Exchanges and Investment Products

**Required:** Yes  
**Doc 1 ref:** §4.11.1  
**Authorization evidence signals:**
- Ad account has Meta written permission for cryptocurrency (account-level flag, if accessible)
- Copy references regulatory license: "licensed by [named regulator]", "registered with [authority]"

**Action if missing:** Warn. Advertiser must submit regulatory license/registration and obtain written permission via Meta Business Suite.

---

### AUTH-07 · SIEP Authorization and Disclaimer

**Required:** Yes (country-specific)  
**Doc 1 ref:** §7  
**Authorization evidence signals:**
- "Paid for by [name]" disclaimer present and non-empty
- Disclaimer includes required information for elected officials/candidates

**Action if missing:** Block (disclaimer is required). Advertiser must complete Meta's SIEP authorization process and add a "Paid for by" disclaimer.  
**Check:** Verify that the "Paid for by" entity name is non-empty and not a placeholder (e.g., not "Paid for by [advertiser]").

---

### AUTH-08 · Age Targeting Compliance (All Restricted Categories)

**Required:** Yes  
**Doc 1 ref:** §9 (Age Targeting Summary in Doc 1)  
**Check:** For every active restricted category requiring 18+ targeting, verify that minimum age in targeting settings ≥ 18.

| Category | Required Minimum Age |
|---|---|
| Alcohol | 18 (higher in some countries — flag for jurisdictional review) |
| Online gambling | 18 |
| Dating services | 18 |
| Addiction treatment | 18 |
| Weight loss / cosmetic procedures | 18 |
| Adult / reproductive health products | 18 |
| Prescription drugs | 18 |
| OTC drugs | 18 |
| Cannabis / CBD | 18 |
| Financial products (credit, loans, insurance) | 18 |
| Weapons (allowed subset) | 18 |

**Action if below 18:** CRITICAL. Block or require targeting adjustment.

---

## Agent 5: Fix Generation (`fix_generation.go`)

**Purpose:** For each violation flagged by Agents 2, 3, or 4, generate a specific, actionable suggested fix. Fixes should be concrete and directly applicable to the flagged copy. Avoid generic advice.

**Output per violation:** original_text (flagged snippet), suggested_fix (rewritten version), explanation (1 sentence, plain language), fix_type (REWRITE / REMOVE / ADD / TARGETING_CHANGE).

---

### Fix Templates by Violation Type

---

#### FIX-001 · Personal Attribute Assertion → Rewrite to Neutral

**Fix type:** REWRITE  
**Pattern:** Rephrase the copy so it describes the product or service without assuming the reader has the attribute.

| Violation Pattern | Suggested Rewrite |
|---|---|
| "Are you struggling with diabetes?" | "Managing blood sugar levels? Our app can help." |
| "Fellow diabetics in Atlanta" | "Blood sugar management tools — available in Atlanta" |
| "Are you a Muslim living in NYC?" | (No valid rewrite for religion + location targeting — this framing should be removed) |
| "Over 50 and worried about your heart?" | "Heart health support for long-term wellness" |
| "Struggling with debt?" | "Debt management solutions — apply in minutes" |
| "Living with chronic pain?" | "Chronic pain support — designed for daily relief" |
| "Gay men in San Francisco — get tested" | "STI testing — available in San Francisco" |

**Explanation template:** "Removed assumption that the reader has [attribute]. Reframed as a product benefit available to anyone who needs it."

---

#### FIX-002 · Body-Shaming / Negative Self-Perception → Reframe to Positive

**Fix type:** REWRITE

| Violation Pattern | Suggested Rewrite |
|---|---|
| "Hate your belly fat?" | "Ready to feel your best?" |
| "Embarrassed by your weight?" | "Supporting your wellness journey" |
| "Sick of your love handles?" | "Designed for your fitness goals" |
| "Your skin is a problem" | "For clearer, healthier-looking skin" |
| "Disgusted by your appearance?" | "Feel confident in your skin" |

**Explanation template:** "Replaced negative self-perception language with positive/aspirational framing, per Meta's prohibition on body-shaming."

---

#### FIX-003 · Guaranteed Outcome → Hedged Claim

**Fix type:** REWRITE

| Violation Pattern | Suggested Rewrite |
|---|---|
| "Guaranteed to lose 10 kg in 30 days" | "Formulated to support weight management goals — individual results vary" |
| "Guaranteed 20% returns" | "Potential returns of up to 20% — subject to market conditions and risk" |
| "100% cure for [condition]" | "Designed to help manage [symptom] — consult your doctor" |
| "You WILL lose weight" | "Our customers report losing weight — results vary" |

**Explanation template:** "Replaced absolute outcome guarantee with hedged language. Added required disclaimer."

---

#### FIX-004 · Medical Cure/Treatment Claim → Structure-Function Claim

**Fix type:** REWRITE

| Violation Pattern | Suggested Rewrite |
|---|---|
| "Cures diabetes" | "Formulated to support healthy blood sugar levels" |
| "Treats cancer" | (Cannot advertise cancer treatment products without Rx authorization) |
| "Prevents Alzheimer's" | "Supports brain health and cognitive function†" + footnote disclaimer |
| "Clinically proven to eliminate pain" | "Clinically studied ingredient — see research at [link]" |

**Explanation template:** "Replaced disease claim with structure-function language per FDA/FTC guidelines and Meta policy. Added disclaimer."

---

#### FIX-005 · Before/After (Weight Loss) → Remove Comparison

**Fix type:** REWRITE or REMOVE  
**Action:** Remove side-by-side comparison framing. Replace with single-result or product-focused imagery.

**Copy rewrite:**
- "See my 3-month before and after results!" → "See how I use [Product] in my daily routine"
- "Compare my transformation" → "My wellness journey with [Product]"

**Image action:** Replace before/after composition with single image of person using product or showing result only.

---

#### FIX-006 · Missing Age Targeting → Targeting Adjustment

**Fix type:** TARGETING_CHANGE  
**Action:** Set minimum age to 18 in Meta Ads Manager targeting settings for the affected ad set.  
**Explanation:** "[Category] advertising requires minimum age targeting of 18+ per Meta policy."

---

#### FIX-007 · Prohibited Financial Product → Remove or Replace

**Fix type:** REMOVE  
**Action:** This product category (payday loans, CFDs, binary options, ICOs, penny auctions) cannot be advertised on Meta under any circumstances. The ad must be withdrawn or the product/service offering replaced with an allowed financial product.

---

#### FIX-008 · Missing Authorization → Process Guidance

**Fix type:** ADD  
**Action per category:**

| Category | Action |
|---|---|
| Gambling | Complete authorization via Meta Business Suite → Authorizations and Verifications tab. Submit regulatory license. |
| Dating | Fill out Meta's dating permission form. Do not run ads until written permission received. |
| Addiction treatment (US) | Obtain LegitScript certification at legitscript.com, then apply to Meta for permission. |
| Prescription drugs (US/CA/NZ) | Obtain LegitScript certification OR complete Meta pharmaceutical manufacturer review. |
| CBD (US) | Obtain LegitScript CBD certification, then apply to Meta for written authorization. |
| Cryptocurrency | Submit regulatory license/registration via Meta Business Suite → Authorizations and Verifications. |
| SIEP | Complete Meta SIEP authorization. Add "Paid for by [name]" disclaimer to all ads. |

---

#### FIX-009 · Missing SIEP Disclaimer → Add Disclaimer

**Fix type:** ADD  
**Action:** Add "Paid for by [Advertiser Name]" to ad copy. Elected officials/candidates must also include name and title.  
**Explanation:** "All ads about social issues, elections, or politics must include a 'Paid for by' disclaimer per Meta's SIEP policy."

---

#### FIX-010 · Fake Urgency → Substantiate or Remove

**Fix type:** REWRITE  
**Action:** Either remove urgency language if it is not factually accurate, or substantiate it (e.g., link to a real limited-time sale with verifiable end date).

| Violation Pattern | Suggested Rewrite |
|---|---|
| "Only 3 left — buy now!" | "Stock is limited — check availability" OR remove if stock is not actually limited |
| "Offer expires tonight!" | If offer is real: "Sale ends [specific date]". If not real: remove entirely |

---

#### FIX-011 · Tobacco / Nicotine Product → Remove or Reframe as Cessation

**Fix type:** REMOVE or REWRITE  
**Action:** If product is a prohibited tobacco/vaping product: ad cannot run on Meta. Remove.  
**Exception — reframe:** If product is genuinely an FDA/WHO-approved cessation product (nicotine patch, nicotine gum), reframe to emphasize the quitting/cessation purpose clearly. Remove any promotional language for the nicotine delivery experience.

---

## Appendix: Rule ID Quick Reference

| Rule ID | Description | Severity | Agent |
|---|---|---|---|
| V-001 | Personal attribute assertion | CRITICAL | Agent 3 |
| V-002 | Profanity | CRITICAL | Agent 2 |
| V-003 | Vaccine discouragement | CRITICAL | Agent 2 |
| V-004 | Misinformation indicators | WARNING | Agent 2 |
| V-005 | Fake urgency / scarcity | WARNING | Agent 2 |
| V-006 | Guaranteed outcome claims | WARNING/CRITICAL | Agent 2 |
| V-007 | Unsubstantiated celebrity endorsement | WARNING | Agent 2 |
| V-008 | Miracle cure / secret language | WARNING/CRITICAL | Agent 2 |
| V-010 | Weight loss before/after side-by-side | CRITICAL | Agent 2 |
| V-011 | Weight loss body fat close-up | CRITICAL | Agent 2 |
| V-012 | Cosmetic anti-aging before/after | CRITICAL | Agent 2 |
| V-013 | Body-shaming / negative body image | CRITICAL | Agent 2 |
| V-014 | Medical cure claim (non-Rx) | CRITICAL | Agent 2 |
| V-015 | Skin whitening / permanent color change | CRITICAL | Agent 2 |
| V-016 | Adult sexual product / service | CRITICAL | Agent 2 |
| V-020 | Prohibited financial product | CRITICAL | Agent 2 |
| V-021 | Investment via DM (US) | CRITICAL | Agent 2 |
| V-022 | Misleading investment return claims | WARNING | Agent 2 |
| V-023 | Missing financial disclosure | WARNING | Agent 2 |
| V-024 | Deceptive student loan claims | WARNING | Agent 2 |
| V-030 | Prohibited tobacco/nicotine product | CRITICAL | Agent 2 |
| V-031 | Alcohol in prohibited country | CRITICAL | Agent 2 |
| V-032 | Alcohol targeting under 18 | CRITICAL | Agent 2 |
| V-040 | Gambling without authorization | WARNING→CRITICAL | Agent 2 |
| V-041 | Gambling targeting under 18 | CRITICAL | Agent 2 |
| V-042 | Gambling in unsupported market | CRITICAL | Agent 2 |
| V-043 | Deceptive gambling claims | CRITICAL | Agent 2 |
| V-050 | Illicit drug promotion | CRITICAL | Agent 2 |
| V-051 | Drug paraphernalia | CRITICAL | Agent 2 |
| V-052 | CBD without authorization (US) | WARNING | Agent 2 |
| V-053 | CBD health claim | CRITICAL | Agent 2 |
| V-054 | THC / psychoactive cannabis | CRITICAL | Agent 2 |
| V-060 | SIEP without authorization | WARNING | Agent 2 |
| V-061 | False voting information | CRITICAL | Agent 2 |
| V-062 | Deepfake political content | CRITICAL | Agent 2 |
| V-070 | Suicide / self-harm encouragement | CRITICAL | Agent 2 |
| V-071 | Eating disorder content | CRITICAL | Agent 2 |
| V-080 | Sexually suggestive language | CRITICAL | Agent 2 |
| V-081 | Excessive focus on body parts | WARNING | Agent 2 |
| V-090 | Shocking / graphic content | CRITICAL | Agent 2 |
| V-100 | Degrading / shaming content | CRITICAL | Agent 2 |
| PA-001 | Health condition assertion | CRITICAL | Agent 3 |
| PA-002 | Financial vulnerability assertion | CRITICAL | Agent 3 |
| PA-003 | Religious belief assertion | CRITICAL | Agent 3 |
| PA-004 | Age assertion | CRITICAL | Agent 3 |
| PA-005 | Sexual orientation / gender identity assertion | CRITICAL | Agent 3 |
| PA-006 | Race or ethnicity assertion | CRITICAL | Agent 3 |
| PA-007 | Voting status / political affiliation | CRITICAL | Agent 3 |
| PA-008 | Criminal record assertion | CRITICAL | Agent 3 |
| PA-009 | Trade union membership | WARNING | Agent 3 |
| PA-010 | Implied personal situation (general) | WARNING | Agent 3 |
| AUTH-01 | Gambling authorization missing | WARNING→CRITICAL | Agent 4 |
| AUTH-02 | Dating authorization missing | WARNING | Agent 4 |
| AUTH-03 | Addiction treatment authorization missing (US) | WARNING | Agent 4 |
| AUTH-04 | Prescription drug authorization missing | WARNING | Agent 4 |
| AUTH-05 | CBD authorization missing (US) | WARNING | Agent 4 |
| AUTH-06 | Crypto authorization missing | WARNING | Agent 4 |
| AUTH-07 | SIEP disclaimer missing | CRITICAL | Agent 4 |
| AUTH-08 | Age targeting below 18 for restricted category | CRITICAL | Agent 4 |

---

*This document is interpretive. All patterns are inferred from Meta's official policies (Document 1). They are not exhaustive — Meta's automated review system may apply additional undisclosed logic. New violation types and edge cases should be added to this document as they are discovered through testing and review feedback. Re-validate against Document 1 whenever that document is updated.*

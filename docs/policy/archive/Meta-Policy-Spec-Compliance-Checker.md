# Meta Ad Policy Specification — Compliance Checker Knowledge Base

**Version:** 1.0 — Draft  
**Based on:** Meta Advertising Standards (transparency.meta.com), collected April 2026  
**Scope:** Phase 1 — Meta (Facebook + Instagram) only  
**Status:** Working document — verify against source URLs before production deployment

---

## How to Use This Document

This specification is the policy knowledge base for the Ad Compliance Checker. It is not a single prompt — it is structured to map onto the multi-agent evaluation pipeline.

| Pipeline Stage | Section(s) Used |
|---|---|
| Pre-filter (keyword/regex) | Section 4 — Copy patterns under each P_ code |
| Phase 0 — Identity Gate | Section 4: P_IMPERSONATION |
| Phase 1 — Category Classifier | Section 3 — Category Taxonomy |
| Phase 2 — Prohibited Agent | Section 4 — P_ violation codes |
| Phase 2 — Auth-Required Agent | Section 5 — A_ violation codes |
| Phase 2 — Restricted Content Agent | Section 6 — R_ violation codes |
| Phase 2 — Personal Attribute Agent | Section 8 — Standalone (applies to ALL ads) |
| Phase 2 — Quality/Deception Agent | Section 7 — M_ violation codes |
| Phase 3 — Age Targeting Check | Section 9 — Age Targeting Rules (deterministic software) |
| Phase 3 — Scoring Engine | Section 10 — Scoring Specification (deterministic software) |
| Fix Generation Agent | Section 11 — Fix Framework Library |

---

## Section 1: Compliance State Taxonomy

Every category and every violation maps to one of four states.

| State | Code Prefix | Meaning | Default checker action |
|---|---|---|---|
| **PROHIBITED** | `P_` | Absolute ban. No authorization path exists. | `not_compliant`. Score = 0. Show why + policy reference. No fix possible. |
| **AUTHORIZATION REQUIRED** | `A_` | Allowed only with Meta written permission or LegitScript certification. | `not_compliant`. Score = 0. Show authorization path. |
| **RESTRICTED** | `R_` | Allowed, but specific content or targeting rules are violated. | `compliant`. Score reflects severity. Show violation + suggested fix. |
| **MODERATE** | `M_` | Quality, deception, or disclosure issue. Reduces score but does not make ad non-compliant. | `compliant`. Score reflects severity. Show violation + suggested fix. |

**Authorization assumption:** The checker defaults to assuming the advertiser does NOT hold Meta authorization for any authorization-required category. If the user indicates they hold authorization, the checker switches to verifying content rules only — the category flag is downgraded and content-specific R_ and M_ rules apply.

---

## Section 2: Violation Code Naming Convention

| Prefix | Tier | User-facing severity | Action |
|---|---|---|---|
| `P_` | Prohibited | High Risk | not_compliant |
| `A_` | Authorization Required | High Risk (unauthorized) / Medium Risk (authorized but content violation) | not_compliant |
| `R_` | Restricted | Medium Risk | compliant (score reflects severity) |
| `M_` | Moderate | Low Risk | compliant (score reflects severity) |

Codes are internal identifiers — they are never shown directly to end users. The user sees: violation title, Meta policy reference link, explanation, and suggested fix.

---

## Section 3: Category Taxonomy

The Category Classifier assigns exactly one category per ad. The category determines:
1. Which violation agents are spawned in Phase 2
2. The default compliance state before any violation scan
3. The age targeting minimum that is checked in Phase 3

### 3.1 Prohibited Categories
Ads in these categories are blocked at category assignment. No violation scan needed.

| Category | Auto-triggered code | Notes |
|---|---|---|
| Weapons & Explosives — Sale | P_WEAPONS_SALE | Video game weapons are NOT this category |
| Illegal Drugs — Sale or Promotion | P_ILLEGAL_DRUGS | CBD without authorization is A_CBD, not this |
| Counterfeit Goods | P_COUNTERFEIT | Authorized resellers are E-commerce — Luxury |
| Tobacco & Nicotine Products | P_TOBACCO_NICOTINE | WHO/FDA-approved cessation products are RESTRICTED |
| Deepfake & AI Identity Tools | P_DEEPFAKE | Entertainment face filters are Software — Games |
| Covert Surveillance / Stalkerware | P_COVERT_SURVEILLANCE | Transparent parental controls are Software — Utilities |
| Adult Explicit Content / Escort Services | P_ADULT_EXPLICIT | Lingerie/swimwear is E-commerce — Apparel |
| Human Exploitation / Trafficking | P_HUMAN_EXPLOITATION | — |
| Payday Loans / Short-term Loans ≤90 days | P_PAYDAY_LOANS | Long-term loans (>90 days) are Financial — Loans |
| Binary Options / CFD Trading | P_BINARY_OPTIONS_CFD | — |
| ICO (Initial Coin Offerings) | P_ICO | — |

### 3.2 Authorization-Required Categories
Default state: flagged. User must indicate Meta authorization to proceed to content-rule checking.

| Category | Required authorization | Age minimum | Eligible countries |
|---|---|---|---|
| Online Gambling & Betting | Meta written permission | 18+ | Country-dependent |
| Dating Platforms | Meta written permission | 18+ | — |
| Cryptocurrency Exchange / Trading / Lending | Meta written permission + regulatory license | — | — |
| CBD Products | LegitScript certification + Meta written permission | 18+ | US only |
| Prescription Drugs | LegitScript certification + Meta written permission | 18+ | US, Canada, New Zealand only |
| Drug & Alcohol Addiction Treatment | LegitScript certification + Meta written permission | 18+ | US only |

### 3.3 Restricted Categories
Allowed, but content and targeting rules apply. The Restricted Content Agent and Age Targeting Check are spawned.

| Category | Age minimum | Key content constraints | Relevant R_ codes |
|---|---|---|---|
| Health — Weight Loss (Mild) | 18+ | No specific claims, no before/after, no body shaming | R_WEIGHT_LOSS_CLAIM, R_BEFORE_AFTER_WEIGHT, R_BODY_SHAMING |
| Health — Medical Aesthetics & Cosmetic Procedures | 18+ | No before/after for wrinkle treatment; before/after allowed for other procedures without negative self-perception | R_BEFORE_AFTER_COSMETIC_WRINKLE, R_BODY_SHAMING |
| Health — Dietary Supplements | 18+ | No disease treatment/cure claims | R_HEALTH_CLAIM |
| Health — Sexual Wellness | 18+ | Health focus only, no sexual pleasure framing | R_SEXUAL_WELLNESS_FRAMING |
| Health — Reproductive Health | 18+ | Informational or medical service focus | R_HEALTH_CLAIM |
| Health — OTC Drugs | 18+ | Local law compliance | R_HEALTH_CLAIM |
| Financial — Credit Cards | 18+ | Disclosures required, no personal data collection | R_FINANCIAL_DISCLOSURE |
| Financial — Long-term Loans (>90 days) | 18+ | Disclosures required, no personal data collection | R_FINANCIAL_DISCLOSURE |
| Financial — Insurance | 18+ | Country-specific authorization may apply | R_FINANCIAL_DISCLOSURE |
| Financial — Investment Products | 18+ | No DM promotion (US), disclosures required | R_FINANCIAL_DISCLOSURE |
| Financial — Mortgages | 18+ | Country-specific authorization may apply | R_FINANCIAL_DISCLOSURE |
| Content — Alcohol | 18+ | Local law compliance; some countries prohibited entirely | R_ALCOHOL_AGE_TARGETING |
| E-commerce — Apparel (Lingerie / Swimwear) | None | No sexually explicit content or suggestive posing | R_SUGGESTIVE_IMAGERY |
| E-commerce — Dropshipping / Unbranded | None | Heightened scam scrutiny | M_ codes apply |
| E-commerce — Luxury & Branded | None | Domain verification for counterfeit risk | P_COUNTERFEIT check |
| Software — Utilities / VPN | None | No scareware tactics | R_SCAREWARE |
| Health — Cannabis / Hemp (non-CBD) | 18+ | Compliant countries only (US/CA/MX), no health claims, <0.3% THC, no CBD | R_HEALTH_CLAIM |

### 3.4 Allowed Categories
No special restrictions. Quality and deception checks (M_ codes) still apply to all ads.

| Category |
|---|
| E-commerce — Apparel & Fashion |
| E-commerce — Beauty & Cosmetics (non-medical) |
| E-commerce — Home & Lifestyle |
| E-commerce — General Store |
| E-commerce — Travel & Ticketing |
| Software — B2B & Productivity |
| Software — Games & Entertainment (non-gambling) |
| Service — Food & Restaurant |
| Service — Local & Home Services |
| Service — Education & Learning |
| Financial — Brand Ads (no product offer) |
| Financial — Educational Financial Content |
| Health — General Wellness |
| Health — Education & Information (no product sales) |
| Cryptocurrency — Education / News / Events |
| Cryptocurrency — NFT (storage only, no trading) |

---

## Section 4: Prohibited Violations (P_ codes)

**Effect of any P_ violation:** `action: block_ad`, `compliance_score: 0`, `risk_category: High Risk`.  
No fix is possible for P_ violations. Fix framework output: *"This content is absolutely prohibited by Meta's advertising policies. This violation cannot be resolved through edits to the ad."*

---

### P_CHILD_SAFETY — Child Sexual Exploitation and Endangerment

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.1

**Hard criteria:** Any ad content that sexually exploits, endangers, or inappropriately sexualizes minors. Any depiction of minors in romantic or sexual contexts.

**Certainty boundary:**
- Educational content about child safety on a verifiable organization domain → NOT P_CHILD_SAFETY
- Any sexualized imagery or language involving or directed at minors → IS P_CHILD_SAFETY. Immediate block.

---

### P_HATE_SPEECH — Hateful Conduct

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.3

**Hard criteria:** Ad attacks or dehumanizes people based on protected characteristics: race, ethnicity, national origin, disability, religious affiliation, caste, sexual orientation, sex, gender identity, serious disease.

**Copy patterns:** Slurs, dehumanizing comparisons, calls for discrimination or harm against a group, extremist organization names or slogans.

**Visual patterns:** Extremist symbols (in a hate context), propaganda-style imagery targeting a group, imagery glorifying violence against an identifiable group.

**Certainty boundary:**
- Historical/educational content on a verifiable educational domain (e.g., WWII content) → NOT P_HATE_SPEECH
- Political criticism of a policy without targeting a group's identity → NOT P_HATE_SPEECH (may be P_POLITICAL_ADVOCACY)
- Any content that dehumanizes, threatens, or calls for harm based on a protected characteristic → IS P_HATE_SPEECH

---

### P_HUMAN_EXPLOITATION — Human Trafficking and Exploitation

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.16

**Hard criteria:** Ad facilitates, coordinates, or promotes the exploitation of humans, including trafficking, forced labor, or sexual exploitation.

**Certainty boundary:**
- Immigration legal services or refugee support organizations on verifiable domains → NOT P_HUMAN_EXPLOITATION
- Any content offering, soliciting, or advertising human trafficking or forced service → IS P_HUMAN_EXPLOITATION

---

### P_WEAPONS_SALE — Weapons, Ammunition, and Explosives

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.15

**Hard criteria:** Ad promotes sale or use of firearms, ammunition, explosives, weapon modification accessories (including suppressors/silencers), or 3D printing files for weapons.

**Copy patterns:** "Buy guns," "Ammo sale," specific firearm model names with purchase CTA, "suppressor," "3D print lower."

**Visual patterns:** Firearms displayed for sale with pricing or CTA, ammunition product shots, weapon modification parts.

**Certainty boundary:**
- Hunting clothing, boots, or optics (no firearms or ammo) → NOT P_WEAPONS_SALE. Categorize as E-commerce — General Store.
- Video game ad featuring guns as gameplay → NOT P_WEAPONS_SALE. Categorize as Software — Games & Entertainment.
- Airsoft or paintball retailer → NOT P_WEAPONS_SALE unless imagery is indistinguishable from real weapons (flag M_SUGGESTIVE).
- Self-defense class or firearms safety education (no product sale) → NOT P_WEAPONS_SALE.
- A gun store, online firearms retailer, or ammunition seller → IS P_WEAPONS_SALE.

---

### P_ILLEGAL_DRUGS — Illegal Drugs and Paraphernalia

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/drugs-pharmaceuticals/  
**Section:** 4A

**Hard criteria:** Ad promotes sale, purchase, trade, or consumption of illicit, recreational, or unsafe drugs. Promotes sale of drug-related paraphernalia. Promotes merchandise depicting high-risk drugs.

**Copy patterns:** "Buy [substance name]," "Research chemicals," "Legal high," "420 delivery," "psychedelics for sale."

**Visual patterns:** Drug product imagery (pills, powders, marijuana buds) in a sales context, paraphernalia product shots (bongs, rolling papers marketed for drug use).

**Certainty boundary:**
- CBD without authorization → NOT P_ILLEGAL_DRUGS. Categorize as A_CBD.
- Smoking cessation products (patches, gum) → NOT P_ILLEGAL_DRUGS. Categorize as Health — General Wellness.
- Dispensary promoting marijuana sales → IS P_ILLEGAL_DRUGS (regardless of local legal status).
- Standard supplements (caffeine, L-theanine branded as nootropics) → NOT P_ILLEGAL_DRUGS. Categorize as Health — Dietary Supplements.
- Reference to drugs in news, political advocacy, or awareness context (no sale or consumption promotion) → NOT P_ILLEGAL_DRUGS.

---

### P_COUNTERFEIT — Counterfeit Goods and IP Infringement

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.17

**Hard criteria:** Ad promotes replica branded goods. Claims "official store" on a non-brand domain. Uses unauthorized logos or trademarks. Sells items mimicking luxury brand design and identity.

**Copy patterns:** "Replica," "1:1 quality," "AAA grade," "Designer inspired," "Official [brand name] store" on wrong domain, "Factory direct [brand name]."

**Visual patterns:** Brand logos on products where the seller domain doesn't match the brand, luxury goods at suspiciously low prices.

**Certainty boundary:**
- Authorized reseller with proper domain and verifiable branding → NOT P_COUNTERFEIT. Categorize as E-commerce — Luxury & Branded.
- "Inspired by" fragrance with no use of the original brand's name or logo → NOT P_COUNTERFEIT.
- Phone cases or accessories referencing Apple/Samsung to indicate compatibility (e.g., "Case for iPhone 15") → NOT P_COUNTERFEIT. Brand name used for compatibility, not impersonation.
- Band merchandise (t-shirts, posters for Nirvana, Metallica, etc.) → NOT P_COUNTERFEIT unless strong indicators of unlicensed goods (very low price + suspicious domain).
- Any ad using a brand's logo on a non-official domain with sales intent → IS P_COUNTERFEIT.

---

### P_TOBACCO_NICOTINE — Tobacco and Nicotine Products

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.14

**Hard criteria:** Ad promotes sale or use of tobacco, cigarettes, cigars, chewing tobacco, nicotine products, e-cigarettes, vaporizers, or products that simulate smoking.

**Copy patterns:** Brand names of tobacco or vaping products with sales CTA, "Buy cigarettes," "Vape now," "Nicotine pouches."

**Certainty boundary:**
- WHO or US FDA-approved smoking cessation products (prescription or OTC) → NOT P_TOBACCO_NICOTINE. Categorize as Health — OTC Drugs (RESTRICTED) or Health — Prescription Drugs (AUTH_REQUIRED).
- Anti-smoking awareness campaigns on verifiable health organization domains → NOT P_TOBACCO_NICOTINE.
- Any vape, e-cigarette, or nicotine product with purchase CTA → IS P_TOBACCO_NICOTINE.

---

### P_VACCINE_DISCOURAGEMENT — Vaccine Discouragement

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.7

**Hard criteria:** Ad discourages people from getting vaccinated, advocates against vaccines, or promotes vaccine misinformation.

**Copy patterns:** "Vaccines cause [harm]," "Don't vaccinate," "The vaccine agenda," anti-vaccine rhetoric paired with a product or service CTA.

**Certainty boundary:**
- Factual news reporting on vaccine debates on a verifiable news domain → NOT P_VACCINE_DISCOURAGEMENT.
- Health supplement marketed as a "natural alternative to vaccines" → IS P_VACCINE_DISCOURAGEMENT.
- Any ad that frames vaccines as harmful or encourages avoidance → IS P_VACCINE_DISCOURAGEMENT.

---

### P_SUICIDE_SELF_HARM — Suicide, Self-Injury, and Eating Disorders

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.13

**Hard criteria:** Ad encourages, glorifies, or trivializes suicide, self-harm, or eating disorders. Mocks victims or survivors. Includes fictional content (memes, illustrations) that promotes these behaviors.

**Certainty boundary:**
- Mental health awareness campaigns with helpline CTAs on verifiable health organization domains → NOT P_SUICIDE_SELF_HARM.
- Any content that frames self-harm or eating disorder behaviors as desirable, funny, or trivial → IS P_SUICIDE_SELF_HARM.

---

### P_BULLYING_HARASSMENT — Bullying and Harassment

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.11

**Hard criteria:** Ad contains attacks meant to degrade or shame public or private individuals. Heightened protection applies to anyone under 18.

**Certainty boundary:**
- Competitive advertising that criticizes a competitor's product (not the person) → NOT P_BULLYING_HARASSMENT.
- Any ad that targets a specific individual with degrading or shaming content → IS P_BULLYING_HARASSMENT.

---

### P_MISINFORMATION — Debunked Misinformation

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.6

**Hard criteria:** Ad contains claims debunked by third-party fact checkers or that violate Meta's Community Standards on misinformation.

**Note for implementation:** This is partially deterministic — a known misinformation claim database should be checked before LLM evaluation. The LLM should flag strong indicators (miracle cure claims, conspiracy framing, false attribution to scientific bodies).

---

### P_PAYDAY_LOANS — Payday Loans and Short-Term Loans ≤90 Days

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/financial-services/  
**Section:** 5A

**Hard criteria:** Ad promotes payday loans (due on next payday), paycheck advances, bail bonds, or any short-term loan with a repayment period of 90 days or less.

**Copy patterns:** "Get cash until payday," "Instant loan today," "Bad credit OK — same day funds," "Bail bonds."

**Certainty boundary:**
- Long-term personal loans (clearly >90 days) → NOT P_PAYDAY_LOANS. Categorize as Financial — Long-term Loans (RESTRICTED).
- BNPL (Buy Now Pay Later) → IS P_PAYDAY_LOANS (short-term credit facility).
- Budgeting app with no credit/loan offering → NOT P_PAYDAY_LOANS.

---

### P_BINARY_OPTIONS_CFD — Binary Options and CFD Trading

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/financial-services/  
**Section:** 5A

**Hard criteria:** Ad promotes binary options (fixed monetary outcome bets on asset price direction) or Contract for Difference (CFD) trading.

**Copy patterns:** "Binary options," "CFD trading," "Fixed return trades," "Bet on market direction."

---

### P_ICO — Initial Coin Offerings

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/cryptocurrency-products-and-services/  
**Section:** 6A

**Hard criteria:** Ad promotes an ICO — raising funds by selling newly issued cryptocurrency tokens.

**Copy patterns:** "Token sale," "ICO launch," "Buy our token," "Whitelist for our token offering," "Presale."

**Certainty boundary:**
- Cryptocurrency exchange or trading platform → NOT P_ICO. Categorize as A_CRYPTO_EXCHANGE (requires authorization).
- Blockchain news or education → NOT P_ICO. Categorize as Cryptocurrency — Education/News (ALLOWED).

---

### P_DEEPFAKE — Deepfake and AI Identity Tools

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content (Covert Surveillance / Identity Manipulation)

**Hard criteria:** Ad promotes face-swapping apps, voice cloning tools, "undress" or "nudify" AI apps, or any tool designed to create deceptive synthetic media of real people.

**Copy patterns:** "Face swap," "Clone any voice," "Undress AI," "Nudify," "Put your face on any video," "Deepfake generator."

**Certainty boundary:**
- Entertainment face filters (Snapchat-style) for fun/gaming → NOT P_DEEPFAKE. Categorize as Software — Games & Entertainment.
- AI art generation tool (original art, not impersonation) → NOT P_DEEPFAKE.
- Text-to-speech with generic AI voices (not impersonating real people) → NOT P_DEEPFAKE.
- Any app marketed to impersonate a specific real person's face or voice → IS P_DEEPFAKE.

---

### P_COVERT_SURVEILLANCE — Covert Surveillance and Stalkerware

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/

**Hard criteria:** Ad promotes hidden cameras, hidden microphones, stalkerware apps, non-consensual device monitoring, or partner tracking tools.

**Copy patterns:** "Spy on," "Track without knowing," "Catch a cheater," "Invisible tracking," "Monitor their phone secretly."

**Certainty boundary:**
- Home security cameras (Ring, Nest) for property security → NOT P_COVERT_SURVEILLANCE.
- Transparent parental control app where child knows it is installed → NOT P_COVERT_SURVEILLANCE. Categorize as Software — Utilities (RESTRICTED).
- GPS tracker for pets or fleet vehicles → NOT P_COVERT_SURVEILLANCE.
- Any app marketed for covert monitoring of a person without their knowledge → IS P_COVERT_SURVEILLANCE.

---

### P_ADULT_EXPLICIT — Adult Nudity and Explicit Sexual Content

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.9

**Hard criteria:** Visible nudity (nipples, genitalia). Simulated sexual acts. Escort or adult sexual services. Cam site promotion. Content focused on sexual pleasure services.

**Visual patterns:** Explicit poses, visible genitalia, sex act depictions, adult service booking interfaces.

**Certainty boundary:**
- Lingerie brand in catalog-standard poses (standing, walking, neutrally seated) → NOT P_ADULT_EXPLICIT. Categorize as E-commerce — Apparel (RESTRICTED).
- Sexual wellness product showing packaging only (no model, no explicit imagery) → NOT P_ADULT_EXPLICIT. Categorize as Health — Sexual Wellness (RESTRICTED, 18+).
- Reproductive health services (IVF, fertility, contraception) with health framing → NOT P_ADULT_EXPLICIT. Categorize as Health — Reproductive Health (RESTRICTED).
- Any visible nudity, simulated sex act, or escort/cam service → IS P_ADULT_EXPLICIT.

---

### P_IMPERSONATION — Unauthorized Identity and Brand Use

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Fraud and Deceptive Practices — 2.4

**Hard criteria:** Unauthorized use of a public figure's name, likeness, or voice. Fake endorsements. Fake news attributions. Brand logo or name used on a non-official domain with sales intent.

**Detection method:** Phase 0 Identity Gate — CV model detects face or logo, software checks domain against brand whitelist.

**Certainty boundary:**
- Legitimate news organization running its own ad → NOT P_IMPERSONATION (verify domain).
- Brand running its own ad from official domain → NOT P_IMPERSONATION.
- Product review site mentioning brands in editorial context without claiming endorsement → NOT P_IMPERSONATION.
- Celebrity or public figure on a non-official domain → IS P_IMPERSONATION.
- "As seen on [outlet]" where LP domain is not that outlet → IS P_IMPERSONATION.

---

## Section 5: Authorization-Required Violations (A_ codes)

**Default behavior (unauthorized):** `action: block_ad`, `risk_category: High Risk`.  
**If user indicates authorization:** Downgrade to content-rule checking — apply relevant R_ and M_ codes only.

Each A_ violation includes:
- What triggers the flag
- The authorization path (shown to user)
- Content rules that apply once authorized

---

### A_GAMBLING — Online Gambling and Betting

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** 7B

**Triggers (without authorization):** Ad promotes any wagering of monetary value, sports betting, online casino, sweepstakes casino, paid fantasy sports with cash prizes, lotteries.

**Copy patterns:** "Bet now," "Free spins," "Deposit bonus," "Jackpot," "Odds," "Cash out," "Sweepstakes casino."

**Visual patterns:** Slot machine imagery, poker/roulette tables, sports betting slips, casino chips.

**Authorization path:** Obtain Meta written permission via Meta Business Help Center. Must comply with applicable local laws. Apply at: https://www.facebook.com/business/help

**Content rules when authorized:**
- Must target 18+ minimum.
- Must comply with all applicable local laws (some countries prohibit gambling ads entirely).
- Must not use excessive urgency or fake scarcity tactics (apply M_URGENCY check).
- Responsible gambling messaging recommended.

**Certainty boundary:**
- Mobile game showing slot-style mechanics but NO real-money wagering → NOT A_GAMBLING. Categorize as Software — Games & Entertainment.
- Free-to-play game with no cash prizes → NOT A_GAMBLING.
- Fantasy sports with paid entry and cash prizes → IS A_GAMBLING.

---

### A_DATING — Dating Platforms and Services

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** 7C

**Triggers (without authorization):** Ad promotes a dating app, matchmaking platform, relationship site, or hookup service.

**Copy patterns:** "Find your match," "Meet singles near you," "Swipe right," "Dating app," "Find love online."

**Authorization path:** Obtain Meta written permission. Must comply with Meta's dating targeting requirements and dating ad guidelines. Apply at: https://www.facebook.com/business/help

**Content rules when authorized:**
- Must follow Meta's dating targeting requirements.
- No sexually suggestive imagery (P_ADULT_EXPLICIT rules still apply).
- Must not use personal attribute assertions in copy (P_PERSONAL_ATTRIBUTE rules apply).

**Certainty boundary:**
- Professional networking app not primarily for dating → NOT A_DATING.
- Relationship counseling or therapy service → NOT A_DATING. Categorize as Health — General Wellness.
- Matrimonial or culturally specific matchmaking service → IS A_DATING.

---

### A_CRYPTO_EXCHANGE — Cryptocurrency Exchange, Trading, and Lending

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/cryptocurrency-products-and-services/  
**Section:** 6A

**Triggers (without authorization):** Ad promotes a cryptocurrency exchange, trading platform, lending/borrowing platform, cryptocurrency wallet with buy/sell/swap/stake features, or cryptocurrency mining services.

**Copy patterns:** "Buy Bitcoin," "Trade ETH," "Crypto exchange," "Staking rewards," "Crypto lending," "Mine crypto."

**Authorization path:** Submit regulatory license or registration proof + obtain Meta written permission via Meta Business Suite (Authorizations and Verifications tab). Reference: https://transparency.meta.com/policies/ad-standards/restricted-goods-services/cryptocurrency-products-and-services/

**Content rules when authorized:**
- No guaranteed return claims (R_GUARANTEED_OUTCOME applies).
- No ICO promotion (P_ICO still applies even when authorized).
- No referral scheme promotion.
- Must not promise high returns on crypto investment.

**Certainty boundary:**
- Blockchain education, events, or news (no trading/exchange CTA) → NOT A_CRYPTO_EXCHANGE. Categorize as Cryptocurrency — Education/News (ALLOWED).
- Storage-only wallet (no buy/sell/swap/stake) → NOT A_CRYPTO_EXCHANGE. Categorize as Cryptocurrency — NFT/Storage (ALLOWED).
- NFT marketplace (not virtual currency) → NOT A_CRYPTO_EXCHANGE. Categorize as Cryptocurrency — NFT (ALLOWED, but check P_COUNTERFEIT if brand-associated).
- Tax services for crypto companies → NOT A_CRYPTO_EXCHANGE. Categorize as Financial — Educational Content (ALLOWED).

---

### A_CBD — CBD Products

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/drugs-pharmaceuticals/  
**Section:** 4D

**Triggers (without authorization):** Ad promotes CBD products without verified LegitScript certification and Meta written permission.

**Authorization path:** Obtain LegitScript certification (https://legitscript.com) + Meta written permission. Eligible in the United States only.

**Content rules when authorized:**
- Must target 18+ only.
- Must NOT make health or medical claims (R_HEALTH_CLAIM applies).
- Must NOT contain THC or psychoactive component claims.
- Product must comply with applicable US laws.

**Certainty boundary:**
- Hemp products with no CBD and <0.3% THC in US, Canada, or Mexico → NOT A_CBD. Allowed with local law compliance (no authorization needed).
- Any CBD product without verified authorization → IS A_CBD.

---

### A_PRESCRIPTION_DRUGS — Prescription Drug Advertising

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/drugs-pharmaceuticals/  
**Section:** 4B

**Triggers (without authorization):** Ad promotes prescription-only drugs by name or drug class, online pharmacies, or telemedicine focused on obtaining specific prescriptions.

**Authorization path:** Pharmaceutical manufacturers: LegitScript certification OR Meta internal review. Online pharmacies/telehealth: LegitScript active certification + Meta authorization. Eligible in US, Canada, New Zealand only.

**Content rules when authorized:**
- Must target 18+ only.
- Must include disclaimer to consult a licensed health professional and obtain a valid prescription.
- Must not promote specific Rx drugs to general audiences without authorization.

**Certainty boundary:**
- General telehealth service promotion (not tied to specific Rx drug) → NOT A_PRESCRIPTION_DRUGS. Categorize as Health — General Wellness (RESTRICTED).
- OTC medication (Tylenol, Advil) → NOT A_PRESCRIPTION_DRUGS. Categorize as Health — OTC Drugs (RESTRICTED).
- Educational or PSA content about prescription drugs (no purchase/consultation CTA) → NOT A_PRESCRIPTION_DRUGS.
- "Get Ozempic online" or "Viagra without a doctor" → IS A_PRESCRIPTION_DRUGS (and likely also P_ICO if no authorization).

---

### A_ADDICTION_TREATMENT — Drug and Alcohol Addiction Treatment

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** 7D

**Triggers (without authorization):** Ad promotes drug or alcohol addiction treatment services in the US without LegitScript certification and Meta permission.

**Authorization path:** LegitScript certification + Meta written permission. US only.

**Content rules when authorized:**
- Must target 18+ only.
- Must not use fear-based or exploitative messaging targeting people in crisis.

---

## Section 6: Restricted Content Violations (R_ codes)

**Effect:** `action: compliant`, `risk_category: Medium Risk`, score reflects severity.  
A fix is possible. Fix framework is provided per violation.

---

### R_WEIGHT_LOSS_CLAIM — Specific or Guaranteed Weight Loss Claims

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/health-wellness/  
**Section:** 3A

**Hard criteria:** Ad makes a specific numeric weight loss claim ("Lose X kg/lbs in Y days") or uses absolute guarantee language for weight loss outcomes without scientific substantiation.

**Copy patterns:** "Lose 20kg in 30 days," "Drop 3 dress sizes," "Guaranteed weight loss," "Proven to melt fat," "100% natural and effective weight loss."

**Certainty boundary:**
- "Get in shape" or "Reach your wellness goals" → NOT R_WEIGHT_LOSS_CLAIM. Aspirational, non-specific.
- "Lose weight naturally with our supplement" (no specific amount or timeframe) → Borderline. Flag M_SUGGESTIVE if "naturally" implies a medical mechanism.
- "Lose 20 pounds in 30 days" → IS R_WEIGHT_LOSS_CLAIM.

**Fix framework:** Remove the specific numeric claim and timeframe. Replace with aspirational but non-specific outcome language. Examples: "Support your weight management goals," "Designed to complement a healthy lifestyle," "Real results with consistent use." Do not imply guarantee. Do not reference a specific timeframe.

---

### R_BEFORE_AFTER_WEIGHT — Before/After Imagery for Weight Loss

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/health-wellness/  
**Section:** 3A

**Hard criteria:** Ad creative contains side-by-side before/after comparison images for weight loss products or services. Close-up imagery pinching or circling body fat.

**Visual patterns:** Two-panel comparison photos showing body transformation, fat-pinching close-ups, red circles on "problem areas."

**Certainty boundary:**
- Before/after for fitness classes (e.g., Pilates, yoga) → NOT R_BEFORE_AFTER_WEIGHT. Explicitly exempted by Meta policy.
- Before/after for cosmetic procedures other than anti-aging/wrinkle → NOT R_BEFORE_AFTER_WEIGHT (check R_BEFORE_AFTER_COSMETIC_WRINKLE separately).
- Side-by-side weight loss transformation → IS R_BEFORE_AFTER_WEIGHT regardless of whether images are realistic.

**Fix framework:** Remove side-by-side comparison imagery. Replace with single "after" state imagery or product-in-use imagery without body transformation framing. Do not show pinching or circling of body areas.

---

### R_BEFORE_AFTER_COSMETIC_WRINKLE — Before/After for Anti-Aging and Wrinkle Treatment

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/health-wellness/  
**Section:** 3A

**Hard criteria:** Ad creative shows side-by-side before/after imagery specifically for wrinkle treatment, Botox, or anti-aging skin procedures.

**Certainty boundary:**
- Close-up of skin improvement for anti-aging (no side-by-side) → NOT R_BEFORE_AFTER_COSMETIC_WRINKLE. Close-ups are allowed for cosmetic/anti-aging.
- Before/after for non-wrinkle cosmetic procedures (rhinoplasty, hair restoration, etc.) → NOT R_BEFORE_AFTER_COSMETIC_WRINKLE. Those procedures may use before/after without negative self-perception.
- Side-by-side wrinkle treatment comparison → IS R_BEFORE_AFTER_COSMETIC_WRINKLE.

**Fix framework:** Remove the side-by-side comparison. Use a single close-up "after" result image instead. Ensure imagery reflects realistic outcomes over time.

---

### R_BODY_SHAMING — Negative Self-Perception and Body Shaming

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/health-wellness/  
**Section:** 3A; Objectionable Content — 8A

**Hard criteria:** Ad uses language or imagery that implies the viewer's current body, appearance, or condition is shameful, wrong, or needs fixing. Exploits insecurities to conform to beauty standards.

**Copy patterns:** "Tired of your belly fat?", "Embarrassed by your skin?", "Stop being ashamed of your body," "Are you still struggling with your weight?", "Do you hate how you look?", language framing the viewer's current state as a problem.

**Visual patterns:** Imagery implying disgust at a body part, exaggerated "problem area" framing, imagery designed to make the viewer feel inadequate about their appearance.

**Certainty boundary:**
- "Feel your best" or "Love your skin" → NOT R_BODY_SHAMING. Positive aspiration.
- Model in athletic wear with a fit physique → NOT R_BODY_SHAMING. Standard aspirational marketing.
- "Stop being ashamed of your body" or "Tired of hiding your belly?" → IS R_BODY_SHAMING.

**Fix framework:** Reframe from problem-focused to aspiration-focused. Remove any language that implies the viewer's current state is shameful or wrong. Replace with positive outcome framing: "Feel confident," "Support your goals," "Designed for you."

---

### R_HEALTH_CLAIM — Unproven Medical or Health Claims

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Objectionable Content — 8B; Drugs/Pharma — 4A, 4D

**Hard criteria:** Ad for a non-prescription product claims to treat, cure, heal, prevent, or diagnose a specific medical condition or disease. Uses pseudo-medical language implying therapeutic effect without substantiation.

**Copy patterns:** "Cures," "Treats," "Heals," "Prevents [named disease]," "Clinically proven to eliminate," "Doctor recommended" (without verifiable source), "Reverses [condition]," "Eliminates [named medical condition]."

**Certainty boundary:**
- "Supports joint health" or "May help with energy levels" → NOT R_HEALTH_CLAIM. General wellness claims with hedging.
- "Reduces arthritis inflammation" → IS R_HEALTH_CLAIM. Names a specific condition with a treatment claim.
- "Clinically proven to reverse hair loss" → IS R_HEALTH_CLAIM. Specific, unverified medical claim.
- "Cures cancer" or "Prevents Alzheimer's" → IS R_HEALTH_CLAIM at maximum severity. Score at lowest end of R_ range.

**Fix framework:** Remove the disease or condition name from the claim. Replace with general wellness language. Examples: "Supports healthy joints" instead of "Treats arthritis." "Promotes hair vitality" instead of "Reverses hair loss." Add a disclaimer if none is present: "These statements have not been evaluated by the FDA. This product is not intended to diagnose, treat, cure, or prevent any disease."

---

### R_GUARANTEED_OUTCOME — Impossible or Absolute Guarantees

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Objectionable Content — 8B

**Hard criteria:** Ad makes absolute outcome claims using "100%," "guaranteed," "instant," or "always" in contexts where such certainty is factually impossible — particularly for health results, financial returns, or employment outcomes.

**Copy patterns:** "100% guaranteed results," "Guaranteed income of $X/month," "Instantly removes," "Always works," "Guaranteed job placement," "Risk-free returns."

**Certainty boundary:**
- "Satisfaction guaranteed or your money back" → NOT R_GUARANTEED_OUTCOME. This is a refund policy, not a performance claim.
- "100% guaranteed to make you $10K/month" → IS R_GUARANTEED_OUTCOME.
- "Our tool helps improve productivity" → NOT R_GUARANTEED_OUTCOME. Hedged, directional claim.

**Fix framework:** Replace absolute language with directional or aspirational language. "Guaranteed" → "Designed to help," "Backed by our satisfaction guarantee." Remove specific income or outcome figures unless verifiable with supporting data.

---

### R_FINANCIAL_DISCLOSURE — Missing Required Financial Disclosures

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/financial-services/  
**Section:** 5B

**Hard criteria:** Ad for a financial product (credit cards, long-term loans, insurance, investment products) is missing legally required disclosures, APR information, or terms and conditions.

**Copy patterns:** Loan or investment CTAs with no rate or fee disclosure, "0% APR" claims with no terms, investment return claims with no disclaimer.

**Additional rule:** Ad must NOT request personally identifiable financial information (bank account numbers, card numbers, routing numbers) directly in the ad.

**Fix framework:** Add legally required disclosures relevant to the financial product. Include APR ranges, fees, and terms where applicable. Add a disclaimer directing users to full terms on the landing page.

---

### R_AGE_TARGETING_BREACH — Targeting Below Minimum Age for Restricted Category

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Age Targeting Requirements

**Hard criteria:** The user's selected age range includes users under 18 for any category that requires 18+ targeting (see Section 9).

**Note for implementation:** This is a deterministic software check in Phase 3 — compare the user's submitted age_min against the category's minimum age requirement from Section 9. No LLM needed.

**Fix framework:** Adjust the target age range minimum to 18+ for this ad category. This is a targeting setting change, not a creative change.

---

### R_SEXUAL_WELLNESS_FRAMING — Adult Sexual Wellness Framing

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/restricted-goods-services/health-wellness/  
**Section:** 3B

**Hard criteria:** Ad for sexual wellness products focuses on sexual pleasure rather than health benefits. Promotes genital procedures focused on sexual pleasure enhancement.

**Copy patterns:** "Enhanced pleasure," "Better orgasms," "Sexual performance enhancement" (as primary framing), "G-spot augmentation."

**Certainty boundary:**
- Erectile dysfunction product framed as a health condition → NOT R_SEXUAL_WELLNESS_FRAMING if health-focused.
- Contraceptive product (condoms) → NOT R_SEXUAL_WELLNESS_FRAMING.
- Pelvic floor trainer framed as health/recovery → NOT R_SEXUAL_WELLNESS_FRAMING.
- Product explicitly focused on sexual pleasure enhancement → IS R_SEXUAL_WELLNESS_FRAMING.

**Fix framework:** Reframe copy to emphasize health, wellness, or medical benefit rather than sexual pleasure. Consult Meta's health and wellness policy for approved framing.

---

### R_SCAREWARE — Fear-Based Software Tactics

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Fraud and Deceptive Practices — 2.4

**Hard criteria:** Ad mimics system-level alerts or notifications to create fear and drive installs or purchases. Fake virus warnings, fake battery/storage indicators, fake system error states.

**Copy patterns:** "Virus detected," "Your phone is infected," "Memory full," "Clean now," "Your data is at risk."

**Visual patterns:** Fake system notification bars, fake virus scan results, simulated device UI with warning overlays, fake battery/storage indicators.

**Certainty boundary:**
- Legitimate antivirus ad using its own branded UI (not mimicking OS) → NOT R_SCAREWARE.
- Phone cleaner with "Speed up your phone" messaging and no fake system alerts → NOT R_SCAREWARE.
- Ad showing a fake Android/iOS notification bar with "3 viruses detected!" → IS R_SCAREWARE.

**Fix framework:** Remove any fake system UI elements. Replace with branded product imagery and messaging. Do not mimic OS-level notifications.

---

## Section 7: Quality and Deception Violations (M_ codes)

**Effect:** `action: compliant`. Score reflects severity — stacking M_ violations reduces it.  
Apply to all ad categories. These do not make an ad non-compliant but reduce the compliance score.

---

### M_URGENCY — Standard Sales Pressure Tactics

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Fraud and Deceptive Practices — 2.4 (soft end)

**Hard criteria:** Ad uses countdown timers, "limited time offer," "today only," "only X left in stock," "one-time offer," seasonal urgency ("Black Friday ends tonight"), or "X% off — while supplies last."

**Important:** These are standard commercial practices. Flag for awareness but score leniently.

**Certainty boundary:**
- "Summer sale — 30% off this week" → IS M_URGENCY. Standard tactic.
- "Shop the new collection" → NOT M_URGENCY. Standard CTA, no pressure.
- Countdown timer with no verifiable end date → IS M_URGENCY. Score slightly lower due to unverifiability.

**Fix framework:** No fix required if used in isolation. If stacked with other violations, consider removing countdown timers or disclosing the actual sale end date for transparency.

---

### M_NARRATIVE — Sales Narratives and Fabricated Social Proof

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Fraud and Deceptive Practices — 2.4

**Hard criteria:** Ad uses closing/liquidation narratives ("store closing after 34 years"), emotional personal stories as primary sales mechanism, fabricated text message or DM screenshots, fake comment section screenshots, fake celebrity endorsement DMs.

**Certainty boundary:**
- Genuine brand founder telling their real founding story → NOT M_NARRATIVE. Authentic brand storytelling.
- Genuine customer video testimonial → NOT M_NARRATIVE. Flag M_DISCLOSURE if no "results may vary" disclaimer.
- Fabricated iMessage screenshots showing friends recommending the product → IS M_NARRATIVE.
- "Store closing after 34 years" → IS M_NARRATIVE. Standard DTC tactic regardless of whether verifiable.

**Fix framework:** Remove fabricated social proof screenshots. Replace with verifiable customer testimonials with proper attribution. If using a closing narrative, consider adding verifiable context.

---

### M_HYGIENE — Low Production Quality

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/

**Hard criteria:** Ad has excessive typos, poor grammar, low-resolution or blurry images, broken layout, gimmicky symbols (excessive emojis, all-caps blocks, symbol spam), watermarked stock images, or generally unprofessional presentation that signals scam risk.

**Certainty boundary:**
- One minor typo in a professionally designed ad → NOT M_HYGIENE.
- Deliberately "raw" or UGC-style ad that is intentionally lo-fi → NOT M_HYGIENE if content is legitimate and style is intentional.
- 3+ typos, low-res imagery, and all-caps headlines → IS M_HYGIENE.
- Watermarked stock photos → IS M_HYGIENE.

**Fix framework:** Correct typos and grammar. Replace low-resolution images with high-quality assets. Remove watermarked stock photos. Reduce all-caps and symbol spam. Ensure layout is clean and professional.

---

### M_SENSATIONALISM — Clickbait and Engagement Bait

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Objectionable Content

**Hard criteria:** Fake "play" button overlaid on a static image. Excessively sensationalist headlines. Outrage bait. Thumbnail or headline that does not match the content.

**Copy patterns:** "You won't believe," "Shocking," "Doctors hate this," "This one trick," "DOCTORS HATE HIM."

**Visual patterns:** Static image with fake play button triangle overlay, misleading thumbnails.

**Certainty boundary:**
- A video ad with a genuine play button → NOT M_SENSATIONALISM.
- "New study finds surprising benefit of walking" → NOT M_SENSATIONALISM. Mildly clickbaity but within norms.
- "DOCTORS HATE HIM! This one trick will SHOCK you!" → IS M_SENSATIONALISM.

**Fix framework:** Remove fake play button overlays. Replace sensationalist headlines with direct, factual product claims. Ensure the headline accurately represents the ad content.

---

### M_DISCLOSURE — Missing Required Disclaimers and Terms

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/

**Hard criteria:** Ad or landing page is missing required disclaimers for its category: "Results may vary" for health/fitness/beauty, FDA disclaimer for supplements, terms and conditions for financial offers, offer end dates.

**Certainty boundary:**
- Clothing sale with no end date → IS M_DISCLOSURE (minor).
- Supplement ad with no FDA disclaimer → IS M_DISCLOSURE.
- Seasonal or event-tied sale with no explicit end date → IS M_DISCLOSURE but note as minor (event implies deadline).
- Simple product ad with clear pricing and no claims → NOT M_DISCLOSURE.

**Fix framework:** Add required disclaimer for the ad category. For supplements: "These statements have not been evaluated by the FDA. This product is not intended to diagnose, treat, cure, or prevent any disease." For health/fitness results: "Individual results may vary." For financial offers: include APR, terms, and link to full conditions.

---

### M_SUGGESTIVE — Suggestive or Borderline Claims

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/

**Hard criteria:** Ad uses pseudo-clinical charts with unlabeled axes, "studies show" without citation, implied performance claims that stop short of an explicit guarantee, or borderline imagery that doesn't meet R_ or P_ threshold.

**Certainty boundary:**
- "Studies show our product works" with no citation → IS M_SUGGESTIVE.
- "Our customers love us" with star ratings → NOT M_SUGGESTIVE. Standard social proof.
- Lingerie ad with mildly flirty posing beyond catalog-standard → IS M_SUGGESTIVE.

**Fix framework:** Replace "studies show" with a specific citation or remove the claim. Replace unlabeled charts with either properly labeled data or remove. Reframe implied performance claims as directional benefits.

---

## Section 8: Personal Attribute Assertions (P_PERSONAL_ATTRIBUTE)

**This section applies to ALL ads regardless of category.** The Personal Attribute Agent runs on every submission.

**Meta policy reference:** https://transparency.meta.com/policies/ad-standards/  
**Section:** Prohibited Content — 2.12

**Effect:** `action: not_compliant`, `compliance_score: 0`, `risk_category: High Risk`.

---

### What the rule prohibits

Meta prohibits ad copy that asserts or implies personal attributes about the people who will see the ad. The prohibition applies to both direct assertions and indirect implications through rhetorical questions or audience framing.

**Protected attributes covered:**
- Race, ethnicity, national origin
- Religion or beliefs
- Age
- Sexual orientation or sexual practices
- Gender identity
- Disability
- Physical or mental health conditions (including named medical conditions)
- Vulnerable financial status
- Voting status
- Trade union membership
- Criminal record
- Name (directly addressing a specific individual)

---

### Violation patterns — copy examples

**Type 1 — Direct assertion:**
- "As a diabetic, you know how hard it is to manage your sugar." ❌
- "For Christians who want to invest ethically." ❌
- "As a member of the LGBTQ+ community, you deserve better healthcare." ❌

**Type 2 — Rhetorical question implying condition:**
- "Are you struggling with debt?" ❌ (implies financial vulnerability)
- "Are you over 50 and worried about your joints?" ❌ (implies age + health condition)
- "Tired of dealing with anxiety?" ❌ (implies mental health condition)
- "Is your credit score holding you back?" ❌ (implies financial vulnerability)

**Type 3 — Audience address implying attribute:**
- "For people dealing with diabetes..." ❌
- "For those battling depression..." ❌
- "For anyone struggling with addiction..." ❌

---

### What is NOT a violation

**Interest or professional role (not a protected attribute):**
- "For fitness enthusiasts" ✓
- "For small business owners" ✓
- "For people interested in healthy living" ✓
- "Love cooking? This is for you." ✓

**Aspirational or outcome-focused framing:**
- "Feel your best every day." ✓
- "Support your energy levels." ✓
- "Designed for an active lifestyle." ✓

**General health category without implying a condition:**
- "For anyone looking to improve their wellness." ✓
- "Supports healthy joints." ✓ (does not imply the reader has a joint condition)

---

### Certainty boundary

The test is: **Does this copy assume the reader has a specific personal attribute?**

- If yes → IS P_PERSONAL_ATTRIBUTE.
- If the copy describes a product benefit without implying the reader currently has a problem or condition → NOT P_PERSONAL_ATTRIBUTE.

"Supports healthy blood sugar" describes a product function. "For people managing their blood sugar" implies the reader has a blood sugar condition. The second is a violation; the first is not.

---

### Fix framework

Rewrite from "you have this condition" framing to "this product supports this function" framing.

| Violation | Compliant rewrite direction |
|---|---|
| "Are you struggling with debt?" | "Designed to help you reach financial freedom." |
| "For people dealing with diabetes" | "Supports healthy blood sugar levels." |
| "Tired of your anxiety?" | "Designed to support calm and focus." |
| "As a Christian investor..." | "Invest according to your values." |
| "Over 50 and worried about joints?" | "Advanced joint support formula." |

---

## Section 9: Age Targeting Compliance Rules

This check is deterministic software — no LLM needed. Compare `age_min` from the submitted check against the minimum required for the assigned category.

| Category | Minimum age_min required |
|---|---|
| Weight loss products and services | 18 |
| Cosmetic products, procedures, or surgeries | 18 |
| Financial products (credit cards, long-term loans, insurance, investment) | 18 |
| Prescription drugs (if authorized) | 18 |
| OTC drugs | 18 |
| Cannabis / CBD products | 18 |
| Alcohol | 18 |
| Online gambling (if authorized) | 18 |
| Reproductive health / sexual wellness | 18 |
| Drug and alcohol addiction treatment (if authorized) | 18 |
| Dating platforms (if authorized) | 18 |
| Dietary supplements | 18 |

**Trigger:** If `age_min` < required minimum → flag `R_AGE_TARGETING_BREACH`.  
**Fix:** Adjust age targeting minimum to 18+. This is a targeting change, not a creative change.

**Exempt from age targeting requirements (no minimum enforced by this checker):**
- Women's hygiene products
- Lingerie and swimwear (must still comply with content rules)
- Sex education with informational focus and no suggestive content
- General food products
- Cosmetic products (creams, makeup, hair products)

---

## Section 10: Scoring Engine Specification

The scoring engine is deterministic software. It receives a list of violation codes from all agents and computes the final output. No LLM involved.

### Input

```
{
  "violations": ["R_WEIGHT_LOSS_CLAIM", "M_DISCLOSURE"],
  "category_state": "RESTRICTED",
  "authorized": false
}
```

### Scoring Waterfall (execute in order — first match wins)

```python
def score(violations, category_state, authorized):

    p_codes = [v for v in violations if v.startswith("P_")]
    a_codes = [v for v in violations if v.startswith("A_")]
    r_codes = [v for v in violations if v.startswith("R_")]
    m_codes = [v for v in violations if v.startswith("M_")]

    # Rule 1: Any P_ violation
    if p_codes:
        return {
            "compliance_score": 0,
            "risk_category": "High Risk",
            "action": "not_compliant"
        }

    # Rule 2: Any A_ violation (unauthorized)
    if a_codes and not authorized:
        return {
            "compliance_score": 0,
            "risk_category": "High Risk",
            "action": "not_compliant"
        }

    # Rule 3: Any R_ violation (no P_ or unauthorized A_)
    if r_codes:
        base = 50  # starting anchor for single R_ violation
        score = base
        score -= 10 * (len(r_codes) - 1)   # each additional R_ deducts 10
        score -= 5 * len(m_codes)            # each M_ stacking deducts 5
        score = max(score, 1)                # floor at 1
        risk = "High Risk" if score < 30 else "Medium Risk"
        return {
            "compliance_score": score,
            "risk_category": risk,
            "action": "compliant"   # not a hard block — user sees violations + fixes and decides
        }

    # Rule 4: Any M_ violation only (no P_, A_, or R_)
    if m_codes:
        base = 75  # starting anchor for single M_ violation
        score = base
        score -= 5 * (len(m_codes) - 1)     # each additional M_ deducts 5
        score = max(score, 45)               # floor at 45
        risk = "Low Risk" if score >= 65 else "Medium Risk"
        return {
            "compliance_score": score,
            "risk_category": risk,
            "action": "compliant"
        }

    # Rule 5: No violations — ALLOWED category
    if category_state == "ALLOWED":
        return {
            "compliance_score": 95,  # refine with quality signals: 85-100
            "risk_category": "Low Risk",
            "action": "compliant"
        }

    # Rule 6: No violations — RESTRICTED or AUTH_REQUIRED (authorized) category
    return {
        "compliance_score": 78,  # refine with quality signals: 71-84
        "risk_category": "Low Risk",
        "action": "compliant"
    }
```

### Score Anchors for Clean Ads (Rules 5 and 6)

| Score | Condition |
|---|---|
| 95–100 | No violations, ALLOWED category, professional creative, clear LP with visible T&Cs |
| 85–94 | No violations, ALLOWED category, minor imperfections (generic creative, LP could be more transparent) |
| 80–84 | No violations, RESTRICTED/AUTHORIZED category, well-produced, appropriate disclaimers present |
| 71–79 | No violations, RESTRICTED/AUTHORIZED category, compliant but minor ambiguities |

---

## Section 11: Fix Framework Library

The Fix Generation Agent receives: the violation code, the original ad copy, and the fix framework below. It generates a specific rewrite based on the framework constraints.

| Violation Code | Fix direction | Cannot include |
|---|---|---|
| P_ (any) | No fix possible. Explain the prohibition and cite the policy. | — |
| A_ (unauthorized) | Explain authorization requirement. Link to Meta authorization process. | Do not suggest creative edits as a workaround. |
| R_WEIGHT_LOSS_CLAIM | Replace specific numeric claims with aspirational non-specific outcome language. | Specific amounts, timeframes, or guarantees. |
| R_BEFORE_AFTER_WEIGHT | Remove side-by-side comparison. Use single-state or product-in-use imagery. | Side-by-side transformation, fat-pinching close-ups. |
| R_BEFORE_AFTER_COSMETIC_WRINKLE | Use single "after" close-up instead of side-by-side. | Side-by-side for wrinkle/anti-aging. |
| R_BODY_SHAMING | Reframe to positive aspiration. Remove problem-focused language. | Any language implying the viewer's current state is shameful or wrong. |
| R_HEALTH_CLAIM | Remove named condition. Replace with general wellness framing. Add FDA disclaimer. | Disease names, "treats," "cures," "prevents," "clinically proven" (without citation). |
| R_GUARANTEED_OUTCOME | Replace absolute language with directional or conditional language. | "Guaranteed," "100%," specific income/outcome figures without substantiation. |
| R_FINANCIAL_DISCLOSURE | Add required disclosures for the financial product type. | Cannot remove the financial offer — only add disclosures. |
| R_AGE_TARGETING_BREACH | Adjust targeting minimum to 18+. This is a targeting setting change only. | No creative change needed. |
| R_SEXUAL_WELLNESS_FRAMING | Reframe to health benefit, not sexual pleasure. | "Enhanced pleasure," "better orgasms," explicit performance framing. |
| R_SCAREWARE | Remove fake system UI elements. Use branded product imagery. | OS-mimicking notification bars, fake virus/battery alerts. |
| P_PERSONAL_ATTRIBUTE | Rewrite from "you have this condition" to "this product supports this function." | Direct or implied assertions about the viewer's protected attributes. |
| M_URGENCY | No creative fix required if isolated. Disclose actual end date for transparency if countdown is present. | — |
| M_NARRATIVE | Replace fabricated screenshots with verifiable testimonials. | Fake text messages, fake DM screenshots, fake comment threads. |
| M_HYGIENE | Correct typos, replace low-res/watermarked images, reduce symbol spam. | — |
| M_SENSATIONALISM | Replace sensationalist headlines with direct product claims. Remove fake play buttons. | "You won't believe," "doctors hate," "one trick," fake play button overlays. |
| M_DISCLOSURE | Add required category-specific disclaimer. | Cannot remove the claim — only add the disclaimer. |
| M_SUGGESTIVE | Add citation for statistics or remove. Replace implied claims with hedged language. | Uncited "studies show," unlabeled performance charts. |

---

*This document is the authoritative policy knowledge base for the Ad Compliance Checker (Phase 1 — Meta). It must be re-verified against Meta's official policy pages before any production deployment. Meta policies change without notice. A policy review process should be scheduled at minimum quarterly.*

*Source URLs: https://transparency.meta.com/policies/ad-standards/ and sub-pages listed per section.*

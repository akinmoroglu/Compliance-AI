package models

type ViolationMeta struct {
	Title    string
	Severity string
}

var ViolationMetadata = map[string]ViolationMeta{
	"P_HATE_SPEECH":         {"Hate Speech", "PROHIBITED"},
	"P_ILLEGAL_DRUGS":       {"Illegal Drug Promotion", "PROHIBITED"},
	"P_IMPERSONATION":       {"Brand or Identity Impersonation", "PROHIBITED"},
	"P_PERSONAL_ATTRIBUTE":  {"Personal Attribute Assertion", "PROHIBITED"},
	"P_BULLYING_HARASSMENT": {"Bullying or Harassment", "PROHIBITED"},
	"A_GAMBLING":            {"Unauthorized Gambling Promotion", "AUTHORIZATION_REQUIRED"},
	"A_CRYPTO_EXCHANGE":     {"Unauthorized Crypto Exchange Promotion", "AUTHORIZATION_REQUIRED"},
	"R_WEIGHT_LOSS_CLAIM":   {"Specific Weight Loss Claim", "RESTRICTED"},
	"R_BODY_SHAMING":        {"Negative Body Image or Body Shaming", "RESTRICTED"},
	"R_HEALTH_CLAIM":        {"Unproven Health or Medical Claim", "RESTRICTED"},
	"R_GUARANTEED_OUTCOME":  {"Guaranteed Financial Outcome", "RESTRICTED"},
	"R_SCAREWARE":           {"Scareware or Fake System Alert", "RESTRICTED"},
	"M_URGENCY":             {"Urgency Tactic", "MODERATE"},
	"M_NARRATIVE":           {"Fabricated Social Proof", "MODERATE"},
	"M_HYGIENE":             {"Low Production Quality", "MODERATE"},
	"M_SENSATIONALISM":      {"Clickbait or Sensationalism", "MODERATE"},
	"M_DISCLOSURE":          {"Missing Required Disclaimer", "MODERATE"},
	
	// Vision Agent codes
	"P_ADULT_EXPLICIT":               {"Explicit Adult Content", "PROHIBITED"},
	"P_VIOLENCE_GRAPHIC":             {"Graphic Violence", "PROHIBITED"},
	"P_WEAPONS_VISUAL":               {"Weapons Display", "PROHIBITED"},
	"P_PROFANITY_VISUAL":             {"Profane Visual Content", "PROHIBITED"},
	"P_CHILD_SAFETY_VISUAL":          {"Child Safety Violation", "PROHIBITED"},
	"R_BEFORE_AFTER_WEIGHT_VISUAL":   {"Before/After Weight Loss Image", "RESTRICTED"},
	"R_BEFORE_AFTER_COSMETIC_VISUAL": {"Before/After Cosmetic Image", "RESTRICTED"},
	"R_BODY_SHAMING_VISUAL":          {"Body Shaming Visual", "RESTRICTED"},
	"R_SKIN_WHITENING_VISUAL":        {"Skin Whitening Visual", "RESTRICTED"},
	"R_SEXUAL_WELLNESS_VISUAL":       {"Sexual Wellness Visual", "RESTRICTED"},
	"M_SUGGESTIVE_VISUAL":            {"Suggestive Imagery", "MODERATE"},
	"M_SHOCK_IMAGERY":                {"Shock Imagery", "MODERATE"},
	"M_SCAREWARE_VISUAL":             {"Scareware Visual", "MODERATE"},

	// LP Agent codes
	"P_COUNTERFEIT":          {"Counterfeit Goods", "PROHIBITED"},
	"P_MISINFORMATION":       {"Dangerous Misinformation", "PROHIBITED"},
	"P_HUMAN_EXPLOITATION":   {"Human Exploitation Content", "PROHIBITED"},
	"R_BEFORE_AFTER_WEIGHT":  {"Before/After Weight Loss Claim", "RESTRICTED"},
	"R_FINANCIAL_DISCLOSURE": {"Missing Financial Disclosure", "RESTRICTED"},

	// Alignment Agent codes
	"ALIGN_PRODUCT_MISMATCH":     {"Ad/LP Product Mismatch", "MODERATE"},
	"ALIGN_CLAIM_MISMATCH":       {"Ad/LP Claim Contradiction", "MODERATE"},
	"ALIGN_PRICE_MISMATCH":       {"Ad/LP Price Contradiction", "MODERATE"},
	"ALIGN_VISUAL_COPY_MISMATCH": {"Visual/Copy Contradiction", "MODERATE"},
}

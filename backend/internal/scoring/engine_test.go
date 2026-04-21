package scoring

import "testing"

func TestScore(t *testing.T) {
	tests := []struct {
		name          string
		violations    []string
		categoryState string
		authorized    bool
		expectedScore int
		expectedRisk  string
		expectedAct   string
	}{
		// Existing 6 test cases updated to v2 action strings
		{
			name:          "P violation - High Risk",
			violations:    []string{"P_HATE_SPEECH"},
			categoryState: "ALLOWED",
			expectedScore: 0,
			expectedRisk:  "High Risk",
			expectedAct:   "PROHIBITED",
		},
		{
			name:          "Unauthorized Category - High Risk",
			violations:    []string{"A_GAMBLING"},
			categoryState: "AUTH_REQUIRED",
			authorized:    false,
			expectedScore: 0,
			expectedRisk:  "High Risk",
			expectedAct:   "AUTHORIZATION_REQUIRED",
		},
		{
			name:          "Restricted violation - Medium Risk",
			violations:    []string{"R_WEIGHT_LOSS_CLAIM"},
			categoryState: "RESTRICTED",
			expectedScore: 50,
			expectedRisk:  "Medium Risk",
			expectedAct:   "RESTRICTED",
		},
		{
			name:          "Moderate violation - Low Risk",
			violations:    []string{"M_URGENCY"},
			categoryState: "ALLOWED",
			expectedScore: 75,
			expectedRisk:  "Low Risk",
			expectedAct:   "ALLOWED_WITH_RESTRICTIONS",
		},
		{
			name:          "Clean Allowed",
			violations:    []string{},
			categoryState: "ALLOWED",
			expectedScore: 95,
			expectedRisk:  "Low Risk",
			expectedAct:   "ALLOWED",
		},
		{
			name:          "Clean Sensitive",
			violations:    []string{},
			categoryState: "RESTRICTED",
			expectedScore: 78,
			expectedRisk:  "Low Risk",
			expectedAct:   "ALLOWED",
		},
		// New v2 test cases
		{
			name:          "R_ + ALIGN_ combination",
			violations:    []string{"R_WEIGHT_LOSS_CLAIM", "ALIGN_CLAIM_MISMATCH"},
			categoryState: "RESTRICTED",
			expectedScore: 45, // 50 - (10 * 0) - (5 * 1)
			expectedRisk:  "Medium Risk",
			expectedAct:   "RESTRICTED",
		},
		{
			name:          "ALIGN_ alone",
			violations:    []string{"ALIGN_CLAIM_MISMATCH"},
			categoryState: "ALLOWED",
			expectedScore: 75, // Treated as M_: 75 - (5 * 0)
			expectedRisk:  "Low Risk",
			expectedAct:   "ALLOWED_WITH_RESTRICTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Score(tt.violations, tt.categoryState, tt.authorized)
			if got.ComplianceScore != tt.expectedScore || got.RiskCategory != tt.expectedRisk || got.Action != tt.expectedAct {
				t.Errorf("Score() = %v, want { %d, %s, %s }", got, tt.expectedScore, tt.expectedRisk, tt.expectedAct)
			}
		})
	}
}

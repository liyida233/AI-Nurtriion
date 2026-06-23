package goal

import "testing"

func TestIsValidationError(t *testing.T) {
	cases := []struct {
		message string
		want    bool
	}{
		{"goal deadline should allow at least one week for meaningful progress tracking", true},
		{"weight loss target appears too aggressive; use a safer weekly target", true},
		{"milestone dueDate must be YYYY-MM-DD", true},
		{"status must be active, completed, paused, or cancelled", true},
		{"database unavailable", false},
	}

	for _, test := range cases {
		if got := isValidationError(test.message); got != test.want {
			t.Fatalf("isValidationError(%q) = %v, want %v", test.message, got, test.want)
		}
	}
}

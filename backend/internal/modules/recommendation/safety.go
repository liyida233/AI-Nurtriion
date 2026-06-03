package recommendation

import "strings"

func IsSafeGeneralWellness(content string) bool {
	blocked := []string{"diagnose", "cure", "starve", "extreme fasting"}
	lowered := strings.ToLower(content)
	for _, word := range blocked {
		if strings.Contains(lowered, word) {
			return false
		}
	}
	return true
}

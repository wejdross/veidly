package helpers

func CutString(s string, sz int) string {
	if len(s) > sz {
		return s[:sz] + "..."
	}
	return s
}

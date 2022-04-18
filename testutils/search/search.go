package search

// searches if needle is present in haystack
func IsPresentIn(needle string, haystack []string) bool {
	for _, thing := range haystack {
		if thing == needle {
			return true
		}
	}

	return false
}

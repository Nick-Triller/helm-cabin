package server

// contains checks if a string slice contains a string value in O(n) time
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

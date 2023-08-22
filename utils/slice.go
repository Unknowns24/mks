package utils

func SliceContainsElement(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

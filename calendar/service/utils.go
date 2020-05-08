package service

func indexInStrings(l []string, s string) int {
	for i, v := range l {
		if s == v {
			return i
		}
	}
	return -1
}

func findInStrings(l []string, s string) bool {
	return indexInStrings(l, s) != -1
}

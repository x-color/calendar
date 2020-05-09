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

func removeAtInStrings(l []string, i int) []string {
	if i == len(l)-1 {
		l = l[:i]
	} else {
		l = append(l[:i], l[i+1:]...)
	}
	return l
}

func removeInStrings(l []string, s string) ([]string, bool) {
	i := indexInStrings(l, s)
	if i == -1 {
		return l, false
	}

	return removeAtInStrings(l, i), true
}

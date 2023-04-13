package search

func addToMap(target, src SingleTableIDmap) {
	for x := range src {
		target[x] = src[x]
	}
}

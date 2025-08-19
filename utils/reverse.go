package utils

func Reverse[T any](arr []T) []T {
	narr := make([]T, len(arr))
	copy(narr, arr)
	for i, j := 0, len(narr)-1; i < j; i, j = i+1, j-1 {
		narr[i], narr[j] = narr[j], narr[i]
	}
	return narr
}

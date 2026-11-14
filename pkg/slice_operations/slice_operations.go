package sliceoperations

import "unsafe"

func SliceToString[T any](s []T) string {
	if len(s) == 0 {
		return ""
	}

	size := len(s) * int(unsafe.Sizeof(s[0]))

	return unsafe.String((*byte)(unsafe.Pointer(&s[0])), size)
}

func CompareSlices[T any](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	if SliceToString(a) != SliceToString(b) {
		return false
	}

	return true
}

func CompareSlicesOfPointers[T any](a []*T, b []*T) bool {
	if len(a) != len(b) {
		return false
	}

	aVals := make([]T, len(a))
	bVals := make([]T, len(b))

	for i := 0; i < len(a); i++ {
		aVals[i] = *a[i]
		bVals[i] = *b[i]
	}

	if SliceToString(aVals) != SliceToString(bVals) {
		return false
	}

	return true
}

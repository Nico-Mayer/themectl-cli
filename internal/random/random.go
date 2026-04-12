package random

import "math/rand/v2"

func Element[T any](slice []T) T {
	return slice[rand.IntN(len(slice))]
}

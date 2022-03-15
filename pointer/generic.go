package pointer

func Pointer[T any](in T) *T {
	return &in
}

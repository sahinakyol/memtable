package main

func Fnv1aHash(s string) uint64 {
	const (
		offset64 uint64 = 14695981039346656037
		prime64  uint64 = 1099511628211
	)

	var hash = offset64

	for i := 0; i < len(s); i++ {
		hash ^= uint64(s[i])
		hash *= prime64
	}
	return hash
}

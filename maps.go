package sequence

// FromMap creates a sequence of pairs, where the first item in the pair is a
// key from the map, and the second item is the value for that key.
// The sequence is Volatile, since multiple iterations could produce different
// results.
func FromMap[K comparable, V any](m map[K]V) Sequence[Pair[K, V]] {
	return GenerateVolatile(func(f func(Pair[K, V]) error) error {
		for k, v := range m {
			if err := f(MakePair(k, v)); err != nil {
				return err
			}
		}
		return nil
	})
}

// IntoMap stores the pairs from the given sequence into the provided map.
// Multiple items with the same key from the sequence will cause the later ones
// to overwrite earlier ones. If an error occurs generating items from the
// sequence the map will be partially updated with every item upto the error.
func IntoMap[M ~map[K]V, K comparable, V any](dst M, s Sequence[Pair[K, V]]) error {
	return EachSimple(s.Sync())(func(p Pair[K, V]) bool {
		k, v := p.AB()
		dst[k] = v
		return true
	})
}

// ToMap creates a map from the given sequence of Pairs, and returning it in a
// Result. The first element in each pair must be comparable, since it must
// become a map key. Earlier items with the same key will be replaced by later
// items in the sequence. If an error occurs, the value of the Result will
// contain all the items processed so far.
func ToMap[K comparable, V any](s Sequence[Pair[K, V]]) Result[map[K]V] {
	m := map[K]V{}
	err := IntoMap(m, s)
	return MakeResult(m, err)
}

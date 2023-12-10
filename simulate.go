package sequence

// Simulate creates a sequence where an initial state is updated by a function
// that produces a new state, and indicates if there is more work to be done.
// The sequence of states is expected to be deterministic given the initial
// state and the step function, if not, the sequence should be marked volatile.
// The output sequence returns the series of states created from the initial
// state to the final state.
func Simulate[T any](initial T, step func(T) (T, bool)) Sequence[T] {
	return Generate(func(f func(T) error) error {
		state := initial
		for more := true; more; state, more = step(state) {
			if err := f(state); err != nil {
				return nil
			}
		}
		return nil
	})
}

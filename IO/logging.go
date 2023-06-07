package IO

import "log"

// Log prints a value using the given format string when the `IO` side effect gets executed
// not that this is a side effect in itself, because the fact of writing to the log
// stream changes the environment
func Log[A any](format string) func(a A) IO[any] {
	return func(a A) IO[any] {
		return func() any {
			log.Printf(format, a)
			return nil
		}
	}
}

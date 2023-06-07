package IO

import "sync"

// IO defines an operation with a side effect. This is signaled by the fact
// that the function does not accept an input value but it returns an output
// this can only mean that 'A' is the result of a side effect because it does not depend on an input
type IO[A any] func() A

// Of returns an `IO` that will always return the same, constant value 'a'
func Of[A any](a A) IO[A] {
	return func() A {
		return a
	}
}

// Ap combines two `IO` side effects (asyncronously). The implementation will spawn one additional goroutine
// to compute one of the side effects. The other one is computed on the current goroutine
func Ap[A, B any](ma IO[A]) func(IO[func(A) B]) IO[B] {
	return func(fab IO[func(A) B]) IO[B] {
		return func() B {
			// use a wait group to synchronize with the go routine that computes one of the values
			// this is a bit mor efficient than using a channel
			var wg sync.WaitGroup
			wg.Add(1)

			// this captures the result of the computation of one side effect
			var a A

			// execute one of the side effects in its own goroutine
			go func() {
				defer wg.Done()
				a = ma()
			}()

			// eagerly compute the other side effect, then wait
			ab := fab()
			wg.Wait()

			// composes the result of both side effects
			return ab(a)
		}
	}
}

// Chain computes two side effects in sequence, both side effects are executed on the current goroutine
// this function is sometimes referred to as `FlatMap`
func Chain[A, B any](f func(A) IO[B]) func(ma IO[A]) IO[B] {
	return func(ma IO[A]) IO[B] {
		return func() B {
			return f(ma())()
		}
	}
}

// Map transforms the result of a side effect via a pure function
func Map[A, B any](f func(A) B) func(ma IO[A]) IO[B] {
	return func(ma IO[A]) IO[B] {
		return func() B {
			return f(ma())
		}
	}
}

// Chain computes two side effects in sequence, but it returns the value of the first
func ChainFirst[A, B any](f func(A) IO[B]) func(ma IO[A]) IO[A] {
	return Chain(func(a A) IO[A] {
		return Map(func(b B) A {
			return a
		})(f(a))
	})
}

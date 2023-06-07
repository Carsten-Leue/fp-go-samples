package IO

// appendToArray is a curried version of `append`
func appendToArray[A any](as []A) func(a A) []A {
	return func(a A) []A {
		return append(as, a)
	}
}

// SequenceArray transforms a sequence of `IO` side effects into an `IO` side effect of a sequence
// Note that the implementation does not make any assumption about the nature of `IO` it delegates
// to the function `Of`, `Map` and `Ap`
// the use of `Ap` will execute the operations in parallel
func SequenceArray[A any](as []IO[A]) IO[[]A] {
	// it's important that the static array has a capacity of `0`, so the first
	// attempt to add a value will make a copy
	current := Of(make([]A, 0)) // IO[[]A]

	// input to the function is an `IO` operation that produces the array to append to
	// output is a function that adds the next value
	add := Map(appendToArray[A]) // func(ma IO[[]A]) IO[func(a A) []A]

	// `fa` and `current` are both `IO` side effects, so this loop
	// will not compute any values. Instead it eagerly constructs an new `IO`
	// side effect that will compute the final array lazily
	for _, fa := range as {
		current = Ap[A, []A](fa)(add(current))
	}

	return current
}

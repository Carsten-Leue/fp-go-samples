package IO

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomInt() int {
	return rand.Intn(100)
}

func TestArray(t *testing.T) {

	// construct an array of `IO` side effects, in this case random numbers
	randomNumbers := make([]IO[int], 10)
	for i := len(randomNumbers) - 1; i >= 0; i-- {
		// note that the numbers are not computed here, `randomInt` is an `IO` side effect
		randomNumbers[i] = randomInt
	}

	// here we lift the array of side effects into a side effect of an array
	// no computation has been done, yet
	numIO := SequenceArray(randomNumbers)

	// compute the numbers here
	num1 := numIO() // []int
	num2 := numIO() // []int

	assert.NotEqual(t, num1, num2)
	assert.Equal(t, len(randomNumbers), len(num1))
	assert.Equal(t, len(randomNumbers), len(num2))

}

func intToString(value int) string {
	return fmt.Sprintf("value: %d", value)
}

func TestMap(t *testing.T) {

	rnd := randomInt // IO[int]

	toString := Map(intToString) // func(ma IO[int]) IO[string]

	// no conversion happened, so far, just setup of the side effect
	rndStrg := toString(rnd)

	// compute the strings here
	s1 := rndStrg() // string
	s2 := rndStrg() // string

	assert.NotEqual(t, s1, s2)
}

func TestLogging(t *testing.T) {

	rnd := randomInt // IO[int]

	rndWithLogging := ChainFirst(Log[int]("randomValue: %d"))(rnd) // IO[int]

	// compute the numbers here
	num1 := rndWithLogging() // int
	num2 := rndWithLogging() // int

	assert.NotEqual(t, num1, num2)
}

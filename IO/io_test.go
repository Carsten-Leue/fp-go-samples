package IO

import (
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
	num1 := numIO()
	num2 := numIO()

	assert.NotEqual(t, num1, num2)
	assert.Equal(t, len(randomNumbers), len(num1))
	assert.Equal(t, len(randomNumbers), len(num2))

}

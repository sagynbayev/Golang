package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (err ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %f", err)
}

func Sqrt(x float64) (float64, error) {

	if x < 0 {
		return x, ErrNegativeSqrt(x)
	}

	z := 1.0
	var a float64
	for math.Abs(a-z) > 1e-9 {
		a = z
		z -= (z*z - x) / (2 * z)
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}

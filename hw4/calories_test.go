package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalories(t *testing.T) {
	weight, height, age, activity, gender := 90.0, 180.0, 19.0, "a2", "male"
	// cCaloryalc := BurnCalories(weight, height, age, activity, gender)
	caloryCalc2 := BurnCalories(weight, height, age, activity, gender)
	assert.Equal(t, 2660.625, caloryCalc2)
}

package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDotProduct(t *testing.T) {
	type testCase struct{
		a Vector3
		b Vector3
		expected float64
	}
	t.Run("dotProduct func", func(t *testing.T) {
		tests := []testCase{
			{a: Vector3{0, 0, 0}, b: Vector3{1, -2, 3}, expected: 0},
			{a: Vector3{1, 5, 3}, b: Vector3{2, 1, -3}, expected: -2},
		}
		for _, test := range tests {
			actual := dotProduct(test.a, test.b)
			assert.Equal(t, test.expected, actual)
		}
	})
}

// tests/solver_test.go
package packsolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactMatch(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := SolvePackDistribution(2000, sizes)
	expected := 2000
	assert.Equal(t, expected, total)
	assert.GreaterOrEqual(t, len(packs), 1)
}

func TestMinimalExcess(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := SolvePackDistribution(2300, sizes)
	assert.GreaterOrEqual(t, total, 2300)
	assert.NotEmpty(t, packs)
}

func TestZeroQuantity(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := SolvePackDistribution(0, sizes)
	assert.Equal(t, 0, total)
	assert.Empty(t, packs)
}

func TestNoSizesAvailable(t *testing.T) {
	sizes := []int{}
	packs, total := SolvePackDistribution(1000, sizes)
	assert.Equal(t, 0, total)
	assert.Empty(t, packs)
}

func TestLargeQuantity(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	packs, total := SolvePackDistribution(12345, sizes)
	assert.GreaterOrEqual(t, total, 12345)
	assert.NotEmpty(t, packs)
}

// internal/packsolver/solver_test.go
package packsolver_test

import (
	"fmt"
	"github.com/rapido-liebre/pack_solver/internal/packsolver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactMatch(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := packsolver.SolvePackDistribution(2000, sizes)
	expected := 2000
	assert.Equal(t, expected, total)
	assert.GreaterOrEqual(t, len(packs), 1)
}

func TestMinimalExcess(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := packsolver.SolvePackDistribution(2300, sizes)
	assert.GreaterOrEqual(t, total, 2300)
	assert.NotEmpty(t, packs)
}

func TestZeroQuantity(t *testing.T) {
	sizes := []int{250, 500, 1000}
	packs, total := packsolver.SolvePackDistribution(0, sizes)
	assert.Equal(t, 0, total)
	assert.Empty(t, packs)
}

func TestNoSizesAvailable(t *testing.T) {
	sizes := []int{}
	packs, total := packsolver.SolvePackDistribution(1000, sizes)
	assert.Equal(t, 0, total)
	assert.Empty(t, packs)
}

func TestLargeQuantity(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}
	packs, total := packsolver.SolvePackDistribution(12345, sizes)
	assert.GreaterOrEqual(t, total, 12345)
	assert.NotEmpty(t, packs)
}

func TestSolvePackDistribution(t *testing.T) {
	sizes := []int{100, 250, 500, 1000}
	quantity := 12001
	packs, total := packsolver.SolvePackDistribution(quantity, sizes)
	assert.NotEmpty(t, packs)
	assert.GreaterOrEqual(t, total, quantity)

	sum := 0
	for _, p := range packs {
		sum += p.Size * p.Count
	}
	assert.Equal(t, total, sum)
}

func TestVeryLargeQuantity(t *testing.T) {
	sizes := []int{23, 31, 53}
	quantity := 500000

	packs, total := packsolver.SolveSmart(quantity, sizes)

	assert.NotEmpty(t, packs)
	assert.GreaterOrEqual(t, total, quantity)

	sum := 0
	for _, p := range packs {
		sum += p.Size * p.Count
	}
	assert.Equal(t, total, sum)
}

func TestCompareAllStrategies(t *testing.T) {
	sizes := []int{23, 31, 53}
	quantity := 500000

	smartPacks, smartTotal := packsolver.SolveSmart(quantity, sizes)
	dpPacks, dpTotal := packsolver.SolvePackDistribution(quantity, sizes)
	greedyPacks, greedyTotal := packsolver.SolveGreedy(quantity, sizes)

	fmt.Println("=== Smart Strategy ===")
	for _, p := range smartPacks {
		fmt.Printf("Pack %d: %d pcs\n", p.Size, p.Count)
	}
	fmt.Printf("Total: %d\n\n", smartTotal)

	fmt.Println("=== Dynamic Programming ===")
	for _, p := range dpPacks {
		fmt.Printf("Pack %d: %d pcs\n", p.Size, p.Count)
	}
	fmt.Printf("Total: %d\n\n", dpTotal)

	fmt.Println("=== Greedy Strategy ===")
	for _, p := range greedyPacks {
		fmt.Printf("Pack %d: %d pcs\n", p.Size, p.Count)
	}
	fmt.Printf("Total: %d\n", greedyTotal)

	assert.GreaterOrEqual(t, smartTotal, quantity)
	assert.GreaterOrEqual(t, dpTotal, quantity)
	assert.GreaterOrEqual(t, greedyTotal, quantity)
}

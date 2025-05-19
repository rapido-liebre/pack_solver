package packsolver

import (
	"math"
	"sort"
)

// PackResult represents one pack size and the number of times it's used.
type PackResult struct {
	Size  int `json:"size"`  // size of the pack
	Count int `json:"count"` // how many times this pack is used
}

// SolveSmart runs greedy and DP and picks the better result based on minimal total.
func SolveSmart(quantity int, sizes []int) ([]PackResult, int) {
	greedy, gTotal := SolveGreedy(quantity, sizes)
	dp, dTotal := SolvePackDistribution(quantity, sizes)

	if dTotal <= gTotal {
		return dp, dTotal
	}
	return greedy, gTotal
}

// SolvePackDistribution uses dynamic programming to find the minimal total quantity of packs
// whose sum is equal or greater than the requested quantity.
func SolvePackDistribution(quantity int, sizes []int) ([]PackResult, int) {
	if len(sizes) == 0 || quantity <= 0 {
		return []PackResult{}, 0
	}

	// Find the largest pack size to set DP search limit
	maxSize := 0
	for _, s := range sizes {
		if s > maxSize {
			maxSize = s
		}
	}

	limit := quantity + maxSize         // allow room for small overage
	dp := make([]int, limit+1)          // dp[i] = min total for i
	packCount := make([][]int, limit+1) // packCount[i] = how many of each size for dp[i]

	for i := 1; i <= limit; i++ {
		dp[i] = math.MaxInt32
	}
	dp[0] = 0
	packCount[0] = make([]int, len(sizes))

	for i := 1; i <= limit; i++ {
		for j, size := range sizes {
			if i >= size && dp[i-size] != math.MaxInt32 {
				if dp[i] > dp[i-size]+size {
					dp[i] = dp[i-size] + size
					packCount[i] = append([]int(nil), packCount[i-size]...)
					packCount[i][j]++
				}
			}
		}
	}

	// Find first valid solution >= quantity
	bestTotal := -1
	for i := quantity; i <= limit; i++ {
		if dp[i] != math.MaxInt32 {
			bestTotal = i
			break
		}
	}
	if bestTotal == -1 {
		return []PackResult{}, 0
	}

	var result []PackResult
	for i, count := range packCount[bestTotal] {
		if count > 0 {
			result = append(result, PackResult{
				Size:  sizes[i],
				Count: count,
			})
		}
	}

	return result, bestTotal
}

// SolveGreedy prefers large packs first, then adjusts to match the quantity.
func SolveGreedy(quantity int, sizes []int) ([]PackResult, int) {
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	remaining := quantity
	packMap := make(map[int]int)

	for _, size := range sizes {
		count := remaining / size
		packMap[size] = count
		remaining -= count * size
	}

	// Try to cover any leftover amount using smallest possible packs
	for remaining > 0 {
		for _, size := range sizes {
			if size >= remaining {
				packMap[size]++
				remaining -= size
				break
			}
		}
		// If nothing fits, break
		if remaining > 0 && packMap[sizes[len(sizes)-1]] == 0 {
			break
		}
	}

	total := 0
	results := []PackResult{}
	for _, size := range sizes {
		count := packMap[size]
		if count > 0 {
			results = append(results, PackResult{Size: size, Count: count})
			total += size * count
		}
	}

	return results, total
}

// SolvePackDistribution2 attempts to find the combination of pack sizes that
// covers the given quantity with the least total number of items.
//
// Parameters:
// - quantity: number of items to be packed
// - sizes: available pack sizes (must be positive integers)
//
// Returns:
// - slice of PackResult (each containing Size and Count)
// - total number of packed items (which may be slightly more than quantity)
//
// Strategy:
// - uses a depth-first search (DFS) to try all combinations
// - explores combinations recursively
// - minimizes total packed items (not number of packs)
func SolvePackDistribution2(quantity int, sizes []int) ([]PackResult, int) {
	var best []PackResult          // best combination found so far
	minTotal := int(^uint(0) >> 1) // set to MaxInt

	// recurse is a recursive DFS function to explore combinations
	var recurse func(index, remaining, currentTotal int, current []PackResult)
	recurse = func(index, remaining, currentTotal int, current []PackResult) {
		// Base case: if remaining is <= 0, we found a valid or overfilled combo
		if remaining <= 0 {
			if currentTotal < minTotal {
				minTotal = currentTotal
				best = make([]PackResult, len(current))
				copy(best, current)
			}
			return
		}

		// If we’ve used all pack sizes and still not satisfied quantity
		if index == len(sizes) {
			return
		}

		packSize := sizes[index]

		// maxCount = how many times we can use this packSize without going too far
		maxCount := (remaining + packSize - 1) / packSize // ceil division

		// Try using this pack size from 0 up to maxCount times
		for count := 0; count <= maxCount; count++ {
			// Create a fresh copy of the current path (to preserve state)
			next := append([]PackResult{}, current...)

			// Only append if we’re actually using this pack size
			if count > 0 {
				next = append(next, PackResult{Size: packSize, Count: count})
			}

			// Recurse to the next pack size
			recurse(index+1, remaining-count*packSize, currentTotal+count*packSize, next)
		}
	}

	// Start DFS from index 0
	recurse(0, quantity, 0, []PackResult{})

	if best == nil {
		return []PackResult{}, 0
	}

	return best, minTotal
}

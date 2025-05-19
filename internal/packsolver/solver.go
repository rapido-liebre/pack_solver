package packsolver

// PackResult represents one pack size and the number of times it's used.
type PackResult struct {
	Size  int // size of the pack
	Count int // how many times this pack is used
}

// SolvePackDistribution attempts to find the combination of pack sizes that
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
func SolvePackDistribution(quantity int, sizes []int) ([]PackResult, int) {
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

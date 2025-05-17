package packsolver

import (
	"sort"
)

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

// SolvePackDistribution calculates the optimal combination of packs to fulfill the order quantity.
// It returns the list of pack sizes with their respective counts and the total quantity packed (which may include excess).
func SolvePackDistribution(quantity int, packSizes []int) ([]Pack, int) {
	// prioritize larger packs
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	bestTotal := 0
	var bestCombo []Pack

	var dfs func(int, []Pack, int)
	dfs = func(remaining int, current []Pack, total int) {
		if total >= quantity && (bestTotal == 0 || total < bestTotal || (total == bestTotal && len(current) < len(bestCombo))) {
			bestTotal = total
			bestCombo = append([]Pack{}, current...)
			return
		}
		if total >= quantity || len(current) > len(packSizes)*2 {
			return
		}
		for _, size := range packSizes {
			dfs(remaining-size, append(current, Pack{Size: size, Count: 1}), total+size)
		}
	}

	dfs(quantity, []Pack{}, 0)

	// collapse duplicates (e.g., 250,250,250 -> 3x250)
	countMap := make(map[int]int)
	for _, p := range bestCombo {
		countMap[p.Size]++
	}
	result := make([]Pack, 0, len(countMap))
	for size, count := range countMap {
		result = append(result, Pack{Size: size, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	return result, bestTotal
}

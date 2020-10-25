package object

import (
	"sort"
	"strings"
)

// Pair is a PanObject key-value pair for obj or map items.
type Pair struct {
	Key   PanObject
	Value PanObject
}

func sortedPairsString(pairs []Pair) string {
	type PairStr struct {
		k string
		v string
	}

	pairStrs := []PairStr{}
	for _, p := range pairs {
		pairStr := PairStr{k: p.Key.Inspect(), v: p.Value.Inspect()}
		pairStrs = append(pairStrs, pairStr)
	}

	sort.Slice(
		pairStrs,
		func(i, j int) bool { return pairStrs[i].k < pairStrs[j].k },
	)

	sortedStrs := []string{}
	for _, p := range pairStrs {
		sortedStrs = append(sortedStrs, p.k+": "+p.v)
	}

	return strings.Join(sortedStrs, ", ")
}

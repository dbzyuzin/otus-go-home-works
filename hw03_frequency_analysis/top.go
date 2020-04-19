package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

func Top10(inp string) []string {
	words := strings.Fields(inp)

	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}

	words = words[:0]
	for word := range frequency {
		words = append(words, word)
	}
	if len(words) <= 10 {
		return words
	}

	sort.Slice(words, func(i, j int) bool {
		return frequency[words[i]] > frequency[words[j]]
	})
	return words[:10]
}

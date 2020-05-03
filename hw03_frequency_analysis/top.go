package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

func Top10(inp string) []string {
	const need = 10
	words := strings.Fields(inp)

	frequency := make(map[string]int)
	for _, word := range words {
		word = strings.Trim(word, "!,.\"-")
		word = strings.ToLower(word)
		if word != "" {
			frequency[word]++
		}
	}

	words = make([]string, len(frequency))
	for word := range frequency {
		words = append(words, word)
	}
	if len(words) <= need {
		return words
	}

	sort.Slice(words, func(i, j int) bool {
		return frequency[words[i]] > frequency[words[j]]
	})

	return words[:need]
}

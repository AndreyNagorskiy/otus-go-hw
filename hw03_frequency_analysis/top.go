package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const WordsLimit = 10

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	wordCounts := make(map[string]int)
	splitWords := strings.Fields(text)

	for _, word := range splitWords {
		formattedWord := formatWord(word)
		if isValidWord(formattedWord) {
			wordCounts[formattedWord]++
		}
	}

	return topWords(wordCounts, WordsLimit)
}

func topWords(wordCounts map[string]int, limit int) []string {
	keys := make([]string, 0, len(wordCounts))

	for k := range wordCounts {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if wordCounts[keys[i]] != wordCounts[keys[j]] {
			return wordCounts[keys[i]] > wordCounts[keys[j]]
		}
		return keys[i] < keys[j]
	})

	if len(keys) < limit {
		limit = len(keys)
	}

	return keys[:limit]
}

func formatWord(w string) string {
	if isHyphenString(w) {
		return w
	}

	re := regexp.MustCompile(`^[^a-zA-Zа-яА-Я0-9]+|[^a-zA-Zа-яА-Я0-9]+$`)
	w = strings.ToLower(strings.TrimSpace(w))
	return re.ReplaceAllString(w, "")
}

func isHyphenString(s string) bool {
	return regexp.MustCompile(`^-+$`).MatchString(s)
}

func isValidWord(w string) bool {
	return w != "-" && w != ""
}

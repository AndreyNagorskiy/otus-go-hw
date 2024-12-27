package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const WordsLimit = 10

var (
	invalidCharsRegex = regexp.MustCompile(`^[^\p{L}\p{N}\p{So}]+|[^\p{L}\p{N}\p{So}]+$`)
	hyphenRegex       = regexp.MustCompile(`^-+$`)
)

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

	res := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		res = append(res, keys[i])
	}

	return res
}

func formatWord(w string) string {
	if isHyphenString(w) {
		return w
	}

	w = strings.ToLower(strings.TrimSpace(w))

	return invalidCharsRegex.ReplaceAllString(w, "")
}

func isHyphenString(s string) bool {
	return hyphenRegex.MatchString(s)
}

func isValidWord(w string) bool {
	return w != "-" && w != ""
}

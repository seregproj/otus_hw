package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var e = regexp.MustCompile(`[\S]+`)

func Top10(str string) []string {
	fields := strings.FieldsFunc(str, func(r rune) bool {
		return unicode.IsPunct(r) && r != []rune("-")[0]
	})

	wordsFrequency := map[string]int{}
	for _, v := range fields {
		subs := e.FindAllString(v, -1)

		for _, v := range subs {
			trimmedStr := strings.Trim(v, "-")

			if len(trimmedStr) > 0 {
				wordsFrequency[strings.ToLower(trimmedStr)]++
			}
		}
	}

	type wordCounts struct {
		w string
		n int
	}

	wordCountsSlice := make([]wordCounts, 0, len(wordsFrequency))
	for k, v := range wordsFrequency {
		wordCountsSlice = append(wordCountsSlice, wordCounts{
			k,
			v,
		})
	}

	sort.Slice(wordCountsSlice, func(i, j int) bool {
		switch {
		case wordCountsSlice[i].n > wordCountsSlice[j].n:
			return true
		case wordCountsSlice[i].n == wordCountsSlice[j].n:
			return wordCountsSlice[i].w < wordCountsSlice[j].w
		default:
			return false
		}
	})

	upperBound := 10
	if len(wordCountsSlice) < upperBound {
		upperBound = len(wordCountsSlice)
	}
	wordCountsSlice = wordCountsSlice[:upperBound]

	topItems := make([]string, 0, 10)
	for _, v := range wordCountsSlice {
		topItems = append(topItems, v.w)
	}

	return topItems
}

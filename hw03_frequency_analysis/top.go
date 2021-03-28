package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var e = regexp.MustCompile(`[^!"#$%&'()*+,./:;<=>?@[\]^_{|}~]+`)

func Top10(str string) []string {
	fields := strings.Fields(str)
	wordsFrequency := map[string]int{}

	for _, v := range fields {
		subs := e.FindAllStringSubmatch(v, -1)

		for _, v := range subs {
			trimmedStr := strings.Trim(v[0], "-")

			if len(trimmedStr) > 0 {
				wordsFrequency[strings.ToLower(trimmedStr)]++
			}
		}
	}

	usesFrequency := map[int][]string{}

	for k, v := range wordsFrequency {
		usesFrequency[v] = append(usesFrequency[v], k)
	}

	usesKeys := make([]int, 0, len(usesFrequency))
	for k := range usesFrequency {
		usesKeys = append(usesKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(usesKeys)))

	topItems := make([]string, 0, 10)
	var upperBound int
	for _, v := range usesKeys {
		cntItems := len(topItems)

		if cntItems == 10 {
			break
		}

		if len(usesFrequency[v]) > (10 - cntItems) {
			upperBound = 10 - cntItems
		} else {
			upperBound = len(usesFrequency[v])
		}

		sort.Strings(usesFrequency[v])
		topItems = append(topItems, usesFrequency[v][:upperBound]...)
	}

	return topItems
}

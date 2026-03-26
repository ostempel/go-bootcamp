package stringutil

import "strings"

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	half := len(s) / 2

	if len(s)%2 == 0 {
		firstHalf := runes[0:half]
		secondHalf := runes[half:]

		return string(firstHalf) == Reverse(string(secondHalf))
	} else {
		firstHalf := runes[0 : half+1]
		secondHalf := runes[half:]

		return string(firstHalf) == Reverse(string(secondHalf))
	}
}

func WordCount(s string) map[string]int {
	wordMap := make(map[string]int)

	words := strings.Fields(strings.ToLower(s))
	for _, w := range words {
		wordMap[w]++
	}

	return wordMap
}

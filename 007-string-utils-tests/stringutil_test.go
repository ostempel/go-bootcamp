package stringutil

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single rune", input: "a", expected: "a"},
		{name: "word string", input: "animal", expected: "lamina"},
		{name: "sentence", input: "Hello, how are you?", expected: "?uoy era woh ,olleH"},
		{name: "unicode", input: "über", expected: "rebü"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Reverse(tt.input)
			if output != tt.expected {
				t.Errorf("invalid reverse of input %q. expected %q but got %q", tt.input, tt.expected, output)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "empty string", input: "", expected: true},
		{name: "single rune", input: "a", expected: true},
		{name: "2 rune palindrome", input: "aa", expected: true},
		{name: "no palindrome, 2 chars", input: "ab", expected: false},
		{name: "unicode", input: "übü", expected: false},
		{name: "3 rune palindrome", input: "aba", expected: true},
		{name: "simple palindrome", input: "anna", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := IsPalindrome(tt.input)
			if tt.expected != output {
				t.Errorf("wrong isPalindrome for %q. expected %t but got %t", tt.input, tt.expected, output)
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{name: "empty string", input: "", expected: map[string]int{}},
		{name: "single char", input: "a", expected: map[string]int{"a": 1}},
		{name: "two separate chars", input: "a b", expected: map[string]int{"a": 1, "b": 1}},
		{name: "two same chars", input: "a a", expected: map[string]int{"a": 2}},
		{name: "word", input: "house", expected: map[string]int{"house": 1}},
		{name: "words", input: "house animal cat", expected: map[string]int{"house": 1, "animal": 1, "cat": 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := WordCount(tt.input)
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("wrong WordCount for %q. expected %v but got %v", tt.input, tt.expected, output)
			}
		})
	}
}

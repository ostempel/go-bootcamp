# Testing

```go
func Test(t *testing.T) {
	tests := []struct{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test logic -> fail with t.Errorf()
		})
	}
}
```

## Table Tests

[]struct for each test case

```go
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
```

## Utilities

- `reflect.DeepEqual` for maps

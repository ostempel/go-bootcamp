# Types

## Enums

- no native support of enums
- iota (numbured) const values for enums

```go
type Status int

const (
	StatusOpen Status = iota
	StatusDone
)

var (
	stateValue = map[Status]string{
		StatusDone: "done",
		StatusOpen: "open",
	}
	stateMap = map[string]Status{
		"done": StatusDone,
		"open": StatusOpen,
	}
)
```

## Interfaces

- are fulfilled implicitly
- determines which functions needs to be implemented by type to be of this interface

```go
// interface
type Priortizer interface {
	Label() string
}

// type
type PriorityNormal struct{}

// type implements interface Prioritizer
func (p PriorityNormal) Label() string {
	return "normal"
}

// type referencing interface
type Task struct {
	ID       string
	Title    string
	Status   Status
	Priority Priortizer // -> can be now PriorityNormal type because it's implementing the interface
}
```

## Strings

- []bytes
- indexings yiels bytes not char/rune

- strings: arbitrary bytes
- string literals: always contain UTF-8

- `range`-loop: encode one rune each iteration -> so not by index

## Codepoint and Rune

- codepoint: unicode standard - refer to item presented by single value
- rune: alias to `int32`
- can be thought as character

`len(string) = Bytes AND IS NOT len([]rune(string)) = chars`

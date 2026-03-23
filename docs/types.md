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

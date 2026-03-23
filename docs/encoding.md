# Encoding and decoding

## JSON

### json.Marshal
- go struct -> `[]byte`

## json.Unmarshal
- `[]byte` -> go struct

## MarshalJSON()
- `Marshaler` interface -> types can marshal themselv into a JSON
- custom marshal logic of type
- e.g:
  - Marshal enum -> must be string
  
  
### Interfaces
- specific types needs to marshal to unique label
- later for unmarshaling determines which type needs to be used

```go
type Priortizer interface {
	Label() string
}

type PriorityNormal struct{}

func (p PriorityNormal) Label() string {
	return "normal"
}

func (p PriorityNormal) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Label())
}
```
  
## UnmarshalJSON()
- `Unmarshaler` interface -> types can unmarshal themselves using a JSON
- custom unmarshal logic of type
- e.g:
  - Unmarshal enum (string in json) -> enum (through enum maps)

### Interfaces in struct
- done through **helper** struct with `string` instead of `interface`
- not straightforward -> what specific type should be used
- unmarshal based on label in json

```go
func (t *Task) UnmarshalJSON(data []byte) error {
	// helper struct which has no interface
	var helper struct {
		ID       string `json:"ID"`
		Title    string `json:"Title"`
		Status   Status `json:"Status"`
		Priority string `json:"Priority"`
	}
	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}

	t.ID = helper.ID
	t.Title = helper.Title
	t.Status = helper.Status

	// set priority type explicitly
	switch helper.Priority {
	case "urgent":
		t.Priority = PriorityUrgent{}
	case "normal":
		t.Priority = PriorityNormal{}
	default:
		return fmt.Errorf("unknown priority: %s", helper.Priority)
	}

	return nil
}
```

- unmarshaling slice of interfaces -> `json.RawMessage`: keep json-bytes for later unmarshaling/concretization
- 2 stage unmarshaling:
  - unmarshal raw
  - loop slice
    - read type
    - unmarshal struct and append
    
```go
func (d *Drawing) UnmarshalJSON(data []byte) error {
	// use RawMessage to keep raw bytes of each element -> later concretize
	var helper struct {
		Shapes []json.RawMessage
	}

	err := json.Unmarshal(data, &helper)
	if err != nil {
		return err
	}

	for _, v := range helper.Shapes {
		var shapeHelper struct {
			Shape string
		}

		err = json.Unmarshal(v, &shapeHelper)
		if err != nil {
			return err
		}

		var dst Shape
		switch shapeHelper.Shape {
		case "circle":
			dst = new(Circle)
		case "rectangle":
			dst = new(Rectangle)
		default:
			return fmt.Errorf("unknown shape: %s", shapeHelper.Shape)
		}

		err := json.Unmarshal(v, dst)
		if err != nil {
			return err
		}

		d.Shapes = append(d.Shapes, dst)
	}

	return nil
}
```

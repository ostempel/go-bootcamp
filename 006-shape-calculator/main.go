package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type Shape interface {
	Area() float64
	String() string
}

type Circle struct {
	Radius int
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(float64(c.Radius), 2)
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle (radius: %d) - Area: %.2f", c.Radius, c.Area())
}

func (c Circle) MarshalJSON() ([]byte, error) {
	var helper struct {
		Radius int
		Shape  string
	}

	helper.Radius = c.Radius
	helper.Shape = "circle"

	data, err := json.Marshal(helper)
	return data, err
}

type Rectangle struct {
	Width  int
	Height int
}

func (r Rectangle) Area() float64 {
	return float64(r.Height) * float64(r.Width)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle (%dx%d) - Area: %.2f", r.Width, r.Height, r.Area())
}

func (r Rectangle) MarshalJSON() ([]byte, error) {
	var helper struct {
		Width  int
		Height int
		Shape  string
	}

	helper.Width = r.Width
	helper.Height = r.Height
	helper.Shape = "rectangle"

	return json.Marshal(helper)
}

type Drawing struct {
	Shapes []Shape
}

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

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("couldn't determine cmd")
	}

	programArgs := args[1:]
	cmd := programArgs[0]
	cmdArgs := programArgs[1:]

	switch cmd {
	case "add":
		if len(cmdArgs) < 2 {
			log.Fatal("add needs at least 2-3 args (depending on shape type)")
		}

		shapeType := cmdArgs[0]
		shapeArgs := cmdArgs[1:]

		err := addShape(shapeType, shapeArgs)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("successfully added shape")
	case "list":
		if len(cmdArgs) != 0 {
			log.Fatal("list takes no args")
		}

		err := listShapes()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown cmd: %s\n", cmd)
	}
}

func addShape(shapeType string, args []string) error {
	d, err := getStore()
	if err != nil {
		return err
	}

	switch shapeType {
	case "circle":
		if len(args) != 1 {
			return fmt.Errorf("circle shape just needs 1 args (radius)")
		}

		radius, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("unable to parse radius: %s", args[0])
		}

		d.Shapes = append(d.Shapes, Circle{
			Radius: radius,
		})
	case "rectangle":
		if len(args) != 2 {
			return fmt.Errorf("rectangle shape needs 2 args (width, height)")
		}

		width, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("unable to parse width: %s", args[0])
		}
		height, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("unable to parse height: %s", args[1])
		}

		d.Shapes = append(d.Shapes, Rectangle{
			Width:  width,
			Height: height,
		})
	default:
		return fmt.Errorf("unknown shape: %s", shapeType)
	}

	return saveStore(d)
}

func listShapes() error {
	d, err := getStore()
	if err != nil {
		return err
	}

	for _, v := range d.Shapes {
		fmt.Println(v)
	}
	return nil
}

func getStore() (*Drawing, error) {
	f, err := os.ReadFile("shapes.json")
	if err != nil {
		if os.IsNotExist(err) {
			return &Drawing{Shapes: []Shape{}}, nil
		}
		return nil, err
	}

	var drawing = &Drawing{}
	err = json.Unmarshal(f, drawing)
	if err != nil {
		return nil, err
	}

	return drawing, nil
}

func saveStore(d *Drawing) error {
	jsonDrawing, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("shapes.json", jsonDrawing, 0644)
}

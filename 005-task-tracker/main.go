package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Status int

const (
	StatusOpen Status = iota
	StatusDone
)

var (
	stateName = map[Status]string{
		StatusDone: "done",
		StatusOpen: "open",
	}
	stateMap = map[string]Status{
		"done": StatusDone,
		"open": StatusOpen,
	}
)

func (s Status) String() string {
	return stateName[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	value, ok := stateMap[status]
	if !ok {
		return fmt.Errorf("unable to unmarshal status %s", status)
	}

	*s = value

	return nil
}

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

type PriorityUrgent struct{}

func (p PriorityUrgent) Label() string {
	return "urgent"
}

func (p PriorityUrgent) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Label())
}

type Task struct {
	ID       string
	Title    string
	Status   Status
	Priority Priortizer
}

type Helper struct {
	Priority string
}

func (t Task) String() string {
	return fmt.Sprintf("[%s] [%s] [%s] %s", t.ID, t.Status, t.Priority.Label(), t.Title)
}

// Don't implement MarshalJSON -> endless recursion if naive implemented
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

func main() {
	args := os.Args

	programmArgs := args[1:]

	if len(programmArgs) == 0 {
		log.Fatalf("no args for task-tracker")
	}

	cmd := programmArgs[0]
	cmdArgs := programmArgs[1:]

	switch cmd {
	case "add":
		if len(cmdArgs) != 2 {
			log.Fatal("needs TITLE and PRIORITY arg")
		}
		var (
			title    = cmdArgs[0]
			priority = cmdArgs[1]
		)

		err := addTask(title, priority)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("successfully added task [%s, %s]", title, priority)
	case "done":
		if len(cmdArgs) != 1 {
			log.Fatal("needs ID of task")
		}
		var id = cmdArgs[0]

		err := doneTask(id)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("successfully finished task %s", id)
	case "list":
		err := listTasks()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("possible commands are: add, done, list")
	}

}

func addTask(title, priority string) error {
	store, err := getStore()
	if err != nil {
		return fmt.Errorf("error getting store: %w", err)
	}

	var p Priortizer
	switch priority {
	case "urgent":
		p = PriorityUrgent{}
	case "normal":
		p = PriorityNormal{}
	default:
		return fmt.Errorf("unknown priority: %s", priority)
	}

	task := &Task{
		ID:       uuid.New().String(),
		Title:    title,
		Status:   StatusOpen,
		Priority: p,
	}
	store[task.ID] = *task

	err = saveStore(store)
	if err != nil {
		return fmt.Errorf("error saving store: %w", err)
	}

	return nil
}

func doneTask(id string) error {
	store, err := getStore()
	if err != nil {
		return fmt.Errorf("error getting store: %w", err)
	}

	task, ok := store[id]
	if !ok {
		return fmt.Errorf("task with id %s not found", id)
	}

	task.Status = StatusDone
	store[id] = task

	err = saveStore(store)
	if err != nil {
		return fmt.Errorf("error saving store: %w", err)
	}

	return nil
}

func listTasks() error {
	store, err := getStore()
	if err != nil {
		return fmt.Errorf("error getting store: %w", err)
	}

	for _, v := range store {
		fmt.Println(v)
	}

	return nil
}

func getStore() (map[string]Task, error) {
	f, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Task), nil
		}
		return nil, err
	}

	var tasks map[string]Task
	err = json.Unmarshal(f, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, err
}

func saveStore(tasks map[string]Task) error {
	jsonTasks, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("tasks.json", jsonTasks, 0644)
}

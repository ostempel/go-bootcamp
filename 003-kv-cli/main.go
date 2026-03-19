package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("not enough args to determine command")
	}

	args = args[1:]
	action := args[0]
	commandArgs := args[1:]

	switch action {
	case "set":
		if len(commandArgs) != 2 {
			log.Fatal("set needs key and value")
		}

		var (
			key   = commandArgs[0]
			value = commandArgs[1]
		)

		err := set(key, value)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("successfully set key %q to value %q", key, value)
	case "get":
		if len(commandArgs) != 1 {
			log.Fatal("get needs key")
		}

		var key = commandArgs[0]

		value, err := get(key)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("result - key: %q - value: %q", key, value)
	case "delete":
		if len(commandArgs) != 1 {
			log.Fatal("delete needs key")
		}

		var key = commandArgs[0]

		err := deleteKey(key)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("successfully deleted key %q", key)
	case "list":
		if len(commandArgs) != 0 {
			log.Fatal("list needs no argument")
		}

		keys, err := list()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("result: %q", keys)
	}
}

func set(key, value string) error {
	store, err := getStore()
	if err != nil {
		return err
	}

	store[key] = value

	err = saveStore(store)
	if err != nil {
		return err
	}

	return nil
}

func get(key string) (string, error) {
	store, err := getStore()
	if err != nil {
		return "", err
	}

	value, ok := store[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return value, nil
}

func deleteKey(key string) error {
	store, err := getStore()
	if err != nil {
		return err
	}

	_, ok := store[key]
	if !ok {
		return errors.New("key not found")
	}

	delete(store, key)

	err = saveStore(store)
	if err != nil {
		return err
	}
	return nil
}

func list() ([]string, error) {
	store, err := getStore()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}

	return keys, nil
}

func getStore() (map[string]string, error) {
	// os.ReadFile returns whole content of file
	f, err := os.ReadFile("store.json")
	if err != nil {
		if os.IsNotExist(err) {
			store := make(map[string]string)
			return store, nil
		}
		return nil, err
	}

	var store map[string]string
	err = json.Unmarshal(f, &store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func saveStore(store map[string]string) error {
	jsonStore, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	// os.WriteFile creates with permission if not existing otherwise truncates and writes
	return os.WriteFile("store.json", jsonStore, 0644)
}

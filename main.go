package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type KVStore struct {
	store map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (k *KVStore) Get(key string) (string, bool) {
	value, ok := k.store[key]
	return value, ok

}

func (k *KVStore) Set(key, value string) {
	k.store[key] = value
}

func (k *KVStore) Delete(key string) {
	delete(k.store, key)
}

func (k *KVStore) PrintValue(key string) {
	value, ok := k.Get(key)
	if ok {
		fmt.Println(value)
	} else {
		fmt.Println("key not found")
	}
}
func Help() {
	fmt.Println("Commands:")
	fmt.Println("set <key> <value>")
	fmt.Println("get <key>")
	fmt.Println("delete <key>")
	fmt.Println("exit")
}

func introduction() {
	fmt.Print("Please type Get then the value you want to return\n")
	fmt.Print("Please type Set then the Key, and value you want to save\n")
	fmt.Print("Please type Delete  then the Key you want to delete\n")
}

func main() {
	store := NewKVStore()
	introduction()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}
		//split the input into the operation(get set delete) the key and value if any
		parts := strings.SplitN(input, " ", 3)
		command := strings.ToLower(parts[0])

		switch command {
		case "get":
			if len(parts) != 2 {
				fmt.Println("usage: get <key>")
				continue
			}
			store.PrintValue(parts[1])

		case "set":
			if len(parts) != 3 {
				fmt.Println("usage: set <key> <value>")
				continue
			}
			store.Set(parts[1], parts[2])
			fmt.Println("OK")

		case "delete":
			if len(parts) != 2 {
				fmt.Println("usage: delete <key>")
				continue
			}
			store.Delete(parts[1])
			fmt.Println("OK")

		case "help":
			Help()

		case "exit":
			return

		default:
			fmt.Println("unknown command")
		}
	}
}

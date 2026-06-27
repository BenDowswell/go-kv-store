package main

func main() {
	store := NewKVStore()
	httpserver(store)
}

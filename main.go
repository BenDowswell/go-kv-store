package main

func main() {
	store := NewKVStore()
	introduction()
	store.RunConsole()
}

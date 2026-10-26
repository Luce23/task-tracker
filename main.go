package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}

	fmt.Println("Аргументы получены:", os.Args[1:])
}

func printUsage() {
	fmt.Println("Команды: add <описание>, list, update <id> <описание>, delete <id>")
}

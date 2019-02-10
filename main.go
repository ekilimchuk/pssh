package main

import (
	"fmt"
	"os"
)

func showUsage() {
	fmt.Println("Usage: ./pssh <run|list> -h")
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run":
		runAction()
	case "list":
		listAction()
	default:
		showUsage()
		os.Exit(1)
	}
}

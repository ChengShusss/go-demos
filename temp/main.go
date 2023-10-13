package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Err Args\n")
		os.Exit(1)
	}

	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	err := os.MkdirAll(os.Args[1], os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create file, err: %v\n", err)
	}
}

package main

import (
	"fmt"
	readline "github.com/chzyer/readline"
	"os"
	"strings"
)

func list_directory(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("ls: error: %s\n", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		arg := strings.Split(line, " ")

		switch {
		case arg[0] == "cd":
			if len(arg) > 1 {
				err := os.Chdir(arg[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case arg[0] == "pwd":
			path, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(path)
		case arg[0] == "ls":
			if len(arg) > 1 {
				list_directory(arg[1])
			} else {
				path, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
				}
				list_directory(path)
			}

		default:
			fmt.Printf("Invalid command: %s\n", line)
		}
	}
}

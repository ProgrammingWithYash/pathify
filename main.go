package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
)

const (
	storageDir      = "C:\\pathify"
	pathsFile       = "C:\\pathify\\paths"
	markedPathsFile = "C:\\pathify\\marked_paths"
)

func main() {
	ensureDirectoryExists(storageDir)

	if len(os.Args) < 2 {
		fmt.Println("Usage: pathify <command> [<args>]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		addPath()
	case "mark":
		markPath()
	case "find":
		findPath(pathsFile)
	case "marked":
		findPath(markedPathsFile)
	case "delete":
		deletePath()
	default:
		fmt.Println("Unknown command:", command)
	}
}

func ensureDirectoryExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}
}

func addPath() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	appendToFile(pathsFile, path)
	fmt.Println("Added path:", path)
}

func markPath() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	appendToFile(markedPathsFile, path)
	fmt.Println("Marked path:", path)
}

func findPath(filePath string) {
	paths, err := readLines(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	selectedPath := fuzzyFind(paths)
	if selectedPath != "" {
		fmt.Print(selectedPath)
	}

	os.Chdir(selectedPath)
}

func deletePath() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: pathify delete <path>")
		return
	}

	pathToDelete := os.Args[2]

	deleteFromFile(pathsFile, pathToDelete)
	deleteFromFile(markedPathsFile, pathToDelete)

	fmt.Println("Deleted path:", pathToDelete)
}

func appendToFile(filePath, text string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(text + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func deleteFromFile(filePath, text string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != text {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func fuzzyFind(paths []string) string {
	fmt.Println("Type to search paths. Press Enter to select. Press Ctrl+C to exit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		matches := fuzzy.Find(input, paths)
		if len(matches) == 0 {
			fmt.Println("No matches found.")
			continue
		}

		for i, match := range matches {
			fmt.Printf("%d: %s\n", i+1, match.Str)
		}

		fmt.Print("Select a number: ")
		scanner.Scan()
		selection := scanner.Text()

		index := -1
		fmt.Sscanf(selection, "%d", &index)
		if index > 0 && index <= len(matches) {
			return matches[index-1].Str
		}

		fmt.Println("Invalid selection.")
	}
}

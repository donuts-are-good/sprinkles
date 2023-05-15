package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a .go file as argument.")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Could not open file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var previousLineWasComment, previousLineWasBlank bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		isCommentLine := strings.HasPrefix(trimmedLine, "//")
		isBlankLine := trimmedLine == ""

		// if the line is blank and the previous line was a comment, skip this line.
		if isBlankLine && previousLineWasComment {
			continue
		}

		if isCommentLine && !previousLineWasComment && !previousLineWasBlank {
			lines = append(lines, "")
		}

		if isCommentLine {
			line = strings.ToLower(line)
		}

		lines = append(lines, line)

		previousLineWasComment = isCommentLine
		previousLineWasBlank = isBlankLine
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	file.Close()
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error overwriting file:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}
	_ = writer.Flush()
}

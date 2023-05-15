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

	if !strings.HasSuffix(filePath, ".go") {
		fmt.Println("Invalid file. Please provide a .go file as an argument.")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Could not open file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var previousLineWasComment, previousLineWasBlank bool
	var inBlockComment bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		isCommentLine := strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "/*") || strings.HasSuffix(trimmedLine, "*/")
		isBlankLine := trimmedLine == ""

		// if the line is blank and the previous line was a comment, skip this line.
		if isBlankLine && previousLineWasComment {
			continue
		}

		if isCommentLine && !previousLineWasComment && !previousLineWasBlank {
			lines = append(lines, "")
		}

		if strings.HasPrefix(trimmedLine, "/*") && strings.HasSuffix(trimmedLine, "*/") {

			// if block comment is on a single line, do not split it
			if !strings.Contains(trimmedLine, "\n") {
				line = strings.ToLower(line)
			} else {
				line = "/*" + strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(trimmedLine), "/*"), "*/") + "*/"
			}
		} else if strings.HasPrefix(trimmedLine, "/*") {
			inBlockComment = true
			line = "\n/*\n"
		} else if strings.HasSuffix(trimmedLine, "*/") {
			inBlockComment = true
			line = "*/"
		} else if inBlockComment {
			line = strings.ToLower(line)
		} else if isCommentLine {
			line = strings.ToLower(line)
		}

		lines = append(lines, line)
		previousLineWasComment = isCommentLine
		previousLineWasBlank = isBlankLine
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	file.Close()
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("error overwriting file:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}
	_ = writer.Flush()
}

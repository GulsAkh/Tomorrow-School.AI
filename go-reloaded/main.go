package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 3 {
		return
	}
	tempFile := os.Args[2]

	if os.Args[1] == os.Args[2] {
		fmt.Println("Error: files names are the same")
		return
	}
	if !strings.HasSuffix(tempFile, ".txt") {
		fmt.Println("Error: The otuput file should have .txt extension.")
		return
	}
	content, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error to open file: ", err)
		return
	}
	defer content.Close()
	scanner := bufio.NewScanner(content)
	var text string
	var flag bool
	for scanner.Scan() {
		line := makeCases(scanner.Text())
		line = makePunct(line)
		line = makeArticle(line)

		if flag {
			text += "\n"
		}
		text += line
		flag = true
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}
	err = os.WriteFile(tempFile, []byte(text), 0o644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

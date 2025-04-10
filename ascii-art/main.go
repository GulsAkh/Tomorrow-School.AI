package main

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed standard.txt
var characters string

//go:embed hash.txt
var originHash string

func createHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func saveHashToFile(hash string) error {
	return os.WriteFile("hash.txt", []byte(hash), 0644)
}

func main() {
	if len(os.Args) != 2 {
		return
	}
	arg := os.Args[1]
	if len(arg) == 0 {
		return
	}

	actualHash := createHash(characters)
	originalHash := strings.TrimSpace(originHash)

	if originalHash == "" {
		if err := saveHashToFile(actualHash); err != nil {
			fmt.Println("Error saving hash: ", err)
			os.Exit(1)
		}
	} else if actualHash != originalHash {
		fmt.Println("File was modified")
		os.Exit(1)
	}
	var asciiArt [8]string
	removeCRLF := strings.ReplaceAll(characters, "\r", "")
	lines := strings.Split(removeCRLF, "\n")
	k := 0
	j := 0
	//adds ascii characters to map
	asciiMap := make(map[int][8]string)
	for k = 32; k <= 126; k++ {
		for i := 0; i < 8 && j < len(lines); i++ {
			if len(lines[j]) == 0 {
				j++
				i--
				continue
			}
			asciiArt[i] = lines[j]
			j++
		}
		asciiMap[k] = asciiArt
		asciiArt = [8]string{}
	}
	handleAscii(arg, asciiMap)
}

// prints ascii art
func handleAscii(str string, asciiTable map[int][8]string) {
	strWithNewline := []string{}
	str = strings.Replace(str, `\n`, "\n", -1)
	strWithNewline = strings.Split(str, "\n")

	if len(strWithNewline) > 0 && strWithNewline[0] == "" {
		strWithNewline = strWithNewline[1:]
	}
	for _, line := range strWithNewline {
		if line == "" {
			fmt.Println()
			continue
		}
		for j := 0; j < 8; j++ {
			for _, el := range line {
				if lines, flag := asciiTable[int(el)]; flag {
					fmt.Print(lines[j])
				}
			}
			fmt.Println()
		}
	}
}

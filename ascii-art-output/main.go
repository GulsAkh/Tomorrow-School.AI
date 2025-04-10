package main

import (
	"crypto/sha256"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

//go:embed standard.txt
var standard string

//go:embed shadow.txt
var shadow string

//go:embed thinkertoy.txt
var thinkertoy string

//go:embed hash.txt
var originHash string

type Banner struct {
	Content string
	Valid   bool
}

var banType = map[string]Banner{
	"standard":   {Content: standard, Valid: true},
	"shadow":     {Content: shadow, Valid: true},
	"thinkertoy": {Content: thinkertoy, Valid: true},
}

func main() {
	var arg, banner, characters string
	var file *os.File
	var nFlag = flag.String("output", "<*.txt>", "file name")
	flag.Parse()

	fmt.Printf("%q\n", *nFlag)
	if *nFlag == "" || *nFlag != "<*.txt>" && !strings.HasSuffix(*nFlag, ".txt") {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
		return
	}

	if len(os.Args) == 2 {
		if *nFlag != "<*.txt>" && strings.HasSuffix(os.Args[1], ".txt") {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
			return
		}
		if *nFlag == "<*.txt>" {
			arg = os.Args[1]
		}
		banner = "standard"
		characters = standard
		file = nil
	} else if len(os.Args) < 4 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
		return
	} else {
		if !strings.HasSuffix(os.Args[1], ".txt") {
			fmt.Println("invalid file extension: must be a .txt file")
			return
		}
		arg = os.Args[2]
		if arg == "" {
			return
		}

		var err error
		if *nFlag != "" {
			file, err = os.Create(strings.TrimSpace(*nFlag))
			if err != nil {
				fmt.Println("Error creating file: ", err)
				return
			}
		}
		defer file.Close()

		banner = strings.TrimSpace(os.Args[3])
		if content, flag := banType[banner]; flag && content.Valid {
			characters = content.Content
		} else {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
			return
		}
	}
	if arg == "" {
		return
	}
	if !isPrinted(arg) {
		fmt.Println("Non printable characters")
		return
	}

	fileModifiedError(characters, banner)
	removeCRLF := strings.ReplaceAll(characters, "\r", "")
	lines := strings.Split(removeCRLF, "\n")
	filledMap := addAsciiToMap(lines)

	// if no file.txt and file is nil
	if file != nil {
		handleAscii(arg, filledMap, file)
		return
	}
	// check if stdout is terminal or pipe
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))
	if isTerminal {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGWINCH)
		go func() {
			for range sigChan {
				getTerminalSize()
				handleAscii(arg, filledMap, nil)
			}
		}()
		handleAscii(arg, filledMap, nil)
		select {}
	} else {
		handleAscii(arg, filledMap, nil)
	}
}

// check non-printable characters
func isPrinted(str string) bool {
	for _, char := range str {
		if char < 32 || char > 126 {
			return false
		}
	}
	return true
}

// adds ascii characters to map
func addAsciiToMap(lines []string) map[int][8]string {
	var asciiArt [8]string
	k := 0
	j := 0
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
	return asciiMap
}

// gets terminal size
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size", err)
		os.Exit(1)
	}
	return width, height
}

// prints ascii art to a created file or terminal
func handleAscii(str string, asciiTable map[int][8]string, file *os.File) {
	strWithNewline := []string{}
	// handle newlines in the string argument
	str = strings.Replace(str, `\n`, "\n", -1)
	strWithNewline = strings.Split(str, "\n")
	if len(strWithNewline) > 0 && strWithNewline[0] == "" {
		strWithNewline = strWithNewline[1:]
	}

	// loop over []string argument splitted by newline
	for _, line := range strWithNewline {
		if line == "" {
			if file == nil {
				fmt.Println()
			} else {
				fmt.Fprintln(file)
			}
			continue
		}
		for j := 0; j < 8; j++ {
			// loop over string
			for _, el := range line {
				// check if rune is available in the map's key
				if lines, flag := asciiTable[int(el)]; flag {

					// prints to the created file
					if file != nil {
						_, err := file.Write([]byte(lines[j]))
						if err != nil {
							fmt.Println("Error writing to file: ", err)
							return
						}

						// prints to terminal if no file.txt
					} else {
						var width int
						currWidth := 0
						// check if  stdout is terminal
						isTerminal := term.IsTerminal(int(os.Stdout.Fd()))
						if isTerminal {
							width, _ = getTerminalSize()
							charWidth := len(lines[0])
							// check if not outside of width of terminal
							if currWidth+charWidth > width {
								break
							}
							fmt.Print(lines[j])
							currWidth += charWidth

							// if stdout is pipe, like go run . "hello" | cat -e
						} else {
							fmt.Print(lines[j])
						}
					}
				}
			}
			if file == nil {
				fmt.Println()
			} else {
				fmt.Fprintln(file)
			}
		}
	}
}

// creates hash
func createHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// check if file is modified
func fileModifiedError(characters, banner string) {
	actualHash := createHash(characters)
	lines := strings.Split(strings.TrimSpace(originHash), "\n")

	var storedBanner, storedHash string
	if len(lines) >= 2 {
		storedBanner = lines[0]
		storedHash = lines[1]
	}
	if storedHash == "" || storedBanner == "" {
		data := fmt.Sprintf("%s\n%s", banner, actualHash)
		if err := os.WriteFile("hash.txt", []byte(data), 0644); err != nil {
			fmt.Println("Error saving hash: ", err)
			os.Exit(1)
		}
	} else if banner != storedBanner {
		data := fmt.Sprintf("%s\n%s", banner, actualHash)
		if err := os.WriteFile("hash.txt", []byte(data), 0644); err != nil {
			fmt.Println("Error saving hash: ", err)
			os.Exit(1)
		}
	} else if banner == storedBanner && actualHash != storedHash {
		fmt.Println("File was modified")
		os.Exit(1)
	}
}

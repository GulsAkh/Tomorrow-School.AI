package main

import (
	"crypto/sha256"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed standard.txt
var standard string

//go:embed shadow.txt
var shadow string

//go:embed thinkertoy.txt
var thinkertoy string

//go:embed hash.txt
var originHash string

var banType = map[string]string{
	"standard":   standard,
	"shadow":     shadow,
	"thinkertoy": thinkertoy,
}

func parseColor(flagValue string) (int, int, int, bool) {
	regRGB := regexp.MustCompile(`^(rgb|RGB)\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*\)$`)
	matches := regRGB.FindStringSubmatch(flagValue)
	fmt.Println(matches)
	if len(matches) == 4 {
		r, _ := strconv.Atoi(matches[1])
		g, _ := strconv.Atoi(matches[2])
		b, _ := strconv.Atoi(matches[3])
		return r, g, b, true
	}
	return 0, 0, 0, false
}

func main() {
	//var arg, characters string
	var color
	var nFlag = flag.Var(&flagVal, "color", "color")
	flag.Parse()
	fmt.Println(*nFlag)
	//format := colorFormat(*nFlag)

	// r, g, b, valid := parseColor(*nFlag)
	// fmt.Println(r, g, b, valid)
	//arg = os.Args[2]
	// fmt.Println(arg)
	// if len(os.Args) == 2 {
	// 	arg = strings.TrimSpace(os.Args[1])
	// 	banner = "standard"
	// 	characters = standard
	// } else if len(os.Args) < 4 {
	// 	fmt.Println(`Usage: go run . [OPTION] [STRING] \nEX: go run . --color=<color> <substring to be colored> "something"`)
	// 	return
	// } else {
	// 	arg = strings.TrimSpace(os.Args[2])
	// 	banner = strings.TrimSpace(os.Args[3])
	// 	if content, flag := banType[banner]; flag {
	// 		characters = content
	// 	}
	// }
	// if len(arg) == 0 {
	// 	return
	// }

	//fileModifiedError(characters, banner)
	// removeCRLF := strings.ReplaceAll(characters, "\r", "")
	// lines := strings.Split(removeCRLF, "\n")
	// filledMap := addAsciiToMap(lines)
	// handleAscii(arg, filledMap, r, g, b)
}

func addAsciiToMap(lines []string) map[int][8]string {
	var asciiArt [8]string
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
	return asciiMap
}

// prints ascii art
func handleAscii(str string, asciiTable map[int][8]string, r, g, b int) {
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
					fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, string(lines[j]))
				}
			}
			fmt.Println()
		}
	}
}

func createHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

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

func colorFormat(flagValue string) string {
	//"#ff0000"
	regHex := regexp.MustCompile(`^#([0-9-a-f-A-F]{6})$`)
	// "rgb(255, 0, 0)"
	regRGB := regexp.MustCompile(`^(rgb|RGB)\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*\)$/gm`)
	// "hsl(0, 100%, 50%)"
	regHSL := regexp.MustCompile(`^(hsl|HSL)\(\s*\d+\s*,\s*\d+%\s*,\s*\d+%\s*\)$`)

	if regHex.MatchString(flagValue) {
		return "hex"
	} else if regRGB.MatchString(flagValue) {
		return "rgb"
	} else if regHSL.MatchString(flagValue) {
		return "hsl"
	} else {
		return "color"
	}
}

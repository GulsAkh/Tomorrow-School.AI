package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Struct to collect invalid matches, errors to print out
type InvalidMatch struct {
	match []string
	error []error
}

var invalid InvalidMatch

// Find matches (low|cap|up|bin|hex, <number>) and recursive call of processText
func makeCases(txt string) string {
	if !isPrint(txt) {
		fmt.Println("Error: found not printable characters")
		return ""
	}
	invalid.match = nil
	invalid.error = nil
	processed, stop := processText(txt)
	want1 := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
	match := want1.MatchString(processed)
	for match && !stop {
		processed, stop = processText(processed)
		match = want1.MatchString(processed)

	}
	return processed
}

// Check non-printable chacters
func isPrint(s string) bool {
	flag := true
	for _, el := range s {
		if !(unicode.IsPrint(el) && el <= 127) {
			flag = false
		}
	}
	return flag
}

// Text formatting
func processText(txt string) (string, bool) {
	want1 := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
	match := want1.FindAllStringIndex(txt, 1)
	data := want1.FindAllStringSubmatch(txt, 1)
	var leftText string
	for i := 0; i < len(match); i++ {
		// Get left text before match for further formatting
		if i == 0 {
			leftText = txt[:match[i][0]]
		} else if i == 1 {
			leftText = txt[match[i-1][1]:match[i][0]]
		}
		// Skip if left text is empty or has only symbols
		re := regexp.MustCompile(`[a-zA-Z-0-9-]+`)
		re1 := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
		if leftText == "" && len(data) == 1 {
			txt = strings.Replace(txt, data[i][0], " ", 1)
			return txt, true
		} else if leftText == "" || i > 1 && !re.MatchString(leftText) || i > 1 && re1.MatchString(leftText) {
			continue
		}
		txt_upd := ""
		leftText_changed := strings.Fields(leftText)

		// Check number is valid, otherwise collect match and error to print out
		lenGoBack, err := checkNumber(i, data)
		if err != nil {
			invalid.match = append(invalid.match, data[i][0])
			invalid.error = append(invalid.error, err)
			check := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
			if len(check.FindAllStringSubmatch(txt, -1)) == 1 {
				txt = strings.Replace(txt, data[i][0], " ", 1)
				if invalid.match != nil {
					fmt.Println(collectInvalid(&invalid))
				}
				return txt, true
			} else {
				txt = strings.Replace(txt, data[i][0], " ", 1)
				continue
			}
		}

		// Check number is not above words number, otherwise collect match and error to print out
		lenWords := findLen(leftText_changed)
		check := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
		if (lenGoBack > lenWords) || (lenGoBack > len(leftText_changed) && len(data) == 1) {
			invalid.match = append(invalid.match, data[i][0])
			invalid.error = append(invalid.error, errors.New("error: number is above words number"))
			if len(check.FindAllStringSubmatch(txt, -1)) == 1 {
				txt = strings.Replace(txt, data[i][0], " ", 1)
				if invalid.match != nil {
					fmt.Println(collectInvalid(&invalid))
				}
				return txt, true
			} else {
				txt = strings.Replace(txt, data[i][0], " ", 1)
				continue
			}
		}
		// Formatting text to uppercase | capitalized | lowcase | binary | hexadecimal
		switch data[i][1] {
		case "up":
			txt_upd = upCase(leftText_changed, lenGoBack)
			txt = strings.Replace(txt, leftText, txt_upd, 1)
			txt = strings.Replace(txt, data[i][0], " ", 1)
		case "cap":
			txt_upd = capCase(leftText_changed, lenGoBack)
			txt = strings.Replace(txt, leftText, txt_upd, 1)
			txt = strings.Replace(txt, data[i][0], " ", 1)
		case "low":
			txt_upd = lowCase(leftText_changed, lenGoBack)
			txt = strings.Replace(txt, leftText, txt_upd, 1)
			txt = strings.Replace(txt, data[i][0], " ", 1)
		case "bin":
			spaceEnter := regexp.MustCompile(`(\W)(])(\d+)`)
			if spaceEnter.MatchString(leftText_changed[len(leftText_changed)-1]) {
				leftText_changed, txt, leftText = checkLastSpace(txt, leftText_changed, leftText)
			}
			txt_upd, err = binCase(leftText_changed, lenGoBack)
			if err != nil {
				invalid.match = append(invalid.match, data[i][0])
				invalid.error = append(invalid.error, err)
				check := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
				if len(check.FindAllStringSubmatch(txt, -1)) == 1 {
					txt = strings.Replace(txt, data[i][0], " ", 1)
					if invalid.match != nil {
						fmt.Println(collectInvalid(&invalid))
						invalid.match = []string{}
						invalid.error = []error{}
					}
					return txt, true
				} else {
					txt = strings.Replace(txt, data[i][0], " ", 1)
					continue
				}
			}
			txt = strings.Replace(txt, leftText, txt_upd, 1)
			txt = strings.Replace(txt, data[i][0], " ", 1)
		case "hex":
			spaceEnter := regexp.MustCompile(`(\W)([0-9a-fA-F]+)`)
			if spaceEnter.MatchString(leftText_changed[len(leftText_changed)-1]) {
				leftText_changed, txt, leftText = checkLastSpace(txt, leftText_changed, leftText)
			}
			txt_upd, err = hexCase(leftText_changed, lenGoBack)

			if err != nil {
				invalid.match = append(invalid.match, data[i][0])
				invalid.error = append(invalid.error, err)
				check := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
				if len(check.FindAllStringSubmatch(txt, -1)) == 1 {
					txt = strings.Replace(txt, data[i][0], " ", 1)
					if invalid.match != nil {
						fmt.Println(collectInvalid(&invalid))
					}
					return txt, true
				} else {
					txt = strings.Replace(txt, data[i][0], " ", 1)
					continue
				}
			}
			txt = strings.Replace(txt, leftText, txt_upd, 1)
			txt = strings.Replace(txt, data[i][0], " ", 1)
		}
	}
	txt = spaceRemoval(txt)
	check := regexp.MustCompile(`\(\s*(low|up|cap|bin|hex)\s*(?:\s*,\s*([^\)]*))?\s*\)`)
	if len(check.FindAllStringSubmatch(txt, -1)) == 1 {
		if invalid.match != nil {
			fmt.Println(collectInvalid(&invalid))
		}
	}
	return txt, false
}

// Collect invalid match, error to print out
func collectInvalid(invalid *InvalidMatch) string {
	if invalid.match == nil {
		invalid.match = []string{}
	}
	if invalid.error == nil {
		invalid.error = []error{}
	}
	var result []string
	for i := 0; i < len(invalid.match); i++ {
		if invalid.error[i] == nil {
			result = append(result, invalid.match[i])
		} else {
			result = append(result, invalid.error[i].Error()+" - "+invalid.match[i])
		}
	}
	return strings.Join(result, "\n")
}

// Check number is valid, otherwise collect invalid match
func checkNumber(i int, data [][]string) (int, error) {
	check := regexp.MustCompile(`\s*(-?[0-9]*\.([0-9]*)?)`)
	var lenGoBack int
	var err error
	data[i][2] = strings.TrimSpace(data[i][2])
	if len(data[i]) > 1 {
		if check.MatchString(data[i][2]) {
			return 0, errors.New("error: invalid float number ")
		}
		if data[i][2] == "" {
			lenGoBack = 1
		} else {
			lenGoBack, err = strconv.Atoi(data[i][2])

			if err != nil {
				return 0, errors.New("error: invalid number ")
			}
			if lenGoBack < 0 {
				return 0, errors.New("error: invalid negative number ")
			}
			if lenGoBack == 0 && err == nil {
				lenGoBack = 1
			}
		}
	}
	return lenGoBack, err
}

// Uppercase formatting
func upCase(leftText_changed []string, lenGoBack int) string {
	txt_upd := ""
	regExp := regexp.MustCompile(`[a-zA-Z]+`)
	for i := len(leftText_changed) - 1; i >= 0; i-- {
		if regExp.MatchString(leftText_changed[i]) && lenGoBack > 0 {
			leftText_changed[i] = strings.ToUpper(leftText_changed[i])
			lenGoBack--
		}
		if i == 0 {
			txt_upd = leftText_changed[i] + txt_upd
		} else {
			txt_upd = " " + leftText_changed[i] + txt_upd
		}
	}
	return txt_upd
}

// Capitalized formatting
func capCase(leftText_changed []string, lenGoBack int) string {
	txt_upd := ""
	regExp := regexp.MustCompile(`[a-zA-Z]+`)
	for i := len(leftText_changed) - 1; i >= 0; i-- {
		if regExp.MatchString(leftText_changed[i]) && lenGoBack > 0 {
			match := regExp.FindStringSubmatch(leftText_changed[i])
			if len(match) > 0 && len(leftText_changed[i]) > 0 {
				transformed := strings.ToUpper(string(match[0][0])) + strings.ToLower(match[0][1:])
				leftText_changed[i] = strings.Replace(leftText_changed[i], match[0], transformed, 1)
			}
			lenGoBack--
		}
		if i == 0 {
			txt_upd = leftText_changed[i] + txt_upd
		} else {
			txt_upd = " " + leftText_changed[i] + txt_upd
		}
	}
	return txt_upd
}

// Lowcase formatting
func lowCase(leftText_changed []string, lenGoBack int) string {
	txt_upd := ""
	regExp := regexp.MustCompile(`[a-zA-Z]+`)
	for i := len(leftText_changed) - 1; i >= 0; i-- {
		if regExp.MatchString(leftText_changed[i]) && lenGoBack > 0 {
			leftText_changed[i] = strings.ToLower(leftText_changed[i])
			lenGoBack--
		}
		if i == 0 {
			txt_upd = leftText_changed[i] + txt_upd
		} else {
			txt_upd = " " + leftText_changed[i] + txt_upd
		}
	}
	return txt_upd
}

// Binary version formatted to decimal
func binCase(leftText []string, lenGoBack int) (string, error) {
	txt_upd := ""
	var num uint64
	var err error
	flag := false
	bin := regexp.MustCompile(`[01]+`)
	digit := regexp.MustCompile(`\d+`)
	for i := len(leftText) - 1; i >= 0; i-- {
		if bin.MatchString(leftText[i]) && lenGoBack != 0 {
			num, err = strconv.ParseUint(leftText[i], 2, 64)
			if err != nil {
				return "", errors.New("error: invalid number ")
			}
			leftText[i] = strconv.FormatUint(uint64(num), 10)
			lenGoBack--
			flag = true
		}
		if digit.MatchString(leftText[i]) && lenGoBack != 0 {
			num, err = strconv.ParseUint(leftText[i], 10, 64)
			if err != nil {
				return "", errors.New("error: invalid number ")
			}
			leftText[i] = strconv.FormatUint(uint64(num), 10)
			lenGoBack--
			flag = true
		}

		if i == 0 {
			txt_upd = leftText[i] + txt_upd
		} else {
			txt_upd = " " + leftText[i] + txt_upd
		}
	}
	if !flag {
		return txt_upd, errors.New("error: invalid number ")
	}
	return txt_upd, err
}

// Hexadecimal version formatted to decimal
func hexCase(leftText []string, lenGoBack int) (string, error) {
	txt_upd := ""
	var num uint64
	var err error
	flag := false
	hex := regexp.MustCompile(`\b(?:0[xX])?[0-9a-fA-F]+\b`)
	for i := len(leftText) - 1; i >= 0; i-- {
		if hex.MatchString(leftText[i]) && lenGoBack != 0 {
			leftSlice := leftText[i]
			if len(leftSlice) > 2 && (leftSlice[:2] == "0x" || leftSlice[:2] == "0X") {
				leftText[i] = leftSlice[2:]
			}
			num, err = strconv.ParseUint(leftText[i], 16, 64)
			if err != nil {
				return "", errors.New("error: invalid number ")
			}
			leftText[i] = strconv.FormatUint(uint64(num), 10)
			lenGoBack--
			flag = true
		}
		if i == 0 {
			txt_upd = leftText[i] + txt_upd
		} else {
			txt_upd = " " + leftText[i] + txt_upd
		}
	}
	if !flag {
		return txt_upd, errors.New("error: invalid number ")
	}
	return txt_upd, err
}

// Add a space for cases <,number>
func checkLastSpace(txt string, leftText_changed []string, leftText string) ([]string, string, string) {
	flag := false
	var updated string
	spaceEnter := regexp.MustCompile(`(\W)([0-9a-fA-F]+|\d+)`)
	if spaceEnter.MatchString(leftText_changed[len(leftText_changed)-1]) {
		lastWord := spaceEnter.ReplaceAllString(leftText_changed[len(leftText_changed)-1], "$1 $2")
		leftText_changed[len(leftText_changed)-1] = lastWord
		updated = strings.Join(leftText_changed, " ")
		flag = true
	}
	if flag {
		txt = strings.Replace(txt, leftText, updated, 1)
		leftText = updated
		leftText_changed = strings.Split(leftText, " ")
	}
	return leftText_changed, txt, leftText
}

// Find length of text before match
func findLen(str []string) int {
	var words []string
	for _, el := range str {
		if el != "" {
			words = append(words, el)
		}
	}
	return len(words)
}

// Remove spaces in match if required
func spaceRemoval(s string) string {
	hex := regexp.MustCompile(`\(\s*h\s*e\s*x\s*\)`)
	bin := regexp.MustCompile(`\(\s*b\s*n\s*n\s*\)`)
	cap := regexp.MustCompile(`\(\s*c\s*a\s*p\s*\)`)
	up := regexp.MustCompile(`\(\s*u\s*p\s*\)`)
	low := regexp.MustCompile(`\(\s*l\s*o\s*w\s*\)`)
	if hex.MatchString(s) {
		s = hex.ReplaceAllString(s, "(hex)")
	} else if bin.MatchString(s) {
		s = bin.ReplaceAllString(s, "(bin)")
	} else if cap.MatchString(s) {
		s = cap.ReplaceAllString(s, "(cap)")
	} else if up.MatchString(s) {
		s = up.ReplaceAllString(s, "(up)")
	} else if low.MatchString(s) {
		s = low.ReplaceAllString(s, "(low)")
	}
	return s
}

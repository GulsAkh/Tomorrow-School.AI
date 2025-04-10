package main

import (
	"regexp"
	"strings"
)

// Punctuation formatting
func makePunct(txt string) string {
	//Add a space after punctuation mark if necessary
	want1 := regexp.MustCompile(`\s*([.,:;!?])\s*`)
	txt = want1.ReplaceAllString(txt, "$1")
	want1 = regexp.MustCompile(`([.,:;!?])([a-z-A-Z-0-9])`)
	txt = want1.ReplaceAllString(txt, "$1 $2")

	txt = strings.Join(strings.Fields(txt), " ")

	//single quote and contractions formatting
	want2 := regexp.MustCompile(`(I|he|she|it|you|we|they|He|She|It|You|We|They|can|may|isn|aren|don|Don|Won|Wouldn|doesn|Doesn|what|What|how|How|There|there|this|This|that|That)\s*(')\s*(ll|t|s|ve|d|m|re)\b`)
	txt = want2.ReplaceAllString(txt, "$1$2$3")
	want := regexp.MustCompile(`(['])\s+([.,!?])`)
	txt = want.ReplaceAllString(txt, "$1$2")

	//brackets formatting
	want = regexp.MustCompile(`(\()\s+`)
	txt = want.ReplaceAllString(txt, "$1")
	want = regexp.MustCompile(`\s+(\))`)
	txt = want.ReplaceAllString(txt, "$1")

	var txtPunct []rune
	openQuote := true
	openDoubleQuote := true
	indexOpen := 0
	indexDoubleOpen := 0

	// Formatting single quotes and double quotes
	for i := 0; i < len(txt); i++ {
		char := rune(txt[i])
		if char == '\'' {
			if (i+2 < len(txt) && (txt[i+1] == 's' || txt[i+1] == 't' || txt[i+1] == 'd' || txt[i+1] == 'm') && txt[i+2] == ' ') ||
				(i+3 < len(txt) && (txt[i+1] == 'l' || txt[i+1] == 'v' || txt[i+1] == 'r' && txt[i+2] == 'l' || txt[i+2] == 'e') && txt[i+3] == ' ') {
				txtPunct = append(txtPunct, char)
				continue
			}
			if !openQuote {
				// closing quote found and mark openQuote as true
				openQuote = true
				tempChar := len(txtPunct) - 1
				// removing a space before closing quote
				if txtPunct[tempChar] == 32 {
					txtPunct = txtPunct[:len(txtPunct)-1]
				}
				// add a closing quote after removal of the space
				txtPunct = append(txtPunct, char)
				// if required add a space after closing quote
				if i+1 < len(txt) && txt[i+1] != ' ' {
					txtPunct = append(txtPunct, ' ')
				}
				// removing a space after opening quote
				if indexOpen+1 < len(txtPunct) && txtPunct[indexOpen+1] == 32 {
					txtPunct = append(txtPunct[:indexOpen+1], txtPunct[indexOpen+2:]...)
				}
				// insert a space before an opening quote
				if indexOpen > 0 && txtPunct[indexOpen-1] != ' ' {
					txtPunct = append(txtPunct[:indexOpen], append([]rune{' '}, txtPunct[indexOpen:]...)...)
				}
			} else {
				// opening quote: set indexOpen with quote start index
				txtPunct = append(txtPunct, char)
				indexOpen = len(txtPunct) - 1
				openQuote = false
			}
		} else if char == '"' {
			if (i+2 < len(txt) && (txt[i+1] == 's' || txt[i+1] == 't' || txt[i+1] == 'd' || txt[i+1] == 'm') && txt[i+2] == ' ') ||
				(i+3 < len(txt) && (txt[i+1] == 'l' || txt[i+1] == 'v' || txt[i+1] == 'r' && txt[i+2] == 'l' || txt[i+2] == 'e') && txt[i+3] == ' ') {
				txtPunct = append(txtPunct, char)
				continue
			}
			if !openDoubleQuote {
				// closing quote found and mark openQuote as true
				openDoubleQuote = true
				tempChar := len(txtPunct) - 1
				// removing a space before closing quote
				if txtPunct[tempChar] == 32 {
					txtPunct = txtPunct[:len(txtPunct)-1]
				}
				// add a closing quote after removal of the space
				txtPunct = append(txtPunct, char)
				// if required add a space after closing quote
				if i+1 < len(txt) && txt[i+1] != ' ' {
					txtPunct = append(txtPunct, ' ')
				}
				// removing a space after opening quote
				if indexDoubleOpen+1 < len(txtPunct) && txtPunct[indexDoubleOpen+1] == 32 {
					txtPunct = append(txtPunct[:indexDoubleOpen+1], txtPunct[indexDoubleOpen+2:]...)
				}
				// insert a space before an opening quote
				if indexDoubleOpen > 0 && txtPunct[indexDoubleOpen-1] != ' ' {
					txtPunct = append(txtPunct[:indexDoubleOpen], append([]rune{' '}, txtPunct[indexDoubleOpen:]...)...)
				}
			} else {
				// opening quote: set indexOpen with quote start index
				txtPunct = append(txtPunct, char)
				indexDoubleOpen = len(txtPunct) - 1
				openDoubleQuote = false
			}
		} else {
			// add non-quote characters
			txtPunct = append(txtPunct, char)
		}
	}
	txt_upd := string(txtPunct)
	want1 = regexp.MustCompile(`([."'])\s+([.",!?])`)
	txt_upd = want1.ReplaceAllString(txt_upd, "$1$2")
	return txt_upd
}

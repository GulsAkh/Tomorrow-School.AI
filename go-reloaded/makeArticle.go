package main

import (
	"regexp"
	"strings"
)

// Articles formatting
func makeArticle(txt string) string {
	conjucUp := map[string]bool{
		"OR":  true,
		"AND": true,
		"BUT": true,
		"NOR": true,
		"FOR": true,
		"SO":  true,
		"YET": true,
	}
	conjucLow := map[string]bool{
		"or":  true,
		"and": true,
		"but": true,
		"nor": true,
		"for": true,
		"so":  true,
		"yet": true,
	}
	silentLow := map[string]bool{
		"hour":       true,
		"honor":      true,
		"honest":     true,
		"heir":       true,
		"herb":       true,
		"honour":     true,
		"hymn":       true,
		"humble":     true,
		"hurt":       true,
		"hinterland": true,
		"horrible":   true,
	}
	silentUp := map[string]bool{
		"HOUR":       true,
		"HONOR":      true,
		"HONEST":     true,
		"HEIR":       true,
		"HERB":       true,
		"HONOUR":     true,
		"HYMN":       true,
		"HUMBLE":     true,
		"HURT":       true,
		"HINTERLAND": true,
		"HORRIBLE":   true,
	}

	// an|An|aN|AN + consonants
	cons := regexp.MustCompile(`(an|An|aN|AN)(\s+)([b-df-hj-np-tv-zB-DF-HJ-NP-TV-Z]\w+)`)
	txt = cons.ReplaceAllStringFunc(txt, func(match string) string {
		article := match[:2]
		rest := match[2:]
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			word := parts[0]
			_, flagLow := conjucLow[word]
			_, flagUp := conjucUp[word]
			if flagLow || flagUp {
				return match
			} else {
				if article == "an" || article == "aN" {
					return "a " + word
				} else if article == "An" || article == "AN" {
					return "A " + word
				}
			}
		}
		return match
	})
	// a|A + VOWEL
	capVow := regexp.MustCompile(`\b(a|A)\b(\s+)([AEIOUHY]\w+)`)
	txt = capVow.ReplaceAllStringFunc(txt, func(match string) string {
		re := regexp.MustCompile(`[a-z]+`)
		article := match[:1]
		rest := match[1:]
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			word := parts[0]
			_, flagLow := conjucLow[word]
			_, flagUp := conjucUp[word]
			if flagLow || flagUp {
				return match
			} else {
				if article == "a" {
					return "an " + word
				} else if article == "A" && re.MatchString(string(word[1:])) {
					return "An " + word
				} else if article == "A" {
					return "AN " + word
				}
			}
		}
		return match
	})
	// a|A + vowel
	vowel := regexp.MustCompile(`\b(a|A)\b(\s+)([aeiouhy]\w+)`)
	txt = vowel.ReplaceAllStringFunc(txt, func(match string) string {
		article := match[:1]
		rest := match[1:]
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			word := parts[0]
			_, flagLow := conjucLow[word]
			_, flagUp := conjucUp[word]
			if flagLow || flagUp {
				return match
			} else {
				if article == "a" {
					return "an " + word
				} else if article == "A" {
					return "An " + word
				}
			}
		}
		return match
	})
	// a HOUR
	hUp := regexp.MustCompile(`\b(a|A)\b\s+([H]\w+)`)
	txt = hUp.ReplaceAllStringFunc(txt, func(match string) string {
		article := match[:1]
		rest := match[1:]
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			word := parts[0]
			_, flagLow := silentLow[word]
			_, flagUp := silentUp[word]
			if flagLow || flagUp {
				if article == "a" {
					return "an " + word
				} else if article == "A" {
					return "AN " + word
				}
			}
		}
		return match
	})
	// a hour
	hLow := regexp.MustCompile(`\b(a|A)\b\s+([h]\w+)`)
	txt = hLow.ReplaceAllStringFunc(txt, func(match string) string {
		article := match[:1]
		rest := match[1:]
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			word := parts[0]
			_, flagLow := silentLow[word]
			_, flagUp := silentUp[word]
			if flagLow || flagUp {
				if article == "a" {
					return "an " + word
				} else if article == "A" {
					return "An " + word
				}
			}
		}
		return match
	})
	return txt
}

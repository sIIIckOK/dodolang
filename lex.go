package main

import (
	"unicode"
)

func lexFile(content string, filePath string) []StringToken {
	contentLen := len(content)
	var tokens []StringToken
	var t StringToken
	t.Loc.FilePath = filePath
	i := 0
	bol := 0
	lineNo := 0
	for unicode.IsSpace(rune(content[i])) && i < contentLen {
		if content[i] == '\n' {
			lineNo++
			bol = i + 1
		}
		i++
	}
	start := i
	for ; i < contentLen; i++ {
		if unicode.IsSpace(rune(content[i])) {
			t.Content = content[start:i]
			t.Loc.Col = uint(start - bol + 1)
			t.Loc.Line = uint(lineNo + 1)
			if i-start > 0 {
				tokens = append(tokens, t)
			}
			for unicode.IsSpace(rune(content[i])) && i < contentLen {
				if content[i] == '\n' {
					lineNo++
					bol = i + 1
				}
				i++
			}
			start = i
		}
		if i < contentLen+1 {
			if content[i] == '/' && content[i+1] == '/' {
				for i < contentLen {
					i++
					if content[i] == '\n' {
						break
					}
				}
				lineNo++
				start = i + 1
				bol = i + 1
			}
		}
	}
	return tokens
}

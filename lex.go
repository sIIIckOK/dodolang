package main

import (
	"strings"
	"unicode"
)

func lexFile(content string, filePath string) []StringToken {
    var tokens []StringToken
    var t StringToken
    t.Loc.FilePath = filePath
    lines := strings.Split(content, "\n")
    for lineNo, _line := range lines {
        if len(_line) == 0 { continue }
        t.Loc.Line = uint(lineNo)
        line := strings.Split(_line, "//")[0]
        ptr := 0
        SkipWhiteSpace(line, &ptr)
        start := ptr
        for ;ptr < len(line); ptr++ {
            if unicode.IsSpace(rune(content[ptr])) {
                t.Content = content[start:ptr]
                if len(t.Content) == 0 {
                    continue
                }
                t.Loc.Col = uint(start)
                tokens = append(tokens, t)
                SkipWhiteSpace(content, &ptr)
                start = ptr+1
            }
        }
    }
    return tokens
}

func SkipWhiteSpace(str string, ptr *int) {
    for _, v := range str {
        if !unicode.IsSpace(v) { break }
        *ptr++
    }
}


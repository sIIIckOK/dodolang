package main

import (
	"strconv"
)

func ParseTokens(strTokens []StringToken) []Token {
	var tokens []Token
	var t Token
	assert(TokenCount == 29, "Exhaustive switch case for ParseToken")
	for i, tok := range strTokens {
		if len(tok.Content) == 0 {
			continue
		}
		t.Loc = tok.Loc
		switch tok.Content {
		case "+":
			t.Type = TokenPlus
			tokens = append(tokens, t)
		case "-":
			t.Type = TokenSub
			tokens = append(tokens, t)
		case "*":
			t.Type = TokenMult
			tokens = append(tokens, t)
		case "divmod":
			t.Type = TokenDivMod
			tokens = append(tokens, t)
		case "print":
			t.Type = TokenPrint
			tokens = append(tokens, t)
		case "swap":
			t.Type = TokenSwap
			tokens = append(tokens, t)
		case "dup":
			t.Type = TokenDup
			tokens = append(tokens, t)
		case "drop":
			t.Type = TokenDrop
			tokens = append(tokens, t)
		case ">":
			t.Type = TokenGt
			tokens = append(tokens, t)
		case ">=":
			t.Type = TokenGe
			tokens = append(tokens, t)
		case "<":
			t.Type = TokenLt
			tokens = append(tokens, t)
		case "<=":
			t.Type = TokenLe
			tokens = append(tokens, t)
		case "=":
			t.Type = TokenEq
			tokens = append(tokens, t)
		case "for":
			t.Type = TokenFor
			tokens = append(tokens, t)
		case "do":
			t.Type = TokenDo
			tokens = append(tokens, t)
		case "if":
			t.Type = TokenIf
			tokens = append(tokens, t)
		case "else":
			t.Type = TokenElse
			tokens = append(tokens, t)
		case "end":
			t.Type = TokenEnd
			tokens = append(tokens, t)
		case "macro":
			t.Type = TokenMacro
			tokens = append(tokens, t)
		case "syscall1":
			t.Type = TokenSyscall1
			tokens = append(tokens, t)
		case "syscall3":
			t.Type = TokenSyscall3
			tokens = append(tokens, t)
		case "rot":
			t.Type = TokenRot
			tokens = append(tokens, t)
		case "var":
			t.Type = TokenVar
			tokens = append(tokens, t)
		case "@":
			t.Type = TokenRead
			tokens = append(tokens, t)
		case "!":
			t.Type = TokenWrite
			tokens = append(tokens, t)
		default:
			if num, err := strconv.ParseUint(tok.Content, 10, 64); err == nil {
				t.Type = TokenInt
				t.Operand = num
				tokens = append(tokens, t)
				continue
			}
			t.Type = TokenWord
			t.Operand = uint64(i)
			tokens = append(tokens, t)
		}
	}
	return tokens
}

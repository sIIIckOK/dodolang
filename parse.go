package main

import (
	"fmt"
	"os"
	"strconv"
)

var tokenKindStr = map[string]TokenType{
	"int":  TokenInt,
	"bool": TokenBool,
	"ptr":  TokenPtr,
}

func parseTokens(strTokens []StringToken, state *CompileState) []Token {
	var (
		mainTokenBuffer    []Token
		macroTokenBuffer   []Token
		currentTokenBuffer *[]Token = &mainTokenBuffer
		macroMode          bool
		macroEndStack      int
		currentMacroName   string
		t                  Token
	)
	if len(strTokens) == 0 {
		return mainTokenBuffer
	}
	t.Loc.FilePath = strTokens[0].Loc.FilePath
	assert(TokenCount == 32, "Exhaustive switch case for ParseToken")

	for i := 0; i < len(strTokens); i++ {
		strTok := strTokens[i]
		if len(strTok.Content) == 0 {
			continue
		}
		mapTok, exists := tokenStr[strTok.Content]
		if macroMode && mapTok == TokenEnd && macroEndStack == 0 {
			t.Type = TokenMacroEnd
			macroTokenBuffer = append(macroTokenBuffer, t)
			globalMacroTable[currentMacroName] = macroTokenBuffer
			currentTokenBuffer = &mainTokenBuffer
			macroTokenBuffer = []Token{}
			macroMode = false
			continue
		}
		if macroMode {
			switch mapTok {
			case TokenFor, TokenIf:
				macroEndStack++
			case TokenEnd:
				macroEndStack--
			case TokenMacro:
				fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
				fmt.Println("macro definition inside of a macro is not supported")
				os.Exit(1)
			default:
			}
		}
		if exists && mapTok == TokenMacro {
			currentTokenBuffer = &macroTokenBuffer
			macroMode = true
			macroEndStack = 0
			currentMacroName = strTokens[i+1].Content
			i++
			continue
		}
		if exists && mapTok == TokenVar {
			var (
				varKind      TokenType
				varName      string
				varNameIndex uint64
			)
			if i+3 >= len(strTokens) {
				fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
				fmt.Println("expected variable definition")
				fmt.Println(
					"variable definition looks like this: \n",
					"  `var <var-name> <var-type> end`\n",
					"eg: \n",
					"  `var x int end`",
				)
				os.Exit(1)
			}
			{
				i++
				strTok = strTokens[i]
				varNameIndex = uint64(i)
				varName = strTokens[varNameIndex].Content
				if tmpT, e := tokenStr[strTok.Content]; e {
					fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
					fmt.Printf(
						"expected TokenWord found keyword %v\n keyword not allowed as variable names\n",
						intrinsicStr[tmpT],
					)
					fmt.Println(
						"variable definition looks like this: \n",
						"  `var <var-name> <var-type> end`\n",
						"eg: \n",
						"  `var x int end`",
					)
					os.Exit(1)
				}
			}
			{
				i++
				strTok = strTokens[i]
				var e bool
				varKind, e = tokenKindStr[strTok.Content]
				if !e {
					fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
					fmt.Printf(
						"expected type found %v\n",
						strTok.Content)
					fmt.Println(
						"variable definition looks like this: \n",
						"  `var <var-name> <var-type> end`\n",
						"eg: \n",
						"  `var x int end`",
					)
					os.Exit(1)
				}
			}
			{
				i++
				strTok = strTokens[i]
				if _, e := tokenStr[strTok.Content]; !e {
					fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
					fmt.Printf(
						"expected TokenEnd found %v\n",
						strTok.Content)
					fmt.Println(
						"variable definition looks like this: \n",
						"  `var <var-name> <var-type> end`\n",
						"eg: \n",
						"  `var x int end`",
					)
					os.Exit(1)
				}
			}
			// TODO(@siiick): swap Token.Kind and Token.Type as `kind` implies things like `end`, `for` etc, while `type` is `int`, `bool` etc
			t.Loc = strTok.Loc
			t.Kind = varKind   // var type
			t.Type = TokenWord // token type
			t.Operand = state.varBufSize
			globalVarsTable[varName] = t
			state.varBufSize += 8
			continue
		}

		if !exists {
			t.Loc = strTok.Loc
			if num, err := strconv.ParseUint(strTok.Content, 10, 64); err == nil {
				t.Type = TokenInt
				t.Operand = num
				*currentTokenBuffer = append(*currentTokenBuffer, t)
				continue
			}
			_, macroFound := globalMacroTable[strTok.Content]
			_, varFound := globalVarsTable[strTok.Content]
			if !macroFound && !varFound {
				fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
				fmt.Printf("Undefined Token `%v`\n", strTok.Content)
				os.Exit(1)
			}
			t.Type = TokenWord
			t.Operand = uint64(i)
			*currentTokenBuffer = append(*currentTokenBuffer, t)
		} else {
			t.Loc = strTok.Loc
			t.Type = mapTok
			*currentTokenBuffer = append(*currentTokenBuffer, t)
		}
	}

	return mainTokenBuffer
}

var tokenStr = map[string]TokenType{
	"+":        TokenPlus,
	"-":        TokenSub,
	"*":        TokenMult,
	"divmod":   TokenDivMod,
	"print":    TokenPrint,
	"swap":     TokenSwap,
	"dup":      TokenDup,
	"drop":     TokenDrop,
	">":        TokenGt,
	">=":       TokenGe,
	"<":        TokenLt,
	"<=":       TokenLe,
	"=":        TokenEq,
	"for":      TokenFor,
	"do":       TokenDo,
	"if":       TokenIf,
	"else":     TokenElse,
	"end":      TokenEnd,
	"macro":    TokenMacro,
	"syscall1": TokenSyscall1,
	"syscall3": TokenSyscall3,
	"rot":      TokenRot,
	"@":        TokenRead,
	"!":        TokenWrite,
	"var":      TokenVar,
	"true":     TokenTrue,
	"false":    TokenFalse,
}

func printTokens(ts []Token) {
	for _, v := range ts {
		typ, found := intrinsicStr[v.Type]
		if found {
			fmt.Printf("%v: %+v\n", typ, v)
		} else {
			fmt.Printf("%v: %+v\n", "[INVALID_TYPE]", v)
		}
	}
}

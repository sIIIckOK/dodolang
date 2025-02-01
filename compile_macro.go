package main

import (
	"fmt"
	"log"
	"os"
)

func compileMacro(f *os.File, strTokens []StringToken, tokens []Token, state *CompileState) {
	var blockStack []TokenType

	assert(TokenCount == 29, "Exhaustive switch case for CompileMacro")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case TokenInt:

			writeStr := compileTokenInt(token)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenPlus:
			writeStr := compileTokenPlus()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenSub:
			writeStr := compileTokenSub()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}

		case TokenMult:

			writeStr := compileTokenMul()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenDivMod:

			writeStr := compileTokenDivMod()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}

		case TokenPrint:

			writeStr := compileTokenPrint()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenSwap:
			writeStr := compileTokenSwap()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenDup:
			writeStr := compileTokenDup()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenDrop:
			writeStr := compileTokenDrop()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenRot:
			writeStr := compileTokenRot()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenGt:
			state.CmpCount++
			writeStr := compileTokenGt(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenGe:
			state.CmpCount++
			writeStr := compileTokenGe(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenLt:
			state.CmpCount++
			writeStr := compileTokenLt(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenLe:
			state.CmpCount++
			writeStr := compileTokenLe(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenEq:
			state.CmpCount++
			writeStr := compileTokenEq(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenFor:
			blockStack = append(blockStack, token.Type)
			writeStr := compileTokenFor(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenDo:
			writeStr := compileTokenDo(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenIf:
			blockStack = append(blockStack, token.Type)
			writeStr := compileTokenIf(state)

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenElse:
			writeStr := compileTokenElse(state)
			state.BranchCount++

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenEnd:
			blockType := blockStack[len(blockStack)-1]
			blockStack = blockStack[:len(blockStack)-1]
			writeStr := compileTokenEnd(state, blockType)
			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenSyscall1:
			writeStr := compileTokenSyscall1()
			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenSyscall3:
			writeStr := compileTokenSyscall3()

			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenVar:
			strTok := strTokens[token.Operand]
			varName := strTok.Content
			var varToken Token
			var found bool
			if varToken, found = globalVarsTable[varName]; !found {
				fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
				fmt.Println("expected variable definition")
				fmt.Println(
					"variable definition looks like this: \n",
					"  `var <var-name> <var-type> end`\n",
					"eg: \n",
					"  `var x int end`",
				)
			}
			writeStr := compileTokenVar(uintptr(varToken.Operand))
			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenRead:
			writeStr := compileTokenRead()
			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
		case TokenWrite:
			writeStr := compileTokenWrite()
			_, err := f.Write([]byte(writeStr))
			if err != nil {
				log.Fatalln(err)
			}
        case TokenWord:
            tokenName := strTokens[tokens[i].Operand].Content
            if _, found := globalMacroTable[tokenName]; found {
                fmt.Println("ERROR:", "macro definition inside a macro is not supported")
            } else if _, found := globalVarsTable[tokenName]; found {
                strTok := strTokens[token.Operand]
                varName := strTok.Content
                var varToken Token
                var found bool
                if varToken, found = globalVarsTable[varName]; !found {
                    fmt.Printf("%v:%v:%v ", strTok.Loc.FilePath, strTok.Loc.Line, strTok.Loc.Col)
                    fmt.Printf("Undefined TokenWord %v\n", varName)
                }
                writeStr := compileTokenVar(uintptr(varToken.Operand))
                _, err := f.Write([]byte(writeStr))
                if err != nil {
                    log.Fatalln(err)
                }
            } else {
                log.Fatalf("undefined `%v`\n", tokenName)
            }
		case TokenMacro:
			fmt.Println("ERROR:", "macro definition inside a macro is not supported")
			os.Exit(1)
		default:
			assert(false, "CompileMacro unreachable")
		}
	}
}

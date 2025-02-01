package main

import (
	"fmt"
	"log"
	"os"
)

func compileTokenInt(token Token) string {
	retStr := "; -- Int Push --\n" +
		fmt.Sprintf("mov rax, %v\n", token.Operand) +
		"push rax\n"

	return retStr
}

func compileTokenPlus() string {
	retStr := "; -- Plus --\n" +
		"pop rax\n" +
		"pop rbx\n" +
		"add rax, rbx\n" +
		"push rax\n"

	return retStr
}

func compileTokenSub() string {
	retStr := "; -- Sub --\n" +
		"pop rax\n" +
		"pop rbx\n" +
		"sub rbx, rax\n" +
		"push rbx\n"

	return retStr
}

func compileTokenMul() string {
	retStr := "; -- Mul --\n" +
		"pop rax\n" +
		"pop rbx\n" +
		"mul rbx\n" +
		"push rax\n"

	return retStr
}

func compileTokenDivMod() string {
	retStr := "; -- DivMod --\n" +
		"xor rdx, rdx\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"div rbx\n" +
		"push rax\n" +
		"push rdx\n" +
		""

	return retStr
}

func compileTokenPrint() string {
	retStr := "; -- Print --\n" +
		"pop rdi\n" +
		"call print\n"

	return retStr
}

func compileTokenSwap() string {
	retStr := "; -- Swap --\n" +
		"pop rax\n" +
		"pop rbx\n" +
		"push rax\n" +
		"push rbx\n"

	return retStr
}

func compileTokenDup() string {
	retStr := "; -- Dup --\n" +
		"pop rax\n" +
		"push rax\n" +
		"push rax\n"

	return retStr
}

func compileTokenDrop() string {
	retStr := "; -- Drop --\n" +
		"pop rax\n"

	return retStr
}

func compileTokenRot() string {
	retStr := "; -- Rot --\n" +
		"pop rcx\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"push rbx\n" +
		"push rcx\n" +
		"push rax\n" +
		""

	return retStr
}

func compileTokenGt(state *CompileState) string {
	retStr := "; -- Gt --\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"cmp rax, rbx\n" +
		fmt.Sprintf("jg gt1_%v\n", state.CmpCount) +
		"push 0\n" +
		fmt.Sprintf("jmp gt2_%v\n", state.CmpCount) +
		fmt.Sprintf("gt1_%v:\n", state.CmpCount) +
		"push 1\n" +
		fmt.Sprintf("gt2_%v:\n", state.CmpCount) +
		""

	return retStr
}

func compileTokenGe(state *CompileState) string {
	retStr := "; -- Ge --\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"cmp rax, rbx\n" +
		fmt.Sprintf("jge ge1_%v\n", state.CmpCount) +
		"push 0\n" +
		fmt.Sprintf("jmp ge2_%v\n", state.CmpCount) +
		fmt.Sprintf("ge1_%v:\n", state.CmpCount) +
		"push 1\n" +
		fmt.Sprintf("ge2_%v:\n", state.CmpCount) +
		""

	return retStr
}

func compileTokenLt(state *CompileState) string {
	retStr := "; -- Lt --\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"cmp rax, rbx\n" +
		fmt.Sprintf("jl lt1_%v\n", state.CmpCount) +
		"push 0\n" +
		fmt.Sprintf("jmp lt2_%v\n", state.CmpCount) +
		fmt.Sprintf("lt1_%v:\n", state.CmpCount) +
		"push 1\n" +
		fmt.Sprintf("lt2_%v:\n", state.CmpCount) +
		""

	return retStr
}

func compileTokenLe(state *CompileState) string {
	retStr := "; -- Le --\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"cmp rax, rbx\n" +
		fmt.Sprintf("jle le1_%v\n", state.CmpCount) +
		"push 0\n" +
		fmt.Sprintf("jmp le2_%v\n", state.CmpCount) +
		fmt.Sprintf("le1_%v:\n", state.CmpCount) +
		"push 1\n" +
		fmt.Sprintf("le2_%v:\n", state.CmpCount) +
		""

	return retStr
}

func compileTokenEq(state *CompileState) string {
	retStr := "; -- Eq --\n" +
		"pop rbx\n" +
		"pop rax\n" +
		"cmp rax, rbx\n" +
		fmt.Sprintf("je eq1_%v\n", state.CmpCount) +
		"push 0\n" +
		fmt.Sprintf("jmp eq2_%v\n", state.CmpCount) +
		fmt.Sprintf("eq1_%v:\n", state.CmpCount) +
		"push 1\n" +
		fmt.Sprintf("eq2_%v:\n", state.CmpCount) +
		""

	return retStr
}

func compileTokenFor(state *CompileState) string {
	state.ForNest++
	retStr := "; -- For --\n" +
		fmt.Sprintf("for_%v_%v:\n", state.ForNest, state.ForCount) +
		""
	return retStr
}

func compileTokenDo(state *CompileState) string {
	retStr := "; -- Do --\n" +
		"pop rax\n" +
		"cmp rax, 0\n" +
		fmt.Sprintf("je forend_%v_%v\n", state.ForNest, state.ForCount) +
		""
	return retStr
}

func compileTokenIf(state *CompileState) string {
	retStr := "; -- If --\n" +
		"pop rax\n" +
		"cmp rax, 0\n" +
		fmt.Sprintf("je ifjmp_%v_%v_%v\n", state.BranchCount, state.IfNest, state.IfCount) +
		""
	state.IfNest++
	return retStr
}

func compileTokenElse(state *CompileState) string {
	retStr := "; -- Else --\n" +
		fmt.Sprintf("jmp ifend_%v_%v\n", state.IfNest, state.IfCount) +
		fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount) +
		""
	return retStr
}

func compileTokenEnd(state *CompileState, blockType TokenType) string {
	var retStr string

	if blockType == TokenIf {
		retStr = "; -- IfEnd --\n" +
			fmt.Sprintf("ifend_%v_%v:\n", state.IfNest, state.IfCount) +
			fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount) +
			""
		state.IfNest--
		if state.IfNest == 0 {
			state.IfCount++
		}
		state.BranchCount = 0
	} else if blockType == TokenFor {
		retStr = "; -- ForEnd --\n" +
			fmt.Sprintf("jmp for_%v_%v\n", state.ForNest, state.ForCount) +
			fmt.Sprintf("forend_%v_%v:\n", state.ForNest, state.ForCount) +
			""
		state.ForNest--
		if state.ForNest == 0 {
			state.ForCount++
		}
	} else if blockType == TokenMacro {
		retStr = "; -- MacroEnd --\n"
	} else {
		assert(false, "unreachable")
	}
	return retStr
}

func compileTokenSyscall1() string {
	retStr := "; -- Syscall1 --\n" +
		"pop rax\n" +
		"pop rdi\n" +
		"syscall\n" +
		""
	return retStr
}

func compileTokenSyscall3() string {
	retStr := "; -- Syscall3 --\n" +
		"pop rax\n" +
		"pop rdi\n" +
		"pop rsi\n" +
		"pop rdx\n" +
		"syscall\n" +
		""
	return retStr
}

func compileTokenVar(offset uintptr) string {
	retStr := "; -- Var --\n" +
		fmt.Sprintf("mov rax, vars_buffer+%v\n", offset) +
		"push rax\n" +
		""
	return retStr
}

func compileTokenRead() string {
	retStr := "; -- Var Read --\n" +
		"pop rbx\n" +
		"mov rax, qword [rbx]\n" +
		"push rax\n" +
		""
	return retStr
}

func compileTokenWrite() string {
	retStr := "; -- Var Write --\n" +
		"pop rax\n" +
		"pop rbx\n" +
		"mov qword [rbx], rax\n"
	return retStr
}

func compileProgram(strTokens []StringToken, tokens []Token, state *CompileState, outPath string) {
	f, err := os.Create(outPath)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	header := `
    ; -- Header --
    BITS 64
    section .text
    
    print_render:
    xor rbx, rbx
    mov rcx, 10 
    .L1:
    xor rdx, rdx
    mov rax, rdi
    div rcx
    mov rdi, rax
    add rdx, '0'
    mov byte [print_buffer + rbx], dl
    cmp rax, 0
    je .exit
    inc rbx
    jmp .L1
    .exit:
    ret
    print_reverse:
    xor rdx, rdx
    mov rcx, 2
    mov rax, rbx
    div rcx
    inc rax
    mov rsi, rbx
    xor r9, r9
    .L1:
    mov cl, byte [print_buffer + rsi]
    mov dil, byte [print_buffer + r9]
    mov byte [print_buffer + rsi], dil
    mov byte [print_buffer + r9], cl
    dec rax
    cmp rax, 0
    je .exit
    inc r9
    dec rsi
    jmp .L1
    .exit:
    inc rbx
    mov byte [print_buffer + rbx], 10
    inc rbx
    ret
    print:
    call print_render
    call print_reverse
    mov rax, 1
    mov rdi, 0
    mov rsi, print_buffer
    mov rdx, rbx
    syscall
    ret

    global _start
    global vars_buffer
    _start: 
    `
	_, err = f.Write([]byte(header))
	if err != nil {
		log.Fatalln(err)
	}
	blockStack := make([]TokenType, 0, 0)

	assert(TokenCount == 29, "Exhaustive switch case for CompileProgram")
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
		case TokenMacro:
			blockStack = append(blockStack, TokenMacro)
			blockCount := 0
			tokenBuffer := []Token{}
			i++
			macroToken := strTokens[tokens[i].Operand]
			for {
				i++
				if tokens[i].Type == TokenFor || tokens[i].Type == TokenIf {
					blockCount++
				} else if tokens[i].Type == TokenEnd {
					if blockCount == 0 {
						break
					}
					blockCount--
				} else if tokens[i].Type == TokenMacro {
					assert(false, "macro definition inside macro definition is not allowed")
				}
				tokenBuffer = append(tokenBuffer, tokens[i])
			}
			globalMacroTable[macroToken.Content] = tokenBuffer
		case TokenVar:
			assert(false, "TokenVar unreachable")
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
			curTok := strTokens[tokens[i].Operand]
			tokenName := curTok.Content
			macroToks, macroFound := globalMacroTable[tokenName]
			varTok, varFound := globalVarsTable[tokenName]
			if macroFound && varFound {
				fmt.Printf("%v:%v:%v ", curTok.Loc.FilePath, curTok.Loc.Line, curTok.Loc.Col)
				fmt.Printf(
					"[COMPILER_ISSUE] Provided TokenWord `%v` is both macro and var, this error should have been caught in the parser\n",
					tokenName,
				)
			}
			if macroFound {
				compileMacro(f, strTokens, macroToks, state)
			} else if varFound {
				writeStr := compileTokenVar(uintptr(varTok.Operand))
				_, err := f.Write([]byte(writeStr))
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				fmt.Printf("%v:%v:%v ", curTok.Loc.FilePath, curTok.Loc.Line, curTok.Loc.Col)
				fmt.Printf("Undefined TokenWord %v\n", tokenName)
			}
		default:
			assert(false, "CompileProgram unreachable")
		}
	}

	footer := "; -- Footer --\n" +
		"mov rax, 60\n" +
		"mov rdi, 0\n" +
		"syscall\n" +
		"section .bss\n" +
		"print_buffer: resb 22\n" +
		fmt.Sprintf("vars_buffer: resb %v\n", state.varBufSize)

	_, err = f.Write([]byte(footer))
	if err != nil {
		log.Fatalln(err)
	}
}

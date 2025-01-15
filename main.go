package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func main() {
    flag.Parse()

    flag.Usage = func() {
        fmt.Printf("Usage of %s:\n", os.Args[0])
        fmt.Println("    build <source file>")
        fmt.Println("        build: compile file")
        flag.PrintDefaults()
    }
    subCom := flag.Arg(0)
    if subCom != "build" && subCom != "run" {
        flag.Usage()
        return
    }
    filePath := flag.Arg(1)
    ext := filepath.Ext(filePath) 
    if ext != ".dodo" {
        fmt.Printf("unknown extension `%v` only valid extension is `.dodo`\n", ext)
    }
    args := flag.Args()
    if len(args) < 2 {
        flag.Usage()
        return
    }

    contentBytes, err := os.ReadFile(filePath)
    if err != nil {
        log.Fatalln(err)
    }
    content := string(contentBytes)+string(rune(0))
    strTokens := LexFile(content, filePath)
    tokens := ParseTokens(strTokens)

    outPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
    if subCom == "build" {
        CompileProgram(strTokens, tokens)
        cmd := []string{"nasm", "-g", "-felf64", outPath+".asm"}
        if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
            log.Fatalln("ERROR:", err, cmd, string(out))
        }
        cmd = []string{"ld", outPath+".o", "-o", outPath}
        if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
            log.Fatalln("ERROR:", err, cmd, string(out))
        }
    } else {
        fmt.Printf("Invalid subcommand `%v`\n", subCom)
        flag.Usage()
    }
}

type Location struct {
    Line     uint
    Col      uint
    FilePath string
}

type StringToken struct {
    Content string
    Loc     Location
}

func LexFile(content string, filePath string) []StringToken {
    var tokens []StringToken
    var t StringToken
    t.Loc.FilePath = filePath
    lines := strings.Split(content, "\n")
    for lineNo, _line := range lines {
        if len(_line) == 0 { continue }
        t.Loc.Line = uint(lineNo)
        commentedOut := strings.Split(_line, "//")[0]
        line := strings.FieldsFunc(commentedOut, unicode.IsSpace)
        for _, tok := range line {
            t.Content = tok
            tokens = append(tokens, t)
        }
    }
    return tokens
}


type TokenType uint
const (
    TokenInt  = iota
    TokenPlus
    TokenSub
    TokenMult
    TokenDivMod
    TokenWord
    TokenPrint
    TokenSwap
    TokenDup
    TokenDrop
    TokenMacro
    TokenEq
    TokenGt
    TokenLt
    TokenLe
    TokenGe
    TokenFor
    TokenDo
    TokenIf
    TokenElse
    TokenEnd
    TokenSyscall1
    TokenSyscall3
    TokenCount
)

type Token struct {
    Type    TokenType
    Loc     Location
    Operand uint64
    NestLvl uint64
}

const stackSize = 1024
type StackUint64 struct {
    Items [stackSize]uint64
    Ptr int
}

func (s *StackUint64) Push(item uint64, loc Location) {
    assert(s.Ptr + 1 < len(s.Items), fmt.Sprintf("%v:%v:%v: stack overflow\n", 
        loc.FilePath, loc.Line, loc.Col))

    s.Items[s.Ptr] = item
    s.Ptr++
}

func (s *StackUint64) Pop(loc Location) uint64 {
    assert(s.Ptr - 1 >= 0, fmt.Sprintf("%v:%v:%v: stack underflow", 
        loc.FilePath, loc.Line, loc.Col))
    item := s.Items[s.Ptr-1]
    s.Ptr--
    return item
}

func (s StackUint64) Len() int {
    return s.Ptr
}
func InterpretProgram(strTokens []StringToken, tokens []Token) {
    stack := StackUint64{}
    assert(TokenCount == 6, "Exhaustive switch case for InterpretProgram")
    for _, token := range tokens {
        loc := token.Loc
        switch token.Type {
        case TokenInt:
            stack.Push(token.Operand, loc)
        case TokenPlus:
            if stack.Len() < 2 {
                log.Fatalf("%v:%v:%v: `+` instrics requires atleast 2 elements on the stack found %v\n", 
                    loc.FilePath, loc.Line, loc.Col, stack.Len())
            } 
            a := stack.Pop(loc) 
            b := stack.Pop(loc) 
            stack.Push(a + b, loc)
        case TokenSub:
            if stack.Len() < 2 {
                log.Fatalf("%v:%v:%v: `-` instrics requires atleast 2 elements on the stack found %v\n", 
                    loc.FilePath, loc.Line, loc.Col, stack.Len())
            } 
            a := stack.Pop(loc) 
            b := stack.Pop(loc) 
            stack.Push(b - a, loc)
        case TokenMult:
            if stack.Len() < 2 {
                log.Fatalf("%v:%v:%v: `*` instrics requires atleast 2 elements on the stack found %v\n", 
                    loc.FilePath, loc.Line, loc.Col, stack.Len())
            } 
            a := stack.Pop(loc)
            b := stack.Pop(loc) 
            stack.Push(a * b, loc)
        case TokenPrint:
            if stack.Len() < 1 {
                log.Fatalf("%v:%v:%v: `print` instrics requires atleast 1 elements on the stack found %v\n",
                    loc.FilePath, loc.Line, loc.Col, stack.Len())
            }
            a := stack.Pop(loc)
            fmt.Println(a)
        }
    }
}

type CompileState struct {
    CmpCount    uint64
    IfCount     uint64
    ForCount    uint64
    ForNest     uint64
    IfNest      uint64
    BranchCount uint64
}

func ParseTokens(strTokens []StringToken) []Token {
    var tokens []Token
    var t Token 
    assert(TokenCount == 23, "Exhaustive switch case for ParseToken")
    for i, tok := range strTokens {
        if rune(tok.Content[0]) == rune(0) {
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

func CompileMacro(f *os.File, strTokens []StringToken, tokens []Token, state *CompileState) {
    var blockStack []TokenType

    assert(TokenCount == 23, "Exhaustive switch case for CompileProgram")
    for i := 0; i < len(tokens); i++ {
        token := tokens[i]
        switch token.Type {
        case TokenInt:

            writeStr := 
            "; -- Int Push --\n"                            +
            fmt.Sprintf("mov rax, %v\n", token.Operand)     +
            "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenPlus:
            
            writeStr := 
                "; -- Plus --\n"    +
                "pop rax\n"         +
                "pop rbx\n"         +
                "add rax, rbx\n"    +
                "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSub:

            writeStr := 
                "; -- Sub --\n"     +
                "pop rax\n"         +
                "pop rbx\n"         +
                "sub rbx, rax\n"    +
                "push rbx\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }

        case TokenMult:

            writeStr := 
                "; -- Mul --\n"     +
                "pop rax\n"         +
                "pop rbx\n"         +
                "mul rbx\n"         +
                "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDivMod:

            writeStr := 
                "; -- DivMod --\n"      +
                "xor rdx, rdx\n"        +
                "pop rax\n"             +
                "pop rbx\n"             +
                "div rbx\n"             +
                "push rdx\n"            +
                "push rax\n"            +
                ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }

        case TokenPrint:
            
            writeStr := 
            "; -- Print --\n" +
            "pop rdi\n"       +
            "call print\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSwap:
            writeStr := 
            "; -- Swap --\n"    +
            "pop rax\n"         +
            "pop rbx\n"         +
            "push rax\n"        +
            "push rbx\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDup:
            writeStr := 
            "; -- Dup --\n"    +
            "pop rax\n"        +
            "push rax\n"        +
            "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDrop:
            writeStr := 
            "; -- Drop --\n"    +
            "pop rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenGt:
            state.CmpCount++
            writeStr := 
            "; -- Gt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jg gt1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp gt2_%v\n", state.CmpCount)   +
            fmt.Sprintf("gt1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("gt2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenGe:
            state.CmpCount++
            writeStr := 
            "; -- Ge --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jge ge1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp ge2_%v\n", state.CmpCount)   +
            fmt.Sprintf("ge1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("ge2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLt:
            state.CmpCount++
            writeStr := 
            "; -- Lt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jl lt1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp lt2_%v\n", state.CmpCount)   +
            fmt.Sprintf("lt1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("lt2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLe:
            state.CmpCount++
            writeStr := 
            "; -- Le --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jle le1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp le2_%v\n", state.CmpCount)   +
            fmt.Sprintf("le1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("le2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenEq:
            state.CmpCount++
            writeStr := 
            "; -- Eq --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("je eq1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp eq2_%v\n", state.CmpCount)  +
            fmt.Sprintf("eq1_%v:\n", state.CmpCount)     +
            "push 1\n"          +
            fmt.Sprintf("eq2_%v:\n", state.CmpCount)     +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenFor:
            blockStack = append(blockStack, token.Type)
            writeStr := 
            "; -- For --\n"                         +
            fmt.Sprintf("for_%v_%v:\n", state.ForNest, state.ForCount) +
            ""
            state.ForNest++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDo:
            writeStr := 
            "; -- Do --\n"                               +
            "pop rax\n"                                  + 
            "cmp rax, 0\n"                               + 
            fmt.Sprintf("je forend_%v_%v\n", state.ForNest-1, state.ForCount) +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenIf:
            blockStack = append(blockStack, token.Type)
            writeStr := 
            "; -- If --\n"                          +
            "pop rax\n"                             +
            "cmp rax, 0\n"                          +
            fmt.Sprintf("je ifjmp_%v_%v_%v\n", state.BranchCount, state.IfNest, state.IfCount)   +
            ""
            state.IfNest++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenElse:
            writeStr := 
            "; -- Else --\n"                               +
            fmt.Sprintf("jmp ifend_%v_%v\n", state.IfNest, state.IfCount)   +
            fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount)   +
            ""
            state.BranchCount++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenEnd:
            blockType := blockStack[len(blockStack)-1]
            blockStack = blockStack[:len(blockStack)-1]
            var writeStr string

            if blockType == TokenIf {
                writeStr = 
                "; -- IfEnd --\n"    +
                fmt.Sprintf("ifend_%v_%v:\n", state.IfNest, state.IfCount)   +
                fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount)   +
                ""
                state.IfNest--
                if state.IfNest == 0 {
                    state.IfCount++
                }
                state.BranchCount = 0
            } else if blockType == TokenFor {
                state.ForNest--
                writeStr = 
                "; -- ForEnd --\n"    +
                fmt.Sprintf("jmp for_%v_%v\n", state.ForNest, state.ForCount)   +
                fmt.Sprintf("forend_%v_%v:\n", state.ForNest, state.ForCount)   +
                ""
                if state.ForNest == 0 {
                    state.ForCount++
                }
            } else if blockType == TokenMacro {
                writeStr = "; -- MacroEnd --\n"
            } else {
                assert(false, "unreachable")
            }
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSyscall1:
            writeStr := 
            "; -- Syscall1 --\n"                               +
            "pop rax\n"                             +
            "pop rdi\n"                             +
            "syscall\n"                             +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSyscall3:
            writeStr := 
            "; -- Syscall3 --\n"                               +
            "pop rax\n"                             +
            "pop rdi\n"                             +
            "pop rdx\n"                             +
            "syscall\n"                             +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenMacro:
            fmt.Println("ERROR:", "macro definition inside a macro is not supported")
            os.Exit(1)
        default:
            assert(false, "CompileMacro unreachable")
        }
    }
}

func CompileProgram(strTokens []StringToken, tokens []Token) {
    if len(tokens) == 0 { 
        log.Fatalln("Empty file not allowed")
    }

    srcPath := tokens[0].Loc.FilePath
    outPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".asm"
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
    _start: 
    `
    _, err = f.Write([]byte(header))
    if err != nil {
        log.Fatalln(err)
    }
    var blockStack []TokenType
    state := CompileState{}

    globalTable := make(map[string][]Token, 100)

    assert(TokenCount == 23, "Exhaustive switch case for CompileProgram")
    for i := 0; i < len(tokens); i++ {
        token := tokens[i]
        switch token.Type {
        case TokenInt:

            writeStr := 
            "; -- Int Push --\n"                            +
            fmt.Sprintf("mov rax, %v\n", token.Operand)     +
            "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenPlus:
            
            writeStr := 
                "; -- Plus --\n"    +
                "pop rax\n"         +
                "pop rbx\n"         +
                "add rax, rbx\n"    +
                "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSub:

            writeStr := 
                "; -- Sub --\n"     +
                "pop rax\n"         +
                "pop rbx\n"         +
                "sub rbx, rax\n"    +
                "push rbx\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }

        case TokenMult:

            writeStr := 
                "; -- Mul --\n"     +
                "pop rax\n"         +
                "pop rbx\n"         +
                "mul rbx\n"         +
                "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDivMod:

            writeStr := 
                "; -- DivMod --\n"      +
                "xor rdx, rdx\n"        +
                "pop rax\n"             +
                "pop rbx\n"             +
                "div rbx\n"             +
                "push rdx\n"            +
                "push rax\n"            +
                ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }

        case TokenPrint:
            
            writeStr := 
            "; -- Print --\n" +
            "pop rdi\n"       +
            "call print\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSwap:
            writeStr := 
            "; -- Swap --\n"    +
            "pop rax\n"         +
            "pop rbx\n"         +
            "push rax\n"        +
            "push rbx\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDup:
            writeStr := 
            "; -- Dup --\n"    +
            "pop rax\n"        +
            "push rax\n"        +
            "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDrop:
            writeStr := 
            "; -- Drop --\n"    +
            "pop rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenGt:
            state.CmpCount++
            writeStr := 
            "; -- Gt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jg gt1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp gt2_%v\n", state.CmpCount)   +
            fmt.Sprintf("gt1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("gt2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenGe:
            state.CmpCount++
            writeStr := 
            "; -- Ge --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jge ge1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp ge2_%v\n", state.CmpCount)   +
            fmt.Sprintf("ge1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("ge2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLt:
            state.CmpCount++
            writeStr := 
            "; -- Lt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jl lt1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp lt2_%v\n", state.CmpCount)   +
            fmt.Sprintf("lt1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("lt2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLe:
            state.CmpCount++
            writeStr := 
            "; -- Le --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jle le1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp le2_%v\n", state.CmpCount)   +
            fmt.Sprintf("le1_%v:\n", state.CmpCount)   +
            "push 1\n"          +
            fmt.Sprintf("le2_%v:\n", state.CmpCount)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenEq:
            state.CmpCount++
            writeStr := 
            "; -- Eq --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("je eq1_%v\n", state.CmpCount)   +
            "push 0\n"          +
            fmt.Sprintf("jmp eq2_%v\n", state.CmpCount)  +
            fmt.Sprintf("eq1_%v:\n", state.CmpCount)     +
            "push 1\n"          +
            fmt.Sprintf("eq2_%v:\n", state.CmpCount)     +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenFor:
            blockStack = append(blockStack, token.Type)
            writeStr := 
            "; -- For --\n"                         +
            fmt.Sprintf("for_%v_%v:\n", state.ForNest, state.ForCount) +
            ""
            state.ForNest++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDo:
            writeStr := 
            "; -- Do --\n"                               +
            "pop rax\n"                                  + 
            "cmp rax, 0\n"                               + 
            fmt.Sprintf("je forend_%v_%v\n", state.ForNest-1, state.ForCount) +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenIf:
            blockStack = append(blockStack, token.Type)
            writeStr := 
            "; -- If --\n"                          +
            "pop rax\n"                             +
            "cmp rax, 0\n"                          +
            fmt.Sprintf("je ifjmp_%v_%v_%v\n", state.BranchCount, state.IfNest, state.IfCount)   +
            ""
            state.IfNest++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenElse:
            writeStr := 
            "; -- Else --\n"                               +
            fmt.Sprintf("jmp ifend_%v_%v\n", state.IfNest, state.IfCount)   +
            fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount)   +
            ""
            state.BranchCount++

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenEnd:
            blockType := blockStack[len(blockStack)-1]
            blockStack = blockStack[:len(blockStack)-1]
            var writeStr string

            if blockType == TokenIf {
                writeStr = 
                "; -- IfEnd --\n"    +
                fmt.Sprintf("ifend_%v_%v:\n", state.IfNest, state.IfCount)   +
                fmt.Sprintf("ifjmp_%v_%v_%v:\n", state.BranchCount, state.IfNest-1, state.IfCount)   +
                ""
                state.IfNest--
                if state.IfNest == 0 {
                    state.IfCount++
                }
                state.BranchCount = 0
            } else if blockType == TokenFor {
                state.ForNest--
                writeStr = 
                "; -- ForEnd --\n"    +
                fmt.Sprintf("jmp for_%v_%v\n", state.ForNest, state.ForCount)   +
                fmt.Sprintf("forend_%v_%v:\n", state.ForNest, state.ForCount)   +
                ""
                if state.ForNest == 0 {
                    state.ForCount++
                }
            } else if blockType == TokenMacro {
                writeStr = "; -- MacroEnd --\n"
            } else {
                assert(false, "unreachable")
            }
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSyscall1:
            writeStr := 
            "; -- Syscall1 --\n"                               +
            "pop rax\n"                             +
            "pop rdi\n"                             +
            "syscall\n"                             +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSyscall3:
            writeStr := 
            "; -- Syscall3 --\n"                               +
            "pop rax\n"                             +
            "pop rdi\n"                             +
            "pop rdx\n"                             +
            "syscall\n"                             +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenMacro:
            blockStack = append(blockStack, TokenMacro)
            blockCount := 0
            tokenBuffer := []Token{}
            i++; macroToken := strTokens[tokens[i].Operand]
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
            globalTable[macroToken.Content] = tokenBuffer
        case TokenWord:
            tokenName := strTokens[tokens[i].Operand].Content
            if rune(tokenName[0]) == rune(0) { continue }
            if toks, found := globalTable[tokenName]; found {
                CompileMacro(f, strTokens, toks, &state)
            } else {
                log.Fatalf("undefined `%v`\n", tokenName)
            }
        default:
            assert(false, "CompileProgram unreachable")
        }
    }

    footer := 
    "; -- Footer --\n"          +
    "mov rax, 60\n"             +
    "mov rdi, 0\n"              +
    "syscall\n"                 +
    "section .bss\n"            +
    "print_buffer: resb 22\n"

    _, err = f.Write([]byte(footer))
    if err != nil {
        log.Fatalln(err)
    }
}

func assert(cond bool, msg string) {
    if !cond {
        log.Println(msg)
        panic(1)
    }
}




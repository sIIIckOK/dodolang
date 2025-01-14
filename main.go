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

func ParseTokens(strTokens []StringToken) []Token {
    var tokens []Token
    var t Token 
    var blockCount uint64
    var branchNo uint64
    var cmpCount uint64
    assert(TokenCount == 23, "Exhaustive switch case for ParseToken")
    for i, tok := range strTokens {
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
            cmpCount++
            t.Type = TokenGt
            t.Operand = cmpCount
            tokens = append(tokens, t)
        case ">=":
            cmpCount++
            t.Type = TokenGe
            t.Operand = cmpCount
            tokens = append(tokens, t)
        case "<":
            cmpCount++
            t.Type = TokenLt
            t.Operand = cmpCount
            tokens = append(tokens, t)
        case "<=":
            cmpCount++
            t.Type = TokenLe
            t.Operand = cmpCount
            tokens = append(tokens, t)
        case "=":
            cmpCount++
            t.Type = TokenEq
            t.Operand = cmpCount
            tokens = append(tokens, t)
        case "for":
            blockCount++
            t.Type = TokenFor
            t.Operand = blockCount
            tokens = append(tokens, t)
        case "do":
            t.Type = TokenDo
            t.Operand = blockCount
            tokens = append(tokens, t)
        case "if":
            blockCount++
            t.Type = TokenIf
            t.Operand = blockCount
            tokens = append(tokens, t)
        case "else":
            t.NestLvl = branchNo
            branchNo++
            t.Type = TokenElse
            t.Operand = blockCount
            tokens = append(tokens, t)
        case "end":
            t.Type = TokenEnd
            t.NestLvl = branchNo
            t.Operand = blockCount 
            tokens = append(tokens, t)
            blockCount--
            branchNo--
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

const stackSize = 1024
type StackUint64 struct {
    Items [stackSize]uint64
    Ptr int
}

func (s *StackUint64) Push(item uint64, loc Location) {
    assert(s.Ptr + 1 < len(s.Items), fmt.Sprintf("%v:%v:%v: stack overflow", 
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
    assert(TokenCount == 23, "Exhaustive switch case for CompileProgram")
    for _, token := range tokens {
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
            writeStr := 
            "; -- Gt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jg gt1_%v\n", token.Operand)   +
            "push 0\n"          +
            fmt.Sprintf("jmp gt2_%v\n", token.Operand)   +
            fmt.Sprintf("gt1_%v:\n", token.Operand)   +
            "push 1\n"          +
            fmt.Sprintf("gt2_%v:\n", token.Operand)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenGe:
            writeStr := 
            "; -- Ge --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jge ge1_%v\n", token.Operand)   +
            "push 0\n"          +
            fmt.Sprintf("jmp ge2_%v\n", token.Operand)   +
            fmt.Sprintf("ge1_%v:\n", token.Operand)   +
            "push 1\n"          +
            fmt.Sprintf("ge2_%v:\n", token.Operand)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLt:
            writeStr := 
            "; -- Lt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jl lt1_%v\n", token.Operand)   +
            "push 0\n"          +
            fmt.Sprintf("jmp lt2_%v\n", token.Operand)   +
            fmt.Sprintf("lt1_%v:\n", token.Operand)   +
            "push 1\n"          +
            fmt.Sprintf("lt2_%v:\n", token.Operand)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenLe:
            writeStr := 
            "; -- Le --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("jle le1_%v\n", token.Operand)   +
            "push 0\n"          +
            fmt.Sprintf("jmp le2_%v\n", token.Operand)   +
            fmt.Sprintf("le1_%v:\n", token.Operand)   +
            "push 1\n"          +
            fmt.Sprintf("le2_%v:\n", token.Operand)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenEq:
            writeStr := 
            "; -- Gt --\n"      +
            "pop rbx\n"         +
            "pop rax\n"         +
            "cmp rax, rbx\n"    +
            fmt.Sprintf("je eq1_%v\n", token.Operand)   +
            "push 0\n"          +
            fmt.Sprintf("jmp eq2_%v\n", token.Operand)  +
            fmt.Sprintf("eq1_%v:\n", token.Operand)     +
            "push 1\n"          +
            fmt.Sprintf("eq2_%v:\n", token.Operand)     +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenFor:
            blockStack = append(blockStack, token.Type)
            writeStr := 
            "; -- For --\n"                         +
            fmt.Sprintf("for_%v:\n", token.Operand) +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenDo:
            writeStr := 
            "; -- Do --\n"                               +
            "pop rax\n"                                  + 
            "cmp rax, 0\n"                               + 
            fmt.Sprintf("je forend_%v\n", token.Operand) +
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
            fmt.Sprintf("je ifjmp_%v_%v\n", token.NestLvl, token.Operand)   +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenElse:
            writeStr := 
            "; -- Else --\n"                               +
            fmt.Sprintf("jmp ifend_%v\n", token.Operand)   +
            fmt.Sprintf("ifjmp_%v_%v:\n", token.NestLvl, token.Operand)   +
            ""

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
                "; -- End --\n"    +
                fmt.Sprintf("ifjmp_%v_%v:\n", token.NestLvl, token.Operand)   +
                fmt.Sprintf("ifend_%v:\n", token.Operand)   +
                ""
            } else if blockType == TokenFor {
                writeStr = 
                "; -- End --\n"    +
                fmt.Sprintf("jmp for_%v\n", token.Operand)   +
                fmt.Sprintf("forend_%v:\n", token.Operand)   +
                ""
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
            "pop rdi\n"                             +
            "pop rax\n"                             +
            "syscall\n"                             +
            ""
            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenSyscall3:
            writeStr := 
            "; -- Syscall3 --\n"                               +
            "pop rdx\n"                             +
            "pop rdi\n"                             +
            "pop rax\n"                             +
            "syscall\n"                             +
            ""

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }
        case TokenWord:
            // assert(false, "compilation of TokenWord not implemented yet")
        default:
            assert(false, "unreachable")
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




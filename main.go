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
        fmt.Println("    [build|run] <source file>")
        fmt.Println("        build: compile")
        fmt.Println("        run:   interpret")
        flag.PrintDefaults()
    }

    subCom := flag.Arg(0)
    if subCom != "build" && subCom != "run" {
        flag.Usage()
        return
    }
    filePath := flag.Arg(1)
    args := flag.Args()
    if len(args) < 2 {
        flag.Usage()
        return
    }

    contentBytes, err := os.ReadFile(filePath)
    if err != nil {
        log.Fatalln(err)
    }
    content := string(contentBytes)
    strTokens := LexFile(content, filePath)
    tokens := ParseTokens(strTokens)

    outPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
    if subCom == "build" {
        CompileProgram(strTokens, tokens)
        cmd := []string{"nasm", "-felf64", outPath+".asm"}
        if _, err := exec.Command(cmd[0], cmd[1:]...).Output(); err != nil {
            log.Fatalln("ERROR:", err, cmd)
        }
        cmd = []string{"ld", outPath+".o", "-o", outPath}
        if _, err := exec.Command(cmd[0], cmd[1:]...).Output(); err != nil {
            log.Fatalln("ERROR:", err, cmd)
        }
    } else if subCom == "run" {
        InterpretProgram(strTokens, tokens)
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
    tokens := []StringToken{}
    t := StringToken{}
    t.Loc.FilePath = filePath
    lines := strings.Split(content, "\n")
    for lineNo, _line := range lines {
        line := SkipWhitespace(_line)
        if len(line) == 0 { continue }
        t.Loc.Line = uint(lineNo)
        var start uint
        for i, v := range line {
            if i == len(line) - 1 {
                t.Loc.Col = start 
                t.Content = line[start:i+1]
                tokens = append(tokens, t)
                start = uint(i) + 1
            }
            if v == ' ' {
                t.Loc.Col = start 
                t.Content = line[start:i]
                tokens = append(tokens, t)
                start = uint(i) + 1
            }
        }
    }
    return tokens
}

func SkipWhitespace(str string) string {
    length := len(str)
    if length != 0 {
        start := 0
        i := 0
        for unicode.IsSpace(rune(str[i])) && i - length - 1 >= 0 {
            i++; start++
        }
        return str[start:]
    }
    return ""
}

type TokenType uint
const (
    TokenInt  = iota
    TokenPlus
    TokenSub
    TokenMult
    TokenWord
    TokenPrint

    TokenCount
)

type Token struct {
    Type    TokenType
    Loc     Location
    Operand uint64
}

func ParseTokens(strTokens []StringToken) []Token {
    var tokens []Token
    var t Token 
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
        case "print":
            t.Type = TokenPrint
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
    if len(tokens) == 0 { return }

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
    .L1:
    mov rax, rdi
    xor rdx, rdx
    mov rcx, 10
    div rcx
    mov rdi, rax
    add rdx, '0' 
    mov byte [print_buffer + rbx], dl
    cmp rax, 0
    je .L2
    inc rbx
    cmp rbx, 21
    jne .L1
    .L2:
    ret
    print_reverse:
    mov rax, rbx
    inc rbx
    mov byte [print_buffer + rbx], 10
    inc rbx
    mov byte [print_buffer + rbx], 0
    dec rbx
    xor rdx, rdx
    mov rcx, 2
    div rcx
    inc rax
    xor rcx, rcx
    .L1:
    mov sil, byte [print_buffer + rcx]
    dec rbx
    mov dil, byte [print_buffer + rbx]
    mov byte [print_buffer + rbx], sil
    mov byte [print_buffer + rcx], dil
    inc rcx
    cmp rcx, rax
    jne .L1
    ret
    print:
    call print_render
    call print_reverse
    mov rax, 1
    mov rdi, 0
    mov rsi, print_buffer
    mov rdx, 22
    syscall
    ret

    global _start
    _start: `
    _, err = f.Write([]byte(header))
    if err != nil {
        log.Fatalln(err)
    }
    assert(TokenCount == 6, "Exhaustive switch case for CompileProgram")
    for _, token := range tokens {
        switch token.Type {
        case TokenInt:

            writeStr := 
            "; -- int --\n"                             +
            fmt.Sprintf("mov rax, %v\n", token.Operand) +
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
                "; -- Sub --\n"    +
                "pop rax\n"         +
                "pop rbx\n"         +
                "sub rbx, rax\n"    +
                "push rax\n"

            _, err := f.Write([]byte(writeStr))
            if err != nil {
                log.Fatalln(err)
            }

        case TokenMult:

            writeStr := 
                "; -- Plus --\n"    +
                "pop rax\n"         +
                "pop rbx\n"         +
                "mul rbx\n"    +
                "push rax\n"

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
        }
    }

    footer := 
    "; -- Footer --\n"          +
    "mov rax, 60\n"             +
    "mov rdi, 0\n"              +
    "syscall\n"                 +
    "section .bss\n"            +
    "print_buffer: resb 21\n"

    _, err = f.Write([]byte(footer))
    if err != nil {
        log.Fatalln(err)
    }
}

func assert(cond bool, msg string) {
    if !cond {
        log.Panicln(msg)
    }
}


package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var globalVarsTable  = make(map[string]uintptr, 100)
var globalMacroTable = make(map[string][]Token, 100)


type Location struct {
    Line     uint
    Col      uint
    FilePath string
}

type StringToken struct {
    Content string
    Loc     Location
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
    TokenRot
    TokenVar
    TokenRead
    TokenWrite

    TokenCount
)

type Token struct {
    Type    TokenType
    Loc     Location
    Operand uint64
    NestLvl uint64
}

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
    strTokens := lexFile(content, filePath)
    tokens := ParseTokens(strTokens)
    typeCheck(strTokens, tokens)

    outPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
    if subCom == "build" {
        compileProgram(strTokens, tokens)
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

type CompileState struct {
    CmpCount    uint64
    IfCount     uint64
    ForCount    uint64
    ForNest     uint64
    IfNest      uint64
    BranchCount uint64
    varOffset   uint64
}

func assert(cond bool, msg string) {
    if !cond {
        log.Println(msg)
        panic(1)
    }
}




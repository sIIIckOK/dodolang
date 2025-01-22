package main

import (
	"fmt"
	"os"
)

type TypeStack []TokenType

func (ts *TypeStack) push(t TokenType) {
    *ts = append(*ts, t)
}

func (ts *TypeStack) pop(loc Location) TokenType {
    if len(*ts) == 0 {
        fmt.Printf("%v:%v:%v: Stack underflow \n",
            loc.FilePath, loc.Line, loc.Col)
    }
    t := (*ts)[len(*ts)-1]
    *ts = append(*ts, (*ts)[:len(*ts)-1]...)
    return t 
}

func (ts TypeStack) len() uint {
    return uint(len(ts))
}


func typeCheck(strTokens []StringToken, tokens []Token) {
    _ = strTokens
    stack := TypeStack{}
    assert(TokenCount == 27, "Exhaustive switch case for typeCheck")
    for _, token := range tokens {
        loc := token.Loc
        switch token.Type {
        case TokenInt:
            stack.push(TokenInt)
        case TokenPlus:
            if stack.len() < 2 {
                printCompilerErrorInstrinsic(token, "expected atleast 2 ints")
            } 
            a := stack.pop(loc)
            if a != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            b := stack.pop(loc) 
            if b != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            stack.push(TokenInt)
        case TokenSub:
            if stack.len() < 2 {
                printCompilerErrorInstrinsic(token, "expected atleast 2 ints")
            } 
            a := stack.pop(loc) 
            if a != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            b := stack.pop(loc)
            if b != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            stack.push(TokenInt)
        case TokenMult:
            if stack.len() < 2 {
                printCompilerErrorInstrinsic(token, "expected atleast 2 ints")
            } 
            a := stack.pop(loc)
            if a != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            b := stack.pop(loc) 
            if b != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            stack.push(TokenInt)
        case TokenPrint:
            if stack.len() < 1 {
                printCompilerErrorInstrinsic(token, "expected atleast 1 int")
            }
            stack.pop(loc)
        case TokenDivMod:
            if stack.len() < 2 {
                printCompilerErrorInstrinsic(token, "expected atleast 2 ints")
            } 
            a := stack.pop(loc)
            if a != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            b := stack.pop(loc) 
            if b != TokenInt {
                printCompilerErrorInstrinsic(token, "takes 2 ints")
            }
            stack.push(TokenInt)
        case TokenSwap:
        case TokenDup:
        case TokenDrop:
        case TokenRot:
        case TokenGt:
        case TokenGe:
        case TokenLt:
        case TokenLe:
        case TokenEq:
        case TokenFor:
        case TokenDo:
        case TokenIf:
        case TokenElse:
        case TokenEnd:
        case TokenSyscall1:
        case TokenSyscall3:
        case TokenMacro:
        case TokenVar:
        case TokenRead:
        case TokenWrite:
        case TokenWord:
        }
    }
}

func printCompilerErrorInstrinsic(token Token, err string) {
    fmt.Printf("%v:%v:%v: `%v` %v\n",
        token.Loc.FilePath, token.Loc.Line, token.Loc.Col, intrinsicStr[token.Type], err)
    os.Exit(1)
}

var intrinsicStr = map[TokenType]string {
    TokenInt: "TokenInt",
    TokenPlus: "TokenPlus",
    TokenSub: "TokenSub",
    TokenMult: "TokenMult",
    TokenPrint: "TokenPrint",
    TokenSwap: "TokenSwap",
    TokenDup: "TokenDup",
    TokenDrop: "TokenDrop",
    TokenRot: "TokenRot",
    TokenGt: "TokenGt",
    TokenGe: "TokenGe",
    TokenLt: "TokenLt",
    TokenLe: "TokenLe",
    TokenEq: "TokenEq",
    TokenFor: "TokenFor",
    TokenDo: "TokenDo",
    TokenIf: "TokenIf",
    TokenElse: "TokenElse",
    TokenEnd: "TokenEnd",
    TokenSyscall1: "TokenSyscall1",
    TokenSyscall3: "TokenSyscall3",
    TokenMacro: "TokenMacro",
    TokenVar: "TokenVar",
    TokenRead: "TokenRead",
    TokenWrite: "TokenWrite",
    TokenWord: "TokenWord",
}




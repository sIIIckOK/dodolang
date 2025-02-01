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
    *ts = (*ts)[:len(*ts)-1]
	return t
}

func (ts TypeStack) len() uint {
	return uint(len(ts))
}

func typeCheck(strTokens []StringToken, tokens []Token) bool {
	_ = strTokens
	var stack TypeStack
	assert(TokenCount == 29, "Exhaustive switch case for typeCheck")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		loc := token.Loc
		switch token.Type {
		case TokenInt:
			stack.push(TokenInt)
		case TokenPlus:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b],
					intrinsicStr[a],
				)
				return false
			}
			stack.push(TokenInt)
		case TokenSub:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b],
					intrinsicStr[a],
				)
				return false
			}
			stack.push(TokenInt)
		case TokenMult:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b],
					intrinsicStr[a],
				)
				return false
			}
			stack.push(TokenInt)
		case TokenPrint:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 int found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			if a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v > on the stack",
					intrinsicStr[a],
				)
				return false
			}
		case TokenDivMod:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(token, "expected atleast 2 ints")
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenInt)
			stack.push(TokenInt)
		case TokenSwap:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 elements found %v elements on the stack",
					stack.len(),
				)
				return false
			}
		case TokenDup:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 elements found %v elements on the stack",
					stack.len(),
				)
				return false
			}
            a := stack.pop(loc)
            stack.push(a)
            stack.push(a)
		case TokenDrop:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 elements found %v elements on the stack",
					stack.len(),
				)
				return false
			}
		case TokenRot:
			if stack.len() < 3 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 3 elements found %v elements on the stack",
					stack.len(),
				)
				return false
			}
		case TokenGt:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenBool)
		case TokenGe:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenBool)
		case TokenLt:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenBool)
		case TokenLe:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenBool)
		case TokenEq:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 ints found %v elements on the stack",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			b := stack.pop(loc)
			if b != TokenInt && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a],
					intrinsicStr[b],
				)
				return false
			}
			stack.push(TokenBool)
		case TokenFor:
		case TokenDo:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 bool got %v elements",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			if a != TokenBool && a != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 bools found < %v > on the stack",
					intrinsicStr[a],
				)
				return false
			}
		case TokenIf:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 bool got %v elements",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			if !(a == TokenBool || a == TokenInt) {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 bools found < %v > on the stack",
					intrinsicStr[a],
				)
				return false
			}
		case TokenElse:
		case TokenEnd:
		case TokenRead:
			if stack.len() < 1 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 ptr found %v elements",
					stack.len(),
				)
				return false
			}
			a := stack.pop(loc)
			if a != TokenPtr {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 ptr found < %v > on the stack",
					intrinsicStr[a],
				)
				return false
			}
			stack.push(TokenInt)
		case TokenWrite:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 ptr and 1 int found %v element(s)",
					stack.len(),
				)
				return false
			}
			b := stack.pop(loc)
			a := stack.pop(loc)
			if a != TokenPtr {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 ptr and 1 int found < %v %v > on the stack",
					intrinsicStr[b],
					intrinsicStr[a],
				)
				return false
			}
		case TokenSyscall1:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 2 elements found %v elements",
					stack.len(),
				)
				return false
			}
		case TokenSyscall3:
			if stack.len() < 4 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 4 elements found %v elements",
					stack.len(),
				)
				return false
			}
		case TokenMacro:
		case TokenVar:
		case TokenWord:
            stack.push(TokenPtr)
		}
	}
	return true
}

func printCompilerErrorInstrinsic(token Token, err string, args ...any) {
	assert(len(intrinsicStr) == TokenCount, "")
	fmtStr := fmt.Sprintf(err, args...)
	fmt.Printf("%v:%v:%v: `%v` %v\n",
		token.Loc.FilePath, token.Loc.Line, token.Loc.Col, intrinsicStr[token.Type], fmtStr)
	os.Exit(1)
}

var intrinsicStr = map[TokenType]string{
	TokenInt:      "TokenInt",
	TokenBool:     "TokenBool",
	TokenPtr:      "TokenPtr",
	TokenPlus:     "TokenPlus",
	TokenSub:      "TokenSub",
	TokenMult:     "TokenMult",
	TokenDivMod:   "TokenDivMod",
	TokenPrint:    "TokenPrint",
	TokenSwap:     "TokenSwap",
	TokenDup:      "TokenDup",
	TokenDrop:     "TokenDrop",
	TokenRot:      "TokenRot",
	TokenGt:       "TokenGt",
	TokenGe:       "TokenGe",
	TokenLt:       "TokenLt",
	TokenLe:       "TokenLe",
	TokenEq:       "TokenEq",
	TokenFor:      "TokenFor",
	TokenDo:       "TokenDo",
	TokenIf:       "TokenIf",
	TokenElse:     "TokenElse",
	TokenEnd:      "TokenEnd",
	TokenSyscall1: "TokenSyscall1",
	TokenSyscall3: "TokenSyscall3",
	TokenMacro:    "TokenMacro",
	TokenVar:      "TokenVar",
	TokenRead:     "TokenRead",
	TokenWrite:    "TokenWrite",
	TokenWord:     "TokenWord",
}

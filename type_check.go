package main

import (
	"fmt"
	"os"
)

type TypeStack []TypeInfo

func (ts *TypeStack) push(t TypeInfo) {
	*ts = append(*ts, t)
}

func (ts *TypeStack) pop(loc Location) TypeInfo {
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

type TypeInfo struct {
	Type TokenType
	Kind TokenType
}

func typeCheck(strTokens []StringToken, tokens []Token) bool {
	_ = strTokens
	var stack TypeStack
	assert(TokenCount == 32, "Exhaustive switch case for typeCheck")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		loc := token.Loc
		switch token.Type {
		case TokenInt:
			stack.push(TypeInfo{TokenInt, TokenInt})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b.Type],
					intrinsicStr[a.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenInt, TokenInt})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b.Type],
					intrinsicStr[a.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenInt, TokenInt})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[b.Type],
					intrinsicStr[a.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenInt, TokenInt})
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
			if !(a.Type == TokenInt || a.Type == TokenPtr || a.Type == TokenBool) {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 ints or ptr found < %v > on the stack",
					intrinsicStr[a.Type],
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 ints found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenInt, TokenInt})
			stack.push(TypeInfo{TokenInt, TokenInt})
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
		case TokenTrue:
			stack.push(TypeInfo{TokenBool, TokenBool})
		case TokenFalse:
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if b.Type != TokenInt && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 2 bools found < %v %v > on the stack",
					intrinsicStr[a.Type],
					intrinsicStr[b.Type],
				)
				return false
			}
			stack.push(TypeInfo{TokenBool, TokenBool})
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
			if a.Type != TokenBool && a.Type != TokenInt {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 bools found < %v > on the stack",
					intrinsicStr[a.Type],
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
			if !(a.Type == TokenBool) {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 bools found < %v > on the stack",
					intrinsicStr[a.Type],
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
			if a.Type != TokenPtr {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 ptr found < %v > on the stack",
					intrinsicStr[a.Type],
				)
				return false
			}
			stack.push(TypeInfo{token.Kind, token.Kind})
		case TokenWrite:
			if stack.len() < 2 {
				printCompilerErrorInstrinsic(
					token,
					"expected atleast 1 ptr and 1 int found %v element(s)",
					stack.len(),
				)
				return false
			}
			value := stack.pop(loc)
			varr := stack.pop(loc)
			if varr.Type != TokenPtr {
				printCompilerErrorInstrinsic(
					token,
					"takes 1 ptr and 1 int found < %v %v > on the stack",
					intrinsicStr[value.Type],
					intrinsicStr[varr.Type],
				)
				return false
			} else if varr.Kind != value.Type {
				printCompilerErrorInstrinsic(
					token,
					"trying to write a value of type %v into type variable %v",
					intrinsicStr[value.Type],
					intrinsicStr[varr.Kind],
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
			stack.push(
				TypeInfo{
					Type: TokenPtr,
					Kind: token.Kind,
				},
			)
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
	TokenMacroEnd: "TokenMacroEnd",
	TokenVar:      "TokenVar",
	TokenRead:     "TokenRead",
	TokenWrite:    "TokenWrite",
	TokenWord:     "TokenWord",
	TokenTrue:     "TokenTrue",
	TokenFalse:    "TokenFalse",
}

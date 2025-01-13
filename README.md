# About the Language
This is a concatenative stack-based programming language written in golang inspired by the likes of [forth](https://www.forth.com/forth/), 
[porth](https://gitlab.com/tsoding/porth) and [other concatenative languages](https://concatenative.org) similar to them.

**NOTE:** This is only a passion project, use this at your own risks, as there is no 
stability and security guaranteed.

## Supported Platforms
- x86 linux only

## Prerequisites
- go compiler
- nasm assembler

## How To Run

```cmd
go build
./dodolang build <file>.dodo
./<built-exe>
```

# Syntax and Features
Consult the `examples/` for up-to-date syntax and features of the language.
Additionally, you can learn more about concatenative languages from here:
- Concatenative language: https://concatenative.org
- Wikipedia: https://en.wikipedia.org/wiki/Concatenative_programming_language

In its current form it is a very limited language, which compiles down to native x86_64 assembly similar to porth, and it uses the Nasm assembler to assemble it down to statically linked bytecode. It is static as it utilizes system-calls' instead of linking to a dynamic library.

# What has to come (this list might change)
- macros and functions
- memory manipulation
- Compile time type checking

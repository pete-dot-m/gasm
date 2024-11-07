package main

import (
	"fmt"
	"gasm/parser"
	"os"
)

func main() {
	file, err := os.Open("hello.asm")
	if err != nil {
		panic(err)
	}

	lexer := parser.NewLexer(file)
	err = lexer.Lex()
	if err != nil {
		panic(err)
	}
	fmt.Println(lexer)
}

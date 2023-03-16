package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type command string

const (
	C_ARITHMETIC command = "C_ARITHMETIC"
	C_PUSH       command = "C_PUSH"
	C_POP        command = "C_POP"
	C_LABEL      command = "C_LABEL"
	C_GOTO       command = "C_GOTO"
	C_IF         command = "C_IF"
	C_FUNCTION   command = "C_FUNCTION"
	C_RETURN     command = "C_RETURN"
	C_CALL       command = "C_CALL"
)

func main() {
	args := os.Args[1:]
	fmt.Printf("args: %+v\n", args)
	parser := newParser(args[0])
	for parser.advance() {
		cmd := parser.commandType()
		line := parser.getLine()
		arg1 := parser.arg1(cmd)
		arg2 := parser.arg2(cmd)
		fmt.Printf("line: %s - cmd: %s arg1: %s arg2: %s\n", line, cmd, arg1, arg2)
	}
}

type parser struct {
	file    *os.File
	scanner *bufio.Scanner
}

func newParser(path string) *parser {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	return &parser{file, scanner}
}

func (p *parser) advance() bool {
	for {
		ok := p.scanner.Scan()
		if !ok {
			return false
		}
		cmd := p.scanner.Text()
		if cmd == "" || cmd[:2] == "//" {
			continue
		}
		return true
	}
}

func (p *parser) getLine() string {
	return p.scanner.Text()
}

func (p *parser) commandType() command {
	cmd := p.scanner.Text()
	if cmd[:3] == "add" {
		return C_ARITHMETIC
	} else if cmd[:4] == "push" {
		return C_PUSH
	} else if cmd[:3] == "pop" {
		return C_POP
	} else {
		log.Fatal("not implemented")
		return ""
	}
}

func (p *parser) arg1(cmd command) string {
	line := p.scanner.Text()
	if cmd == C_ARITHMETIC {
		return line
	} else {
		return strings.Split(line, " ")[1]
	}
}

func (p *parser) arg2(cmd command) string {
	line := p.scanner.Text()
	if cmd == C_PUSH || cmd == C_POP || cmd == C_FUNCTION || cmd == C_CALL {
		return strings.Split(line, " ")[2]
	} else {
		return ""
	}
}

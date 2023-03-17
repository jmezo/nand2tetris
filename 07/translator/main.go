package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

type segment string

const (
	argument segment = "argument"
	local    segment = "local"
	static   segment = "static"
	constant segment = "constant"
	this     segment = "this"
	that     segment = "that"
	pointer  segment = "pointer"
	temp     segment = "temp"
)

func main() {
	args := os.Args[1:]
	fmt.Printf("args: %+v\n", args)
	parser := newParser(args[0])
	codeWriter := newCodeWriter("out.asm")
	for parser.advance() {
		cmd := parser.commandType()
		line := parser.getLine()
		arg1 := parser.arg1(cmd)
		arg2 := parser.arg2(cmd)
		fmt.Printf("line: %s - cmd: %s arg1: %s arg2: %d\n", line, cmd, arg1, arg2)
		if cmd == C_ARITHMETIC {
			codeWriter.writeArithmetic(arg1)
		} else {
			codeWriter.writePushPop(cmd, arg1, arg2)
		}
	}
	codeWriter.close()
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
	fullCmd := p.scanner.Text()
	cmd := fullCmd[:3]
	if cmd == "add" || cmd == "sub" || cmd == "neg" || cmd == "eq" || cmd == "gt" || cmd == "lt" || cmd == "and" || cmd == "or" || cmd == "not" {
		return C_ARITHMETIC
	} else if fullCmd[:4] == "push" {
		return C_PUSH
	} else if fullCmd[:3] == "pop" {
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

func (p *parser) arg2(cmd command) int {
	line := p.scanner.Text()
	if cmd == C_PUSH || cmd == C_POP || cmd == C_FUNCTION || cmd == C_CALL {
		argStr := strings.Split(line, " ")[2]
		arg, err := strconv.Atoi(argStr)
		if err != nil {
			log.Fatalf("arg2: %s is not a number on line: %s", argStr, line)
		}
		return arg
	} else {
		return -1
	}
}

type codeWriter struct {
	file       *os.File
	vmFileName string
}

func newCodeWriter(fileName string) *codeWriter {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return &codeWriter{file, ""}
}

func (c *codeWriter) setFileName(fileName string) {
	c.vmFileName = fileName
}

func (c *codeWriter) writeArithmetic(cmd string) {
	// TODO
}

func (c *codeWriter) writePushPop(cmd command, segment string, index int) {
	// TODO
}

func (c *codeWriter) close() {
	c.file.Close()
}

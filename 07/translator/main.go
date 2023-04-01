package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type command string

const stackPointerDefault = 256

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
	filePath := args[0]
	fileName := filePath
	var parsers []*parser
	// if filename has .vm extension, then it's a single file
	if len(filePath) > 3 && filePath[len(filePath)-3:] == ".vm" {
		fileName = filePath[:len(filePath)-3]
		parsers = append(parsers, newParser(filePath))
	} else {
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".vm") {
				parsers = append(parsers, newParser(filePath+"/"+file.Name()))
			}
		}

	}

	codeWriter := newCodeWriter(fileName + ".asm")
	codeWriter.initStack()
	for _, parser := range parsers {
		for parser.advance() {
			cmd := parser.commandType()
			line := parser.getLine()
			arg1 := parser.arg1(cmd)
			arg2 := parser.arg2(cmd)
			fmt.Printf("line: %s - cmd: %s arg1: %s arg2: %d\n",
				line, cmd, arg1, arg2)
			if cmd == C_ARITHMETIC {
				codeWriter.writeArithmetic(arg1)
			} else if cmd == C_PUSH || cmd == C_POP {
				codeWriter.writePushPop(cmd, arg1, arg2)
			} else {
				log.Fatal("not implemented")
			}
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
	cmd := p.scanner.Text()
	if cmd == "add" || cmd == "sub" || cmd == "neg" || cmd == "eq" || cmd == "gt" || cmd == "lt" || cmd == "and" || cmd == "or" || cmd == "not" {
		return C_ARITHMETIC
	} else if len(cmd) > 4 && cmd[:4] == "push" {
		return C_PUSH
	} else if len(cmd) > 3 && cmd[:3] == "pop" {
		return C_POP
	} else {
		log.Fatalf("not implemented: %s", cmd)
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
	stackIndex int
}

func newCodeWriter(fileName string) *codeWriter {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return &codeWriter{file, "", stackPointerDefault}
}

func (c *codeWriter) initStack() {
	c.file.WriteString("// init stack\n")
	c.file.WriteString("@256\nD=A\n@SP\nM=D\n\n")
}

func (c *codeWriter) setFileName(fileName string) {
	c.vmFileName = fileName
}

func (c *codeWriter) writeArithmetic(cmd string) {
	getStackTop := "@SP\nM=M-1\nA=M\n"
	setStackTopToD := getStackTop + "D=M\n" + "M=0\n"
	incrStack := "@SP\nM=M+1\n"
	var asmC string
	switch cmd {
	case "add":
		c.file.WriteString("// add\n")
		asmC = "@SP\nM=M-1\nA=M\nD=M\nM=0\n@SP\nM=M-1\nA=M\nM=M+D\n@SP\nM=M+1\n\n"
	case "sub":
		c.file.WriteString("// sub\n")
		asmC = "@SP\nM=M-1\nA=M\nD=M\nM=0\n@SP\nM=M-1\nA=M\nM=M-D\n@SP\nM=M+1\n\n"
	case "neg":
		c.file.WriteString("// neg\n")
		asmC = "@SP\nM=M-1\nA=M\nM=-M\n@SP\nM=M+1\n\n"
	case "eq":
		c.file.WriteString("// eq\n")
		asmC = setStackTopToD + getStackTop +
			"D=M-D\n" +
			"@EQ\n" +
			"D;JEQ\n" +
			"M=0\n" +
			"@END\n" +
			"(EQ)\n" +
			"M=-1\n" +
			"(END)\n" +
			incrStack + "\n"
	case "gt": // x > y
		c.file.WriteString("// gt\n")
		asmC = setStackTopToD + getStackTop +
			"D=M-D\n" +
			"@EQ\n" +
			"D;JGT\n" +
			"M=0\n" +
			"@END\n" +
			"(EQ)\n" +
			"M=-1\n" +
			"(END)\n" +
			incrStack + "\n"
	case "lt": // x < y
		c.file.WriteString("// lt\n")
		asmC = setStackTopToD + getStackTop +
			"D=M-D\n" +
			"@EQ\n" +
			"D;JLT\n" +
			"M=0\n" +
			"@END\n" +
			"(EQ)\n" +
			"M=-1\n" +
			"(END)\n" +
			incrStack + "\n"
	case "and":
		c.file.WriteString("// and\n")
		asmC = setStackTopToD + getStackTop + "M=D&M\n" + incrStack + "\n"
	case "or":
		c.file.WriteString("// or\n")
		asmC = setStackTopToD + getStackTop + "M=D|M\n" + incrStack + "\n"
	case "not":
		c.file.WriteString("// not\n")
		asmC = getStackTop + "M=!M\n" + incrStack + "\n"
	default:
		log.Fatal("not implemented")
	}
	c.file.WriteString(asmC)
}

func (c *codeWriter) writePushPop(cmd command, segment string, index int) {
	switch cmd {
	case C_PUSH:
		switch segment {
		case "constant":
			c.file.WriteString(fmt.Sprintf("// push constant %d\n", index))
			asm := "@%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
			c.file.WriteString(fmt.Sprintf(asm, index))
		default:
			log.Fatal("not implemented")
		}
	case C_POP:
		log.Fatal("not implemented")
	default:
		log.Fatal("not implemented")
	}
}

func (c *codeWriter) close() {
	c.file.Close()
}

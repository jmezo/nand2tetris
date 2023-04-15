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
	asmPath := filePath + ".asm"
	var parsers []*parser
	// if filename has .vm extension, then it's a single file
	if len(filePath) > 3 && filePath[len(filePath)-3:] == ".vm" {
		fileName := filePath[:len(filePath)-3]
		asmPath = fileName + ".asm"
		parsers = append(parsers, newParser(filePath))
	} else {
		fileName := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
		asmPath = filePath + "/" + fileName + ".asm"
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
	fmt.Println("fileName: ", asmPath)

	codeWriter := newCodeWriter(asmPath)
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
	cmdCount   int
}

func newCodeWriter(fileName string) *codeWriter {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return &codeWriter{file, "", stackPointerDefault, 0}
}

func (c *codeWriter) initStack() {
	cmd := "// init stack\n"
	cmd += "@256\nD=A\n@SP\nM=D\n\n"
	c.writeCommand(cmd)
}

func (c *codeWriter) setFileName(fileName string) {
	c.vmFileName = fileName
}

func (c *codeWriter) writeCommand(cmd string) {
	c.file.WriteString(cmd)
	c.cmdCount++
}

func (c *codeWriter) writeArithmetic(cmd string) {
	cmdCount := strconv.Itoa(c.cmdCount)
	getStackTop := "@SP\n" +
		"M=M-1\n" +
		"A=M\n"

	incrStack := "@SP\n" +
		"M=M+1\n"

	alu2ParamCommand := getStackTop +
		"D=M\n" +
		"@SP\n" +
		"M=M-1\n" +
		"A=M\n" +
		"M=%s\n" + // M=M+D, M=M-D, M=M&D, M=M|D
		incrStack + "\n"

	alu1ParamCommand := getStackTop +
		"M=%s\n" + // M=-M, M=!M
		incrStack + "\n"

	cmpCommand := getStackTop +
		"D+M\n" +
		getStackTop +
		"D=M-D\n" +
		"@CMD" + cmdCount + "\n" +
		"D;%s\n" + // JEQ, JGT, JLT
		"@SP\n" +
		"A=M\n" +
		"M=0\n" +
		"@END" + cmdCount + "\n" +
		"0;JMP\n" +
		"(CMD" + cmdCount + ")\n" +
		"@SP\n" +
		"A=M\n" +
		"M=-1\n" +
		"(END" + cmdCount + ")\n" +
		incrStack + "\n"

	var asmC string
	switch cmd {
	case "add":
		asmC = "// add\n"
		asmC += fmt.Sprintf(alu2ParamCommand, "M+D")
	case "sub":
		asmC = "// sub\n"
		asmC += fmt.Sprintf(alu2ParamCommand, "M-D")
	case "neg":
		asmC = "// neg\n"
		asmC += fmt.Sprintf(alu1ParamCommand, "-M")
	case "eq":
		asmC = "// eq\n"
		asmC += fmt.Sprintf(cmpCommand, "JEQ")
	case "gt": // x > y
		asmC = "// gt\n"
		asmC += fmt.Sprintf(cmpCommand, "JGT")
	case "lt": // x < y
		asmC = "// lt\n"
		asmC += fmt.Sprintf(cmpCommand, "JLT")
	case "and":
		asmC = "// and\n"
		asmC += fmt.Sprintf(alu2ParamCommand, "M&D")
	case "or":
		asmC = "// or\n"
		asmC += fmt.Sprintf(alu2ParamCommand, "M|D")
	case "not":
		asmC = "// not\n"
		asmC += fmt.Sprintf(alu1ParamCommand, "!M")
	default:
		log.Fatal("not implemented")
	}
	c.writeCommand(asmC)
}

func (c *codeWriter) writePushPop(cmd command, segment string, index int) {
	pushDToStack := "@SP\n" +
		"A=M\n" +
		"M=D\n" +
		"@SP\n" +
		"M=M+1\n"

	pushSegmentToStack := "@%s\n" +
		"D=M\n" +
		"@%d\n" +
		"A=A+D\n" +
		"D=M\n" +
		pushDToStack

	pushRamToStack := "@%s\n" +
		"D=A\n" +
		"@%d\n" +
		"A=A+D\n" +
		"D=M\n" +
		pushDToStack

	popStackToD := "@SP\n" +
		"M=M-1\n" +
		"A=M\n" +
		"D=M\n"

	popStackToRam := "@%s\n" +
		"D=A\n" +
		"@%d\n" +
		"D=A+D\n" +
		"@R13\n" +
		"M=D\n" +
		popStackToD +
		"@R13\n" +
		"A=M\n" +
		"M=D\n"

	popStackToSegment := "@%s\n" +
		"D=M\n" +
		"@%d\n" +
		"D=A+D\n" +
		"@R13\n" +
		"M=D\n" +
		popStackToD +
		"@R13\n" +
		"A=M\n" +
		"M=D\n"

	switch cmd {
	case C_PUSH:
		switch segment {
		case "constant":
			cmd := fmt.Sprintf("// push constant %d\n", index)
			cmd += fmt.Sprintf("@%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n", index)
			c.writeCommand(cmd)
		case "local":
			cmd := fmt.Sprintf("// push local %d\n", index)
			cmd += fmt.Sprintf(pushSegmentToStack, "LCL", index)
			c.writeCommand(cmd)
		case "argument":
			cmd := fmt.Sprintf("// push argument %d\n", index)
			cmd += fmt.Sprintf(pushSegmentToStack, "ARG", index)
			c.writeCommand(cmd)
		case "this":
			cmd := fmt.Sprintf("// push this %d\n", index)
			cmd += fmt.Sprintf(pushSegmentToStack, "THIS", index)
			c.writeCommand(cmd)
		case "that":
			cmd := fmt.Sprintf("// push that %d\n", index)
			cmd += fmt.Sprintf(pushSegmentToStack, "THAT", index)
			c.writeCommand(cmd)
		case "pointer":
			cmd := fmt.Sprintf("// push pointer %d\n", index)
			cmd += fmt.Sprintf(pushRamToStack, "THIS", index)
			c.writeCommand(cmd)
		case "temp":
			cmd := fmt.Sprintf("// push temp %d\n", index)
			cmd += fmt.Sprintf(pushRamToStack, "R5", index)
			c.writeCommand(cmd)
		default:
			log.Fatal("not implemented")
		}
	case C_POP:
		switch segment {
		case "local":
			cmd := fmt.Sprintf("// pop local %d\n", index)
			cmd += fmt.Sprintf(popStackToSegment, "LCL", index)
			c.writeCommand(cmd)
		case "argument":
			cmd := fmt.Sprintf("// pop argument %d\n", index)
			cmd += fmt.Sprintf(popStackToSegment, "ARG", index)
			c.writeCommand(cmd)
		case "this":
			cmd := fmt.Sprintf("// pop this %d\n", index)
			cmd += fmt.Sprintf(popStackToSegment, "THIS", index)
			c.writeCommand(cmd)
		case "that":
			cmd := fmt.Sprintf("// pop that %d\n", index)
			cmd += fmt.Sprintf(popStackToSegment, "THAT", index)
			c.writeCommand(cmd)
		case "pointer":
			cmd := fmt.Sprintf("// pop pointer %d\n", index)
			cmd += fmt.Sprintf(popStackToRam, "THIS", index)
			c.writeCommand(cmd)
		case "temp":
			cmd := fmt.Sprintf("// pop temp %d\n", index)
			cmd += fmt.Sprintf(popStackToRam, "R5", index)
			c.writeCommand(cmd)
		default:
			log.Fatal("not implemented")
		}
	default:
		log.Fatal("not implemented")
	}
}

func (c *codeWriter) close() {
	c.file.Close()
}

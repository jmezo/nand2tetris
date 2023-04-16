// init stack
@256
D=A
@SP
M=D

// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1

// eq
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD3
D;JEQ
@SP
A=M
M=0
@END3
0;JMP
(CMD3)
@SP
A=M
M=-1
(END3)
@SP
M=M+1

// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 16
@16
D=A
@SP
A=M
M=D
@SP
M=M+1

// eq
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD6
D;JEQ
@SP
A=M
M=0
@END6
0;JMP
(CMD6)
@SP
A=M
M=-1
(END6)
@SP
M=M+1

// push constant 16
@16
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1

// eq
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD9
D;JEQ
@SP
A=M
M=0
@END9
0;JMP
(CMD9)
@SP
A=M
M=-1
(END9)
@SP
M=M+1

// push constant 892
@892
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1

// lt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD12
D;JLT
@SP
A=M
M=0
@END12
0;JMP
(CMD12)
@SP
A=M
M=-1
(END12)
@SP
M=M+1

// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 892
@892
D=A
@SP
A=M
M=D
@SP
M=M+1

// lt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD15
D;JLT
@SP
A=M
M=0
@END15
0;JMP
(CMD15)
@SP
A=M
M=-1
(END15)
@SP
M=M+1

// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1

// lt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD18
D;JLT
@SP
A=M
M=0
@END18
0;JMP
(CMD18)
@SP
A=M
M=-1
(END18)
@SP
M=M+1

// push constant 32767
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1

// gt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD21
D;JGT
@SP
A=M
M=0
@END21
0;JMP
(CMD21)
@SP
A=M
M=-1
(END21)
@SP
M=M+1

// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 32767
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1

// gt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD24
D;JGT
@SP
A=M
M=0
@END24
0;JMP
(CMD24)
@SP
A=M
M=-1
(END24)
@SP
M=M+1

// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1

// gt
@SP
M=M-1
A=M
D+M
@SP
M=M-1
A=M
D=M-D
@CMD27
D;JGT
@SP
A=M
M=0
@END27
0;JMP
(CMD27)
@SP
A=M
M=-1
(END27)
@SP
M=M+1

// push constant 57
@57
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 31
@31
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 53
@53
D=A
@SP
A=M
M=D
@SP
M=M+1

// add
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M+D
@SP
M=M+1

// push constant 112
@112
D=A
@SP
A=M
M=D
@SP
M=M+1

// sub
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1

// neg
@SP
M=M-1
A=M
M=-M
@SP
M=M+1

// and
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M&D
@SP
M=M+1

// push constant 82
@82
D=A
@SP
A=M
M=D
@SP
M=M+1

// or
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M|D
@SP
M=M+1

// not
@SP
M=M-1
A=M
M=!M
@SP
M=M+1


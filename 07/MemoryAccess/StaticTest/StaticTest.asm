// init stack
@256
D=A
@SP
M=D

// push constant 111
@111
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 333
@333
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 888
@888
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop static 8
@SP
M=M-1
A=M
D=M
@static.8
M=D
// pop static 3
@SP
M=M-1
A=M
D=M
@static.3
M=D
// pop static 1
@SP
M=M-1
A=M
D=M
@static.1
M=D
// push static 3
@static.3
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@static.1
D=M
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

// push static 8
@static.8
D=M
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


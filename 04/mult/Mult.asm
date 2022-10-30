// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// Put your code here.


// clear variables
@i
M=0
@z
M=0
@R2
M=0

// if R0 = 0 then jump to end
@R0
D=M
@END
D;JEQ
// if R1 = 0 then jump to end
@R1
D=M
@END
D;JEQ
// if R1 = 1 then set R2 to R0
D=D-1
@r1is1
D;JEQ

// z = x * y
// set i = R1 - 1
@R1
D=M-1
@i
M=D
// set z = R0
@R0
D=M
@z
M=D

(LOOP)
// set z += R0
@z
D=M
@R0
D=D+M
@z
M=D

// set i--
@i
M=M-1
D=M
@LOOP
D;JNE

// set R2 = z
@z
D=M
@R2
M=D
// end
(END)
@END
0;JMP

(r1is1)
@R0
D=M
@R2
M=D
@END
0;JMP

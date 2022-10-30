// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

// put the @screen pointer into @pixel register
@SCREEN
D=A
@pixel
M=D

(checkKeyboard)
// put value of keyboard into D
@24576
D=M
// if key is pressed start blackening screen
@blackScreen
D;JNE
// else start whitening screen
@whiteScreen
0;JMP

(blackScreen)
	// set current pixel black
	@pixel
	A=M
	M=-1

	// if @pixel is at the last pixel, jump back to @checkKeyboard
	@pixel
	D=M
	@24576 // last register
	D=A-D
	@checkKeyboard
	D;JLE
	// else next pixel
	@pixel
	M=M+1
	// go back to keyboard check
	@checkKeyboard
	0;JMP

(whiteScreen)
	// set current pixel white
	@pixel
	A=M
	M=0

	// if @pixel is at the first pixel, jump back to @checkKeyboard
	@pixel
	D=M
	@SCREEN
	D=D-A
	@checkKeyboard
	D;JLE

	// else set to previous pixel
	@pixel
	M=M-1
	// go back to keyboard check
	@checkKeyboard
	0;JMP

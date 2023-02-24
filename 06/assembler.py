import sys
import typing

A_COMMAND = "A_COMMAND"
C_COMMAND = "C_COMMAND"
L_COMMAND = "L_COMMAND"

def hasMoreCommands(line: str):
    return line != ""

def advance(f: typing.IO) -> str:
    return f.readline()

def stripLine(line: str) -> str:
    res = ""
    i = 0
    while i < len(line):
        if line[i] == " " or line[i] == "\t" or line[i] == "\r" or line[i] == "\n":
            i += 1
            continue
        if i + 1 < len(line) and line[i] == "/" and line[i + 1] == "/":
            break
        res += line[i]
        i += 1
    return res

def commandType(line: str) -> str:
    if line[0] == "@":
        return A_COMMAND
    elif line[0] == "(":
        return L_COMMAND
    else:
        return C_COMMAND

def symbol(line: str) -> str:
    nameBeginIndex = 0
    nameEndIndex = 0
    i = 0
    while i < len(line):
        if line[i] == ")":
            nameEndIndex = i
        i += 1
    return line[nameBeginIndex:nameEndIndex]

def dest(line: str) -> str:
    equalIndex = 0
    i = 0
    while i < len(line):
        if line[i] == "=":
            equalIndex = i
            break
        i += 1
    return line[:equalIndex]

def comp(line: str) -> str:
    destIndex = -1
    jumpIndex = -1 
    i = 0
    while i < len(line):
        if line[i] == "=":
            destIndex = i
        if line[i] == ";":
            jumpIndex = i
        i += 1
    if jumpIndex == -1:
        jumpIndex = len(line)
    return line[destIndex + 1:jumpIndex]

def jump(line: str) -> str:
    semicolonIndex = 0
    i = 0
    while i < len(line):
        if line[i] == ";":
            semicolonIndex = i
            break
        i += 1
    if semicolonIndex == 0:
        return ""
    return line[semicolonIndex + 1:]

def dest_code(dest: str) -> str:
    destMap = {
            "": "000",
            "M": "001",
            "D": "010",
            "MD": "011",
            "A": "100",
            "AM": "101",
            "AD": "110",
            "AMD": "111"
            }
    return destMap[dest]

def comp_code(comp: str) -> str:
    compMap = {
            "0": "0101010",
            "1": "0111111",
            "-1": "0111010",
            "D": "0001100",
            "A": "0110000",
            "!D": "0001101",
            "!A": "0110001",
            "-D": "0001111",
            "-A": "0110011",
            "D+1": "0011111",
            "A+1": "0110111",
            "D-1": "0001110",
            "A-1": "0110010",
            "D+A": "0000010",
            "D-A": "0010011",
            "A-D": "0000111",
            "D&A": "0000000",
            "D|A": "0010101",
            "M": "1110000",
            "!M": "1110001",
            "-M": "1110011",
            "M+1": "1110111",
            "M-1": "1110010",
            "D+M": "1000010",
            "D-M": "1010011",
            "M-D": "1000111",
            "D&M": "1000000",
            "D|M": "1010101"
            }
    return compMap[comp]

def jump_code(jump: str) -> str:
    jumpMap = {
            "": "000",
            "JGT": "001",
            "JEQ": "010",
            "JGE": "011",
            "JLT": "100",
            "JNE": "101",
            "JLE": "110",
            "JMP": "111"
            }
    return jumpMap[jump]

def main():
    # get args
    args = sys.argv[1:]
    bin_commands = []
    file_name_with_ext = args[0]
    file_name = file_name_with_ext.split(".")[0]

    with open(file_name_with_ext, 'r') as f:
        while True:
            nstr = advance(f)
            if not hasMoreCommands(nstr):
                break
            nstr = stripLine(nstr)
            if nstr == "":
                continue
            cmdTzpe = commandType(nstr)
            if cmdTzpe == A_COMMAND:
                bin_cmd = "{0:016b}".format(int(nstr[1:]))
                bin_commands.append(bin_cmd)
            elif cmdTzpe == C_COMMAND:
                dst = dest(nstr)
                dst_code = dest_code(dst)
                cmp = comp(nstr)
                cmp_code = comp_code(cmp)
                jmp = jump(nstr)
                jmp_code = jump_code(jmp)
                # bin_cmd = "111" +  " " + cmp_code +  " " + dst_code +  " " + jmp_code
                bin_cmd = "111" + cmp_code + dst_code + jmp_code
                bin_commands.append(bin_cmd)
            elif cmdTzpe == L_COMMAND:
                print("L_COMMAND: " + nstr)

    for line in bin_commands:
        print(line)
    with open(file_name + ".hack", 'w') as f:
        for line in bin_commands:
            f.write("%s\n" % line)

if __name__ == "__main__":
    main()

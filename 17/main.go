package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type registers struct {
	a int
	b int
	c int
}

func parseComputer(scanner *bufio.Scanner) (registers, []int) {
	registers := registers{}
	program := []int{}
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Register A: %d", &registers.a)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Register B: %d", &registers.b)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Register C: %d", &registers.c)
	scanner.Scan()
	scanner.Scan()
	programStr := strings.Split(strings.Split(scanner.Text(), ": ")[1], ",")
	for _, s := range programStr {
		p, _ := strconv.Atoi(s)
		program = append(program, p)
	}

	return registers, program
}

func formatOutput(out []int) string {
	line := strconv.Itoa(out[0])
	for _, n := range out[1:] {
		line += ","
		line += strconv.Itoa(n)
	}
	return line
}

func powInt(base, exp int) int {
	result := 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}
	return result
}

func combo(n int, reg registers) int {
	operand := n
	switch n {
	case 4:
		operand = reg.a
	case 5:
		operand = reg.b
	case 6:
		operand = reg.c
	case 7:
		panic("program is invalid")
	}
	return operand
}

func adv(reg registers, operand int) registers {
	reg.a = reg.a / powInt(2, operand)
	return reg
}

func bxl(reg registers, operand int) registers {
	reg.b = reg.b ^ operand
	return reg
}

func bst(reg registers, operand int) registers {
	reg.b = operand % 8
	return reg
}

func jnz(reg registers, pointer int, operand int) int {
	if reg.a != 0 {
		return operand
	}
	return pointer
}

func bxc(reg registers) registers {
	reg.b = reg.b ^ reg.c
	return reg
}

func out(operand int) int {
	return operand % 8
}

func bdv(reg registers, operand int) registers {
	reg.b = reg.a / powInt(2, operand)
	return reg
}

func cdv(reg registers, operand int) registers {
	reg.c = reg.a / powInt(2, operand)
	return reg
}

func performProgram(reg registers, program []int) []int {
	outs := []int{}
	pointer := 0
	halt := false
	for !halt && pointer >= 0 && pointer < len(program) {
		switch program[pointer] {
		case 0:
			reg = adv(reg, combo(program[pointer+1], reg))
			pointer += 2
		case 1:
			reg = bxl(reg, program[pointer+1])
			pointer += 2
		case 2:
			reg = bst(reg, combo(program[pointer+1], reg))
			pointer += 2
		case 3:
			prevPointer := pointer
			pointer = jnz(reg, pointer, program[pointer+1])
			if pointer == prevPointer {
				halt = true
			}
		case 4:
			reg = bxc(reg)
			pointer += 2
		case 5:
			outs = append(outs, out(combo(program[pointer+1], reg)))
			pointer += 2
		case 6:
			reg = bdv(reg, combo(program[pointer+1], reg))
			pointer += 2
		case 7:
			reg = cdv(reg, combo(program[pointer+1], reg))
			pointer += 2
		}
	}

	return outs
}

func findFirstProgramMatch(program []int, start int, match []int) int {
	outs := []int{}
	i := start
	for !slices.Equal(match, outs) {
		outs = performProgram(registers{i, 0, 0}, program)
		i++
	}
	return i - 1
}

// think this could be smarter but it's smart enough
func findRegStartToCreateProgram(program []int) int {
	v := 0
	for i := len(program) - 1; i >= 0; i-- {
		match := program[i:]
		v = findFirstProgramMatch(program, v*8, match)
	}
	return v
}

func main() {
	start := time.Now()
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	regs, program := parseComputer(scanner)
	outs := performProgram(regs, program)

	fmt.Println("part 1:", formatOutput(outs))
	fmt.Println("part 2:", findRegStartToCreateProgram(program))

	log.Printf("Time taken: %s", time.Since(start))
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type pos struct {
	x int
	y int
}

func parseWarehouse(scanner *bufio.Scanner) (pos, map[pos]rune, []rune) {
	robot := pos{}
	warehouse := map[pos]rune{}
	instr := []rune{}
	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		if len(line) > 0 {
			if line[0] == '#' {
				for x, c := range line {
					if c == '@' {
						robot = pos{x: x, y: y}
					} else if c != '.' {
						warehouse[pos{x: x, y: y}] = c
					}
				}
			} else {
				instr = append(instr, []rune(line)...)
			}
		}
	}
	return robot, warehouse, instr
}

func getStepFromInstr(instr rune) pos {
	switch instr {
	case '<':
		return pos{x: -1, y: 0}
	case '>':
		return pos{x: 1, y: 0}
	case '^':
		return pos{x: 0, y: -1}
	case 'v':
		return pos{x: 0, y: 1}
	default:
		panic("unrecognised instruction")
	}
}

func performInstr(robot pos, warehouse map[pos]rune, instr rune) (pos, map[pos]rune) {
	step := getStepFromInstr(instr)
	nextPos := pos{x: robot.x + step.x, y: robot.y + step.y}
	if warehouse[nextPos] == 0 {
		robot = nextPos
	} else if warehouse[nextPos] == 'O' {
		currentPos := pos{x: nextPos.x + step.x, y: nextPos.y + step.y}
		for warehouse[currentPos] == 'O' {
			currentPos = pos{x: currentPos.x + step.x, y: currentPos.y + step.y}
		}
		if warehouse[currentPos] == 0 {
			robot = nextPos
			warehouse[currentPos] = 'O'
			warehouse[nextPos] = 0
		}
	}

	return robot, warehouse
}

func sumBoxGPS(warehouse map[pos]rune) int {
	total := 0
	for pos, c := range warehouse {
		if c == 'O' {
			total += pos.x + pos.y*100
		}
	}
	return total
}

func sumPosAfterInstr(robot pos, warehouse map[pos]rune, instr []rune) int {
	for _, i := range instr {
		robot, warehouse = performInstr(robot, warehouse, i)
	}

	return sumBoxGPS(warehouse)
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
	robot, warehouse, instr := parseWarehouse(scanner)
	fmt.Println("Part 1:", sumPosAfterInstr(robot, warehouse, instr))

	log.Printf("Time taken: %s", time.Since(start))
}

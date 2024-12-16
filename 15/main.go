package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
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

func sumBoxGPS(warehouse map[pos]rune, box rune) int {
	total := 0
	for pos, c := range warehouse {
		if c == box {
			total += pos.x + pos.y*100
		}
	}
	return total
}

func sumPosAfterInstr(robot pos, warehouse map[pos]rune, instr []rune) int {
	for _, i := range instr {
		robot, warehouse = performInstr(robot, warehouse, i)
	}

	return sumBoxGPS(warehouse, 'O')
}

func widenWarehouse(warehouse map[pos]rune) map[pos]rune {
	wideWarehouse := map[pos]rune{}
	for p, c := range warehouse {
		if c == 'O' {
			wideWarehouse[pos{x: p.x * 2, y: p.y}] = '['
			wideWarehouse[pos{x: p.x*2 + 1, y: p.y}] = ']'
		} else {
			wideWarehouse[pos{x: p.x * 2, y: p.y}] = c
			wideWarehouse[pos{x: p.x*2 + 1, y: p.y}] = c
		}
	}

	return wideWarehouse
}

func pushHorizontal(robot pos, warehouse map[pos]rune, step pos) (pos, map[pos]rune) {
	row := []pos{}
	currentPos := pos{x: robot.x + step.x, y: robot.y + step.y}
	for warehouse[currentPos] == '[' || warehouse[currentPos] == ']' {
		currentPos = pos{x: currentPos.x + step.x, y: currentPos.y + step.y}
		row = append(row, currentPos)
	}

	if warehouse[currentPos] == 0 {
		robot = pos{x: robot.x + step.x, y: robot.y + step.y}
		for i := len(row) - 1; i >= 0; i-- {
			warehouse[row[i]] = warehouse[pos{x: row[i].x - step.x, y: row[i].y}]
		}
		warehouse[robot] = 0
	}

	return robot, warehouse
}

func allTouchingSpace(rows [][]pos, warehouse map[pos]rune, step pos) bool {
	allTouchingSpace := true
	for _, p := range rows[len(rows)-1] {
		allTouchingSpace = allTouchingSpace && warehouse[pos{x: p.x, y: p.y + step.y}] == 0
	}
	return allTouchingSpace
}

func touchingWall(rows [][]pos, warehouse map[pos]rune, step pos) bool {
	touchingWall := false
	for _, p := range rows[len(rows)-1] {
		touchingWall = touchingWall || warehouse[pos{x: p.x, y: p.y + step.y}] == '#'
	}
	return touchingWall
}

func pushVertical(robot pos, warehouse map[pos]rune, step pos) (pos, map[pos]rune) {
	rows := [][]pos{}
	currentPos := pos{x: robot.x + step.x, y: robot.y + step.y}
	if warehouse[currentPos] == '[' {
		rows = append(rows, []pos{currentPos, {x: currentPos.x + 1, y: currentPos.y}})
	} else {
		rows = append(rows, []pos{currentPos, {x: currentPos.x - 1, y: currentPos.y}})
	}

	for !touchingWall(rows, warehouse, step) && !allTouchingSpace(rows, warehouse, step) {
		newRow := []pos{}
		for _, p := range rows[len(rows)-1] {
			nextP := pos{x: p.x, y: p.y + step.y}
			if warehouse[nextP] == '[' || warehouse[nextP] == ']' {
				newRow = append(newRow, nextP)
			}
		}
		sort.Slice(newRow, func(i, j int) bool { return newRow[i].x < newRow[j].x })
		for _, p := range newRow {
			if warehouse[p] == ']' && !slices.Contains(newRow, pos{p.x - 1, p.y}) {
				newRow = append(newRow, pos{x: p.x - 1, y: p.y})
			}
			if warehouse[p] == '[' && !slices.Contains(newRow, pos{p.x + 1, p.y}) {
				newRow = append(newRow, pos{x: p.x + 1, y: p.y})
			}
		}
		rows = append(rows, newRow)
	}

	if allTouchingSpace(rows, warehouse, step) {
		for i := len(rows) - 1; i >= 0; i-- {
			for _, p := range rows[i] {
				warehouse[pos{x: p.x, y: p.y + step.y}] = warehouse[p]
				warehouse[p] = 0
			}
		}
		robot = pos{robot.x, robot.y + step.y}
	}

	return robot, warehouse
}

func performInstrWide(robot pos, warehouse map[pos]rune, instr rune) (pos, map[pos]rune) {
	step := getStepFromInstr(instr)
	nextPos := pos{x: robot.x + step.x, y: robot.y + step.y}
	if warehouse[nextPos] == 0 {
		robot = nextPos
	} else if warehouse[nextPos] == '[' || warehouse[nextPos] == ']' {
		switch instr {
		case '<', '>':
			robot, warehouse = pushHorizontal(robot, warehouse, step)
		case '^', 'v':
			robot, warehouse = pushVertical(robot, warehouse, step)
		}
	}

	return robot, warehouse
}

func sumPosAfterInstrWide(robot pos, warehouse map[pos]rune, instr []rune) int {
	for _, i := range instr {
		robot, warehouse = performInstrWide(robot, warehouse, i)
	}

	return sumBoxGPS(warehouse, '[')
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
	wideWarehouse := widenWarehouse(warehouse)
	wideRobot := pos{x: robot.x * 2, y: robot.y}
	fmt.Println("Part 1:", sumPosAfterInstr(robot, warehouse, instr))
	fmt.Println("Part 2:", sumPosAfterInstrWide(wideRobot, wideWarehouse, instr))

	log.Printf("Time taken: %s", time.Since(start))
}

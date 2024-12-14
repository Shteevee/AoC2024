package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const X_LIMIT = 100
const Y_LIMIT = 102

type vector struct {
	x int
	y int
}

type robot struct {
	pos  vector
	velo vector
}

func parseRobots(scanner *bufio.Scanner) []robot {
	robots := []robot{}
	for scanner.Scan() {
		var px, py, vx, vy int
		fmt.Sscanf(scanner.Text(), "p=%d,%d v=%d,%d", &px, &py, &vx, &vy)
		robots = append(robots, robot{pos: vector{x: px, y: py}, velo: vector{x: vx, y: vy}})
	}
	return robots
}

func updateRobotPos(robot robot) robot {
	newRobot := robot
	newX, newY := robot.pos.x+robot.velo.x, robot.pos.y+robot.velo.y
	xOffset, yOffset := X_LIMIT+1, Y_LIMIT+1

	switch {
	case newX > X_LIMIT:
		newRobot.pos.x = newX - xOffset
	case newX < 0:
		newRobot.pos.x = newX + xOffset
	default:
		newRobot.pos.x = newX
	}

	switch {
	case newY > Y_LIMIT:
		newRobot.pos.y = newY - yOffset
	case newY < 0:
		newRobot.pos.y = newY + yOffset
	default:
		newRobot.pos.y = newY
	}

	return newRobot
}

func countQuadrant(robots []robot, lower vector, upper vector) int {
	count := 0
	for _, robot := range robots {
		if robot.pos.x >= lower.x &&
			robot.pos.x <= upper.x &&
			robot.pos.y >= lower.y &&
			robot.pos.y <= upper.y {
			count += 1
		}
	}
	return count
}

func floor(robots []robot) [Y_LIMIT + 1][X_LIMIT + 1]int {
	floor := [Y_LIMIT + 1][X_LIMIT + 1]int{}
	for _, robot := range robots {
		floor[robot.pos.y][robot.pos.x] += 1
	}
	return floor
}

func calcSafetyFactor(robots []robot) int {
	midX, midY := X_LIMIT/2, Y_LIMIT/2
	topLeft := countQuadrant(robots, vector{x: 0, y: 0}, vector{x: midX - 1, y: midY - 1})
	topRight := countQuadrant(robots, vector{x: midX + 1, y: 0}, vector{x: X_LIMIT, y: midY - 1})
	botLeft := countQuadrant(robots, vector{x: 0, y: midY + 1}, vector{x: midX - 1, y: Y_LIMIT})
	botRight := countQuadrant(robots, vector{x: midX + 1, y: midY + 1}, vector{x: X_LIMIT, y: Y_LIMIT})
	return topLeft * topRight * botLeft * botRight
}

func findSafetyFactor(robots []robot, time int) int {
	for i := 0; i < time; i++ {
		for j := range robots {
			robots[j] = updateRobotPos(robots[j])
		}
	}

	return calcSafetyFactor(robots)
}

func christmasTreeMatch(robots []robot) bool {
	posMap := map[vector]bool{}
	for _, robot := range robots {
		posMap[robot.pos] = true
	}

	treePattern := []vector{{x: 0, y: 1}, {x: 1, y: 1}, {x: -1, y: 1}, {x: -2, y: 2}, {x: -1, y: 2}, {x: 0, y: 2}, {x: 1, y: 2}, {x: 2, y: 2}}
	for _, robot := range robots {
		match := true
		for _, pos := range treePattern {
			match = match && posMap[vector{x: robot.pos.x + pos.x, y: robot.pos.y + pos.y}]
		}
		if match {
			return true
		}
	}
	return false
}

func prettyPrintFloor(floor [Y_LIMIT + 1][X_LIMIT + 1]int) {
	for _, row := range floor {
		str := ""
		for _, n := range row {
			if n == 0 {
				str += " "
			} else {
				str += "*"
			}
		}
		fmt.Println(str)
	}
}

func findChristmasTree(robots []robot, time int) {
	for i := 1; i < time; i++ {
		for j := range robots {
			robots[j] = updateRobotPos(robots[j])
		}

		if christmasTreeMatch(robots) {
			fmt.Println(i)
			prettyPrintFloor(floor(robots))
		}
	}
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
	robots := parseRobots(scanner)
	robotsCopy := append([]robot{}, robots...)

	fmt.Println("Part 1:", findSafetyFactor(robots, 100))
	fmt.Println("Part 2:")
	findChristmasTree(robotsCopy, 10000)

	log.Printf("Time taken: %s", time.Since(start))
}

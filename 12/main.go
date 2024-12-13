package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type point struct {
	x int
	y int
}

type plant struct {
	t       rune
	visited bool
}

func parseGardens(scanner *bufio.Scanner) map[point]*plant {
	gardens := map[point]*plant{}
	for y := 0; scanner.Scan(); y++ {
		for x, c := range scanner.Text() {
			gardens[point{y: y, x: x}] = &plant{t: c, visited: false}
		}
	}
	return gardens
}

func neighbors(p point) []point {
	return []point{
		{x: p.x - 1, y: p.y},
		{x: p.x + 1, y: p.y},
		{x: p.x, y: p.y - 1},
		{x: p.x, y: p.y + 1},
	}
}

func checkPointToVisit(curr point, next point, gardens map[point]*plant) bool {
	return gardens[next] != nil &&
		gardens[next].t == gardens[curr].t &&
		!gardens[next].visited
}

func visit(p point, gardens map[point]*plant, area []point) []point {
	area = append(area, p)
	gardens[p].visited = true
	neighbors := neighbors(p)

	for _, n := range neighbors {
		if checkPointToVisit(p, n, gardens) {
			area = visit(n, gardens, area)
		}
	}

	return area
}

func findAreas(gardens map[point]*plant) [][]point {
	areas := [][]point{}

	for k, v := range gardens {
		if !v.visited {
			areas = append(areas, visit(k, gardens, []point{}))
		}
	}

	return areas
}

func findPerimeter(area []point, gardens map[point]*plant) int {
	perimeter := 0
	for _, p := range area {
		neighbors := neighbors(p)
		for _, n := range neighbors {
			if gardens[n] == nil || gardens[n].t != gardens[p].t {
				perimeter += 1
			}
		}
	}

	return perimeter
}

func calcFenceCost(gardens map[point]*plant) int {
	areas := findAreas(gardens)

	total := 0
	for _, area := range areas {
		total += len(area) * findPerimeter(area, gardens)
	}

	return total
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
	gardens := parseGardens(scanner)

	fmt.Println("Part 1:", calcFenceCost(gardens))

	log.Printf("Time taken: %s", time.Since(start))
}

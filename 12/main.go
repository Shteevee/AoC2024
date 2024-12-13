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

func findPerimeterLength(area []point, gardens map[point]*plant) int {
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

func calcFenceCost(areas [][]point, gardens map[point]*plant) int {
	total := 0
	for _, area := range areas {
		total += len(area) * findPerimeterLength(area, gardens)
	}

	return total
}

func findAreaSet(area []point, gardens map[point]*plant) map[point]bool {
	areaSet := map[point]bool{}
	for _, p := range area {
		areaSet[p] = true
	}

	return areaSet
}

func countCorners(area map[point]bool) int {
	corners := 0

	for p := range area {
		left, right := point{x: p.x - 1, y: p.y}, point{x: p.x + 1, y: p.y}
		up, down := point{x: p.x, y: p.y - 1}, point{x: p.x, y: p.y + 1}
		topLeft, topRight := point{x: p.x - 1, y: p.y - 1}, point{x: p.x + 1, y: p.y - 1}
		botLeft, botRight := point{x: p.x - 1, y: p.y + 1}, point{x: p.x + 1, y: p.y + 1}

		if area[left] && area[up] && !area[topLeft] {
			corners += 1
		}

		if area[right] && area[up] && !area[topRight] {
			corners += 1
		}

		if area[left] && area[down] && !area[botLeft] {
			corners += 1
		}

		if area[right] && area[down] && !area[botRight] {
			corners += 1
		}

		if !area[right] && !area[down] {
			corners += 1
		}

		if !area[left] && !area[down] {
			corners += 1
		}

		if !area[right] && !area[up] {
			corners += 1
		}

		if !area[left] && !area[up] {
			corners += 1
		}
	}

	return corners
}

func calcDiscountFenceCost(areas [][]point, gardens map[point]*plant) int {
	total := 0
	for _, area := range areas {
		areaSet := findAreaSet(area, gardens)
		total += len(area) * countCorners(areaSet)
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
	areas := findAreas(gardens)

	fmt.Println("Part 1:", calcFenceCost(areas, gardens))
	fmt.Println("Part 2:", calcDiscountFenceCost(areas, gardens))

	log.Printf("Time taken: %s", time.Since(start))
}

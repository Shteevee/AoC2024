package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"time"
)

type pos struct {
	x int
	y int
}

type qItem struct {
	pos   pos
	dist  int
	index int
	prev  *qItem
}

const (
	UPPER_MAZE_BOUND = 70
)

type PriorityQueue []*qItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*qItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func parseBytes(scanner *bufio.Scanner) []pos {
	bytes := []pos{}
	for scanner.Scan() {
		p := pos{}
		fmt.Sscanf(scanner.Text(), "%d,%d", &p.x, &p.y)
		bytes = append(bytes, p)
	}

	return bytes
}

func buildWalls(ps []pos, sim int) map[pos]bool {
	walls := map[pos]bool{}
	for i := 0; i < sim; i++ {
		walls[ps[i]] = true
	}

	for x := -1; x <= UPPER_MAZE_BOUND+1; x++ {
		for y := -1; y <= UPPER_MAZE_BOUND+1; y++ {
			if x == -1 || y == -1 || x == UPPER_MAZE_BOUND+1 || y == UPPER_MAZE_BOUND+1 {
				walls[pos{x, y}] = true
			}
		}
	}
	return walls
}

func nextNeighbors(qItem qItem, walls map[pos]bool) []pos {
	x, y := qItem.pos.x, qItem.pos.y

	i := 0
	neighbors := []pos{
		{x - 1, y},
		{x + 1, y},
		{x, y - 1},
		{x, y + 1},
	}
	for _, n := range neighbors {
		if !walls[n] {
			neighbors[i] = n
			i++
		}
	}

	return neighbors[:i]
}

func findLowestScoreThroughMaze(start pos, end pos, walls map[pos]bool) (*qItem, bool) {
	q := &PriorityQueue{&qItem{pos: start}}
	heap.Init(q)
	visited := map[pos]bool{}

	for len(*q) > 0 {
		u := heap.Pop(q).(*qItem)
		switch {
		case visited[u.pos]:
			continue
		case u.pos == end:
			return u, true
		default:
			visited[u.pos] = true
			neighbors := nextNeighbors(*u, walls)
			for _, n := range neighbors {
				heap.Push(q, &qItem{pos: n, dist: u.dist + 1, prev: u})
			}
		}
	}

	return &qItem{}, false
}

func findFirstBlockingByte(bytes []pos) pos {
	lowestBlock := len(bytes) - 1
	highestPass := 0
	n := 0
	for lowestBlock-highestPass != 1 {
		n = (lowestBlock + highestPass) / 2
		_, pass := findLowestScoreThroughMaze(
			pos{0, 0},
			pos{UPPER_MAZE_BOUND, UPPER_MAZE_BOUND},
			buildWalls(bytes, n),
		)

		if pass {
			highestPass = n
		} else {
			lowestBlock = n
		}
	}
	return bytes[lowestBlock-1]
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
	bytes := parseBytes(scanner)
	walls := buildWalls(bytes, 1024)
	end, _ := findLowestScoreThroughMaze(pos{0, 0}, pos{UPPER_MAZE_BOUND, UPPER_MAZE_BOUND}, walls)
	fmt.Println("Part 1:", end.dist)
	fmt.Println("Part 2:", findFirstBlockingByte(bytes))

	log.Printf("Time taken: %s", time.Since(start))
}

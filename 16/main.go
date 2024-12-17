package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"maps"
	"os"
	"time"
)

type pos struct {
	x int
	y int
}

type directedPos struct {
	pos pos
	dir rune
}

type qItem struct {
	pos   pos
	dir   rune
	dist  int
	index int
	prev  *qItem
}

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

func parseMaze(scanner *bufio.Scanner) (pos, pos, map[pos]bool) {
	start, end := pos{}, pos{}
	walls := map[pos]bool{}
	for y := 0; scanner.Scan(); y++ {
		for x, c := range scanner.Text() {
			switch c {
			case '#':
				walls[pos{x, y}] = true
			case 'S':
				start = pos{x, y}
			case 'E':
				end = pos{x, y}
			}
		}
	}

	return start, end, walls
}

func oppositeDir(curr rune, next rune) bool {
	opposite := false
	switch curr {
	case '<':
		opposite = next == '>'
	case '>':
		opposite = next == '<'
	case '^':
		opposite = next == 'v'
	case 'v':
		opposite = next == '^'
	}
	return opposite
}

func nextNeighbors(qItem qItem, walls map[pos]bool) []directedPos {
	x, y := qItem.pos.x, qItem.pos.y

	i := 0
	neighbors := []directedPos{
		{pos: pos{x - 1, y}, dir: '<'},
		{pos: pos{x + 1, y}, dir: '>'},
		{pos: pos{x, y - 1}, dir: '^'},
		{pos: pos{x, y + 1}, dir: 'v'},
	}
	for _, n := range neighbors {
		if !walls[n.pos] && !oppositeDir(qItem.dir, n.dir) {
			neighbors[i] = n
			i++
		}
	}

	return neighbors[:i]
}

func calcDist(prev qItem, next directedPos) int {
	dist := prev.dist + 1
	if prev.dir != next.dir {
		dist += 1000
	}
	return dist
}

func findLowestScoreThroughMaze(start pos, end pos, walls map[pos]bool) (*qItem, bool) {
	q := &PriorityQueue{&qItem{pos: start, dir: '>'}}
	heap.Init(q)
	visited := map[directedPos]bool{}

	for len(*q) > 0 {
		u := heap.Pop(q).(*qItem)
		dPos := directedPos{u.pos, u.dir}
		switch {
		case visited[dPos]:
			continue
		case u.pos == end:
			return u, true
		default:
			visited[dPos] = true
			neighbors := nextNeighbors(*u, walls)
			for _, n := range neighbors {
				dist := calcDist(*u, n)
				heap.Push(q, &qItem{pos: n.pos, dir: n.dir, dist: dist, prev: u})
			}
		}
	}

	return &qItem{}, false
}

func findUniqTilesOfAllBestPaths(start pos, walls map[pos]bool, end *qItem) map[pos]bool {
	uniqTiles := map[pos]bool{}
	shortestPath := end.dist
	path := []pos{}
	next := end
	for next.prev != nil {
		path = append(path, next.pos)
		uniqTiles[next.pos] = true
		next = next.prev
	}

	for _, p := range path {
		newWalls := maps.Clone(walls)
		newWalls[p] = true
		if e, ok := findLowestScoreThroughMaze(start, end.pos, newWalls); ok && e.dist == shortestPath {
			for e.prev != nil {
				uniqTiles[e.pos] = true
				e = e.prev
			}
			uniqTiles[e.pos] = true
		}
	}

	return uniqTiles
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
	s, e, walls := parseMaze(scanner)

	end, _ := findLowestScoreThroughMaze(s, e, walls)
	uniqTiles := findUniqTilesOfAllBestPaths(s, walls, end)
	fmt.Println("Part 1:", end.dist)
	fmt.Println("Part 2:", len(uniqTiles))

	log.Printf("Time taken: %s", time.Since(start))
}

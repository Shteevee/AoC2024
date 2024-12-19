package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

func parseTowels(scanner *bufio.Scanner) (map[rune][][]rune, [][]rune) {
	patterns := map[rune][][]rune{}
	designs := [][]rune{}
	scanner.Scan()
	towelPatterns := strings.Split(scanner.Text(), ", ")
	for _, p := range towelPatterns {
		rs := []rune(p)
		patterns[rs[0]] = append(patterns[rs[0]], rs)
	}
	scanner.Scan()
	for scanner.Scan() {
		designs = append(designs, []rune(scanner.Text()))
	}

	return patterns, designs
}

func designPossible(patterns map[rune][][]rune, design []rune) bool {
	if len(design) == 0 {
		return true
	}
	possible := false
	for _, pattern := range patterns[design[0]] {
		if len(design) >= len(pattern) && slices.Equal(pattern, design[0:len(pattern)]) {
			possible = possible || designPossible(patterns, design[len(pattern):])
		}
	}
	return possible
}

func countPossibleDesigns(patterns map[rune][][]rune, designs [][]rune) int {
	possible := 0
	for _, design := range designs {
		if designPossible(patterns, design) {
			possible += 1
		}
	}
	return possible
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
	patterns, designs := parseTowels(scanner)

	fmt.Println("Part 1:", countPossibleDesigns(patterns, designs))

	log.Printf("Time taken: %s", time.Since(start))
}

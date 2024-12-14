package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const BUTTON_A_COST = 3
const CONVERSION_ERROR = 10000000000000

type vector struct {
	x int
	y int
}

type game struct {
	a     vector
	b     vector
	prize vector
}

func parseVector(line string, separator string) vector {
	split := strings.Split(line, ", ")
	x, _ := strconv.Atoi(strings.Split(split[0], separator)[1])
	y, _ := strconv.Atoi(strings.Split(split[1], separator)[1])
	return vector{x: x, y: y}
}

func parseCraneGames(scanner *bufio.Scanner) []game {
	games := []game{}
	g := game{}
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			games = append(games, g)
			g = game{}
		} else {
			split := strings.Split(scanner.Text(), ": ")
			switch split[0] {
			case "Button A":
				g.a = parseVector(split[1], "+")
			case "Button B":
				g.b = parseVector(split[1], "+")
			case "Prize":
				g.prize = parseVector(split[1], "=")
			}
		}
	}
	games = append(games, g)
	return games
}

func calcTokensForPrize(game game) int {
	ax, ay := game.a.x, game.a.y
	bx, by := game.b.x, game.b.y
	px, py := game.prize.x, game.prize.y

	a := (px*by - py*bx) / (ax*by - ay*bx)
	b := (px - a*ax) / bx

	if (px-a*ax)%bx != 0 || (px*by-py*bx)%(ax*by-ay*bx) != 0 {
		return 0
	}

	return a*BUTTON_A_COST + b
}

func sumMinimumTokens(games []game) int {
	total := 0
	for _, game := range games {
		total += calcTokensForPrize(game)
	}

	return total
}

func sumMinimumTokensWithError(games []game) int {
	total := 0
	for _, game := range games {
		game.prize.x += CONVERSION_ERROR
		game.prize.y += CONVERSION_ERROR
		total += calcTokensForPrize(game)
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
	games := parseCraneGames(scanner)

	fmt.Println("Part 1:", sumMinimumTokens(games))
	fmt.Println("Part 2:", sumMinimumTokensWithError(games))

	log.Printf("Time taken: %s", time.Since(start))
}

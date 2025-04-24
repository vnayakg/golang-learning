package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

const (
	totalGames    = 10
	winningScore  = 100
	diceFaceCount = 6
	passValue     = 1
)

type Player struct {
	id                             string
	totalScore, holdCapacity, wins int
}

func (p *Player) resetTotalScore() {
	p.totalScore = 0
}

func (p *Player) playTurn(roll func() int) {
	turnScore := 0

	for turnScore < p.holdCapacity {
		diceValue := roll()
		if diceValue == passValue {
			return
		}
		turnScore += diceValue
	}
	p.totalScore += turnScore
}

func rollDice() int {
	return rand.IntN(diceFaceCount) + 1
}

func play(playerOne, playerTwo *Player) {
	for range totalGames {
		currentPlayer := playerOne

		for {
			if currentPlayer.totalScore >= winningScore {
				currentPlayer.wins++
				break
			}
			currentPlayer.playTurn(rollDice)
			currentPlayer = switchPlayer(currentPlayer, playerOne, playerTwo)
		}
		playerOne.resetTotalScore()
		playerTwo.resetTotalScore()
	}
}

func switchPlayer(currentPlayer, p1, p2 *Player) *Player {
	if currentPlayer.id == p1.id {
		return p2
	}
	return p1
}

func formatResult(p1, p2 *Player) string {
	return fmt.Sprintf(
		"Holding at %d wins: %d/%d (%.1f%%) vs Holding at %d wins: %d/%d (%.1f%%)",
		p1.holdCapacity, p1.wins, totalGames, float32(p1.wins)*100/totalGames,
		p2.holdCapacity, p2.wins, totalGames, float32(p2.wins)*100/totalGames,
	)
}

func parseArgs() (int, int, error) {
	if len(os.Args) < 3 {
		return 0, 0, fmt.Errorf("missing command line arguments")
	}

	p1Cap, err1 := strconv.Atoi(os.Args[1])
	p2Cap, err2 := strconv.Atoi(os.Args[2])

	if err1 != nil {
		return 0, 0, fmt.Errorf("invalid hold capacity for player 1: %v", os.Args[1])
	}
	if err2 != nil {
		return 0, 0, fmt.Errorf("invalid hold capacity for player 2: %v", os.Args[2])
	}
	if p1Cap < 1 || p2Cap < 1 {
		return 0, 0, fmt.Errorf("hold capacities must be at least 1")
	}

	return p1Cap, p2Cap, nil
}

func main() {

	p1Cap, p2Cap, err := parseArgs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	playerOne := Player{id: "1", holdCapacity: p1Cap}
	playerTwo := Player{id: "2", holdCapacity: p2Cap}

	play(&playerOne, &playerTwo)
	fmt.Println(formatResult(&playerOne, &playerTwo))
}

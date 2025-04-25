package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
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

var playFixedStrategyAgainstFixedStrategy = func(p1HoldCapacity int, p2HoldCapacity int) (*Player, *Player) {
	playerOne := Player{id: "1", holdCapacity: p1HoldCapacity}
	playerTwo := Player{id: "2", holdCapacity: p2HoldCapacity}

	play(&playerOne, &playerTwo)
	return &playerOne, &playerTwo
}

var playFixedStrategyAgainstVariableStrategy = func(p1HoldCapacity int, p2HoldRange Range) []*[2]*Player {
	var results []*[2]*Player
	for p2HoldCapacity := p2HoldRange.Start; p2HoldCapacity <= p2HoldRange.End; p2HoldCapacity++ {
		if p1HoldCapacity == p2HoldCapacity {
			continue
		}

		playerOne := Player{id: "1", holdCapacity: p1HoldCapacity}
		playerTwo := Player{id: "2", holdCapacity: p2HoldCapacity}

		play(&playerOne, &playerTwo)
		results = append(results, &[2]*Player{&playerOne, &playerTwo})
	}
	return results
}

var playVariableStrategyAgainstVariableStrategy = func(p1HoldRange Range, p2HoldRange Range) []*[2]*Player {
	var results []*[2]*Player

	for p1HoldCapacity := p1HoldRange.Start; p1HoldCapacity <= p1HoldRange.End; p1HoldCapacity++ {
		result := playFixedStrategyAgainstVariableStrategy(p1HoldCapacity, p2HoldRange)
		results = append(results, result...)
	}
	return results
}

var play = func(playerOne, playerTwo *Player) {
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

type Range struct {
	Start int
	End   int
}

type ParsedArg struct {
	IsRange bool
	Value   int
	Range   Range
}

func parseArg(arg string) (*ParsedArg, error) {
	if strings.Contains(arg, "-") {
		parts := strings.Split(arg, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", arg)
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid range numbers in: %s", arg)
		}
		if start <= 0 || end <= 0 || start > 100 || end > 100 {
			return nil, fmt.Errorf("invalid range numbers in: %s", arg)
		}
		return &ParsedArg{
			IsRange: true,
			Range:   Range{Start: start, End: end},
		}, nil
	}

	val, err := strconv.Atoi(arg)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", arg)
	}
	if val <= 0 || val > 100 {
		return nil, fmt.Errorf("invalid range numbers in: %s", arg)
	}
	return &ParsedArg{
		IsRange: false,
		Value:   val,
	}, nil
}

func playStrategies(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: ./pig <number|range> <number|range>")
	}

	arg1, err1 := parseArg(args[0])
	arg2, err2 := parseArg(args[1])

	if err1 != nil || err2 != nil {
		return fmt.Errorf("error: %v, %v", err1, err2)
	}

	switch {
	case !arg1.IsRange && !arg2.IsRange:
		playerOne, playerTwo := playFixedStrategyAgainstFixedStrategy(arg1.Value, arg2.Value)
		fmt.Println(formatResult(playerOne, playerTwo))
	case !arg1.IsRange && arg2.IsRange:
		results := playFixedStrategyAgainstVariableStrategy(arg1.Value, arg2.Range)
		for _, result := range results {
			fmt.Println(formatResult(result[0], result[1]))
		}
	case arg1.IsRange && arg2.IsRange:
		results := playVariableStrategyAgainstVariableStrategy(arg1.Range, arg2.Range)
		fmt.Println(formatVariableStrategyResults(results, arg1.Range))
	}
	return nil
}
func formatVariableStrategyResults(results []*[2]*Player, p1HoldRange Range) []string {
	resultsByP1 := make(map[int][]*[2]*Player)

	for _, result := range results {
		p1HoldCapacity := result[0].holdCapacity
		resultsByP1[p1HoldCapacity] = append(resultsByP1[p1HoldCapacity], result)
	}
	var formattedResults []string
	for p1HoldCapacity := p1HoldRange.Start; p1HoldCapacity <= p1HoldRange.End; p1HoldCapacity++ {
		resultsForP1 := resultsByP1[p1HoldCapacity]
		if len(resultsForP1) == 0 {
			continue
		}
		var playerOneWins int = 0
		var playerTwoWins int = 0

		for _, res := range resultsForP1 {
			playerOneWins += res[0].wins
			playerTwoWins += res[1].wins
		}

		totalGames := playerOneWins + playerTwoWins
		formattedResult := fmt.Sprintf("Result: Wins, losses staying at k = %v: %v/%v (%0.1f%%), %v/%v (%.1f%%)\n",
			p1HoldCapacity, playerOneWins, totalGames, float32(playerOneWins)*100/float32(totalGames),
			playerTwoWins, totalGames, float32(playerTwoWins)*100/float32(totalGames))

		formattedResults = append(formattedResults, formattedResult)
	}
	return formattedResults
}

func main() {
	if err := playStrategies(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

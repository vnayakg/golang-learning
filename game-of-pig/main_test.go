package main

import (
	"os"
	"testing"
)

func mockRoller(rolls []int) func() int {
	index := 0
	return func() int {
		if index >= len(rolls) {
			return 1
		}
		rollValue := rolls[index]
		index++
		return rollValue
	}
}

func TestPlayTurn_ShouldPassWhenDiceValueIsOne(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{1}))

	if player.totalScore != 0 {
		t.Errorf("expected total score to be 0, got %v", player.totalScore)
	}
}

func TestPlayTurn_ShouldDiscardAccumulatedScoreWhenDiceValueIsOne(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 1}))

	if player.totalScore != 0 {
		t.Errorf("expected total score to be 0, got %v", player.totalScore)
	}
}

func TestPlayTurn_ShouldAccumulateScoreTillExactHoldCapacity(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 4}))

	if player.totalScore != 10 {
		t.Errorf("expected total score to be 10, got %v", player.totalScore)
	}
}

func TestPlayTurn_ShouldAccumulateScoreTillHoldCapacity(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 6}))

	if player.totalScore != 12 {
		t.Errorf("expected total score to be 12, got %v", player.totalScore)
	}
}

func TestSwitchPlayer_ShouldSwitchCurrentPlayer(t *testing.T) {
	playerOne := Player{id: "1"}
	playerTwo := Player{id: "2"}

	result := switchPlayer(&playerOne, &playerOne, &playerTwo)

	if result != &playerTwo {
		t.Errorf("expected: %v, got %v", playerTwo, result)
	}
}

func TestResetTotalScore_ShouldResetScoreOfPlayer(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10, totalScore: 99}

	player.resetTotalScore()

	if player.totalScore != 0 {
		t.Errorf("expected totalScore: 0, got %v", player.totalScore)
	}
}

func TestFormatResult(t *testing.T) {
	playerOne := &Player{id: "1", holdCapacity: 20, wins: 6}
	playerTwo := &Player{id: "2", holdCapacity: 10, wins: 4}

	result := formatResult(playerOne, playerTwo)
	expected := "Holding at 20 wins: 6/10 (60.0%) vs Holding at 10 wins: 4/10 (40.0%)"

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestPlay_ShouldPlayGivenNumberOfGames(t *testing.T) {
	playerOne := &Player{id: "1", holdCapacity: 20, wins: 0}
	playerTwo := &Player{id: "2", holdCapacity: 10, wins: 0}

	play(playerOne, playerTwo)

	totalGames := playerOne.wins + playerTwo.wins
	if totalGames != 10 {
		t.Errorf("expected %v, got %v", 10, totalGames)
	}
}

func TestParseArgs_ShouldNotGiveErrorForValidArgs(t *testing.T) {
	os.Args = []string{"cmd", "5", "8"}

	p1, p2, err := parseArgs()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if p1 != 5 || p2 != 8 {
		t.Errorf("Expected 5 and 8, got %d and %d", p1, p2)
	}
}

func TestParseArgs_ShouldGiveErrorForInvalidArgs(t *testing.T) {
	cases := [][]string{
		{"program"},
		{"program", "5"},
		{"program", "x", "10"},
		{"program", "5", "x"},
		{"program", "0", "10"},
		{"program", "5", "0"},
	}

	for _, args := range cases {
		os.Args = args

		_, _, err := parseArgs()
		if err == nil {
			t.Errorf("Expected error for args: %v", args)
		}
	}
}

package main

import (
	"reflect"
	"testing"
)

func assertEqual[T comparable](t *testing.T, got, want T, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

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

	assertEqual(t, player.totalScore, 0, "expected total score to be 0")
}

func TestPlayTurn_ShouldDiscardAccumulatedScoreWhenDiceValueIsOne(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 1}))

	assertEqual(t, player.totalScore, 0, "expected total score to be 0")
}

func TestPlayTurn_ShouldAccumulateScoreTillExactHoldCapacity(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 4}))

	assertEqual(t, player.totalScore, 10, "expected total score to be 10")
}

func TestPlayTurn_ShouldAccumulateScoreTillHoldCapacity(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10}

	player.playTurn(mockRoller([]int{6, 6}))

	assertEqual(t, player.totalScore, 12, "expected total score to be 12")
}

func TestSwitchPlayer_ShouldSwitchCurrentPlayer(t *testing.T) {
	playerOne := Player{id: "1"}
	playerTwo := Player{id: "2"}

	result := switchPlayer(&playerOne, &playerOne, &playerTwo)

	assertEqual(t, result, &playerTwo, "shouldSwitchPlayer")
}

func TestResetTotalScore_ShouldResetScoreOfPlayer(t *testing.T) {
	player := Player{id: "1", holdCapacity: 10, totalScore: 99}

	player.resetTotalScore()

	assertEqual(t, player.totalScore, 0, "expected totalScore: 0")
}

func TestFormatResult(t *testing.T) {
	playerOne := &Player{id: "1", holdCapacity: 20, wins: 6}
	playerTwo := &Player{id: "2", holdCapacity: 10, wins: 4}

	result := formatResult(playerOne, playerTwo)
	expected := "Holding at 20 wins: 6/10 (60.0%) vs Holding at 10 wins: 4/10 (40.0%)"

	assertEqual(t, result, expected, "")
}

func TestPlay_ShouldPlayGivenNumberOfGames(t *testing.T) {
	playerOne := &Player{id: "1", holdCapacity: 20, wins: 0}
	playerTwo := &Player{id: "2", holdCapacity: 10, wins: 0}

	play(playerOne, playerTwo)
	totalGames := playerOne.wins + playerTwo.wins

	assertEqual(t, totalGames, 10, "")
}

func TestParseArgs_ShouldParseForValidArgs(t *testing.T) {
	testCases := []struct {
		input    string
		expected *ParsedArg
	}{
		{"5", &ParsedArg{Value: 5, IsRange: false}},
		{"5-100", &ParsedArg{IsRange: true, Range: Range{5, 100}}},
	}

	for _, testCase := range testCases {
		actual, err := parseArg(testCase.input)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, actual)
		}
	}
}

func TestParseArgs_ShouldGiveErrorForInvalidArgs(t *testing.T) {
	testCases := []struct {
		input string
	}{{"a"}, {"0"}, {"0-1"}, {"1-"}, {"-"}, {"-10"}, {"0-100"}, {"1-101"}}

	for _, testCase := range testCases {
		if _, err := parseArg(testCase.input); err == nil {
			t.Errorf("Expected error for %v, got %v", testCase.input, err)
		}
	}
}

func TestRun_FixedVsFixed(t *testing.T) {
	called := false
	//cleanup mocks
	originalPlayFixedStrategyAgainstFixedStrategy := playFixedStrategyAgainstFixedStrategy
	defer func() { playFixedStrategyAgainstFixedStrategy = originalPlayFixedStrategyAgainstFixedStrategy }()
	playFixedStrategyAgainstFixedStrategy = func(a, b int) (*Player, *Player) {
		called = true
		if a != 5 || b != 6 {
			t.Errorf("unexpected arguments: got %d and %d", a, b)
		}
		return &Player{}, &Player{}
	}

	err := playStrategies([]string{"5", "6"})

	if err != nil || !called {
		t.Error("playFixedStrategyAgainstFixedStrategy was not called")
	}
}

func TestRun_FixedVsVariable(t *testing.T) {
	called := false
	//cleanup mocks
	originalPlayFixedStrategyAgainstVariableStrategy := playFixedStrategyAgainstVariableStrategy
	defer func() { playFixedStrategyAgainstVariableStrategy = originalPlayFixedStrategyAgainstVariableStrategy }()
	playFixedStrategyAgainstVariableStrategy = func(a int, b Range) []*[2]*Player {
		called = true
		if (a != 5 || b != Range{1, 100}) {
			t.Errorf("unexpected arguments: got %d and %d", a, b)
		}
		return []*[2]*Player{}
	}

	err := playStrategies([]string{"5", "1-100"})

	if err != nil || !called {
		t.Error("playFixedStrategyAgainstFixedStrategy was not called")
	}
}

func TestRun_VariableVsVariable(t *testing.T) {
	called := false
	//cleanup mocks
	originalPlayVariableStrategyAgainstVariableStrategy := playVariableStrategyAgainstVariableStrategy
	defer func() {
		playVariableStrategyAgainstVariableStrategy = originalPlayVariableStrategyAgainstVariableStrategy
	}()
	playVariableStrategyAgainstVariableStrategy = func(a, b Range) []*[2]*Player {
		called = true
		if (a != Range{1, 100} || b != Range{1, 100}) {
			t.Errorf("unexpected arguments: got %d and %d", a, b)
		}
		return []*[2]*Player{}
	}

	err := playStrategies([]string{"1-100", "1-100"})

	if err != nil || !called {
		t.Error("playVariableStrategyAgainstFixedStrategy was not called")
	}
}

func TestPlayFixedVsFixed_ReturnsCorrectPlayers(t *testing.T) {
	called := false
	//cleanup mocks
	originalPlay := play
	defer func() { play = originalPlay }()
	play = func(p1, p2 *Player) {
		called = true
		p1.wins = 1
		p2.wins = 99
	}

	p1, p2 := playFixedStrategyAgainstFixedStrategy(5, 6)

	if !called {
		t.Fatal("play was not called")
	}
	assertEqual(t, p1.wins, 1, "unexpected hold capacities")
	assertEqual(t, p2.wins, 99, "unexpected hold capacities")
}

func TestPlayFixedVsVariable_ReturnsCorrectPlayerPairs(t *testing.T) {
	var calledWith [][2]int
	//cleanup mocks
	originalPlay := play
	defer func() { play = originalPlay }()
	play = func(p1, p2 *Player) {
		calledWith = append(calledWith, [2]int{p1.holdCapacity, p2.holdCapacity})
	}

	pairs := playFixedStrategyAgainstVariableStrategy(5, Range{Start: 1, End: 2})

	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	expected := [][2]int{
		{5, 1},
		{5, 2},
	}
	for i, pair := range pairs {
		got := [2]int{pair[0].holdCapacity, pair[1].holdCapacity}
		if got != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], got)
		}
	}
	if !reflect.DeepEqual(calledWith, expected) {
		t.Errorf("play called with unexpected values. got %v, want %v", calledWith, expected)
	}
}

func TestPlayVariableVsVariable_ReturnsCorrectResult(t *testing.T) {
	var calledWith [][2]int
	//cleanup mocks
	originalPlay := play
	defer func() { play = originalPlay }()
	play = func(p1, p2 *Player) {
		calledWith = append(calledWith, [2]int{p1.holdCapacity, p2.holdCapacity})
	}

	pairs := playVariableStrategyAgainstVariableStrategy(Range{1, 2}, Range{1, 3})

	if len(pairs) != 4 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	expected := [][2]int{
		{1, 2},
		{1, 3},
		{2, 1},
		{2, 3},
	}
	for i, pair := range pairs {
		got := [2]int{pair[0].holdCapacity, pair[1].holdCapacity}
		if got != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], got)
		}
	}
	if !reflect.DeepEqual(calledWith, expected) {
		t.Errorf("play called with unexpected values. got %v, want %v", calledWith, expected)
	}
}

func TestFormatVariableStrategyResults(t *testing.T) {
	player1 := &Player{id: "1", holdCapacity: 1, wins: 7}
	player2 := &Player{id: "2", holdCapacity: 2, wins: 3}
	player3 := &Player{id: "1", holdCapacity: 1, wins: 6}
	player4 := &Player{id: "2", holdCapacity: 3, wins: 4}

	results := []*[2]*Player{
		{player1, player2},
		{player3, player4},
	}

	statStrings := formatVariableStrategyResults(results, Range{Start: 1, End: 1})

	if len(statStrings) != 1 {
		t.Errorf("Expected 1 stat string, got %d", len(statStrings))
	}

	expectedString := "Result: Wins, losses staying at k = 1: 13/20 (65.0%), 7/20 (35.0%)\n"
	if statStrings[0] != expectedString {
		t.Errorf("Expected string: %s, got: %s", expectedString, statStrings[0])
	}
}

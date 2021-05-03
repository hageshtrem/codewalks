package pig

import "math/rand"

// A score includes scores accumulated in previous turns for each player,
// as well as the points scored by the current player in this turn.
type Score struct {
	Player, Opponent, ThisTurn uint
}

// An Action is a player`s action
type Action int

//go:generate stringer -type=Action
const (
	// Roll one more time
	Roll Action = iota + 1
	// Stay and pass the dice
	Stay
)

// DiceValue is possible dice value
type DiceValue int

// Possible dice value
//go:generate stringer -type=DiceValue
const (
	One DiceValue = iota + 1
	Two
	Three
	Four
	Five
	Six
)

// Dice is something that can return a DiceValue.
type Dice interface { // or type DiceRoll func() DiceValue
	// Roll returnsa DiceValue.
	Roll() DiceValue
}

type dice struct{}

func (dice) Roll() DiceValue {
	return DiceValue(rand.Intn(6) + 1)
}

// NewRandomDice returns a dice with random behaviour.
func NewRandomDice() Dice {
	return dice{}
}

// Game is a pig game.
type Game struct {
	dice Dice
	win  uint
}

// NewGame receives a dice instance and a win score and returns a Game.
func NewGame(dice Dice, win uint) Game {
	return Game{dice, win}
}

// roll returns the (result, turnIsOver) outcome of simulating a die roll.
// If the roll value is 1, then thisTurn score is abandoned, and the players'
// roles swap.  Otherwise, the roll value is added to thisTurn.
func (g Game) roll(s Score) (Score, bool) {
	outcome := uint(g.dice.Roll())
	if outcome == 1 {
		return Score{s.Opponent, s.Player, 0}, true
	}
	return Score{s.Player, s.Opponent, s.ThisTurn + outcome}, false
}

// stay returns the result outcome of staying.
// thisTurn score is added to the player's score, and the players' roles swap.
func (Game) stay(s Score) Score {
	return Score{s.Opponent, s.Player + s.ThisTurn, 0}
}

// Play runs a game and returns a winner.
func (g Game) Play(player1, player2 Player) Player {
	players := []Player{player1, player2}
	currentPlayer := rand.Intn(2) // Randomly decide who plays first
	var score Score
	var turnIsOver bool
	for score.Player+score.ThisTurn < g.win {
		switch players[currentPlayer].MakeChoice(score) {
		case Roll:
			if score, turnIsOver = g.roll(score); turnIsOver {
				currentPlayer = (currentPlayer + 1) % 2
			}
		case Stay:
			score = g.stay(score)
		}
	}
	return players[currentPlayer]
}

// Player makes a choice in the game by the Score.
type Player struct {
	Name     string
	strategy Strategy
}

// NewPlayer returns a Player that will play with Strategy.
func NewPlayer(name string, strategy Strategy) Player {
	return Player{name, strategy}
}

// MakeChoice returnc an Action.
func (p Player) MakeChoice(score Score) Action {
	return p.strategy(score)
}

// Strategy chooses an action for any given score.
type Strategy func(Score) Action

// StayAtK returns a strategy that rolls until ThisTurn is at least k, then stays.
func StayAtK(k uint) Strategy {
	return func(score Score) Action {
		if score.ThisTurn < k {
			return Roll
		}
		return Stay
	}
}

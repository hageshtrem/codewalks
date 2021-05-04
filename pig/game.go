package pig

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

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
type Game interface {
	Play(player1, player2 Player) Player
}

type game struct {
	dice Dice
	win  uint
}

// NewGame receives a dice instance and a win score and returns a Game.
func NewGame(dice Dice, win uint) Game {
	return game{dice, win}
}

// roll returns the (result, turnIsOver) outcome of simulating a die roll.
// If the roll value is 1, then thisTurn score is abandoned, and the players'
// roles swap.  Otherwise, the roll value is added to thisTurn.
func (g game) roll(s Score) (Score, bool) {
	outcome := uint(g.dice.Roll())
	if outcome == uint(One) {
		return Score{Player: s.Opponent, Opponent: s.Player, ThisTurn: 0}, true
	}
	return Score{Player: s.Player, Opponent: s.Opponent, ThisTurn: s.ThisTurn + outcome}, false
}

// stay returns the result outcome of staying.
// thisTurn score is added to the player's score, and the players' roles swap.
func (game) stay(s Score) Score {
	return Score{Player: s.Opponent, Opponent: s.Player + s.ThisTurn, ThisTurn: 0}
}

// Play runs a game and returns a winner. First player plays first.
func (g game) Play(player1, player2 Player) Player {
	players := coupleOfPlayers{player1, player2}
	currentPlayer := players.first()
	var score Score
	var turnIsOver bool
	for score.Player+score.ThisTurn < g.win {
		switch currentPlayer.MakeChoice(score) {
		case Roll:
			if score, turnIsOver = g.roll(score); turnIsOver {
				currentPlayer = players.swap(currentPlayer)
			}
		case Stay:
			score = g.stay(score)
			currentPlayer = players.swap(currentPlayer)
		}
	}
	return currentPlayer
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

type coupleOfPlayers struct {
	p1, p2 Player
}

func (c coupleOfPlayers) first() Player {
	return c.p1
}

func (c coupleOfPlayers) swap(p Player) Player {
	if p.Name != c.p1.Name {
		return c.p1
	}
	return c.p2
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

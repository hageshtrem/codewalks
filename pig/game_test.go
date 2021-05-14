package pig_test

import (
	"reflect"
	"testing"

	"github.com/hageshtrem/codewalks/pig"
	"github.com/hageshtrem/codewalks/pig/tournament"
)

// The winning score in a game of Pig
const win = 100

// sixToOneDice is implementation of the pig.Dice with a predicted roll result.
// Each Roll call returns values in the sequence 6, 5, 4, 3, 2, 1, 6, ...
type sixToOneDice int

func (d *sixToOneDice) Roll() pig.DiceValue {
	val := pig.DiceValue(((*d + 5) % 6) + 1)
	*d = (*d - 1) % 6
	return val
}

func BenchmarkSixToOneDice(b *testing.B) {
	var dice sixToOneDice
	var v pig.DiceValue
	for i := 0; i < b.N; i++ {
		v = dice.Roll()
	}
	b.Log(v)
}

func TestSixToOneDice(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
		values     []pig.DiceValue
	}{
		{
			name:       "3",
			iterations: 3,
			values:     []pig.DiceValue{pig.Six, pig.Five, pig.Four},
		},
		{
			name:       "7",
			iterations: 7,
			values:     []pig.DiceValue{pig.Six, pig.Five, pig.Four, pig.Three, pig.Two, pig.One, pig.Six},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var dice sixToOneDice
			for i := 0; i < tt.iterations; i++ {
				if v := dice.Roll(); v != tt.values[i] {
					t.Errorf("want: %v, got: %v", tt.values[i], v)
				}
			}
		})
	}
}

func TestGame(t *testing.T) {
	// In both competitions "20 vs 10" and "10 vs 20", player 1 wins because he plays first.
	tests := []struct {
		name             string
		player1, player2 pig.Player
		winnerName       string
	}{
		{
			name:       "20 vs 10",
			player1:    pig.NewPlayer("p1", pig.StayAtK(20)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(10)),
			winnerName: "p1",
		},
		{
			name:       "10 vs 20",
			player1:    pig.NewPlayer("p1", pig.StayAtK(10)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(20)),
			winnerName: "p1",
		},
		{
			name:       "1 vs 3",
			player1:    pig.NewPlayer("p1", pig.StayAtK(1)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(3)),
			winnerName: "p1",
		},
		{
			name:       "3 vs 1",
			player1:    pig.NewPlayer("p1", pig.StayAtK(3)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(1)),
			winnerName: "p2",
		},
		{
			name:       "1 vs 2",
			player1:    pig.NewPlayer("p1", pig.StayAtK(1)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(2)),
			winnerName: "p1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dice := new(sixToOneDice)
			game := pig.NewGame(dice, win)
			if winner := game.Play(tt.player1, tt.player2); winner.Name != tt.winnerName {
				t.Errorf("winner: %s, but must be: %s", winner.Name, tt.winnerName)
			}
		})
	}
}

func TestStayAtK(t *testing.T) {
	tests := []struct {
		name   string
		stayAt uint
		want   []pig.Action
	}{
		{
			name:   "StayAt 1",
			stayAt: 1,
			want:   []pig.Action{pig.Roll, pig.Stay},
		},
		{
			name:   "StayAt 3",
			stayAt: 3,
			want:   []pig.Action{pig.Roll, pig.Stay},
		},
		{
			name:   "StayAt 16",
			stayAt: 16,
			want:   []pig.Action{pig.Roll, pig.Roll, pig.Roll, pig.Roll, pig.Stay},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			strategy := pig.StayAtK(tt.stayAt)
			var score pig.Score
			var dice sixToOneDice
			for i, ta := range tt.want {
				if action := strategy(score); ta != action {
					t.Errorf("iteration: %d, want: %v, got: %v", i, ta, action)
				}
				score.ThisTurn += uint(dice.Roll())
			}
		})
	}
}

// firstWinGame implements pig.Game.
type firstWinGame struct{}

// Play returns player1 as winner.
func (firstWinGame) Play(player1, _ pig.Player) pig.Player {
	return player1
}

func TestTournament(t *testing.T) {
	var game firstWinGame
	tests := []struct {
		name           string
		gamesPerSeries uint
		players        uint
		wins           []uint
	}{
		{
			name:           "4 players",
			gamesPerSeries: 10,
			players:        4,
			wins:           []uint{30, 20, 10, 0},
		},
		{
			name:           "10 players",
			gamesPerSeries: 10,
			players:        10,
			wins:           []uint{90, 80, 70, 60, 50, 40, 30, 20, 10, 0},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tour := tournament.NewTournament(game, tt.gamesPerSeries)
			wins, _ := tour.RoundRobin(tournament.CreatePlayers(tt.players))
			if !reflect.DeepEqual(tt.wins, wins) {
				t.Errorf("\nwant: %v\ngot:  %v", tt.wins, wins)
			}
		})
	}
}

func BenchmarkTournament(b *testing.B) {
	game := pig.NewGame(pig.NewRandomDice(), win)
	tour := tournament.NewTournament(game, 10)
	players := tournament.CreatePlayers(10)
	var wins []uint
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wins, _ = tour.RoundRobin(players)
	}
	b.Log(wins)
}

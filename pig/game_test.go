package pig_test

import (
	"testing"

	"github.com/hageshtrem/codewalks/pig"
)

// treeTwoOneDice is implementation of the pig.Dice with a predicted roll result.
// Each Roll call returns values in the sequence 3, 2, 1, 3, 2, 1, ...
type treeTwoOneDice int

func (d *treeTwoOneDice) Roll() pig.DiceValue {
	defer func() { *d = (*d - 1) % 3 }()
	return pig.DiceValue(((*d + 2) % 3) + 1)
}

func TestGame(t *testing.T) {
	tests := []struct {
		name             string
		player1, player2 pig.Player
		winnerName       string
	}{
		{
			name:       "2 vs 3",
			player1:    pig.NewPlayer("p1", pig.StayAtK(2)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(3)),
			winnerName: "p2",
		},
		{
			name:       "3 vs 2",
			player1:    pig.NewPlayer("p1", pig.StayAtK(3)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(2)),
			winnerName: "p1",
		},
		{
			name:       "1 vs 3",
			player1:    pig.NewPlayer("p1", pig.StayAtK(1)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(3)),
			winnerName: "p2",
		},
		{
			name:       "3 vs 1",
			player1:    pig.NewPlayer("p1", pig.StayAtK(3)),
			player2:    pig.NewPlayer("p2", pig.StayAtK(1)),
			winnerName: "p1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dice := new(treeTwoOneDice)
			game := pig.NewGame(dice, 10)
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
			want:   []pig.Action{pig.Roll, pig.Roll, pig.Roll, pig.Stay},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			strategy := pig.StayAtK(tt.stayAt)
			var score pig.Score
			for i, ta := range tt.want {
				if action := strategy(score); ta != action {
					t.Errorf("iteration: %d, want: %v, got: %v", i, ta, action)
				}
				score.ThisTurn++
			}
		})
	}
}

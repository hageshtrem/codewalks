package pig

import (
	"strconv"
)

// Tournament represents game tournament.
type Tournament struct {
	game           Game
	gamesPerSeries uint
}

// NewTournamet receives the Game instance and the number of games per series to
// simulate. It returns instance of tournament.
func NewTournament(game Game, gamesPerSeries uint) Tournament {
	return Tournament{game, gamesPerSeries}
}

// RoundRobin simulates a series of games between every pair of players.
// It returns scores for each player and count of games per player.
func (t Tournament) RoundRobin(players []Player) ([]uint, uint) {
	wins := make([]uint, len(players))
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			for k := 0; k < int(t.gamesPerSeries); k++ {
				winner := t.game.Play(players[i], players[j])
				if winner.Name == players[i].Name {
					wins[i]++
				} else {
					wins[j]++
				}
			}
		}
	}
	gamesPerPlayer := t.gamesPerSeries * uint(len(players)-1) // no self play
	return wins, gamesPerPlayer
}

// CreatePlayers returns slice of n Players.
func CreatePlayers(n uint) []Player {
	players := make([]Player, 0, n)
	for i := 0; i < int(n); i++ {
		p := NewPlayer(strconv.Itoa(i), StayAtK(uint(i+1)))
		players = append(players, p)
	}
	return players
}

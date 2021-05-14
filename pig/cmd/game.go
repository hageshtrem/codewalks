package main

import (
	"fmt"

	"github.com/hageshtrem/codewalks/pig"
	"github.com/hageshtrem/codewalks/pig/tournament"
)

const (
	win            = 100 // The winning score in a game of Pig
	gamesPerSeries = 10  // The number of games per series to simulate
)

// ratioString takes a list of integer values and returns a string that lists
// each value and its percentage of the sum of all values.
// e.g., ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"
func ratioString(vals ...uint) string {
	var total uint
	for _, val := range vals {
		total += val
	}
	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}

func main() {
	game := pig.NewGame(pig.NewRandomDice(), win)
	players := tournament.CreatePlayers(win)
	tournament := tournament.NewTournament(game, gamesPerSeries)
	wins, games := tournament.RoundRobin(players)

	for i := range players {
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			i+1, ratioString(wins[i], games-wins[i]))

	}
}

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hageshtrem/codewalks/markov/original"
)

func main() {
	// Register command-line flags.
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := original.NewChain(*prefixLen) // Initialize a new Chain.
	c.Build(os.Stdin)                  // Build chains from standard input.
	text := c.Generate(*numWords)      // Generate text.
	fmt.Println(text)                  // Write text to standard output.
}

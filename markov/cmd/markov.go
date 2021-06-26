package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hageshtrem/codewalks/markov"
)

func main() {
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")
	seed := flag.String("seed", "", "first few words")
	flag.Parse()

	c := markov.NewChain(uint(*prefixLen))
	c.Build(os.Stdin)

	var sb strings.Builder
	if err := c.Generate(&sb, *numWords, *seed); err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(sb.String())
}

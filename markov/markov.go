package markov

import (
	"bufio"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/scanner"
)

// Chain contains model.
type Chain struct {
	corpus    map[string]relationSlice
	prefixLen uint
}

// NewChain returns a new Chain with prefixes of prefixLen tokens.
func NewChain(prefixLen uint) *Chain {
	return &Chain{make(map[string]relationSlice), prefixLen}
}

// Build reads text from the provided Reader and parses it.
func (c *Chain) Build(r io.Reader) {
	defer c.sort() // sort tokens at the end of building

	var s scanner.Scanner
	s.Init(r)

	p := make(prefix, c.prefixLen)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		next := s.TokenText()
		rs := c.corpus[p.String()]
		rs.addValue(next)
		c.corpus[p.String()] = rs
		p.shift(next)
	}
}

// Generate generates text of len n with seed as first words and
// writes it to Writer.
func (c *Chain) Generate(w io.Writer, n int, seed string) error {
	bw := bufio.NewWriter(w)

	seedTokens := splitIntoTokens(seed)
	p := make(prefix, c.prefixLen)
	if len(seedTokens) <= len(p) {
		for _, t := range seedTokens {
			p.shift(t)
		}
	} else {
		for i := range p {
			p.shift(seedTokens[i])
		}
	}

	if _, err := bw.WriteString(p.String()); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if rs, ok := c.corpus[p.String()]; ok {
			if len(rs) == 0 {
				continue
			}
			// get the most common
			v := rs[len(rs)-1].val
			if !isPunctuationMark(v) {
				if _, err := bw.WriteString(" "); err != nil {
					return err
				}
			}
			if _, err := bw.WriteString(v); err != nil {
				return err
			}
			p.shift(v)
		}
	}

	return bw.Flush()
}

func (c *Chain) sort() {
	for k := range c.corpus {
		rs := c.corpus[k]
		rs.Sort()
		c.corpus[k] = rs
	}
}

func isPunctuationMark(v string) bool {
	r := regexp.MustCompile(`\pP`)
	return r.MatchString(v)
}

type relation struct {
	val   string
	count int
}

func newRelation(val string) relation {
	return relation{val, 1}
}

func (r *relation) inc() {
	r.count++
}

type relationSlice []relation

func (rs *relationSlice) addValue(val string) {
	var r *relation
	for i := range *rs {
		if (*rs)[i].val == val {
			r = &(*rs)[i]
		}
	}
	if r != nil {
		r.inc()
	} else {
		*rs = append(*rs, newRelation(val))
	}
}

func (rs *relationSlice) Len() int {
	return len(*rs)
}

func (rs *relationSlice) Less(i, j int) bool {
	return (*rs)[i].count < (*rs)[j].count
}

func (rs *relationSlice) Swap(i, j int) {
	(*rs)[i], (*rs)[j] = (*rs)[j], (*rs)[i]
}

func (rs *relationSlice) Sort() {
	sort.Sort(rs)
}

func splitIntoTokens(seed string) []string {
	var s scanner.Scanner
	s.Init(strings.NewReader(seed))
	tokens := make([]string, 0)
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		tokens = append(tokens, s.TokenText())
	}
	return tokens
}

// prefix is a Markov chain prefix of one or more words.
type prefix []string

// String returns the prefix as a string (for use as a map key).
func (p prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the prefix and appends the given word.
func (p prefix) shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

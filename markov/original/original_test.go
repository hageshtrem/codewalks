package original_test

import (
	"testing"

	"github.com/hageshtrem/codewalks/markov/original"
)

func BenchmarkPrefix(b *testing.B) {
	pref := make(original.Prefix, 2)
	for i := 0; i < b.N; i++ {
		pref.Shift("AABBCCDD")
	}
}

func TestPrefix(t *testing.T) {
	pref := make(original.Prefix, 2)
	pref.Shift("AA")
	pref.Shift("BB")
	pref.Shift("CC")
	s := pref.String()
	t.Log(s)
}

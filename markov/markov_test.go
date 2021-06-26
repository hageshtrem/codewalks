package markov

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

var relationSliceTests = []struct {
	vals   []string
	counts []int
}{
	{
		vals:   []string{"A", "B", "C", "D"},
		counts: []int{1, 1, 1, 1},
	},
	{
		vals:   []string{"AAA", "AAA", "AAA", "B", "CC", "CC", "D"},
		counts: []int{3, 1, 2, 1},
	},
}

func TestRelationSliceAddValue(t *testing.T) {
	for _, tc := range relationSliceTests {
		rs := make(relationSlice, 0, len(tc.vals))
		for _, v := range tc.vals {
			rs.addValue(v)
		}
		if len(rs) != len(tc.counts) {
			t.Errorf("len of relationSlice: %d, must: %d", len(rs), len(tc.counts))
		}
		for i := range rs {
			if rs[i].count != tc.counts[i] {
				t.Errorf("rs: %+v, must contains count: %d", rs[i], tc.counts[i])
			}
		}
	}
}

func TestRelationSliceSort(t *testing.T) {
	for _, tc := range relationSliceTests {
		rs := make(relationSlice, 0, len(tc.vals))
		for _, v := range tc.vals {
			rs.addValue(v)
		}
		var counts sort.IntSlice = tc.counts
		counts.Sort()
		rs.Sort()
		if len(rs) != len(counts) {
			t.Errorf("len of relationSlice: %d, must: %d", len(rs), len(counts))
		}
		for i := range rs {
			if rs[i].count != counts[i] {
				t.Errorf("rs: %+v, must contains count: %d", rs[i], counts[i])
			}
		}
	}
}

var markovChainTests = []struct {
	name         string
	prefixLen    int
	src          string
	partOfCorpus map[string]relationSlice
	genLimit     int
	genSeed      string
	generated    string
}{
	{
		name:      "{I am not a number...} prefix 1",
		prefixLen: 1,
		src:       "I am not a number! I am a free man! I will go away!",
		partOfCorpus: map[string]relationSlice{
			"":   []relation{{val: "I", count: 1}},
			"I":  []relation{{val: "will", count: 1}, {val: "am", count: 2}},
			"a":  []relation{{val: "number!", count: 1}, {val: "free", count: 1}},
			"am": []relation{{val: "not", count: 1}, {val: "a", count: 1}},
		},
		genLimit:  9,
		genSeed:   "I will",
		generated: "I am a free man! I am a free man!",
	},
	{
		name:      "{I am not a number...} prefix 2",
		prefixLen: 2,
		src:       "I am not a number! I am a free man! I will go away!",
		partOfCorpus: map[string]relationSlice{
			" ":     []relation{{val: "I", count: 1}}, // space because of prefix.String
			"not a": []relation{{val: "number!", count: 1}},
			"I am":  []relation{{val: "not", count: 1}, {val: "a", count: 1}},
		},
		genLimit:  9,
		genSeed:   "I will",
		generated: "I will go away!",
	},
}

func TestBuild(t *testing.T) {
	for _, tc := range markovChainTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			chain := NewChain(uint(tc.prefixLen))
			chain.Build(strings.NewReader(tc.src))
			for k, v := range tc.partOfCorpus {
				got := chain.corpus[k]
				if !reflect.DeepEqual(got, v) {
					t.Errorf("prefix: %s, suffix: {want: %v, got %v}", k, v, got)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	for _, tc := range markovChainTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			chain := NewChain(uint(tc.prefixLen))
			chain.Build(strings.NewReader(tc.src))
			var sb strings.Builder
			if err := chain.Generate(&sb, tc.genLimit, tc.genSeed); err != nil {
				t.Error(err)
			}
			result := sb.String()
			if result != tc.generated {
				t.Errorf("want: '%s', got: '%s'", tc.generated, result)
			}
		})
	}
}

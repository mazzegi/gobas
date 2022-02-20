package lex

type bracketPair struct {
	open  rune
	close rune
}

var bracketPairs = []bracketPair{
	{'"', '"'},
	{'\'', '\''},
	{'(', ')'},
	{'{', '}'},
	{'[', ']'},
}

type bracketier struct {
	open map[rune]int
}

func newBracketier() *bracketier {
	return &bracketier{
		open: map[rune]int{},
	}
}

func (b *bracketier) on(r rune) {
	for _, bp := range bracketPairs {
		if r == bp.close {
			if b.open[bp.open] > 0 {
				b.open[bp.open]--
				continue
			}
		}
		if r == bp.open {
			b.open[r]++
		}
	}
}

func (b *bracketier) inBrackets() bool {
	for _, n := range b.open {
		if n > 0 {
			return true
		}
	}
	return false
}

func splitSkipBrackets(s string, sep rune) []string {
	if s == "" {
		return []string{}
	}
	brc := newBracketier()
	sl := []string{""}
	inQuotes := false
	for _, r := range s {
		brc.on(r)

		if r == '"' {
			if inQuotes {
				inQuotes = false
			} else {
				inQuotes = true
			}
		}

		if !brc.inBrackets() && r == sep {
			sl = append(sl, "")
			continue
		}
		sl[len(sl)-1] += string(r)
	}
	return sl
}

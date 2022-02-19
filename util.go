package gobas

import (
	"constraints"
	"strconv"
	"strings"
)

func trimWhite(s string) string {
	return strings.Trim(s, " \r\n\t")
}

func splitOutsideQuotes(s string, sep rune) []string {
	if s == "" {
		return []string{}
	}
	sl := []string{""}
	inQuotes := false
	for _, r := range s {
		if r == '"' {
			if inQuotes {
				inQuotes = false
			} else {
				inQuotes = true
			}
		}

		if !inQuotes && r == sep {
			sl = append(sl, "")
			continue
		}
		sl[len(sl)-1] += string(r)
	}
	return sl
}

func parseInts[T constraints.Integer](s string, sep rune) ([]T, error) {
	var ns []T
	sl := splitOutsideQuotes(s, sep)
	for _, sn := range sl {
		n, err := strconv.ParseInt(trimWhite(sn), 10, 64)
		if err != nil {
			return nil, err
		}
		ns = append(ns, T(n))
	}
	return ns, nil
}

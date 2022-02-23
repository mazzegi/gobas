package gobas

import "strings"

func mustParseInput(raw string) INPUT {
	inp := INPUT{}

	var varsRaw string
	var inQuotes bool
	for i, r := range raw {
		if !inQuotes && r == ' ' {
			continue
		}
		if r == '"' {
			if !inQuotes {
				inQuotes = true
			} else {
				inQuotes = false
			}
			continue
		}
		if inQuotes {
			inp.Msg += string(r)
			continue
		}

		if r == ';' {
			inp.Semicolon = true
			varsRaw = raw[i+1:]
		} else {
			varsRaw = raw[i:]
		}
		break
	}

	inp.Vars = strings.Split(varsRaw, ",")
	for i, v := range inp.Vars {
		inp.Vars[i] = strings.TrimSpace(v)
	}

	return inp
}

package gobas

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type rawLine struct {
	num  uint32
	text string
}

func rawRead(r io.Reader) ([]rawLine, error) {
	var rls []rawLine
	scanner := bufio.NewScanner(r)
	var lno int = -1
	for scanner.Scan() {
		lno++
		ln := strings.Trim(scanner.Text(), " \r\n\t")
		if ln == "" {
			continue
		}
		snum, text, ok := strings.Cut(ln, " ")
		if !ok {
			return nil, errors.Errorf("invalid input-line %d (no line-number separator): %q", lno, ln)
		}
		num, err := strconv.ParseUint(snum, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "scanning line-number in input-line %d", lno)
		}
		rls = append(rls, rawLine{
			num:  uint32(num),
			text: strings.Trim(text, " \r\n\t"),
		})
	}
	return rls, nil
}

func rawReadFile(file string) ([]rawLine, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "open file %q", file)
	}
	defer f.Close()
	return rawRead(f)
}

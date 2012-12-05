package text

import (
	"strings"
)

// An implementation of Plan 9's tokenize (see
// http://plan9.bell-labs.com/magic/man2html/2/getfields)
//
// Tokenize is similar to strings.Fields – an input string is split
// into fields separated by whitespace. Additionally, single quotes
// are interpreted and do not appear in the output. In a quoted part
// of the string, whitespace will not create a new field, and two
// consequtive single quotes will result in one quote in the output.
func Tokenize(s string) (fields []string) {
	if n := do(nil, s); n > 0 { // at first run count fields only
		fields = make([]string, n)
		do(fields, s)
	}
	return
}

func do(fields []string, s string) (n int) {
	var (
		countOnly = fields == nil

		qf      []string
		quoting = false
		wasq    = false

		i0 = -1

		addField = func(f string) {
			if len(qf) == 0 {
				fields[n] = f
			} else {
				if len(f) > 0 {
					qf = append(qf, f)
				}
				fields[n] = strings.Join(qf, "")
				qf = qf[:0]
			}
		}
	)

	for i, r := range s {
		switch r {
		case ' ', '\t', '\r', '\n':
			if !quoting && i0 != -1 {
				if !countOnly {
					addField(s[i0:i])
				}
				n++
				i0 = -1
			}

		case '\'':
			if !quoting {
				if wasq {
					i0--
				}
				if i0 >= 0 && i-i0 > 0 {
					if !countOnly {
						qf = append(qf, s[i0:i])
					}
				}
				i0 = i + 1
				quoting = true
			} else {
				if i-i0 > 0 {
					if !countOnly {
						qf = append(qf, s[i0:i])
					}
				}
				i0 = i + 1
				quoting = false
				wasq = true
				continue
			}
		default:
			if i0 == -1 {
				i0 = i
			}
		}
		wasq = false
	}
	if i0 != -1 {
		if !countOnly {
			addField(s[i0:])
		}
		n++
	}
	return
}

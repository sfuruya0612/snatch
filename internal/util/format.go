package util

import "strconv"

//import "strings"

func Formatln(list ...[]string) string {
	return Format(list...) + "\n"
}

func Format(list ...[]string) string {
	f := ""
	for _, c := range list {
		f = f + format(c)
	}

	return f
}

func format(list []string) string {
	max := 0
	for _, name := range list {
		if len(name) < max {
			continue
		}
		max = len(name)
	}

	return "%-" + strconv.Itoa(max+1) + "s"
}

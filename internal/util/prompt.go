package util

import "github.com/manifoldco/promptui"

func Prompt(elements []string, label string)
	searcher := func(input string, index int) bool {
		cluster := strings.ToLower(elements[index])
		return strings.Contains(cluster, input)
	}

	return promptui.Select{
		Label:    label,
		Items:    elements,
		Size:     50,
		Searcher: searcher,
	}
}

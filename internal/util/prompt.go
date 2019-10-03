package util

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
)

func Prompt(elements []string, label string) (string, error) {
	searcher := func(input string, index int) bool {
		lower := strings.ToLower(elements[index])
		return strings.Contains(lower, input)
	}

	prompt := promptui.Select{
		Label:    label,
		Items:    elements,
		Size:     50,
		Searcher: searcher,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	return result, nil
}

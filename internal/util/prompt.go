package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func Prompt(elements []string, label string) (string, error) {
	if len(elements) < 1 {
		return "", fmt.Errorf("Elements is empty")
	}

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

func Confirm(element string) bool {
	fmt.Printf("Choose %v ? (Y/N) ", element)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	switch scanner.Text() {
	case "Y", "y", "YES", "Yes", "yes":
		return true
	case "N", "n", "NO", "No", "no":
		return false
	}

	fmt.Printf("No match input pattern")
	return false
}

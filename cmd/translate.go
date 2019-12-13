package cmd

import (
	"fmt"
	"unicode"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/translate"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/urfave/cli"
)

// Language source and target language code struct
type Language struct {
	Source string
	Target string
	Text   string
}

func Translate(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	text := c.String("text")

	l := chkLanguage(text)

	input := &translate.TextInput{
		SourceLanguageCode: aws.String(l.Source),
		TargetLanguageCode: aws.String(l.Target),
		Text:               aws.String(text),
	}

	if err := input.Validate(); err != nil {
		return fmt.Errorf("%v", err)
	}

	client := saws.NewTransSess(profile, region)
	output, err := client.TranslateText(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("[ %v -> %v ]\n", *output.SourceLanguageCode, *output.TargetLanguageCode)
	fmt.Printf("Before: %v\n", text)
	fmt.Printf("After:  %v\n", *output.TranslatedText)

	return nil
}

func chkLanguage(text string) Language {
	l := Language{
		Source: "en",
		Target: "ja",
	}

	for _, t := range text {
		if unicode.In(t, unicode.Hiragana) {
			l = Language{
				Source: "ja",
				Target: "en",
			}
			return l
		}

		if unicode.In(t, unicode.Katakana) {
			l = Language{
				Source: "ja",
				Target: "en",
			}
			return l
		}
	}

	return l
}

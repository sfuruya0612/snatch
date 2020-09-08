package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/translate"
)

// Translate client struct
type Translate struct {
	Client *translate.Translate
}

// NewTransSess return Translate struct initialized
func NewTransSess(profile, region string) *Translate {
	return &Translate{
		Client: translate.New(GetSession(profile, region)),
	}
}

// TranslateText return Stacks
// input translate.TextInput
func (c *Translate) TranslateText(input *translate.TextInput) (*translate.TextOutput, error) {
	output, err := c.Client.Text(input)
	if err != nil {
		return nil, fmt.Errorf("translate text: %v", err)
	}

	return output, nil
}

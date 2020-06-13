package change

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

const UiApi = 1

type surveyContext struct {
	entry *Entry
}

func (c *surveyContext) WriteAnswer(field string, value interface{}) error {
	var err error

	switch field {
	case "title":
		c.entry.Title = value.(string)
	case "author":
		c.entry.Author = value.(string)
	case "reason":
		c.entry.Reason, err = ParseReason(value.(core.OptionAnswer).Value)
	case "comp":
		c.entry.Component = value.(string)
	case "req":
		if val, ok := value.(string); ok && len(val) > 0 {
			c.entry.MergeRequest, err = strconv.Atoi(val)
		}
	default:
		err = fmt.Errorf("unable to process answer for %q", field)
	}

	return err
}

type EntryPrompt struct {
	TargetUiApi int
}

func (p *EntryPrompt) Run(entry *Entry) (bool, error) {
	context := surveyContext{entry}
	questions := p.buildSurvey(entry)
	err := survey.Ask(questions, &context)

	if err == terminal.InterruptErr {
		return false, nil
	}

	return (err == nil), err
}

func (p *EntryPrompt) buildSurvey(entry *Entry) []*survey.Question {
	targetApi := TargetUiApi()
	questions := []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "Changelog title:",
				Default: entry.Title,
			},
			Validate: survey.Required,
		},
	}

	if targetApi >= 1 {
		questions = append(questions,
			&survey.Question{
				Name: "reason",
				Prompt: &survey.Select{
					Message: "Reason:",
					Options: ReasonNames[1:],
					Default: ReasonNames[1],
				},
				Validate: survey.Required,
			},
			&survey.Question{
				Name: "author",
				Prompt: &survey.Input{
					Message: "Author:",
					Help:    "Name to accredit this contribution to",
					Default: entry.Author,
				},
				Validate: survey.Required,
			},
			&survey.Question{
				Name: "comp",
				Prompt: &survey.Input{
					Message: "Affected application component:",
					Help:    "Area of the project this change affects or empty if general change",
					Default: entry.Component,
				},
			},
			&survey.Question{
				Name: "req",
				Prompt: &survey.Input{
					Message: "Associated merge request:",
					Help:    "Numeric ID of the pull/merge request",
					Default: entry.MergeRequestString(),
				},
			},
		)
	}

	return questions
}

func TargetUiApi() int {
	val := os.Getenv("CHANGELOG_ASSEMBLER_UI_API")

	if len(val) == 0 {
		return UiApi
	}

	if api, err := strconv.Atoi(val); err == nil {
		return api
	}

	return UiApi
}

func NewEntryPrompt(targetUiApi int) *EntryPrompt {
	entryPrompt := &EntryPrompt{
		TargetUiApi: targetUiApi,
	}

	return entryPrompt
}

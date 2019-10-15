package common

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

var Log logrus.Logger

func PromptForInput(questionText string, allowEmpty bool, defaultAnswer string)string{
	prompt := promptui.Prompt{
		Label:    questionText,
		Default: defaultAnswer,
	}
	result, err := prompt.Run()

	if err != nil {
		Log.Fatal(err)
	}
	return result
}

func TextPromptForSelection(questionText string, options []string)string{
	prompt := promptui.Select{
		Label: questionText,
		Items: options,
	}
	_, result, err := prompt.Run()

	if err != nil {
		Log.Fatal(err)
	}
	return result
}


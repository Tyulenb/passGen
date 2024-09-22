package main

import (
	"errors"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

var StartMenu = &survey.Select{
	Message: "Choose an option",
	Options: []string{"Sign in", "Sign up"},
}

var MainSelect = &survey.Select{
	Message: "Choose an option",
	Options: []string{"Generate New Password", "Get My Passwords", "Exit"},
}

var GeneratePassQuestions = []*survey.Question{
	{
		Name: "key-digit",
		Prompt: &survey.Input{
			Message: "Enter your key digit[1 to 9]",
		},
		Validate: func(ans interface{}) error {
			input, ok := strconv.Atoi(ans.(string))
			if ok != nil || input > 10 || input < 1 {
				return errors.New("you should enter a digit in this input")
			}
			return nil
		},
	},
	{
		Name: "Length",
		Prompt: &survey.Input{
			Message: "Enter length of your password [from 8 to 256]",
		},
		Validate: func(ans interface{}) error {
			input, ok := strconv.Atoi(ans.(string))
			if ok != nil || input > 256 || input < 8 {
				return errors.New("password length should be from 8 to 256")
			}
			return nil
		},
	},
	{
		Name: "Service",
		Prompt: &survey.Input{
			Message: "Provide the service where you intend to use this password.",
			Default: "",
			Help:    "Otherwise, it won't be stored in your password list.",
		},
	},
}

var Authorization = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Enter your user name",
			Help:    "Your registration name",
		},
		Validate: survey.Required,
	},
	{
		Name: "key",
		Prompt: &survey.Input{
			Message: "Enter your key",
			Help:    "You were given it during registration",
		},
		Validate: survey.Required,
	},
}

var Registration = &survey.Input{
	Message: "Please enter your name it will serve as your login for future authorizations",
	Help:    "Keep it safe, or you won't be able to access your password storage!",
}

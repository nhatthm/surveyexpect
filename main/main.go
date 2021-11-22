package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	p := &survey.Select{Message: "Select a country", Options: []string{"Germany", "Vietnam"}}

	var answer string

	err := survey.AskOne(p, &answer)

	fmt.Println(err)
}

package cliinput

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func InputProjectName() string {
	prompt := promptui.Prompt{
		Label: "Enter Project Name",
	}
	res, err := prompt.Run()
	if err != nil {
		fmt.Println("error is ", err)
		return ""
	}
	if res == "" {
		return "my-express-project"
	}
	return res
}

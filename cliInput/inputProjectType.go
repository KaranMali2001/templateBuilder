package cliinput

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func InputProjectType() string {
	items := []string{"Express with Typescript", "Express with Javascript"}
	prompt := promptui.Select{
		Label: "Enter project Name",
		Items: items,
	}
	_, res, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return "my-new-express-project"
	}
	return res

}

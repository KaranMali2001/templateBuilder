package cliinput

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func InputDb() (string, string) {
	items := []string{"postgresql", "mongodb", "None"}
	prompt := promptui.Select{
		Label: "choose Database",
		Items: items,
	}
	_, res, err := prompt.Run()
	if err != nil {
		fmt.Println("error ", err)
	}
	if res == "None" {
		return res, ""
	}
	orm := inputOrm()
	return res, orm
}
func inputOrm() string {
	items := []string{"Prisma", "Will add it Later"}
	prompt := promptui.Select{
		Label: "Choose ORM",
		Items: items,
	}
	_, res, err := prompt.Run()
	if err != nil {
		fmt.Println("error while orm", err)
	}
	return res
}

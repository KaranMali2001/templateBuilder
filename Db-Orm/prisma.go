package dborm

import (
	"fmt"
	"os"
	"os/exec"

	projectstructure "github.com/KaranMali2001/templateBuilder/projectStructure"
)

func InitPrisma(dbString string, db string, projectName string, projectType string, orm string) error {
	filePath := projectName + "/.env"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("DATABASE_URL=%s", dbString))
	if err != nil {
		return err
	}
	filePath = projectName + "/prisma/schema.prisma"
	content := `generator client {
provider = "prisma-client-js"
}

datasource db {
provider = "mongodb"
url      = env("DATABASE_URL")
}

	model Users {
id   String @id @map("_id")
name String
}
	`

	err = projectstructure.CreateFileWithContent(filePath, content)
	if err != nil {
		return err
	}
	if projectType == "Express with Typescript" && orm == "Prisma" {
		err = os.Chdir(projectName)
		if err != nil {

			return err
		}
		cmd := exec.Command("npm", "install", "prisma", "@prisma/client", "--save-dev")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {

			return err
		}
		cmd = exec.Command("npx", "prisma", "init", "--datasource-provider", db)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("error while initalizing Prisma")
		}
		cmd = exec.Command("npx", "prisma", "generate")

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("error while initalizing Prisma")
		}

	}

	return nil
}

package main

import (
	"fmt"
	"os"

	dborm "github.com/KaranMali2001/templateBuilder/Db-Orm"
	cliinput "github.com/KaranMali2001/templateBuilder/cliInput"
	projectstructure "github.com/KaranMali2001/templateBuilder/projectStructure"
	"github.com/spf13/cobra"
)

func main() {
	var rootCommand = &cobra.Command{
		Use:   "express-ts",
		Short: "tool to generate express application",
	}
	var initCommand = &cobra.Command{
		Use:   "init",
		Short: "initialize new typescript project",
		Run: func(cmd *cobra.Command, args []string) {
			projectName := cliinput.InputProjectName()
			projectType := cliinput.InputProjectType()

			dbName, orm := cliinput.InputDb()
			fmt.Println("project type is", projectType)
			fmt.Println("orm is ", orm)
			fmt.Printf("Creating Express TypeScript project %s with %s database...\n", projectName, dbName)
			if err := os.Mkdir(projectName, 0775); err != nil {
				fmt.Println("error while creating dir", err)
				return
			}
			//without any db config
			if projectType == "Express with Typescript" {
				err := projectstructure.TsProjectStructure(projectName)
				if err != nil {
					fmt.Println("error while creating project structure", err)
					return
				}
			} else {
				//create javascript project
				err := projectstructure.JsProjectStrcuture(projectName)
				if err != nil {
					fmt.Println("error ", err)
					return
				}
			}

			if dbName == "postgresql" {

				imageName := "postgres"
				fmt.Println("image name is ", imageName)
				dbString, err := dborm.RunDockerContainer(imageName, projectName)
				if err != nil {
					fmt.Println("error while creating docker container", err)
					return
				}
				//init primsa
				err = dborm.InitPrisma(dbString, dbName, projectName, projectType, orm)
				if err != nil {
					fmt.Println("error while init primsa is", err)
					return
				}

			} else if dbName == "mongodb" {
				imageName := "mongo"

				fmt.Println("image name is ", imageName)
				dbString, err := dborm.RunDockerContainer(imageName, projectName)
				if err != nil {
					fmt.Println("error while creating docker container", err)
					return
				}
				err = dborm.InitPrisma(dbString, dbName, projectName, projectType, orm)
				if err != nil {
					fmt.Println("error while init primsa is", err)
					return
				}
				//run docker container with volume same as name of project
			}
		},
	}
	rootCommand.AddCommand(initCommand)
	rootCommand.Execute()

}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pkgName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an Potato Automated Testing Project",
	Long: `Initialize (potato init) will create a new project with the appropriate structure.
	
	  * If a name is provided, a directory with that name will be created in the current directory;
	  * If no name is provided, the current directory will be assumed;
	`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := initializeProject(args)
		cobra.CheckErr(err)
		fmt.Printf("Your Potato Automated Testing Project is ready at\n%s\n", projectPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&pkgName, "pkg-name", "", "fully qualified pkg name")
	cobra.CheckErr(initCmd.MarkFlagRequired("pkg-name"))
}

func initializeProject(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			wd = fmt.Sprintf("%s/%s", wd, args[0])
		}
	}

	project := &Project{
		AbsolutePath: wd,
		PkgName:      pkgName,
		Viper:        viper.GetBool("useViper"),
		ProjectName:  path.Base(pkgName),
	}
	if err := project.Create(); err != nil {
		return "", err
	}
	return project.AbsolutePath, nil
}

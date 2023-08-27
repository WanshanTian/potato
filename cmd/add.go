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
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addSuite = &cobra.Command{
	Use:     "add [testsuite name]",
	Aliases: []string{"testsuite, suite"},
	Short:   "Add a testsuite to a Automated Testing Project",
	Long: `Add (potato add) will create a new testsuite with the appropriate structure.
	
	If you want your testsuite to be executed, pass in the testsuite name
	with an initial uppercase letter.
	
	Example: potato add Network -> resulting in a new Network/network.go`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("add needs a name for the testsuite"))
		}

		wd, err := os.Getwd()
		cobra.CheckErr(err)
		if path.IsAbs(args[0]) {
			fmt.Println("please enter a relative path")
			os.Exit(1)
		}
		suiteName := validateSuiteName(args[0])
		suite := &Testsuite{
			SuiteName:            suiteName,
			ToLowerSuiteBaseName: strings.ToLower(path.Base(suiteName)),
			Project: &Project{
				AbsolutePath: wd,
			},
		}
		cobra.CheckErr(suite.Create())
		fmt.Printf("%s created at %s\n", suite.SuiteName, suite.AbsolutePath)
	},
}

func init() {
	rootCmd.AddCommand(addSuite)

}
func validateSuiteName(source string) string {
	origin := source
	source = path.Base(source)
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l {
		if source[i] == '-' || source[i] == '_' {
			if output == "" {
				output = source[:i]
			}

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 {
				break
			}

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' {
				i++
				continue
			}

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
			output += string(source[i])
		}
		i++
	}

	if output == "" {
		return origin // source is initially valid name.
	}
	return path.Join(path.Dir(origin), output)
}

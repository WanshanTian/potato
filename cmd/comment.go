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
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:  "comment",
	Long: `Automatically generate the comment of testcases and testsuites.`,
	Run: func(cmd *cobra.Command, args []string) {
		var dst string
		wd, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		if strings.Contains(wd, TestsuitesDirName) {
			dst = strings.Split(wd, TestsuitesDirName)[0]
		} else {
			if _, err := os.Stat(path.Join(wd, TestsuitesDirName)); os.IsNotExist(err) {
				fmt.Println(err.Error() + "\nPlease execute comment under the rootpath of potato project or the path of testsuites")
				os.Exit(1)
			}
			dst = path.Join(wd, TestsuitesDirName)
		}
		if !strings.Contains(dst, TestsuitesDirName) {
			fmt.Printf("The absolute path is %s, please enter the correct relative path of testsuite", dst)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

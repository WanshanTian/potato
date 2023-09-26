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
	"log"
	"path"

	"github.com/spf13/cobra"
	"github.com/txy2023/potato/utils"
)

type CommentInfo struct {
	Testsuite string
	*Project
}

var Commentinfomation = new(CommentInfo)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:  "comment",
	Long: `Automatically generate the comment of testcases and testsuites.`,
	Run: func(cmd *cobra.Command, args []string) {
		dst, err := utils.GetTestSuiteAbsoluteRootDir(TestsuitesDirName)
		if err != nil {
			log.Panic(err)
		}
		com, err := utils.GetAllTestSuitesComment(dst)
		if err != nil {
			log.Panic(err)
		}
		prettysuitecom := utils.PrettySuiteComment(com)
		Commentinfomation.Testsuite = prettysuitecom
		utils.CommentWrite(Commentinfomation, path.Dir(dst))
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
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
	"github.com/txy2023/potato/tpl"
	"github.com/txy2023/potato/utils"
)

type CommentInfo struct {
	Testsuite      string
	TestsuiteCount int
	Testcase       string
	TestcaseCount  int
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
		// testsuite
		com, err := utils.GetAllTestSuitesComment(dst)
		if err != nil {
			log.Panic(err)
		}
		testsuiteNum := utils.GetTestSuitesNum(com)
		prettysuitecom := utils.PrettySuiteComment(com)
		// testcase
		comCase, err := utils.GetAllTestCasesComment(dst)
		if err != nil {
			log.Panic(err)
		}
		testcaseeNum := utils.GetTestCasesNum(comCase)
		prettycasecom := utils.PrettyCaseComment(comCase)
		// write
		Commentinfomation.Testsuite = prettysuitecom
		Commentinfomation.TestsuiteCount = testsuiteNum
		Commentinfomation.Testcase = prettycasecom
		Commentinfomation.TestcaseCount = testcaseeNum
		utils.CommentWrite(Commentinfomation, path.Dir(dst), string(tpl.TestSuiteCommentTemplate()), string(tpl.TestCaseCommentTemplate()))
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}

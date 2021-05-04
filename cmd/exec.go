/*
Copyright © 2021 SuperOrbital, LLC <info@superorbital.io>

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

	"github.com/spf13/cobra"
)

func AddExec(rootCmd *cobra.Command) {
	// execCmd represents the exec command
	var execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Execute a command with the correct AWS credentials",
		Long:  `This will connect to the cludo server, retrieve a temporary token and then execute your command with the AWS credential variables set correctly.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exec something")
		},
	}

	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

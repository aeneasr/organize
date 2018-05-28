// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"path/filepath"
	"os"
	"io/ioutil"
	"path"
	"strings"
	"github.com/pkg/errors"
)

func must(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

var offset = 0

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "move <from> <to>",
	Short: "Recursively scans files from src, and renames and moves them to dst",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(cmd.UsageString())
			return
		}

		files, err := ioutil.ReadDir(args[1])
		must(err)
		offset = len(files)

		must(copyDir(args[0], args[1]))
	},
}

func copyDir(dir string, to string)error{
	fmt.Printf("Entering directory %s\n", dir)
	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		fmt.Printf("Processing %s\n", p)
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() && strings.ToLower(dir) == strings.ToLower(p) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		fileName := fmt.Sprintf("%d_%d_%d-%d_%d-%d%s",
			info.ModTime().Year(),
			int(info.ModTime().Month()),
			info.ModTime().Day(),
			info.ModTime().Hour(),
			info.ModTime().Minute(),
			offset,
			filepath.Ext(p),
		)

		cpto := path.Join(to, fileName)
		fmt.Printf("Moving file from %s to %s\n", p, cpto)

		offset++
		if err := os.Rename(p, cpto); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}
func init() {
	RootCmd.AddCommand(dirCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dirCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

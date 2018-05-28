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
	"github.com/spf13/cobra"
	"fmt"
	"path/filepath"
	"os"
	"github.com/pkg/errors"
	"io"
	"encoding/hex"
	"crypto/sha256"
)

// dupesCmd represents the dupes command
var dupesCmd = &cobra.Command{
	Use:   "dupes <src> <dst>",
	Short: "Recursively finds dupes in src using sha256 and moves them to dst",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println(cmd.UsageString())
			return
		}

		var dupes int
		var md5sums = make(map[string][]string)

		must(filepath.Walk(args[0], func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.WithStack(err)
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(p)
			if err != nil {
				return errors.WithStack(err)

			}

			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				f.Close()
				return errors.WithStack(err)
			}
			f.Close()

			sum := hex.EncodeToString(h.Sum(nil))
			md5sums[sum] = append(md5sums[sum], p)
			if len(md5sums[sum]) > 1 {
				if err := os.Rename(p, filepath.Join(args[1], fmt.Sprintf("%s_%s", sum, filepath.Base(p)))); err != nil {
					return errors.WithStack(err)
				}
				fmt.Printf("Found dupes for %s: %v\n", sum, md5sums[sum])
				dupes ++
			}

			return nil
		}))

		fmt.Printf("Found %d dupes\n", dupes)
	},
}

func init() {
	RootCmd.AddCommand(dupesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dupesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dupesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

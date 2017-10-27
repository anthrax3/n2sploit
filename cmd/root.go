// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"os"
	"os/exec"
	"strings"

	"github.com/kkirsche/nmap-searchsploit/libnmap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	fuzzy bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nmap-searchsploit",
	Short: "Used to take Nmap data and search Exploit-DB via SearchSploit",
	Long: `Used to take Nmap data and search Exploit-DB via SearchSploit for
exploits related to the service`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			n, err := libnmap.Parse(arg)
			if err != nil {
				continue
			}
			for _, port := range n.Host.Ports.Port {
				version := port.Service.AttrVersion
				if fuzzy {
					s := strings.Split(port.Service.AttrVersion, ".")
					vl := s[:1]
					version = strings.Join(vl, ".")
				}
				q := fmt.Sprintf("%s %s", port.Service.AttrProduct, version)

				logrus.Infof("Searching searchsploit for \"%s\"...", q)
				out, err := exec.Command("searchsploit", q).Output()
				if err != nil {
					logrus.WithError(err).Errorln("Failed to search...")
				}

				logrus.Println(string(out))
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&fuzzy, "fuzzy", "f", false, "Fuzzy version searching")
}

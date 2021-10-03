/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/olekukonko/tablewriter"
	"github.com/tomokazukozuma/polkadot-cli/lib/encode"

	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			log.Fatalf("Failed Get create: %s", err.Error())
		}
		publicKey, ss58Prefix, err := encode.DecodeAddress(address)
		if err != nil {
			log.Fatalf("Failed Get create: %s", err.Error())
		}
		data := [][]string{
			{"Address", fmt.Sprintf("%s", address)},
			{"PublicKey", fmt.Sprintf("0x%x", publicKey)},
			{"ss58Prefix", fmt.Sprintf("%d", ss58Prefix)},
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	addressCmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	decodeCmd.Flags().String("address", "", "address")
}

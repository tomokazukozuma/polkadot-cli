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
	"crypto/ed25519"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tomokazukozuma/polkadot-cli/lib/encode"
)

// addressCmd represents the address command
var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := cmd.Flags().GetBool("create")
		if err != nil {
			log.Fatalf("Failed Get create: %s", err.Error())
		}
		if c {
			publicKey, privateKey, err := ed25519.GenerateKey(nil)
			if err != nil {
				log.Fatalf("Failed GenerateKey", err)
			}

			ss58Prefix, err := cmd.Flags().GetInt8("ss58Prefix")
			if err != nil {
				log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
			}
			address := encode.EncodeAddress(publicKey, ss58Prefix)

			data := [][]string{
				{"PrivateKey", fmt.Sprintf("%x", privateKey)},
				{"PublicKey", fmt.Sprintf("%x", publicKey)},
				{"Address", fmt.Sprintf("%s", address)},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(true)
			table.AppendBulk(data)
			table.Render()
		}

		d, err := cmd.Flags().GetBool("decode")
		if err != nil {
			log.Fatalf("Failed Get create: %s", err.Error())
		}
		if d {
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
				{"PublicKey", fmt.Sprintf("%x", publicKey)},
				{"ss58Prefix", fmt.Sprintf("%d", ss58Prefix)},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(true)
			table.AppendBulk(data)
			table.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(addressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addressCmd.Flags().Bool("create", false, "create address")
	addressCmd.Flags().Int8("ss58Prefix", 0, "SS58Prefix 0: Polkadot, 2: Kusama, 42: Westend")
	addressCmd.Flags().Bool("decode", false, "decode address")
	addressCmd.Flags().String("address", "", "address")
}

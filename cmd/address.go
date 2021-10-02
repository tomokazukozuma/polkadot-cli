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

	"github.com/akamensky/base58"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/blake2b"
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
		c, _ := cmd.Flags().GetBool("create")
		if c {
			publicKey, privateKey, err := ed25519.GenerateKey(nil)
			if err != nil {
				log.Fatalf("Failed GenerateKey", err)
			}

			ss58Prefix, _ := cmd.Flags().GetInt8("ss58Prefix")
			address := EncodeAddress(publicKey, ss58Prefix)

			data := [][]string{
				{"PrivateKey", fmt.Sprintf("%x", privateKey)},
				{"PublicKey", fmt.Sprintf("%x", publicKey)},
				{"Address", fmt.Sprintf("%s", address)},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(true)
			table.AppendBulk(data)
			table.Render()

		} else {
			fmt.Println("address called")
		}
	},
}

type option struct {
	create bool
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
	addressCmd.Flags().Bool("create", false,"create address")
	addressCmd.Flags().Int8("ss58Prefix", 0,"SS58Prefix 0: Polkadot, 2: Kusama, 42: Westend")
}

var (
	prefix = []byte("SS58PRE")
)

func EncodeAddress(pubKey []byte, ss58Prefix int8) string {
	var raw []byte
	raw = append([]byte{byte(ss58Prefix)}, pubKey...)
	checksum := blake2b.Sum512(append(prefix, raw...))
	address := base58.Encode(append(raw, checksum[0:2]...))
	return address
}
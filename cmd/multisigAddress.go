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
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/blake2b"
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// multisigAddressCmd represents the multisigAddress command
var multisigAddressCmd = &cobra.Command{
	Use:   "multisigAddress",
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
			addresses, err := cmd.Flags().GetStringSlice("addresses")
			if err != nil {
				log.Fatalf("Failed Get addresses: %s", err.Error())
			}
			threshold, err := cmd.Flags().GetUint8("threshold")
			if err != nil {
				log.Fatalf("Failed Get threshold: %s", err.Error())
			}
			ss58Prefix, err := cmd.Flags().GetInt8("ss58Prefix")
			if err != nil {
				log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
			}

			var publicKeys [][]byte
			for _, address := range addresses {
				publicKey, _, _:= DecodeAddress(address)
				publicKeys = append(publicKeys, publicKey)
			}
			sort.Slice(publicKeys, func(i,j int) bool {
				return bytes.Compare(publicKeys[i], publicKeys[j]) < 0
			})
			var payload []byte
			var multiSigPrefix = []byte("modlpy/utilisuba")
			payload = append(payload, multiSigPrefix...)
			var keyLength = uint8(len(publicKeys)) << 2
			payload = append(payload, keyLength)
			for _, pubKey := range publicKeys {
				payload = append(payload, pubKey...)
			}
			payload = append(payload, threshold, 0x00)
			multisigPubKey := blake2b.Sum256(payload)
			address := EncodeAddress(multisigPubKey[:], ss58Prefix)
			data := [][]string{
				{"Multisig PublicKey", fmt.Sprintf("%x", multisigPubKey)},
				{"Address", fmt.Sprintf("%s", address)},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(true)
			table.AppendBulk(data)
			table.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(multisigAddressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// multisigAddressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// multisigAddressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	multisigAddressCmd.Flags().Bool("create", false,"create address")
	multisigAddressCmd.Flags().StringSlice("addresses", nil,"addresses")
	multisigAddressCmd.Flags().Uint8("threshold", 0,"threshold")
	multisigAddressCmd.Flags().Int8("ss58Prefix", 0,"SS58Prefix 0: Polkadot, 2: Kusama, 42: Westend")
}

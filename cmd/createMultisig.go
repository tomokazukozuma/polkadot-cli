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
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tomokazukozuma/polkadot-cli/lib/encode"
	"golang.org/x/crypto/blake2b"
)

// createMultisigCmd represents the createMultisig command
var createMultisigCmd = &cobra.Command{
	Use:   "createMultisig",
	Short: "Create multisig address",
	Long: `Create multisig address. For example:

address createMultisig --threshold=2 --addresses="1VhsQ5adREGuorYyrKacR5KB4XCkYbCr7YunQW5pAPgiVP9,15FGVSb62LVw4saLnK43PHT1N7fpjcdzTwJYsfuquNfhTeT4".`,
	Run: func(cmd *cobra.Command, args []string) {
		addresses, err := cmd.Flags().GetStringSlice("addresses")
		if err != nil {
			log.Fatalf("Failed Get addresses: %s", err.Error())
		}
		threshold, err := cmd.Flags().GetUint16("threshold")
		if err != nil {
			log.Fatalf("Failed Get threshold: %s", err.Error())
		}
		ss58Prefix, err := cmd.Flags().GetInt8("ss58Prefix")
		if err != nil {
			log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
		}

		var publicKeys [][]byte
		for _, address := range addresses {
			publicKey, _, _ := encode.DecodeAddress(address)
			publicKeys = append(publicKeys, publicKey)
		}
		sort.Slice(publicKeys, func(i, j int) bool {
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
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, threshold)
		payload = append(payload, b...)
		multisigPubKey := blake2b.Sum256(payload)
		address := encode.EncodeAddress(multisigPubKey[:], ss58Prefix)
		data := [][]string{
			{"Multisig PublicKey", fmt.Sprintf("0x%x", multisigPubKey)},
			{"Threshold", fmt.Sprintf("%d", threshold)},
			{"Address", fmt.Sprintf("%s", address)},
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	addressCmd.AddCommand(createMultisigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createMultisigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createMultisigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createMultisigCmd.Flags().StringSlice("addresses", nil, "addresses")
	createMultisigCmd.Flags().Uint16("threshold", 0, "threshold")
	createMultisigCmd.Flags().Int8("ss58Prefix", 0, "SS58Prefix 0: Polkadot, 2: Kusama, 42: Westend")
}

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
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tomokazukozuma/polkadot-cli/lib/bip39"
	"github.com/tomokazukozuma/polkadot-cli/lib/encode"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ss58 address",
	Long: `create ss58 address. For example:

address create --ss58Prefix=0 --mnemonic="food ring street ... shield"`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			privateKey []byte
			publicKey  []byte
		)
		stringMnemonic, err := cmd.Flags().GetString("mnemonic")
		if err != nil {
			log.Fatalf("Failed Get mnemonic: %s", err.Error())
		}
		log.Printf("stringMnemonic: %s", stringMnemonic)
		stringPrivateKey, err := cmd.Flags().GetString("privateKey")
		if err != nil {
			log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
		}
		stringPublicKey, err := cmd.Flags().GetString("publicKey")
		if err != nil {
			log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
		}
		var data [][]string
		if stringMnemonic != "" {
			stringPassphrase, err := cmd.Flags().GetString("passphrase")
			if err != nil {
				log.Fatalf("Failed Get passphrase: %s", err.Error())
			}
			entropy, err := bip39.MnemonicToEntropy(stringMnemonic)
			if err != nil {
				log.Fatalf("Failed MnemonicToEntropy: %s", err.Error())
			}
			seed, err := bip39.MnemonicToSeed(stringMnemonic, stringPassphrase)
			if err != nil {
				log.Fatalf("Failed MnemonicToSeed: %s", err.Error())
			}
			extendedPrivateKey := ed25519.NewKeyFromSeed(seed[:32])
			privateKey = extendedPrivateKey[:32]
			publicKey = extendedPrivateKey[32:]

			data = [][]string{
				{"entropy", fmt.Sprintf("%x", entropy)},
				{"mnemonic", fmt.Sprintf("%s", stringMnemonic)},
			}
		} else if stringPrivateKey != "" {
			p, err := hex.DecodeString(stringPrivateKey)
			if err != nil {
				log.Fatalf("Failed Decode privateKey: %s", err.Error())
			}
			extendedPrivateKey := ed25519.NewKeyFromSeed(p)
			privateKey = extendedPrivateKey[:32]
			publicKey = extendedPrivateKey[32:]
		} else if stringPublicKey != "" {
			publicKey, err = hex.DecodeString(stringPublicKey)
			if err != nil {
				log.Fatalf("Failed Decode privateKey: %s", err.Error())
			}
		} else {
			stringPassphrase, err := cmd.Flags().GetString("passphrase")
			if err != nil {
				log.Fatalf("Failed Get passphrase: %s", err.Error())
			}
			mnemonic, entropy, err := bip39.GenerateMnemonic()
			if err != nil {
				log.Fatalf("Failed GenerateMnemonic: %s", err.Error())
			}
			seed, err := bip39.MnemonicToSeed(strings.Join(mnemonic, " "), stringPassphrase)
			if err != nil {
				log.Fatalf("Failed MnemonicToSeed: %s", err.Error())
			}
			extendedPrivateKey := ed25519.NewKeyFromSeed(seed[:32])
			privateKey = extendedPrivateKey[:32]
			publicKey = extendedPrivateKey[32:]
			data = [][]string{
				{"entropy", fmt.Sprintf("%x", entropy)},
				{"mnemonic", fmt.Sprintf("%s", strings.Join(mnemonic, " "))},
			}
		}

		ss58Prefix, err := cmd.Flags().GetInt8("ss58Prefix")
		if err != nil {
			log.Fatalf("Failed Get ss58Prefix: %s", err.Error())
		}
		address := encode.EncodeAddress(publicKey, ss58Prefix)

		data = append(data, [][]string{
			{"PrivateKey", fmt.Sprintf("0x%x", privateKey)},
			{"PublicKey", fmt.Sprintf("0x%x", publicKey)},
			{"Address", fmt.Sprintf("%s", address)},
		}...)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	addressCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().String("mnemonic", "", "Mnemonic")
	createCmd.Flags().String("passphrase", "", "passphrase")
	createCmd.Flags().String("privateKey", "", "private key")
	createCmd.Flags().String("publicKey", "", "public key")
	createCmd.Flags().Int8("ss58Prefix", 0, "SS58Prefix 0: Polkadot, 2: Kusama, 42: Westend")
}

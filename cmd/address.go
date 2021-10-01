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
	"crypto/rand"
	"fmt"
	"log"

	"github.com/akamensky/base58"
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
			log.Println("address create called")
			publicKey, privateKey, err := ed25519.GenerateKey(nil)
			if err != nil {
				log.Fatalf("Failed GenerateKey", err)
			}
			address := EncodeAddress(publicKey)
			log.Printf("privateKey: %x", privateKey)
			log.Printf("publicKey: %x", publicKey)
			log.Printf("address: %s", address)

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
	var o = &option{}
	addressCmd.Flags().BoolVar(&o.create, "create", false,"create address")
}

func GenerateRandom(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var (
	prefix = []byte("SS58PRE")
)
func EncodeAddress(pubKey []byte) string {
	var raw []byte
	addressType := []byte{0x00}
	//if isTestnet {
	//	// Westend is 42
	//	addressType = []byte{0x2A}
	//}
	raw = append(addressType, pubKey...)
	checksum := blake2b.Sum512(append(prefix, raw...))
	address := base58.Encode(append(raw, checksum[0:2]...))
	return address
}
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/urfave/cli/v2"
)

// Superchain registry url: https://raw.githubusercontent.com/ethereum-optimism/superchain-registry/refs/heads/main

func CreateApp() *cli.App {
	return &cli.App{
		Name:  "superchain-registry-cli",
		Usage: "A tool for interacting with the superchain-registry",
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version of the CLI",
				Action: func(c *cli.Context) error {
					fmt.Printf("superchain-registry-cli version %s\n", Version)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all chains in the superchain",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Enable verbose output",
					},
					&cli.BoolFlag{
						Name:    "testnet",
						Aliases: []string{"t"},
						Usage:   "Filter for testnet chains",
					},
				},
				Action: func(c *cli.Context) error {
					testnet := c.Bool("testnet")

					for _, chain := range superchain.OPChains {
						if testnet && chain.Superchain != "sepolia" || !testnet && chain.Superchain != "mainnet" {
							continue
						}

						if chain.SuperchainLevel == superchain.Standard {
							fmt.Printf("Chain: %s\n", chain.Chain)
							fmt.Printf("Network: %s\n", chain.Name)
							if c.Bool("verbose") {
								fmt.Printf("  Identifier: %s\n", chain.Identifier())
								fmt.Printf("  Chain ID: %d\n", chain.ChainID)
								fmt.Printf("  RPC: %s\n", chain.PublicRPC)
								fmt.Printf("  Explorer: %s\n", chain.Explorer)
							}
						}
					}
					return nil
				},
			},
			{
				Name:                   "get-addresses",
				Aliases:                []string{"ga"},
				Usage:                  "Gets addresses for a given chain",
				HideHelpCommand:        true,
				UseShortOptionHandling: true,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Enable verbose output",
					},
					&cli.BoolFlag{
						Name:    "testnet",
						Aliases: []string{"t"},
						Usage:   "Filter for testnet chains",
					},
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"a"},
						Usage:   "Address to find",
					},
					&cli.StringFlag{
						Name:    "address-name",
						Aliases: []string{"an"},
						Usage:   "Address name to find",
					},
					&cli.StringFlag{
						Name:    "chain",
						Aliases: []string{"c"},
						Usage:   "Chain to filter by",
					},
				},
				Action: func(c *cli.Context) error {
					if !c.IsSet("address") && !c.IsSet("chain") && !c.IsSet("address-name") {
						return cli.ShowCommandHelp(c, "get-addresses")
					}

					// TODO: validate address
					address := c.String("address")
					chain := c.String("chain")
					testnet := c.Bool("testnet")
					addressName := c.String("address-name")

					findChain(superchain.OPChains, chain, address, addressName, testnet)
					return nil
				},
			},
		},
	}
}

func main() {
	app := CreateApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type NamedAddress struct {
	Name    string
	Address superchain.Address
}

func ConvertAddressListToNamedAddresses(addressList superchain.AddressList) []NamedAddress {
	var namedAddresses []NamedAddress

	val := reflect.ValueOf(addressList)
	typ := reflect.TypeOf(addressList)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if field.Type() == reflect.TypeOf(superchain.Address{}) {
			namedAddresses = append(namedAddresses, NamedAddress{
				Name:    fieldName,
				Address: field.Interface().(superchain.Address),
			})
		}
	}

	return namedAddresses
}

func findChain(opChains map[uint64]*superchain.ChainConfig, chainName, addressToFind, addressNameToFind string, isTestnet bool) {
	for _, chain := range opChains {
		if chainName != "" && chain.Chain != chainName {
			continue
		}

		if isTestnet && chain.Superchain != "sepolia" || !isTestnet && chain.Superchain != "mainnet" {
			continue
		}

		printChainInfo := func(addressName, address string) {
			if address == "0x0000000000000000000000000000000000000000" {
				fmt.Printf("  %s: %s\n", addressName, "<n/a>")
			} else {
				fmt.Print(CreateHyperlinkedAddress(addressName, GetEtherscanURL(address, isTestnet)))
			}
		}
		namedAddresses := ConvertAddressListToNamedAddresses(chain.Addresses)

		if addressToFind == "" {
			fmt.Printf("Chain: %s\n", chain.Chain)
			fmt.Printf("Network: %s\n", chain.Name)

			for _, namedAddress := range namedAddresses {
				if addressNameToFind == "" || strings.Contains(strings.ToLower(namedAddress.Name), strings.ToLower(addressNameToFind)) {
					printChainInfo(namedAddress.Name, namedAddress.Address.String())
				}
			}
		} else {
			for _, namedAddress := range namedAddresses {
				if namedAddress.Address.String() == addressToFind {
					fmt.Printf("Chain: %s\n", chain.Chain)
					fmt.Printf("Network: %s\n", chain.Name)
					fmt.Printf("  Name: %s\n", namedAddress.Name)
					printChainInfo(namedAddress.Name, namedAddress.Address.String())
				}
			}
		}
	}
}

func CreateHyperlinkedAddress(addressName string, etherscanAddressURL string) string {
	addressPart := etherscanAddressURL[len(etherscanAddressURL)-42:]
	return fmt.Sprintf("  %s: \033]8;;%s\033\\%s\033]8;;\033\\\n", addressName, etherscanAddressURL, addressPart)
}

func GetEtherscanURL(address string, isTestnet bool) string {
	baseURL := "https://etherscan.io/address/%s"
	if isTestnet {
		baseURL = "https://sepolia.etherscan.io/address/%s"
	}
	return fmt.Sprintf(baseURL, address)
}

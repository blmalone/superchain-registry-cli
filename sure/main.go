package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

// Superchain registry url: https://raw.githubusercontent.com/ethereum-optimism/superchain-registry/refs/heads/main
func main() {
	app := &cli.App{
		Name:  "superchain-registry-cli",
		Usage: "A tool for interacting with the superchain-registry",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all chains in the superchain",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "superchain-registry-url",
						Aliases: []string{"scr"},
						Usage:   "Specify the superchain registry URL",
						EnvVars: []string{"SUPERCHAIN_REGISTRY_URL"},
					},
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
					superchainRegistryURL := c.String("superchain-registry-url")
					if superchainRegistryURL == "" {
						fmt.Println("SUPERCHAIN_REGISTRY_URL is not set. Please set it or use the --superchain-registry-url flag.")
						return nil
					}

					chainListURL := superchainRegistryURL + "/chainList.toml"

					verbose := c.Bool("verbose")
					testnet := c.Bool("testnet")

					body, err := httpGet(chainListURL)
					if err != nil {
						fmt.Printf("Failed to perform httpGet: %v\n", err)
						return err
					}

					var chainList ChainList
					if err := toml.Unmarshal(body, &chainList); err != nil {
						fmt.Printf("Failed to parse TOML file: %v\n", err)
						return err
					}

					if len(chainList.Chains) == 0 {
						fmt.Println("No chains found in the TOML file.")
					} else {
						for _, chain := range chainList.Chains {
							if chain.SuperchainLevel == STANDARD_CHAIN {
								if (testnet && chain.Parent.Chain == SEPOLIA) || (!testnet && chain.Parent.Chain == MAINNET) {
									fmt.Printf("  Name: %s\n", chain.Name)
									if verbose {
										fmt.Printf("  Identifier: %s\n", chain.Identifier)
										fmt.Printf("  Chain ID: %d\n", chain.ChainID)
										fmt.Printf("  RPC: %v\n", chain.RPC)
										fmt.Printf("  Explorers: %v\n", chain.Explorers)
									}
									fmt.Println()
								}
							}
						}
					}

					return nil
				},
			},
			{
				Name:    "get-addresses",
				Aliases: []string{"ga"},
				Usage:   "Gets addresses for a given chain",
				// HideHelp: true,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "superchain-registry-url",
						Aliases: []string{"scr"},
						Usage:   "Specify the superchain registry URL",
						EnvVars: []string{"SUPERCHAIN_REGISTRY_URL"},
					},
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
						Usage:   "address to find",
					},
					&cli.StringFlag{
						Name:    "address-name",
						Aliases: []string{"an"},
						Usage:   "address name to find",
					},
					&cli.StringFlag{
						Name:    "network",
						Aliases: []string{"n"},
						Usage:   "network to filter by",
					},
				},
				Action: func(c *cli.Context) error {

					superchainRegistryURL := c.String("superchain-registry-url")
					if superchainRegistryURL == "" {
						fmt.Println("SUPERCHAIN_REGISTRY_URL is not set. Please set it or use the --superchain-registry-url flag.")
						return nil
					}

					if !c.IsSet("address") && !c.IsSet("network") && !c.IsSet("address-name") {
						err := cli.ShowCommandHelp(c, "get-addresses")
						return err
					}

					// TODO: validate address
					address := c.String("address")
					network := c.String("network")
					testnet := c.Bool("testnet")
					addressName := c.String("address-name")

					superchainConfigs := superchainRegistryURL + "/superchain/configs/configs.json"

					body, _ := httpGet(superchainConfigs)

					var superchains Superchains
					err := json.Unmarshal([]byte(body), &superchains)
					if err != nil {
						println(err)
						return err
					}

					superchain := superchains.Superchains[0]
					if testnet {
						superchain = superchains.Superchains[1]
					}

					findChain(superchain, network, address, addressName, testnet)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func findChain(superchain Superchain, chainName string, addressToFind string, addressNameToFind string, isTestnet bool) (string, *ChainInfo) {
	for _, chain := range superchain.ChainInfos {
		if chainName != "" && chain.Chain != chainName {
			continue
		}

		printChainInfo := func(addressName, address string) {
			etherscanURL := getEtherscanURL(address, isTestnet)
			fmt.Printf("  %s: \033]8;;%s\033\\%s\033]8;;\033\\\n", addressName, etherscanURL, address)
		}

		if addressToFind == "" {
			fmt.Printf("Chain: %s\n", chain.Chain)
			fmt.Printf("Network: %s\n", superchain.Name)

			for addressName, address := range chain.Addresses {
				if addressNameToFind == "" || strings.Contains(strings.ToLower(addressName), strings.ToLower(addressNameToFind)) {
					printChainInfo(addressName, address)
				}
			}
			return "", &chain
		}

		for addressName, address := range chain.Addresses {
			if address == addressToFind {
				fmt.Printf("Chain: %s\n", chain.Chain)
				fmt.Printf("Network: %s\n", superchain.Name)
				fmt.Printf("  Name: %s\n", addressName)
				printChainInfo(addressName, address)
				return addressName, &chain
			}
		}
	}
	return "", nil
}

func getEtherscanURL(address string, isTestnet bool) string {
	if isTestnet {
		return fmt.Sprintf("https://sepolia.etherscan.io/address/%s", address)
	}
	return fmt.Sprintf("https://etherscan.io/address/%s", address)
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

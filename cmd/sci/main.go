package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

// Superchain registry url: https://raw.githubusercontent.com/ethereum-optimism/superchain-registry/refs/heads/main
func main() {
	app := &cli.App{
		Name:  "superchain-insights",
		Usage: "A tool for gathering Superchain insights",
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
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

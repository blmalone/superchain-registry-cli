package main

import (
	"fmt"
	"os"

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
						// Simple chain name from superchain-registry corresponds to the file names in this directory: https://github.com/ethereum-optimism/superchain-registry/tree/821930ffec82ed6095130951e25bc1322190e0d9/superchain/configs/mainnet
						Usage: "Chain to filter by - simple chain name from superchain-registry (e.g. op, zora, base) or full chain name from https://github.com/ethereum-lists/chains (e.g. 'Metal L2', 'OP Mainnet')",
					},
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Output data as JSON",
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
					verbose := c.Bool("verbose")
					addressName := c.String("address-name")
					isJson := c.Bool("json")

					GetAddresses(superchain.OPChains, chain, address, addressName, testnet, verbose, isJson)
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

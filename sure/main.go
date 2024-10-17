package main

import (
	"fmt"
	"os"
	"sort"

	"text/tabwriter"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/urfave/cli/v2"
)

func CreateApp() *cli.App {
	return &cli.App{
		Name:  "sure",
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
				Usage:   "List all chains in the superchain in superchain registry",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "target",
						Aliases: []string{"tg"},
						Usage:   "Target chain to filter by (mainnet, sepolia, sepolia-dev-0)",
						Value:   "mainnet", // Default to mainnet
					},
				},
				Action: func(c *cli.Context) error {
					target := c.String("target")

					// Sort the chains by name
					chains := make([]*superchain.ChainConfig, 0, len(superchain.OPChains))
					for _, chain := range superchain.OPChains {
						chains = append(chains, chain)
					}
					sort.Slice(chains, func(i, j int) bool {
						return chains[i].Name < chains[j].Name
					})

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
					fmt.Fprintln(w, "Chain name\tShort name\tChain ID\tRPC URL\tSuperchain Level\tStandard Chain Candidate")
					fmt.Fprintln(w, "----------\t----------\t--------\t-------\t----------------\t------------------------")

					for _, chain := range chains {
						if chain.Superchain != target {
							continue
						}
						fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%d\t%t\n", chain.Name, chain.Chain, chain.ChainID, chain.PublicRPC, chain.SuperchainLevel, chain.StandardChainCandidate)
					}

					w.Flush()
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
					&cli.StringFlag{
						Name:    "target",
						Aliases: []string{"tg"},
						Usage:   "Target chain to filter by (mainnet, sepolia, sepolia-dev-0)",
						Value:   "mainnet",
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
					superchainTarget := c.String("target")
					addressName := c.String("address-name")
					isJson := c.Bool("json")

					GetAddresses(superchain.OPChains, chain, address, addressName, superchainTarget, isJson)
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

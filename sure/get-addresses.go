package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/urfave/cli/v2"
)

func GetAddresses(ctx *cli.Context, opChains map[uint64]*superchain.ChainConfig, chainName string, addressToFind string, addressNameToFind string, superchainTarget string, isJson bool) error {
	jsonResult := make(map[string]interface{})

	relevantSuperchain := getRelevantSuperchain(superchainTarget)
	chainExists := false

	chains := make([]*superchain.ChainConfig, 0, len(superchain.OPChains))
	for _, chain := range superchain.OPChains {
		chains = append(chains, chain)
	}
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].Name < chains[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	if !isJson {
		fmt.Fprintln(w, "Network\tAddress Name                    \tAddress")
		fmt.Fprintln(w, "-------\t--------------------------------\t-------")
	}

	for _, chain := range chains {
		if !isChainMatching(chain, chainName, relevantSuperchain) {
			continue // Skip chains that do not match the criteria
		}

		chainExists = true
		namedAddresses := ConvertAddressListToNamedAddresses(chain.Addresses)

		if addressToFind == "" {
			collectAddressNameSearchResults(relevantSuperchain, chain, namedAddresses, addressNameToFind, isJson, jsonResult, w)
		} else {
			collectAddressSearchResults(relevantSuperchain, chain, namedAddresses, addressToFind, isJson, jsonResult, w)
		}
	}

	if !chainExists {
		fmt.Fprintf(os.Stderr, "Error: Chain '%s' not found\n\n", chainName)
		os.Exit(1)
	}

	if isJson {
		outputJsonResults(jsonResult)
	} else {
		w.Flush() // Flush the writer to output the table
	}
	return nil
}

// Helper function to determine the relevant superchain
func getRelevantSuperchain(superchainTarget string) *superchain.Superchain {
	superchain := superchain.Superchains[superchainTarget]
	if superchain == nil {
		fmt.Fprintf(os.Stderr, "Error: Superchain target %s not found\n\n", superchainTarget)
		os.Exit(1)
	}
	return superchain
}

// Helper function to output JSON results
func outputJsonResults(jsonResult map[string]interface{}) {
	jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func isChainMatching(chain *superchain.ChainConfig, chainName string, relevantSuperchain *superchain.Superchain) bool {
	if chainName != "" && !strings.EqualFold(chain.Chain, chainName) && !strings.EqualFold(chain.Name, chainName) {
		return false
	}
	const (
		mainnet     = "mainnet"
		sepolia     = "sepolia"
		sepoliaDev0 = "sepolia-dev-0"
	)

	if (relevantSuperchain.Superchain == mainnet) && (chain.Superchain == sepolia || chain.Superchain == sepoliaDev0) {
		return false
	}
	if (relevantSuperchain.Superchain == sepolia) && (chain.Superchain == mainnet || chain.Superchain == sepoliaDev0) {
		return false
	}
	if (relevantSuperchain.Superchain == sepoliaDev0) && (chain.Superchain == mainnet || chain.Superchain == sepolia) {
		return false
	}
	return true
}

func collectAddressNameSearchResults(relevantSuperchain *superchain.Superchain, chain *superchain.ChainConfig, namedAddresses []NamedAddress, addressNameToFind string, isJson bool, jsonResults map[string]interface{}, w *tabwriter.Writer) {
	jsonResults[chain.Chain] = make(map[string]interface{})
	chainMap, _ := jsonResults[chain.Chain].(map[string]interface{})
	chainMap["network"] = chain.Name
	chainMap["addrs"] = make(map[string]interface{})
	addressMap := chainMap["addrs"].(map[string]interface{})
	for _, namedAddress := range namedAddresses {
		if addressNameToFind == "" || strings.Contains(strings.ToLower(namedAddress.Name), strings.ToLower(addressNameToFind)) {
			if namedAddress.Name == "SuperchainConfig" { // TODO: This is a hack to get the superchain config address
				namedAddress.Address = *relevantSuperchain.Config.SuperchainConfigAddr
			}
			if isJson {
				addressMap[namedAddress.Name] = namedAddress.Address.String()
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t\t\n", chain.Name, namedAddress.Name, FormatAddress(namedAddress.Address.String(), isTestnetSuperchain(chain.Superchain)))
			}
		}
	}
}

func collectAddressSearchResults(relevantSuperchain *superchain.Superchain, chain *superchain.ChainConfig, namedAddresses []NamedAddress, addressToFind string, isJson bool, jsonResults map[string]interface{}, w *tabwriter.Writer) {
	addressesProperty := "addrs"
	jsonResults["network"] = chain.Name
	jsonResults["chain"] = chain.Chain
	for _, namedAddress := range namedAddresses {
		if namedAddress.Name == "SuperchainConfig" { // TODO: This is a hack to get the superchain config address
			namedAddress.Address = *relevantSuperchain.Config.SuperchainConfigAddr
		}
		if namedAddress.Address.String() == addressToFind {
			if isJson {
				if _, exists := jsonResults[addressesProperty]; !exists {
					jsonResults[addressesProperty] = make(map[string]interface{})
				}
				chainData := jsonResults[addressesProperty].(map[string]interface{})
				chainData[namedAddress.Name] = namedAddress.Address.String()
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t\t\n", chain.Name, namedAddress.Name, FormatAddress(namedAddress.Address.String(), isTestnetSuperchain(chain.Superchain)))
			}
		}
	}
}

func FormatAddress(address string, isTestnet bool) string {
	if address == "0x0000000000000000000000000000000000000000" {
		return fmt.Sprintf("  \033]8;;%s\033\\%s\033]8;;\033\\\n", "N/A", "N/A")
	}
	return CreateHyperlinkedAddress(GetEtherscanURL(address, isTestnet))
}

func CreateHyperlinkedAddress(etherscanAddressURL string) string {
	addressPart := etherscanAddressURL[len(etherscanAddressURL)-42:]
	return fmt.Sprintf("  \033]8;;%s\033\\%s\033]8;;\033\\\n", etherscanAddressURL, addressPart)
}

func GetEtherscanURL(address string, isTestnet bool) string {
	baseURL := "https://etherscan.io/address/%s"
	if isTestnet {
		baseURL = "https://sepolia.etherscan.io/address/%s"
	}
	return fmt.Sprintf(baseURL, address)
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

func isTestnetSuperchain(superchainName string) bool {
	return superchainName == "sepolia" || superchainName == "sepolia-dev-0"
}

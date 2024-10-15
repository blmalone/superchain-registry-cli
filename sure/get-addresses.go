package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum-optimism/superchain-registry/superchain"
)

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

func GetAddresses(opChains map[uint64]*superchain.ChainConfig, chainName, addressToFind, addressNameToFind string, isTestnet bool, isVerbose bool, isJson bool) {
	jsonResult := make(map[string]interface{})

	for _, chain := range opChains {
		if !isChainMatching(chain, chainName, isTestnet) {
			continue // chain doesn't match criteria, skip this chain
		}

		namedAddresses := ConvertAddressListToNamedAddresses(chain.Addresses)

		if addressToFind == "" {
			collectAddressNameSearchResults(chain, namedAddresses, addressNameToFind, isTestnet, isJson, jsonResult)
		} else {
			collectAddressSearchResults(chain, namedAddresses, addressToFind, isTestnet, isJson, jsonResult)
		}
	}

	if isJson {
		jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	}
}

func isChainMatching(chain *superchain.ChainConfig, chainName string, isTestnet bool) bool {
	if chainName != "" && !strings.EqualFold(chain.Chain, chainName) && !strings.EqualFold(chain.Name, chainName) {
		return false
	}
	if isTestnet && chain.Superchain != "sepolia" || !isTestnet && chain.Superchain != "mainnet" {
		return false
	}
	return true
}

func collectAddressNameSearchResults(chain *superchain.ChainConfig, namedAddresses []NamedAddress, addressNameToFind string, isTestnet bool, isJson bool, jsonResults map[string]interface{}) {
	if !isJson {
		printChainAndNetwork(chain)
	}
	jsonResults[chain.Chain] = make(map[string]interface{})
	chainMap, _ := jsonResults[chain.Chain].(map[string]interface{})
	chainMap["network"] = chain.Name
	chainMap["addrs"] = make(map[string]interface{})
	addressMap := chainMap["addrs"].(map[string]interface{})
	for _, namedAddress := range namedAddresses {
		if addressNameToFind == "" || strings.Contains(strings.ToLower(namedAddress.Name), strings.ToLower(addressNameToFind)) {
			if isJson {
				addressMap[namedAddress.Name] = namedAddress.Address.String()
			} else {
				printChainInfo(namedAddress.Name, namedAddress.Address.String(), isTestnet)
			}
		}
	}
}

func collectAddressSearchResults(chain *superchain.ChainConfig, namedAddresses []NamedAddress, addressToFind string, isTestnet bool, isJson bool, jsonResults map[string]interface{}) {
	addressesProperty := "addrs"
	jsonResults["network"] = chain.Name
	jsonResults["chain"] = chain.Chain
	for _, namedAddress := range namedAddresses {
		if namedAddress.Address.String() == addressToFind {
			if isJson {
				if _, exists := jsonResults[addressesProperty]; !exists {
					jsonResults[addressesProperty] = make(map[string]interface{})
				}
				chainData := jsonResults[addressesProperty].(map[string]interface{})
				chainData[namedAddress.Name] = namedAddress.Address.String()
			} else {
				printChainAndNetwork(chain)
				printChainInfo(namedAddress.Name, namedAddress.Address.String(), isTestnet)
			}
		}
	}
}

func printChainInfo(addressName, address string, isTestnet bool) {
	if address == "0x0000000000000000000000000000000000000000" {
		fmt.Printf("  %s: %s\n", addressName, "<n/a>")
	} else {
		fmt.Print(CreateHyperlinkedAddress(addressName, GetEtherscanURL(address, isTestnet)))
	}
}

func printChainAndNetwork(chain *superchain.ChainConfig) {
	fmt.Printf("Chain: %s\n", chain.Chain)
	fmt.Printf("Network: %s\n", chain.Name)
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

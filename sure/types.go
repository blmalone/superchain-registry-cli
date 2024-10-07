package main

type Parent struct {
	Type  string `toml:"type"`
	Chain string `toml:"chain"`
}

// Define the Chain struct that includes the Parent struct
type Chain struct {
	Name                 string   `toml:"name"`
	Identifier           string   `toml:"identifier"`
	ChainID              int      `toml:"chain_id"`
	RPC                  []string `toml:"rpc"`
	Explorers            []string `toml:"explorers"`
	SuperchainLevel      int      `toml:"superchain_level"`
	DataAvailabilityType string   `toml:"data_availability_type"`
	Parent               Parent   `toml:"parent"`
}

// Define the ChainList struct to match the structure of the TOML file
type ChainList struct {
	Chains []Chain `toml:"chains"`
}

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

type Config struct {
	Name                 string `json:"Name"`
	L1                   L1     `json:"L1"`
	ProtocolVersionsAddr string `json:"ProtocolVersionsAddr"`
	SuperchainConfigAddr string `json:"SuperchainConfigAddr"`
}

type L1 struct {
	ChainID   int    `json:"ChainID"`
	PublicRPC string `json:"PublicRPC"`
	Explorer  string `json:"Explorer"`
}

type ChainInfo struct {
	Name                   string            `json:"Name"`
	L2ChainID              int               `json:"l2_chain_id"`
	PublicRPC              string            `json:"PublicRPC"`
	SequencerRPC           string            `json:"SequencerRPC"`
	Explorer               string            `json:"Explorer"`
	SuperchainLevel        int               `json:"SuperchainLevel"`
	StandardChainCandidate bool              `json:"StandardChainCandidate"`
	SuperchainTime         *int              `json:"SuperchainTime"`
	BatchInboxAddress      string            `json:"batch_inbox_address"`
	Superchain             string            `json:"Superchain"`
	Chain                  string            `json:"Chain"`
	CanyonTime             int               `json:"canyon_time"`
	DeltaTime              int               `json:"delta_time"`
	EcotoneTime            int               `json:"ecotone_time"`
	FjordTime              int               `json:"fjord_time"`
	GraniteTime            int               `json:"granite_time"`
	BlockTime              int               `json:"block_time"`
	SeqWindowSize          int               `json:"seq_window_size"`
	MaxSequencerDrift      int               `json:"max_sequencer_drift"`
	DataAvailabilityType   string            `json:"DataAvailabilityType"`
	Optimism               Optimism          `json:"optimism"`
	GasPayingToken         *string           `json:"GasPayingToken"`
	Genesis                Genesis           `json:"genesis"`
	Addresses              map[string]string `json:"Addresses"`
}

type Optimism struct {
	EIP1559Elasticity        int `json:"eip1559Elasticity"`
	EIP1559Denominator       int `json:"eip1559Denominator"`
	EIP1559DenominatorCanyon int `json:"eip1559DenominatorCanyon"`
}

type Genesis struct {
	L1           GenesisL1    `json:"l1"`
	L2           GenesisL2    `json:"l2"`
	L2Time       int          `json:"l2_time"`
	SystemConfig SystemConfig `json:"system_config"`
}

type GenesisL1 struct {
	Hash   string `json:"hash"`
	Number int    `json:"number"`
}

type GenesisL2 struct {
	Hash   string `json:"hash"`
	Number int    `json:"number"`
}

type SystemConfig struct {
	BatcherAddr string `json:"batcherAddr"`
	Overhead    string `json:"overhead"`
	Scalar      string `json:"scalar"`
	GasLimit    int    `json:"gasLimit"`
}

type Superchain struct {
	Name       string      `json:"name"`
	Config     Config      `json:"config"`
	ChainInfos []ChainInfo `json:"chains"`
}

type Superchains struct {
	Superchains []Superchain `json:"superchains"`
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func formatAddress(name, address string, isTestnet bool) string {
	return CreateHyperlinkedAddress(name, GetEtherscanURL(address, isTestnet))
}

func runTest(t *testing.T, name string, args []string, wantOutput []string, wantErr bool) {
	t.Run(name, func(t *testing.T) {
		var output bytes.Buffer

		// Save the original os.Stdout
		oldStdout := os.Stdout

		// Redirect stdout to the buffer
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the CLI command
		app := CreateApp()
		err := app.Run(args)

		// Close the writer and restore os.Stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output into the buffer
		_, _ = output.ReadFrom(r)

		if (err != nil) != wantErr {
			t.Errorf("got error = %v, wantErr %v", err, wantErr)
		}

		gotOutput := strings.TrimSpace(output.String())
		for _, want := range wantOutput {
			if !strings.Contains(gotOutput, strings.TrimSpace(want)) {
				t.Errorf("got output = %v, want %v", gotOutput, want)
			}
		}
	})
}

func TestListChains(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "List all chains in the superchain",
			args:       []string{"sure", "list"},
			wantOutput: []string{"Chain: op\nNetwork: OP Mainnet"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		runTest(t, tt.name, tt.args, tt.wantOutput, tt.wantErr)
	}
}

func TestGetAddresses(t *testing.T) {
	opMainnetAddressManager := formatAddress("AddressManager", "0xdE1FCfB0851916CA5101820A69b13a4E276bd81F", false)
	opMainnetL1CrossDomainMessengerProxy := formatAddress("L1CrossDomainMessengerProxy", "0x25ace71c97B33Cc4729CF772ae268934F7ab5fA1", false)
	allChainAddressesOPMainnet := fmt.Sprintf("Chain: op\nNetwork: OP Mainnet\n%s%s", opMainnetAddressManager, opMainnetL1CrossDomainMessengerProxy)

	opSepoliaAddressManager := formatAddress("AddressManager", "0x9bFE9c5609311DF1c011c47642253B78a4f33F4B", true)
	opSepoliaL1CrossDomainMessengerProxy := formatAddress("L1CrossDomainMessengerProxy", "0x58Cc85b8D04EA49cC6DBd3CbFFd00B4B8D6cb3ef", true)
	allChainAddressesOPSepolia := fmt.Sprintf("Chain: op\nNetwork: OP Sepolia Testnet\n%s%s", opSepoliaAddressManager, opSepoliaL1CrossDomainMessengerProxy)

	l1CrossDomainMessengerAddressOPMainnet := fmt.Sprintf("Chain: op\nNetwork: OP Mainnet\n%s", opMainnetL1CrossDomainMessengerProxy)
	zoraMainnetL1CrossDomainMessengerProxy := formatAddress("L1CrossDomainMessengerProxy", "0xdC40a14d9abd6F410226f1E6de71aE03441ca506", false)
	l1CrossDomainMessengerAddressZora := fmt.Sprintf("Chain: zora\nNetwork: Zora\n%s", zoraMainnetL1CrossDomainMessengerProxy)

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "Find all chain addresses",
			args:       []string{"sure", "get-addresses", "--chain", "op"},
			wantOutput: []string{allChainAddressesOPMainnet},
			wantErr:    false,
		},
		{
			name:       "Find all testnet chain addresses",
			args:       []string{"sure", "get-addresses", "--chain", "op", "-t"},
			wantOutput: []string{allChainAddressesOPSepolia},
			wantErr:    false,
		},
		{
			name:       "Find specific address by name",
			args:       []string{"sure", "get-addresses", "--chain", "op", "-an", "L1CrossDomainMessengerProxy"},
			wantOutput: []string{fmt.Sprintf("Chain: op\nNetwork: OP Mainnet\n%s", formatAddress("L1CrossDomainMessengerProxy", "0x25ace71c97B33Cc4729CF772ae268934F7ab5fA1", false))},
			wantErr:    false,
		},
		{
			name:       "Find specific address by name across all chains",
			args:       []string{"sure", "get-addresses", "-an", "L1Cross"},
			wantOutput: []string{l1CrossDomainMessengerAddressOPMainnet, l1CrossDomainMessengerAddressZora}, // should have more than one change
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		runTest(t, tt.name, tt.args, tt.wantOutput, tt.wantErr)
	}
}

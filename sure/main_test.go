package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

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
	expectedHeader := []string{
		"Chain name", "Short name", "Chain ID", "RPC URL", "Superchain Level", "Standard Chain Candidate",
		"Base", "base", "8453", "https://mainnet.base.org", "0", "true",
		"Binary Mainnet", "tbn", "624", "https://rpc.zero.thebinaryholdings.com", "0", "false",
		"Zora", "zora", "7777777", "https://rpc.zora.energy", "0", "true",
	}

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "List all chains in the superchain",
			args:       []string{"sure", "list"},
			wantOutput: expectedHeader,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		runTest(t, tt.name, tt.args, tt.wantOutput, tt.wantErr)
	}
}

func TestGetAddresses(t *testing.T) {

	expectedHeader := []string{
		"Network", "Address Name", "Address",
	}
	opfirstRow := []string{
		"OP Mainnet", "AddressManager", FormatAddress("0xdE1FCfB0851916CA5101820A69b13a4E276bd81F", false),
	}
	oplastRow := []string{
		"OP Mainnet", "DAChallengeAddress", FormatAddress("0x0000000000000000000000000000000000000000", false),
	}

	opSepoliaFirstRow := []string{
		"OP Sepolia Testnet", "AddressManager", FormatAddress("0x9bFE9c5609311DF1c011c47642253B78a4f33F4B", true),
	}
	opSepoliaLastRow := []string{
		"OP Sepolia Testnet", "DAChallengeAddress", FormatAddress("0x0000000000000000000000000000000000000000", true),
	}

	specificAddressOpRow := []string{
		"OP Mainnet", "L1CrossDomainMessengerProxy", FormatAddress("0x25ace71c97B33Cc4729CF772ae268934F7ab5fA1", false),
	}
	specificAddressZoraRow := []string{
		"Zora", "L1CrossDomainMessengerProxy", FormatAddress("0xdC40a14d9abd6F410226f1E6de71aE03441ca506", false),
	}

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "Find all chain addresses",
			args:       []string{"sure", "get-addresses", "--chain", "op"},
			wantOutput: append(append(expectedHeader, opfirstRow...), oplastRow...),
			wantErr:    false,
		},
		{
			name:       "Find all testnet chain addresses",
			args:       []string{"sure", "get-addresses", "--chain", "op", "-tg", "sepolia"},
			wantOutput: append(append(expectedHeader, opSepoliaFirstRow...), opSepoliaLastRow...),
			wantErr:    false,
		},
		{
			name:       "Find specific address by name",
			args:       []string{"sure", "get-addresses", "--chain", "op", "-an", "L1CrossDomainMessengerProxy"},
			wantOutput: append(expectedHeader, specificAddressOpRow...),
			wantErr:    false,
		},
		{
			name:       "Find specific address by name across all chains",
			args:       []string{"sure", "get-addresses", "-an", "L1Cross"},
			wantOutput: append(append(expectedHeader, specificAddressZoraRow...), specificAddressOpRow...),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		runTest(t, tt.name, tt.args, tt.wantOutput, tt.wantErr)
	}
}

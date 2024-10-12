# superchain-registry-cli

A command line tool that wraps Optimisms [superchain-registry](https://github.com/ethereum-optimism/superchain-registry) repository.

## Introduction

The `superchain-registry-cli` is designed to simplify interactions with the superchain-registry, providing a user-friendly command line interface for developers and operators.

## Prerequisites

- Ensure you have [Go](https://formulae.brew.sh/formula/go) installed on your system.

## Installation

- Install `superchain-registry-cli`:
    - `go install github.com/blmalone/superchain-registry-cli/sure@latest`
 
## Usage

After installation, run `sure --help` for a full breakdown of the available functionality.

## Examples

- List all chains in the superchain (default is mainnet): `sure ls` or `sure list` or `sure ls --testnet`
- Get all addresses for a chain: `sure ga -c op` or `sure get-addresses -c op` or `sure ga -c op -t`
- Get a specific address by name (fuzzy match on `--address-name`): `sure ga -c zora -an l1` or `go run sure/*.go ga -c zora -an L1CrossDomainMessengerProxy`
- Get all addresses for a given name across the superchain (fuzzy match on `--address-name`): `sure ga -an L1StandardBridge`

## Contributing

Contributions are welcome! Fork this repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# superchain-registry-cli

<p align="center">
    <img src="./superchain-registry-cli-dalle.png" alt="Generated Dall-e given this README as a prompt" width="400"/ title="Generated Dall-e given this README as a prompt">
</p>

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

- List all chains in the superchain (default is mainnet): 
    ```bash 
        sure ls
        sure list
        sure ls --testnet
    ``` 
- Get all addresses for a chain: 
    ```bash 
        sure get-addresses --chain op
        sure ga -c op
        sure ga -c op -t
    ```
- Get a specific address by name (fuzzy match on `--address-name`): 
    ```bash
        sure ga -c zora --address-name l1 # all addresses containing "l1" - not case sensitive
        sure ga -c zora -an L1CrossDomainMessengerProxy
    ```

- Get all addresses for a given name across the superchain (fuzzy match on `--address-name`): 
    ```bash
        sure ga -an L1StandardBridge
    ```

- Usage with [cast](https://book.getfoundry.sh/cast/):
    ```bash
        cast call $(sure ga -c op -an L1Standard --json | jq -r '.addrs.L1StandardBridgeProxy') "version()(string)"

        # When you know there will be only one address returned
        cast call $(sure ga -c op -an L1Standard --json | jq -r '.addrs | to_entries | .[0].value') "version()(string)"
    ```

## Contributing

Contributions are welcome! Fork this repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

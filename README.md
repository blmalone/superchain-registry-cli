# superchain-insights
A command line tool that wraps Optimisms [superchain-registry](https://github.com/ethereum-optimism/superchain-registry) repository.


## Installation

- Install [Go](https://formulae.brew.sh/formula/go)
- Install `superchain-insights`
    - `go install github.com/blmalone/superchain-insights/cmd/sci@latest`
- Set `SUPERCHAIN_REGISTRY_URL` environment variable
    - e.g. `export SUPERCHAIN_REGISTRY_URL=https://raw.githubusercontent.com/ethereum-optimism/superchain-registry/refs/heads/main`
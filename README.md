# Vega Protocol - CLI Market Maker

## Getting Started

1. Build with `go build`
2. Create a file in the same directory as the binary called `.secret`
3. Generate a new [BIP39 passphrase](https://iancoleman.io/bip39/) and save it in the `.secret` file
4. Edit `markets.json` so that it contains the markets you want to trade
5. Run with `./vega-cli-mm`

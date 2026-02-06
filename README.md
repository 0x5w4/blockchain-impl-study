# blockchain-impl-study

Usage
-----

Prerequisites:

- Go 1.18+ installed

Build:

```bash
go build ./...
```

Run CLI:

Create a new blockchain (creates data under `./tmp/blocks_node_1`):

```bash
./blockchain-impl-study createblockchain -address YOUR_ADDRESS
```

Add a block (creates a coinbase tx and mines a block):

```bash
./blockchain-impl-study addblock -data "some reward message"
```

Print the chain:

```bash
./blockchain-impl-study printchain
```

Create a wallet (saves wallets to `wallet_%s.dat` using node id):

```bash
./blockchain-impl-study createwallet
```

Notes:

- The project stores blockchain data in `./tmp/blocks_node_1` by default.
- This repo uses a Bitcoin-like transaction and block serialization for learning purposes.
- If you change node id or data directories, adjust commands accordingly.
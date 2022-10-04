# go-zerokit-rln

Go wrappers for [zerokit's RLN](https://github.com/vacp2p/zerokit)

### Building this library
```
git clone https://github.com/status-im/go-zerokit-rln
cd go-zerokit-rln
git submodule init
git submodule update --recursive
make
```

To generate smaller static libraries, before `make`, edit `./zerokit/rln/Cargo.toml` and use `branch = "no-ethers-core"` for `ark-circom`

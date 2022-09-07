# go-zerokit-rln

Go wrappers for [zerokit's RLN](https://github.com/vacp2p/zerokit)

### Building this library

```
make
```

Some architectures are not available in cross unless they're locally build. This [PR](https://github.com/cross-rs/cross/pull/591) will update ubuntu base version on cross. But while it's merged, build them locally. To build them locally execute the following instructions (adapted from [here](https://github.com/cross-rs/cross/wiki/FAQ#newer-linux-versions)):

```
git clone --single-branch --depth 1 --branch increment_versions https://github.com/Alexhuszagh/cross
cd cross
cargo build-docker-image x86_64-pc-windows-gnu
cargo build-docker-image aarch64-unknown-linux-gnu
cargo build-docker-image x86_64-unknown-linux-gnu
cargo build-docker-image arm-unknown-linux-gnueabi
cargo build-docker-image i686-unknown-linux-gnu
cargo build-docker-image arm-unknown-linux-gnueabihf
cargo build-docker-image mips-unknown-linux-gnu
cargo build-docker-image mips64-unknown-linux-gnuabi64
cargo build-docker-image mips64el-unknown-linux-gnuabi64
cargo build-docker-image mipsel-unknown-linux-gnu
```

`i686-pc-windows-gnu`, and `mips` / `mips64` are currently not supported 
# go-zerokit-rln

Go wrappers for [zerokit's RLN](https://github.com/vacp2p/zerokit)



### Updating vacp2p/zerokit

To overcome the limit of 500mb github has for repositories, go-zerokit-rln depends on 3 projects:
- https://github.com/waku-org/go-zerokit-rln-apple
- https://github.com/waku-org/go-zerokit-rln-arm
- https://github.com/waku-org/go-zerokit-rln-x86_64

Zerokit must be updated in these 3 repositories. The instructions are the same for each of the architectures,
except for `-apple` which require macos to be executed. You need to have docker and rust installed.

```bash
export GO_RLN_ARCH=x86_64       # Replace this for x86_64, arm or apple
export ZEROKIT_COMMIT=master    # Use a commit, branch or tag

git clone https://github.com/waku-org/go-zerokit-rln_${GO_RLN_ARCH}
cd go-zerokit-rln-${GO_RLN_ARCH}
git submodule init
git submodule update --recursive
cd zerokit
git pull
git checkout ${ZEROKIT_COMMIT}
cd ..
make
git add zerokit
git add libs/*/librln.a
git commit -m "chore: bump zerokit"
git push
```

Once you execute the previous commands for each one of the architectures, update go.mod:
```bash
cd /path/to/go-zerokit-rln
go get github.com/waku-org/go-zerokit-rln-apple@latest
go get github.com/waku-org/go-zerokit-rln-arm@latest
go get github.com/waku-org/go-zerokit-rln-x86_64@latest
git checkout master
git add go.mod
git add go.sum
git commit -m "chore: bump zerokit"
git push
```

And later in go-waku, update the go-zerokit-rln dependency with
```
cd /path/to/go-waku
git fetch
git checkout -b `date +"%Y%m%d%H%M%S"-bump-zerokit` origin/master
go get github.com/waku-org/go-zerokit-rln@latest
git add go.mod
git add go.sum
git commit -m "chore: bump go-zerokit-rln"
git push
````
And create a PR


# Adding a new architecture

1. Depending on the architecture/platform you want to add, clone one of these repositories.
- https://github.com/waku-org/go-zerokit-rln-apple
- https://github.com/waku-org/go-zerokit-rln-arm
- https://github.com/waku-org/go-zerokit-rln-x86_64
2. Edit `./scripts/build.sh` to add the new architecture. Refer to this list of [supported targets](https://github.com/cross-rs/cross#supported-targets)
3. Edit `./rln/link.go` to add a **cgo** build tag for the new architecture
4. Commit the changes and then execute
```bash
export GO_RLN_ARCH=x86_64       # Replace this for the platform you just updated: x86_64, arm or apple
cd /path/to/go-zerokit-rln
git checkout master
git pull
go get github.com/waku-org/go-zerokit-rln-${GO_RLN_ARCH}@latest
git add go.mod
git add go.sum
git commit -m "chore: bump go-zerokit-rln-${GO_RLN_ARCH}"
git push
```

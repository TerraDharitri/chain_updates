# drt-go-chain-tools-accounts-manager

The go implementation for the drt-go-chain-tools-accounts-manager

- This application will be responsible to fetch all dharitri-accounts that have staked an amount of REWA tokens. 
After the accounts are fetched from API it will process all the information, and it will index the new data 
in a new Elaticsearch index.

- The new Elastisearch index will contain all the accounts that have balance and also information 
about the staked balance and energy.

### Sources of accounts with stake

- This go client will fetch information from:
    1. Validators system smart contract
    2. Delegation manager system smart contracts
    3. Legacy delegation smart contract
    4. Energy smart contract
    

### Installation and running


#### Step 1: install & configure go:

The installation of go should proceed as shown in official golang 
installation guide https://golang.org/doc/install . In order to run the node, minimum golang 
version should be 1.12.4.


#### Step 2: clone the repository and build the binary:

```
 $ git clone https://github.com/TerraDharitri/drt-go-chain-tools-accounts-manager.git
 $ cd accounts-manager-go/cmd/manager
 $ GO111MODULE=on go mod vendor
 $ go build
```

#### Step 3: run manager
```
 $ ./manager --config="pathToConfig/config.toml"
```

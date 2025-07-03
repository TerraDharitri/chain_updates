# Command Line Interface

## Overview

**drtpy** exposes a number of CLI **commands**, organized within **groups**.


```
$ drtpy --help
usage: drtpy [-h] [-v] [--verbose] COMMAND-GROUP [-h] COMMAND ...

-----------
DESCRIPTION
-----------
drtpy is part of the dharitri-sdk and consists of Command Line Tools and Python SDK
for interacting with the Blockchain (in general) and with Smart Contracts (in particular).

drtpy targets a broad audience of users and developers.

See:
 - https://docs.dharitri.org/sdk-and-tools/sdk-py
 - https://docs.dharitri.org/sdk-and-tools/sdk-py/drtpy-cli


COMMAND GROUPS:
  {config-wallet,contract,tx,validator,ledger,wallet,validator-wallet,deps,config,localnet,data,staking-provider,dns,faucet,multisig,governance,config-env,get}

TOP-LEVEL OPTIONS:
  -h, --help            show this help message and exit
  -v, --version         show program's version number and exit
  --verbose
  --log-level {debug,info,warning,error}
                        default: info

----------------------
COMMAND GROUPS summary
----------------------
config-wallet                  Configure Dharitri CLI to use a default wallet.
contract                       Deploy, upgrade and interact with Smart Contracts
tx                             Create and broadcast Transactions
validator                      Stake, UnStake, UnBond, Unjail and other actions useful for Validators
ledger                         Get Ledger App addresses and version
wallet                         Create wallet, derive secret key from mnemonic, bech32 address helpers etc.
validator-wallet               Create a validator wallet, sign and verify messages and convert a validator wallet to a hex secret key.
deps                           Manage dependencies or dharitri-sdk modules
config                         Configure Dharitri CLI (default values etc.)
localnet                       Set up, start and control localnets
data                           Data manipulation omnitool
staking-provider               Staking provider omnitool
dns                            Operations related to the Domain Name Service
faucet                         Get xREWA on Devnet or Testnet
multisig                       Deploy and interact with the Multisig Smart Contract
governance                     Propose, vote and interact with the governance contract.
config-env                     Configure Dharitri CLI to use specific environment values.
get                            Get info from the network.

```
## Group **Contract**


```
$ drtpy contract --help
usage: drtpy contract COMMAND [-h] ...

Deploy, upgrade and interact with Smart Contracts

COMMANDS:
  {deploy,call,upgrade,query,verify,unverify,reproducible-build,build}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
deploy                         Deploy a Smart Contract.
call                           Interact with a Smart Contract (execute function).
upgrade                        Upgrade a previously-deployed Smart Contract.
query                          Query a Smart Contract (call a pure function)
verify                         Verify the authenticity of the code of a deployed Smart Contract
unverify                       Unverify a previously verified Smart Contract
reproducible-build             Build a Smart Contract and get the same output as a previously built Smart Contract
build                          Build a Smart Contract project. This command is DISABLED.

```
### Contract.Deploy


```
$ drtpy contract deploy --help
usage: drtpy contract deploy [-h] ...

Deploy a Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash",
    "contractAddress": "the address of the contract",
    "transactionOnNetwork": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "simulation": {
        "execution": {
            "...": "..."
        },
        "cost": {
            "...": "..."
        }
    }
}

options:
  -h, --help                                     show this help message and exit
  --bytecode BYTECODE                            the file containing the WASM bytecode
  --abi ABI                                      the ABI file of the Smart Contract
  --metadata-not-upgradeable                     ‼ mark the contract as NOT upgradeable (default: upgradeable)
  --metadata-not-readable                        ‼ mark the contract as NOT readable (default: readable)
  --metadata-payable                             ‼ mark the contract as payable (default: not payable)
  --metadata-payable-by-sc                       ‼ mark the contract as payable by SC (default: not payable by SC)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --arguments ARGUMENTS [ARGUMENTS ...]          arguments for the contract transaction, as [number, bech32-address,
                                                 ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                 0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                a json file containing the arguments. ONLY if abi file is provided.
                                                 E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### Contract.Call


```
$ drtpy contract call --help
usage: drtpy contract call [-h] ...

Interact with a Smart Contract (execute function).

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash",
    "contractAddress": "the address of the contract",
    "transactionOnNetwork": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "simulation": {
        "execution": {
            "...": "..."
        },
        "cost": {
            "...": "..."
        }
    }
}

positional arguments:
  contract                                        🖄 the bech32 address of the Smart Contract

options:
  -h, --help                                      show this help message and exit
  --abi ABI                                       the ABI file of the Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --function FUNCTION                             the function to call
  --arguments ARGUMENTS [ARGUMENTS ...]           arguments for the contract transaction, as [number, bech32-address,
                                                  ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                  0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                 a json file containing the arguments. ONLY if abi file is provided.
                                                  E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --token-transfers TOKEN_TRANSFERS [TOKEN_TRANSFERS ...]
                                                  token transfers for transfer & execute, as [token, amount] E.g.
                                                  --token-transfers NFT-123456-0a 1 DCDT-987654 100000000
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)

```
### Contract.Upgrade


```
$ drtpy contract upgrade --help
usage: drtpy contract upgrade [-h] ...

Upgrade a previously-deployed Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash",
    "contractAddress": "the address of the contract",
    "transactionOnNetwork": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "simulation": {
        "execution": {
            "...": "..."
        },
        "cost": {
            "...": "..."
        }
    }
}

positional arguments:
  contract                                       🖄 the bech32 address of the Smart Contract

options:
  -h, --help                                     show this help message and exit
  --abi ABI                                      the ABI file of the Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --bytecode BYTECODE                            the file containing the WASM bytecode
  --metadata-not-upgradeable                     ‼ mark the contract as NOT upgradeable (default: upgradeable)
  --metadata-not-readable                        ‼ mark the contract as NOT readable (default: readable)
  --metadata-payable                             ‼ mark the contract as payable (default: not payable)
  --metadata-payable-by-sc                       ‼ mark the contract as payable by SC (default: not payable by SC)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --arguments ARGUMENTS [ARGUMENTS ...]          arguments for the contract transaction, as [number, bech32-address,
                                                 ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                 0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                a json file containing the arguments. ONLY if abi file is provided.
                                                 E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### Contract.Query


```
$ drtpy contract query --help
usage: drtpy contract query [-h] ...

Query a Smart Contract (call a pure function)

positional arguments:
  contract                               🖄 the bech32 address of the Smart Contract

options:
  -h, --help                             show this help message and exit
  --abi ABI                              the ABI file of the Smart Contract
  --proxy PROXY                          🔗 the URL of the proxy
  --function FUNCTION                    the function to call
  --arguments ARGUMENTS [ARGUMENTS ...]  arguments for the contract transaction, as [number, bech32-address, ascii
                                         string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000 0xabba
                                         str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE        a json file containing the arguments. ONLY if abi file is provided. E.g. [{
                                         'to': 'erd1...', 'amount': 10000000000 }]

```
### Contract.Verify


```
$ drtpy contract verify --help
usage: drtpy contract verify [-h] ...

Verify the authenticity of the code of a deployed Smart Contract

positional arguments:
  contract                                   🖄 the bech32 address of the Smart Contract

options:
  -h, --help                                 show this help message and exit
  --packaged-src PACKAGED_SRC                JSON file containing the source code of the contract
  --verifier-url VERIFIER_URL                the url of the service that validates the contract
  --docker-image DOCKER_IMAGE                the docker image used for the build
  --contract-variant CONTRACT_VARIANT        in case of a multicontract, specify the contract variant you want to verify
  --sender SENDER                            the alias of the wallet set in the address config
  --pem PEM                                  🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                          🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                        DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter the
                                             password.
  --ledger                                   🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type mnemonic
                                             or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME          🖄 the username of the sender
  --hrp HRP                                  The hrp used to convert the address to its bech32 representation
  --skip-confirmation, -y                    can be used to skip the confirmation prompt

```
### Contract.Unverify


```
$ drtpy contract unverify --help
usage: drtpy contract unverify [-h] ...

Unverify a previously verified Smart Contract

positional arguments:
  contract                                   🖄 the bech32 address of the Smart Contract

options:
  -h, --help                                 show this help message and exit
  --code-hash CODE_HASH                      the code hash of the contract
  --verifier-url VERIFIER_URL                the url of the service that validates the contract
  --sender SENDER                            the alias of the wallet set in the address config
  --pem PEM                                  🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                          🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                        DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter the
                                             password.
  --ledger                                   🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type mnemonic
                                             or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME          🖄 the username of the sender
  --hrp HRP                                  The hrp used to convert the address to its bech32 representation

```
### Contract.ReproducibleBuild


```
$ drtpy contract reproducible-build --help
usage: drtpy contract reproducible-build [-h] ...

Build a Smart Contract and get the same output as a previously built Smart Contract

positional arguments:
  project                              the project directory (default: current directory)

options:
  -h, --help                           show this help message and exit
  --debug                              set debug flag (default: False)
  --no-optimization                    bypass optimizations (for clang) (default: False)
  --no-wasm-opt                        do not optimize wasm files after the build (default: False)
  --cargo-target-dir CARGO_TARGET_DIR  for rust projects, forward the parameter to Cargo
  --wasm-symbols                       for rust projects, does not strip the symbols from the wasm output. Useful for
                                       analysing the bytecode. Creates larger wasm files. Avoid in production (default:
                                       False)
  --wasm-name WASM_NAME                for rust projects, optionally specify the name of the wasm bytecode output file
  --wasm-suffix WASM_SUFFIX            for rust projects, optionally specify the suffix of the wasm bytecode output file
  --docker-image DOCKER_IMAGE          the docker image tag used to build the contract
  --contract CONTRACT                  contract to build (contract name, as found in Cargo.toml)
  --no-docker-interactive
  --no-docker-tty
  --no-default-platform                do not set DOCKER_DEFAULT_PLATFORM environment variable to 'linux/amd64'

```
## Group **Transactions**


```
$ drtpy tx --help
usage: drtpy tx COMMAND [-h] ...

Create and broadcast Transactions

COMMANDS:
  {new,send,sign,relay}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
new                            Create a new transaction.
send                           Send a previously saved transaction.
sign                           Sign a previously saved transaction.
relay                          Relay a previously saved transaction.

```
### Transactions.New


```
$ drtpy tx new --help
usage: drtpy tx new [-h] ...

Create a new transaction.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --receiver RECEIVER                             🖄 the address of the receiver
  --receiver-username RECEIVER_USERNAME           🖄 the username of the receiver
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --data DATA                                     the payload, or 'memo' of the transaction (default: )
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --data-file DATA_FILE                           a file containing transaction data
  --token-transfers TOKEN_TRANSFERS [TOKEN_TRANSFERS ...]
                                                  token transfers for transfer & execute, as [token, amount] E.g.
                                                  --token-transfers NFT-123456-0a 1 DCDT-987654 100000000
  --outfile OUTFILE                               where to save the output (signed transaction, hash) (default: stdout)
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --proxy PROXY                                   🔗 the URL of the proxy
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Transactions.Send


```
$ drtpy tx send --help
usage: drtpy tx send [-h] ...

Send a previously saved transaction.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help         show this help message and exit
  --infile INFILE    input file (a previously saved transaction)
  --outfile OUTFILE  where to save the output (the hash) (default: stdout)
  --proxy PROXY      🔗 the URL of the proxy

```
### Transactions.Sign


```
$ drtpy tx sign --help
usage: drtpy tx sign [-h] ...

Sign a previously saved transaction.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --infile INFILE                                input file (a previously saved transaction)
  --outfile OUTFILE                              where to save the output (the signed transaction) (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --proxy PROXY                                  🔗 the URL of the proxy
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### Transactions.Relay


```
$ drtpy tx relay --help
usage: drtpy tx relay [-h] ...

Relay a previously saved transaction.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                   show this help message and exit
  --relayer-pem RELAYER_PEM                    🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                               the password.
  --relayer-ledger                             🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type mnemonic
                                               or Ledger devices (default: 0)
  --infile INFILE                              input file (a previously saved transaction)
  --outfile OUTFILE                            where to save the output (the relayer signed transaction) (default:
                                               stdout)
  --send                                       ✓ whether to broadcast the transaction (default: False)
  --simulate                                   whether to simulate the transaction (default: False)
  --proxy PROXY                                🔗 the URL of the proxy

```
## Group **Validator**


```
$ drtpy validator --help
usage: drtpy validator COMMAND [-h] ...

Stake, UnStake, UnBond, Unjail and other actions useful for Validators

COMMANDS:
  {stake,unstake,unjail,unbond,change-reward-address,claim,unstake-nodes,unstake-tokens,unbond-nodes,unbond-tokens,clean-registered-data,restake-unstaked-nodes}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
stake                          Stake value into the Network
unstake                        Unstake value
unjail                         Unjail a Validator Node
unbond                         Unbond tokens for a bls key
change-reward-address          Change the reward address
claim                          Claim rewards
unstake-nodes                  Unstake-nodes will unstake nodes for provided bls keys
unstake-tokens                 This command will un-stake the given amount (if value is greater than the existing topUp value, it will unStake one or several nodes)
unbond-nodes                   It will unBond nodes
unbond-tokens                  It will unBond tokens, if provided value is bigger that topUp value will unBond nodes
clean-registered-data          Deletes duplicated keys from registered data
restake-unstaked-nodes         It will reStake UnStaked nodes

```
### Validator.Stake


```
$ drtpy validator stake --help
usage: drtpy validator stake [-h] ...

Stake value into the Network

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --reward-address REWARD_ADDRESS                the reward address
  --validators-pem VALIDATORS_PEM                a PEM file describing the nodes; can contain multiple nodes
  --top-up                                       Stake value for top up

```
### Validator.Unstake


```
$ drtpy validator unstake --help
usage: drtpy validator unstake [-h] ...

Unstake value

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
### Validator.Unjail


```
$ drtpy validator unjail --help
usage: drtpy validator unjail [-h] ...

Unjail a Validator Node

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
### Validator.Unbond


```
$ drtpy validator unbond --help
usage: drtpy validator unbond [-h] ...

Unbond tokens for a bls key

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
### Validator.ChangeRewardAddress


```
$ drtpy validator change-reward-address --help
usage: drtpy validator change-reward-address [-h] ...

Change the reward address

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --reward-address REWARD_ADDRESS                the new reward address

```
### Validator.Claim


```
$ drtpy validator claim --help
usage: drtpy validator claim [-h] ...

Claim rewards

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### Validator.UnstakeNodes


```
$ drtpy validator unstake-nodes --help
usage: drtpy validator unstake-nodes [-h] ...

Unstake-nodes will unstake nodes for provided bls keys

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
### Validator.UnstakeTokens


```
$ drtpy validator unstake-tokens --help
usage: drtpy validator unstake-tokens [-h] ...

This command will un-stake the given amount (if value is greater than the existing topUp value, it will unStake one or several nodes)

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --unstake-value UNSTAKE_VALUE                  the unstake value

```
### Validator.UnbondNodes


```
$ drtpy validator unbond-nodes --help
usage: drtpy validator unbond-nodes [-h] ...

It will unBond nodes

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
### Validator.UnbondTokens


```
$ drtpy validator unbond-tokens --help
usage: drtpy validator unbond-tokens [-h] ...

It will unBond tokens, if provided value is bigger that topUp value will unBond nodes

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --unbond-value UNBOND_VALUE                    the unbond value

```
### Validator.CleanRegisteredData


```
$ drtpy validator clean-registered-data --help
usage: drtpy validator clean-registered-data [-h] ...

Deletes duplicated keys from registered data

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### Validator.RestakeUnstakedNodes


```
$ drtpy validator restake-unstaked-nodes --help
usage: drtpy validator restake-unstaked-nodes [-h] ...

It will reStake UnStaked nodes

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --nodes-public-keys NODES_PUBLIC_KEYS          the public keys of the nodes as CSV (addrA,addrB)

```
## Group **StakingProvider**


```
$ drtpy staking-provider --help
usage: drtpy staking-provider COMMAND [-h] ...

Staking provider omnitool

COMMANDS:
  {create-new-delegation-contract,get-contract-address,add-nodes,remove-nodes,stake-nodes,unbond-nodes,unstake-nodes,unjail-nodes,delegate,claim-rewards,redelegate-rewards,undelegate,withdraw,change-service-fee,modify-delegation-cap,automatic-activation,redelegate-cap,set-metadata,make-delegation-contract-from-validator}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
create-new-delegation-contract Create a new delegation system smart contract, transferred value must be greater than baseIssuingCost + min deposit value
get-contract-address           Get create contract address by transaction hash
add-nodes                      Add new nodes must be called by the contract owner
remove-nodes                   Remove nodes must be called by the contract owner
stake-nodes                    Stake nodes must be called by the contract owner
unbond-nodes                   Unbond nodes must be called by the contract owner
unstake-nodes                  Unstake nodes must be called by the contract owner
unjail-nodes                   Unjail nodes must be called by the contract owner
delegate                       Delegate funds to a delegation contract
claim-rewards                  Claim the rewards earned for delegating
redelegate-rewards             Redelegate the rewards earned for delegating
undelegate                     Undelegate funds from a delegation contract
withdraw                       Withdraw funds from a delegation contract
change-service-fee             Change service fee must be called by the contract owner
modify-delegation-cap          Modify delegation cap must be called by the contract owner
automatic-activation           Automatic activation must be called by the contract owner
redelegate-cap                 Redelegate cap must be called by the contract owner
set-metadata                   Set metadata must be called by the contract owner
make-delegation-contract-from-validator Create a delegation contract from validator data. Must be called by the node operator

```
### StakingProvider.CreateNewDelegationContract


```
$ drtpy staking-provider create-new-delegation-contract --help
usage: drtpy staking-provider create-new-delegation-contract [-h] ...

Create a new delegation system smart contract, transferred value must be greater than baseIssuingCost + min deposit value

options:
  -h, --help                                     show this help message and exit
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --total-delegation-cap TOTAL_DELEGATION_CAP    the total delegation contract capacity
  --service-fee SERVICE_FEE                      the delegation contract service fee

```
### StakingProvider.GetContractAddress


```
$ drtpy staking-provider get-contract-address --help
usage: drtpy staking-provider get-contract-address [-h] ...

Get create contract address by transaction hash

options:
  -h, --help                       show this help message and exit
  --create-tx-hash CREATE_TX_HASH  the hash
  --proxy PROXY                    🔗 the URL of the proxy

```
### StakingProvider.AddNodes


```
$ drtpy staking-provider add-nodes --help
usage: drtpy staking-provider add-nodes [-h] ...

Add new nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      bech32 address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.RemoveNodes


```
$ drtpy staking-provider remove-nodes --help
usage: drtpy staking-provider remove-nodes [-h] ...

Remove nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --bls-keys BLS_KEYS                            a list with the bls keys of the nodes as CSV (addrA,addrB)
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.StakeNodes


```
$ drtpy staking-provider stake-nodes --help
usage: drtpy staking-provider stake-nodes [-h] ...

Stake nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --bls-keys BLS_KEYS                            a list with the bls keys of the nodes as CSV (addrA,addrB)
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      bech32 address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.UnbondNodes


```
$ drtpy staking-provider unbond-nodes --help
usage: drtpy staking-provider unbond-nodes [-h] ...

Unbond nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --bls-keys BLS_KEYS                            a list with the bls keys of the nodes as CSV (addrA,addrB)
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.UnstakeNodes


```
$ drtpy staking-provider unstake-nodes --help
usage: drtpy staking-provider unstake-nodes [-h] ...

Unstake nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --bls-keys BLS_KEYS                            a list with the bls keys of the nodes as CSV (addrA,addrB)
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.UnjailNodes


```
$ drtpy staking-provider unjail-nodes --help
usage: drtpy staking-provider unjail-nodes [-h] ...

Unjail nodes must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --bls-keys BLS_KEYS                            a list with the bls keys of the nodes as CSV (addrA,addrB)
  --validators-pem VALIDATORS_PEM                a PEM file holding the BLS keys; can contain multiple nodes
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.Delegate


```
$ drtpy staking-provider delegate --help
usage: drtpy staking-provider delegate [-h] ...

Delegate funds to a delegation contract

options:
  -h, --help                                     show this help message and exit
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.ClaimRewards


```
$ drtpy staking-provider claim-rewards --help
usage: drtpy staking-provider claim-rewards [-h] ...

Claim the rewards earned for delegating

options:
  -h, --help                                     show this help message and exit
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.RedelegateRewards


```
$ drtpy staking-provider redelegate-rewards --help
usage: drtpy staking-provider redelegate-rewards [-h] ...

Redelegate the rewards earned for delegating

options:
  -h, --help                                     show this help message and exit
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.Undelegate


```
$ drtpy staking-provider undelegate --help
usage: drtpy staking-provider undelegate [-h] ...

Undelegate funds from a delegation contract

options:
  -h, --help                                     show this help message and exit
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.Withdraw


```
$ drtpy staking-provider withdraw --help
usage: drtpy staking-provider withdraw [-h] ...

Withdraw funds from a delegation contract

options:
  -h, --help                                     show this help message and exit
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.ChangeServiceFee


```
$ drtpy staking-provider change-service-fee --help
usage: drtpy staking-provider change-service-fee [-h] ...

Change service fee must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --service-fee SERVICE_FEE                      new service fee value
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.ModifyDelegationCap


```
$ drtpy staking-provider modify-delegation-cap --help
usage: drtpy staking-provider modify-delegation-cap [-h] ...

Modify delegation cap must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --delegation-cap DELEGATION_CAP                new delegation contract capacity
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.AutomaticActivation


```
$ drtpy staking-provider automatic-activation --help
usage: drtpy staking-provider automatic-activation [-h] ...

Automatic activation must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --set                                          set automatic activation True
  --unset                                        set automatic activation False
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.RedelegateCap


```
$ drtpy staking-provider redelegate-cap --help
usage: drtpy staking-provider redelegate-cap [-h] ...

Redelegate cap must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --set                                          set redelegate cap True
  --unset                                        set redelegate cap False
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.SetMetadata


```
$ drtpy staking-provider set-metadata --help
usage: drtpy staking-provider set-metadata [-h] ...

Set metadata must be called by the contract owner

options:
  -h, --help                                     show this help message and exit
  --name NAME                                    name field in staking provider metadata
  --website WEBSITE                              website field in staking provider metadata
  --identifier IDENTIFIER                        identifier field in staking provider metadata
  --delegation-contract DELEGATION_CONTRACT      address of the delegation contract
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
### StakingProvider.MakeDelegationContractFromValidator


```
$ drtpy staking-provider make-delegation-contract-from-validator --help
usage: drtpy staking-provider make-delegation-contract-from-validator [-h] ...

Create a delegation contract from validator data. Must be called by the node operator

options:
  -h, --help                                     show this help message and exit
  --max-cap MAX_CAP                              total delegation cap in REWA, fully denominated. Use value 0 for
                                                 uncapped
  --fee FEE                                      service fee as hundredths of percents. (e.g. a service fee of 37.45
                                                 percent is expressed by the integer 3745)
  --proxy PROXY                                  🔗 the URL of the proxy
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --outfile OUTFILE                              where to save the output (signed transaction, hash) (default: stdout)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)

```
## Group **Wallet**


```
$ drtpy wallet --help
usage: drtpy wallet COMMAND [-h] ...

Create wallet, derive secret key from mnemonic, bech32 address helpers etc.

COMMANDS:
  {new,convert,bech32,sign-message,verify-message}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
new                            Create a new wallet and print its mnemonic; optionally save as password-protected JSON (recommended) or PEM (not recommended)
convert                        Convert a wallet from one format to another
bech32                         Helper for encoding and decoding bech32 addresses
sign-message                   Sign a message
verify-message                 Verify a previously signed message

```
### Wallet.New


```
$ drtpy wallet new --help
usage: drtpy wallet new [-h] ...

Create a new wallet and print its mnemonic; optionally save as password-protected JSON (recommended) or PEM (not recommended)

options:
  -h, --help                                      show this help message and exit
  --format {raw-mnemonic,keystore-mnemonic,keystore-secret-key,pem}
                                                  the format of the generated wallet file (default: None)
  --outfile OUTFILE                               the output path and base file name for the generated wallet files
                                                  (default: None)
  --address-hrp ADDRESS_HRP                       the human-readable part of the address, when format is keystore-
                                                  secret-key or pem (default: erd)
  --shard SHARD                                   the shard in which the address will be generated; (default: random)

```
### Wallet.Convert


```
$ drtpy wallet convert --help
usage: drtpy wallet convert [-h] ...

Convert a wallet from one format to another

options:
  -h, --help                                      show this help message and exit
  --infile INFILE                                 path to the input file
  --outfile OUTFILE                               path to the output file
  --in-format {raw-mnemonic,keystore-mnemonic,keystore-secret-key,pem}
                                                  the format of the input file
  --out-format {raw-mnemonic,keystore-mnemonic,keystore-secret-key,pem,address-bech32,address-hex,secret-key}
                                                  the format of the output file
  --address-index ADDRESS_INDEX                   the address index, if input format is raw-mnemonic, keystore-mnemonic
                                                  or pem (with multiple entries) and the output format is keystore-
                                                  secret-key or pem
  --address-hrp ADDRESS_HRP                       the human-readable part of the address, when the output format is
                                                  keystore-secret-key or pem (default: erd)

```
### Wallet.Bech32


```
$ drtpy wallet bech32 --help
usage: drtpy wallet bech32 [-h] ...

Helper for encoding and decoding bech32 addresses

positional arguments:
  value       the value to encode or decode

options:
  -h, --help  show this help message and exit
  --encode    whether to encode
  --decode    whether to decode
  --hrp HRP   the human readable part; only used for encoding to bech32 (default: erd)

```
### Wallet.SignMessage


```
$ drtpy wallet sign-message --help
usage: drtpy wallet sign-message [-h] ...

Sign a message

options:
  -h, --help                                 show this help message and exit
  --message MESSAGE                          the message you want to sign
  --sender SENDER                            the alias of the wallet set in the address config
  --pem PEM                                  🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                          🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                        DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter the
                                             password.
  --ledger                                   🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type mnemonic
                                             or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME          🖄 the username of the sender
  --hrp HRP                                  The hrp used to convert the address to its bech32 representation

```
### Wallet.VerifyMessage


```
$ drtpy wallet verify-message --help
usage: drtpy wallet verify-message [-h] ...

Verify a previously signed message

options:
  -h, --help             show this help message and exit
  --address ADDRESS      the bech32 address of the signer
  --message MESSAGE      the previously signed message(readable text, as it was signed)
  --signature SIGNATURE  the signature in hex format

```
## Group **ValidatorWallet**


```
$ drtpy validator-wallet --help
usage: drtpy validator-wallet COMMAND [-h] ...

Create a validator wallet, sign and verify messages and convert a validator wallet to a hex secret key.

COMMANDS:
  {new,sign-message,verify-message-signature,convert}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
new                            Create a new validator wallet and save it as a PEM file.
sign-message                   Sign a message.
verify-message-signature       Verify a previously signed message.
convert                        Convert a validator pem file to a hex secret key.

```
### Wallet.New


```
$ drtpy validator-wallet new --help
usage: drtpy validator-wallet new [-h] ...

Create a new validator wallet and save it as a PEM file.

options:
  -h, --help         show this help message and exit
  --outfile OUTFILE  the output path and file name for the generated wallet

```
### Wallet.Convert


```
$ drtpy validator-wallet convert --help
usage: drtpy validator-wallet convert [-h] ...

Convert a validator pem file to a hex secret key.

options:
  -h, --help       show this help message and exit
  --infile INFILE  the pem file of the wallet
  --index INDEX    the index of the validator in case the file contains multiple validators (default: 0)

```
### Wallet.SignMessage


```
$ drtpy validator-wallet sign-message --help
usage: drtpy validator-wallet sign-message [-h] ...

Sign a message.

options:
  -h, --help         show this help message and exit
  --message MESSAGE  the message you want to sign
  --pem PEM          the path to a validator pem file
  --index INDEX      the index of the validator in case the file contains multiple validators (default: 0)

```
### Wallet.VerifyMessage


```
$ drtpy validator-wallet verify-message-signature --help
usage: drtpy validator-wallet verify-message-signature [-h] ...

Verify a previously signed message.

options:
  -h, --help             show this help message and exit
  --pubkey PUBKEY        the hex string representing the validator's public key
  --message MESSAGE      the previously signed message(readable text, as it was signed)
  --signature SIGNATURE  the signature in hex format

```
## Group **Localnet**


```
$ drtpy localnet --help
usage: drtpy localnet COMMAND [-h] ...

Set up, start and control localnets

COMMANDS:
  {setup,new,prerequisites,build,start,config,clean}

OPTIONS:
  -h, --help            show this help message and exit

```
### Localnet.Setup


```
$ drtpy localnet setup --help
usage: drtpy localnet setup [-h] ...

Set up a localnet (runs 'prerequisites', 'build' and 'config' in one go)

options:
  -h, --help               show this help message and exit
  --configfile CONFIGFILE  An optional configuration file describing the localnet

```
### Localnet.New


```
$ drtpy localnet new --help
usage: drtpy localnet new [-h] ...

Create a new localnet configuration

options:
  -h, --help               show this help message and exit
  --configfile CONFIGFILE  An optional configuration file describing the localnet

```
### Localnet.Prerequisites


```
$ drtpy localnet prerequisites --help
usage: drtpy localnet prerequisites [-h] ...

Download and verify the prerequisites for running a localnet

options:
  -h, --help               show this help message and exit
  --configfile CONFIGFILE  An optional configuration file describing the localnet

```
### Localnet.Build


```
$ drtpy localnet build --help
usage: drtpy localnet build [-h] ...

Build necessary software for running a localnet

options:
  -h, --help                                      show this help message and exit
  --configfile CONFIGFILE                         An optional configuration file describing the localnet
  --software {node,seednode,proxy} [{node,seednode,proxy} ...]
                                                  The software to build (default: ['node', 'seednode', 'proxy'])

```
### Localnet.Config


```
$ drtpy localnet config --help
usage: drtpy localnet config [-h] ...

Configure a localnet (required before starting it the first time or after clean)

options:
  -h, --help               show this help message and exit
  --configfile CONFIGFILE  An optional configuration file describing the localnet

```
### Localnet.Start


```
$ drtpy localnet start --help
usage: drtpy localnet start [-h] ...

Start a localnet

options:
  -h, --help                               show this help message and exit
  --configfile CONFIGFILE                  An optional configuration file describing the localnet
  --stop-after-seconds STOP_AFTER_SECONDS  Stop the localnet after a given number of seconds (default: 31536000)

```
### Localnet.Clean


```
$ drtpy localnet clean --help
usage: drtpy localnet clean [-h] ...

Erase the currently configured localnet (must be already stopped)

options:
  -h, --help               show this help message and exit
  --configfile CONFIGFILE  An optional configuration file describing the localnet

```
## Group **Dependencies**


```
$ drtpy deps --help
usage: drtpy deps COMMAND [-h] ...

Manage dependencies or dharitri-sdk modules

COMMANDS:
  {install,check}

OPTIONS:
  -h, --help       show this help message and exit

----------------
COMMANDS summary
----------------
install                        Install dependencies or dharitri-sdk modules.
check                          Check whether a dependency is installed.

```
### Dependencies.Install


```
$ drtpy deps install --help
usage: drtpy deps install [-h] ...

Install dependencies or dharitri-sdk modules.

positional arguments:
  {all,golang,testwallets}  the dependency to install

options:
  -h, --help                show this help message and exit
  --overwrite               whether to overwrite an existing installation

```
### Dependencies.Check


```
$ drtpy deps check --help
usage: drtpy deps check [-h] ...

Check whether a dependency is installed.

positional arguments:
  {all,golang,testwallets}  the dependency to check

options:
  -h, --help                show this help message and exit

```
## Group **Configuration**


```
$ drtpy config --help
usage: drtpy config COMMAND [-h] ...

Configure Dharitri CLI (default values etc.)

COMMANDS:
  {dump,get,set,delete,new,switch,list,reset}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
dump                           Dumps the active configuration.
get                            Gets a configuration value from the active configuration.
set                            Sets a configuration value for the active configuration.
delete                         Deletes a configuration value from the active configuration.
new                            Creates a new configuration and sets it as the active configuration.
switch                         Switch to a different config.
list                           List available configs
reset                          Deletes the config file. Default config will be used.

```
### Configuration.Dump


```
$ drtpy config dump --help
usage: drtpy config dump [-h] ...

Dumps the active configuration.

options:
  -h, --help  show this help message and exit
  --defaults  dump defaults instead of local config

```
### Configuration.Get


```
$ drtpy config get --help
usage: drtpy config get [-h] ...

Gets a configuration value from the active configuration.

positional arguments:
  name        the name of the configuration entry

options:
  -h, --help  show this help message and exit

```
### Configuration.Set


```
$ drtpy config set --help
usage: drtpy config set [-h] ...

Sets a configuration value for the active configuration.

positional arguments:
  name        the name of the configuration entry
  value       the new value

options:
  -h, --help  show this help message and exit

```
### Configuration.New


```
$ drtpy config new --help
usage: drtpy config new [-h] ...

Creates a new configuration and sets it as the active configuration.

positional arguments:
  name                 the name of the configuration entry

options:
  -h, --help           show this help message and exit
  --template TEMPLATE  template from which to create the new config

```
### Configuration.Switch


```
$ drtpy config switch --help
usage: drtpy config switch [-h] ...

Switch to a different config.

positional arguments:
  name        the name of the configuration entry

options:
  -h, --help  show this help message and exit

```
### Configuration.List


```
$ drtpy config list --help
usage: drtpy config list [-h] ...

List available configs

options:
  -h, --help  show this help message and exit

```
### Configuration.Reset


```
$ drtpy config reset --help
usage: drtpy config reset [-h] ...

Deletes the config file. Default config will be used.

options:
  -h, --help  show this help message and exit

```
## Group **Data**


```
$ drtpy data --help
usage: drtpy data COMMAND [-h] ...

Data manipulation omnitool

COMMANDS:
  {parse,store,load}

OPTIONS:
  -h, --help          show this help message and exit

----------------
COMMANDS summary
----------------
parse                          Parses values from a given file
store                          Stores a key-value pair within a partition
load                           Loads a key-value pair from a storage partition

```
### Data.Dump


```
$ drtpy data parse --help
usage: drtpy data parse [-h] ...

Parses values from a given file

options:
  -h, --help               show this help message and exit
  --file FILE              path of the file to parse
  --expression EXPRESSION  the Python-Dictionary expression to evaluate in order to extract the data

```
### Data.Store


```
$ drtpy data store --help
usage: drtpy data store [-h] ...

Stores a key-value pair within a partition

options:
  -h, --help             show this help message and exit
  --key KEY              the key
  --value VALUE          the value to save
  --partition PARTITION  the storage partition (default: *)
  --use-global           use the global storage (default: False)

```
### Data.Load


```
$ drtpy data load --help
usage: drtpy data load [-h] ...

Loads a key-value pair from a storage partition

options:
  -h, --help             show this help message and exit
  --key KEY              the key
  --partition PARTITION  the storage partition (default: *)
  --use-global           use the global storage (default: False)

```
## Group **Faucet**


```
$ drtpy faucet --help
usage: drtpy faucet COMMAND [-h] ...

Get xREWA on Devnet or Testnet

COMMANDS:
  {request}

OPTIONS:
  -h, --help  show this help message and exit

----------------
COMMANDS summary
----------------
request                        Request xREWA.

```
### Faucet.Request


```
$ drtpy faucet request --help
usage: drtpy faucet request [-h] ...

Request xREWA.

options:
  -h, --help                                 show this help message and exit
  --sender SENDER                            the alias of the wallet set in the address config
  --pem PEM                                  🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                          🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                        DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter the
                                             password.
  --ledger                                   🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type mnemonic
                                             or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME          🖄 the username of the sender
  --hrp HRP                                  The hrp used to convert the address to its bech32 representation
  --chain {D,T}                              the chain identifier
  --api API                                  custom api url for the native auth client
  --wallet-url WALLET_URL                    custom wallet url to call the faucet from

```
## Group **Multisig**


```
$ drtpy multisig --help
usage: drtpy multisig COMMAND [-h] ...

Deploy and interact with the Multisig Smart Contract

COMMANDS:
  {deploy,deposit,discard-action,discard-batch,add-board-member,add-proposer,remove-user,change-quorum,transfer-and-execute,transfer-and-execute-dcdt,async-call,deploy-from-source,upgrade-from-source,sign-action,sign-batch,sign-and-perform,sign-batch-and-perform,unsign-action,unsign-batch,unsign-for-outdated-members,perform-action,perform-batch,get-quorum,get-num-board-members,get-num-groups,get-num-proposers,get-action-group,get-last-action-group-id,get-action-last-index,is-signed-by,is-quorum-reached,get-pending-actions,get-user-role,get-board-members,get-proposers,get-action-data,get-action-signers,get-action-signers-count,get-action-valid-signers-count,parse-propose-action}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
deploy                         Deploy a Multisig Smart Contract.
deposit                        Deposit native tokens (REWA) or DCDT tokens into a Multisig Smart Contract.
discard-action                 Discard a proposed action. Signatures must be removed first via `unsign`.
discard-batch                  Discard all the actions for the specified IDs.
add-board-member               Propose adding a new board member.
add-proposer                   Propose adding a new proposer.
remove-user                    Propose removing a user from the Multisig Smart Contract.
change-quorum                  Propose changing the quorum of the Multisig Smart Contract.
transfer-and-execute           Propose transferring REWA and optionally calling a smart contract.
transfer-and-execute-dcdt      Propose transferring DCDTs and optionally calling a smart contract.
async-call                     Propose a transaction in which the contract will perform an async call.
deploy-from-source             Propose a smart contract deploy from a previously deployed smart contract.
upgrade-from-source            Propose a smart contract upgrade from a previously deployed smart contract.
sign-action                    Sign a proposed action.
sign-batch                     Sign a batch of actions.
sign-and-perform               Sign a proposed action and perform it. Works only if quorum is reached.
sign-batch-and-perform         Sign a batch of actions and perform them. Works only if quorum is reached.
unsign-action                  Unsign a proposed action.
unsign-batch                   Unsign a batch of actions.
unsign-for-outdated-members    Unsign an action for outdated board members.
perform-action                 Perform an action that has reached quorum.
perform-batch                  Perform a batch of actions that has reached quorum.
get-quorum                     Perform a smart contract query to get the quorum.
get-num-board-members          Perform a smart contract query to get the number of board members.
get-num-groups                 Perform a smart contract query to get the number of groups.
get-num-proposers              Perform a smart contract query to get the number of proposers.
get-action-group               Perform a smart contract query to get the actions in a group.
get-last-action-group-id       Perform a smart contract query to get the id of the last action in a group.
get-action-last-index          Perform a smart contract query to get the index of the last action.
is-signed-by                   Perform a smart contract query to check if an action is signed by a user.
is-quorum-reached              Perform a smart contract query to check if an action has reached quorum.
get-pending-actions            Perform a smart contract query to get the pending actions full info.
get-user-role                  Perform a smart contract query to get the role of a user.
get-board-members              Perform a smart contract query to get all the board members.
get-proposers                  Perform a smart contract query to get all the proposers.
get-action-data                Perform a smart contract query to get the data of an action.
get-action-signers             Perform a smart contract query to get the signers of an action.
get-action-signers-count       Perform a smart contract query to get the number of signers of an action.
get-action-valid-signers-count Perform a smart contract query to get the number of valid signers of an action.
parse-propose-action           Parses the propose action transaction to extract proposal ID.

```
### Multisig.Deploy


```
$ drtpy multisig deploy --help
usage: drtpy multisig deploy [-h] ...

Deploy a Multisig Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --bytecode BYTECODE                             the file containing the WASM bytecode
  --quorum QUORUM                                 the number of signatures required to approve a proposal
  --board-members BOARD_MEMBERS [BOARD_MEMBERS ...]
                                                  the bech32 addresses of the board members
  --metadata-not-upgradeable                      ‼ mark the contract as NOT upgradeable (default: upgradeable)
  --metadata-not-readable                         ‼ mark the contract as NOT readable (default: readable)
  --metadata-payable                              ‼ mark the contract as payable (default: not payable)
  --metadata-payable-by-sc                        ‼ mark the contract as payable by SC (default: not payable by SC)
  --abi ABI                                       the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Multisig.Deposit


```
$ drtpy multisig deposit --help
usage: drtpy multisig deposit [-h] ...

Deposit native tokens (REWA) or DCDT tokens into a Multisig Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --token-transfers TOKEN_TRANSFERS [TOKEN_TRANSFERS ...]
                                                  token transfers for transfer & execute, as [token, amount] E.g.
                                                  --token-transfers NFT-123456-0a 1 DCDT-987654 100000000
  --contract CONTRACT                             🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                       the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Multisig.DiscardAction


```
$ drtpy multisig discard-action --help
usage: drtpy multisig discard-action [-h] ...

Discard a proposed action. Signatures must be removed first via `unsign`.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action ACTION                                the id of the action
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.DiscardBatch


```
$ drtpy multisig discard-batch --help
usage: drtpy multisig discard-batch [-h] ...

Discard all the actions for the specified IDs.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action-ids ACTION_IDS [ACTION_IDS ...]       the IDs of the actions to discard
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.AddBoardMember


```
$ drtpy multisig add-board-member --help
usage: drtpy multisig add-board-member [-h] ...

Propose adding a new board member.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --board-member BOARD_MEMBER                    the bech32 address of the proposed board member
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.AddProposer


```
$ drtpy multisig add-proposer --help
usage: drtpy multisig add-proposer [-h] ...

Propose adding a new proposer.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --proposer PROPOSER                            the bech32 address of the proposed proposer
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.RemoveUser


```
$ drtpy multisig remove-user --help
usage: drtpy multisig remove-user [-h] ...

Propose removing a user from the Multisig Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --user USER                                    the bech32 address of the proposed user to be removed
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.ChangeQuorum


```
$ drtpy multisig change-quorum --help
usage: drtpy multisig change-quorum [-h] ...

Propose changing the quorum of the Multisig Smart Contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --quorum QUORUM                                the size of the new quorum (number of signatures required to approve a
                                                 proposal)
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.TransferAndExecute


```
$ drtpy multisig transfer-and-execute --help
usage: drtpy multisig transfer-and-execute [-h] ...

Propose transferring REWA and optionally calling a smart contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --opt-gas-limit OPT_GAS_LIMIT                  optional gas limit for the async call
  --contract-abi CONTRACT_ABI                    the ABI file of the contract to call
  --function FUNCTION                            the function to call
  --arguments ARGUMENTS [ARGUMENTS ...]          arguments for the contract transaction, as [number, bech32-address,
                                                 ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                 0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                a json file containing the arguments. ONLY if abi file is provided.
                                                 E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --receiver RECEIVER                            🖄 the address of the receiver
  --receiver-username RECEIVER_USERNAME          🖄 the username of the receiver
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.TransferAndExecuteDcdt


```
$ drtpy multisig transfer-and-execute-dcdt --help
usage: drtpy multisig transfer-and-execute-dcdt [-h] ...

Propose transferring DCDTs and optionally calling a smart contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --token-transfers TOKEN_TRANSFERS [TOKEN_TRANSFERS ...]
                                                  token transfers for transfer & execute, as [token, amount] E.g.
                                                  --token-transfers NFT-123456-0a 1 DCDT-987654 100000000
  --opt-gas-limit OPT_GAS_LIMIT                   optional gas limit for the async call
  --contract-abi CONTRACT_ABI                     the ABI file of the contract to call
  --function FUNCTION                             the function to call
  --arguments ARGUMENTS [ARGUMENTS ...]           arguments for the contract transaction, as [number, bech32-address,
                                                  ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                  0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                 a json file containing the arguments. ONLY if abi file is provided.
                                                  E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --contract CONTRACT                             🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                       the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --receiver RECEIVER                             🖄 the address of the receiver
  --receiver-username RECEIVER_USERNAME           🖄 the username of the receiver
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Multisig.AsyncCall


```
$ drtpy multisig async-call --help
usage: drtpy multisig async-call [-h] ...

Propose a transaction in which the contract will perform an async call.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --token-transfers TOKEN_TRANSFERS [TOKEN_TRANSFERS ...]
                                                  token transfers for transfer & execute, as [token, amount] E.g.
                                                  --token-transfers NFT-123456-0a 1 DCDT-987654 100000000
  --opt-gas-limit OPT_GAS_LIMIT                   optional gas limit for the async call
  --contract-abi CONTRACT_ABI                     the ABI file of the contract to call
  --function FUNCTION                             the function to call
  --arguments ARGUMENTS [ARGUMENTS ...]           arguments for the contract transaction, as [number, bech32-address,
                                                  ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                  0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                 a json file containing the arguments. ONLY if abi file is provided.
                                                  E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --contract CONTRACT                             🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                       the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --receiver RECEIVER                             🖄 the address of the receiver
  --receiver-username RECEIVER_USERNAME           🖄 the username of the receiver
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Multisig.DeployFromSource


```
$ drtpy multisig deploy-from-source --help
usage: drtpy multisig deploy-from-source [-h] ...

Propose a smart contract deploy from a previously deployed smart contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --contract-to-copy CONTRACT_TO_COPY            the bech32 address of the contract to copy
  --contract-abi CONTRACT_ABI                    the ABI file of the contract to copy
  --metadata-not-upgradeable                     ‼ mark the contract as NOT upgradeable (default: upgradeable)
  --metadata-not-readable                        ‼ mark the contract as NOT readable (default: readable)
  --metadata-payable                             ‼ mark the contract as payable (default: not payable)
  --metadata-payable-by-sc                       ‼ mark the contract as payable by SC (default: not payable by SC)
  --arguments ARGUMENTS [ARGUMENTS ...]          arguments for the contract transaction, as [number, bech32-address,
                                                 ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                 0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                a json file containing the arguments. ONLY if abi file is provided.
                                                 E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.UpgradeFromSource


```
$ drtpy multisig upgrade-from-source --help
usage: drtpy multisig upgrade-from-source [-h] ...

Propose a smart contract upgrade from a previously deployed smart contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --contract-to-upgrade CONTRACT_TO_UPGRADE      the bech32 address of the contract to upgrade
  --contract-to-copy CONTRACT_TO_COPY            the bech32 address of the contract to copy
  --contract-abi CONTRACT_ABI                    the ABI file of the contract to copy
  --metadata-not-upgradeable                     ‼ mark the contract as NOT upgradeable (default: upgradeable)
  --metadata-not-readable                        ‼ mark the contract as NOT readable (default: readable)
  --metadata-payable                             ‼ mark the contract as payable (default: not payable)
  --metadata-payable-by-sc                       ‼ mark the contract as payable by SC (default: not payable by SC)
  --arguments ARGUMENTS [ARGUMENTS ...]          arguments for the contract transaction, as [number, bech32-address,
                                                 ascii string, boolean] or hex-encoded. E.g. --arguments 42 0x64 1000
                                                 0xabba str:TOK-a1c2ef true addr:erd1[..]
  --arguments-file ARGUMENTS_FILE                a json file containing the arguments. ONLY if abi file is provided.
                                                 E.g. [{ 'to': 'erd1...', 'amount': 10000000000 }]
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.SignAction


```
$ drtpy multisig sign-action --help
usage: drtpy multisig sign-action [-h] ...

Sign a proposed action.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action ACTION                                the id of the action
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.SignBatch


```
$ drtpy multisig sign-batch --help
usage: drtpy multisig sign-batch [-h] ...

Sign a batch of actions.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --batch BATCH                                  the id of the batch to sign
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.SignAndPerform


```
$ drtpy multisig sign-and-perform --help
usage: drtpy multisig sign-and-perform [-h] ...

Sign a proposed action and perform it. Works only if quorum is reached.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action ACTION                                the id of the action
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.SignBatchAndPerform


```
$ drtpy multisig sign-batch-and-perform --help
usage: drtpy multisig sign-batch-and-perform [-h] ...

Sign a batch of actions and perform them. Works only if quorum is reached.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --batch BATCH                                  the id of the batch to sign
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.UnsignAction


```
$ drtpy multisig unsign-action --help
usage: drtpy multisig unsign-action [-h] ...

Unsign a proposed action.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action ACTION                                the id of the action
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.UnsignBatch


```
$ drtpy multisig unsign-batch --help
usage: drtpy multisig unsign-batch [-h] ...

Unsign a batch of actions.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --batch BATCH                                  the id of the batch to unsign
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.UnsignForOutdatedMembers


```
$ drtpy multisig unsign-for-outdated-members --help
usage: drtpy multisig unsign-for-outdated-members [-h] ...

Unsign an action for outdated board members.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                      show this help message and exit
  --action ACTION                                 the id of the action
  --outdated-members OUTDATED_MEMBERS [OUTDATED_MEMBERS ...]
                                                  IDs of the outdated board members
  --contract CONTRACT                             🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                       the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                               where to save the output (default: stdout)
  --sender SENDER                                 the alias of the wallet set in the address config
  --pem PEM                                       🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                               🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --ledger                                        🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX       🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME               🖄 the username of the sender
  --hrp HRP                                       The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                   🔗 the URL of the proxy
  --nonce NONCE                                   # the nonce for the transaction. If not provided, is fetched from the
                                                  network.
  --gas-price GAS_PRICE                           ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                           ⛽ the gas limit
  --value VALUE                                   the value to transfer (default: 0)
  --chain CHAIN                                   the chain identifier
  --version VERSION                               the transaction version (default: 2)
  --options OPTIONS                               the transaction options (default: 0)
  --relayer RELAYER                               the bech32 address of the relayer
  --guardian GUARDIAN                             the bech32 address of the guardian
  --send                                          ✓ whether to broadcast the transaction (default: False)
  --simulate                                      whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL     the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE           the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                     🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE             🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE           DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --guardian-ledger                               🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX   🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                       🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE               🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE             DEPRECATED, do not use it anymore. Instead, you'll be prompted to
                                                  enter the password.
  --relayer-ledger                                🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX     🔑 the address index; can be used for PEM files, keyfiles of type
                                                  mnemonic or Ledger devices (default: 0)
  --wait-result                                   signal to wait for the transaction result - only valid if --send is
                                                  set
  --timeout TIMEOUT                               max num of seconds to wait for result - only valid if --wait-result is
                                                  set

```
### Multisig.PerformAction


```
$ drtpy multisig perform-action --help
usage: drtpy multisig perform-action [-h] ...

Perform an action that has reached quorum.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --action ACTION                                the id of the action
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.PerformBatch


```
$ drtpy multisig perform-batch --help
usage: drtpy multisig perform-batch [-h] ...

Perform a batch of actions that has reached quorum.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --batch BATCH                                  the id of the batch to perform
  --contract CONTRACT                            🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI                                      the ABI file of the Multisig Smart Contract
  --outfile OUTFILE                              where to save the output (default: stdout)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Multisig.GetQuorum


```
$ drtpy multisig get-quorum --help
usage: drtpy multisig get-quorum [-h] ...

Perform a smart contract query to get the quorum.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetNumBoardMembers


```
$ drtpy multisig get-num-board-members --help
usage: drtpy multisig get-num-board-members [-h] ...

Perform a smart contract query to get the number of board members.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetNumGroups


```
$ drtpy multisig get-num-groups --help
usage: drtpy multisig get-num-groups [-h] ...

Perform a smart contract query to get the number of groups.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetNumProposers


```
$ drtpy multisig get-num-proposers --help
usage: drtpy multisig get-num-proposers [-h] ...

Perform a smart contract query to get the number of proposers.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetActionGroup


```
$ drtpy multisig get-action-group --help
usage: drtpy multisig get-action-group [-h] ...

Perform a smart contract query to get the actions in a group.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --group GROUP        the group id
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetLastActionGroupId


```
$ drtpy multisig get-last-action-group-id --help
usage: drtpy multisig get-last-action-group-id [-h] ...

Perform a smart contract query to get the id of the last action in a group.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetLastActionLastIndex


```
$ drtpy multisig get-action-last-index --help
usage: drtpy multisig get-action-last-index [-h] ...

Perform a smart contract query to get the index of the last action.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.IsSignedBy


```
$ drtpy multisig is-signed-by --help
usage: drtpy multisig is-signed-by [-h] ...

Perform a smart contract query to check if an action is signed by a user.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --user USER          the bech32 address of the user
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.IsQuorumReached


```
$ drtpy multisig is-quorum-reached --help
usage: drtpy multisig is-quorum-reached [-h] ...

Perform a smart contract query to check if an action has reached quorum.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetPendingActions


```
$ drtpy multisig get-pending-actions --help
usage: drtpy multisig get-pending-actions [-h] ...

Perform a smart contract query to get the pending actions full info.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetUserRole


```
$ drtpy multisig get-user-role --help
usage: drtpy multisig get-user-role [-h] ...

Perform a smart contract query to get the role of a user.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --user USER          the bech32 address of the user
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetBoardMemebers


```
$ drtpy multisig get-board-members --help
usage: drtpy multisig get-board-members [-h] ...

Perform a smart contract query to get all the board members.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetProposers


```
$ drtpy multisig get-proposers --help
usage: drtpy multisig get-proposers [-h] ...

Perform a smart contract query to get all the proposers.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetActionData


```
$ drtpy multisig get-action-data --help
usage: drtpy multisig get-action-data [-h] ...

Perform a smart contract query to get the data of an action.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetActionSigners


```
$ drtpy multisig get-action-signers --help
usage: drtpy multisig get-action-signers [-h] ...

Perform a smart contract query to get the signers of an action.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetActionSignersCount


```
$ drtpy multisig get-action-signers-count --help
usage: drtpy multisig get-action-signers-count [-h] ...

Perform a smart contract query to get the number of signers of an action.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.GetActionValidSignersCount


```
$ drtpy multisig get-action-valid-signers-count --help
usage: drtpy multisig get-action-valid-signers-count [-h] ...

Perform a smart contract query to get the number of valid signers of an action.

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  🖄 the bech32 address of the Multisig Smart Contract
  --abi ABI            the ABI file of the Multisig Smart Contract
  --action ACTION      the id of the action
  --proxy PROXY        🔗 the URL of the proxy

```
### Multisig.ParseProposeAction


```
$ drtpy multisig parse-propose-action --help
usage: drtpy multisig parse-propose-action [-h] ...

Parses the propose action transaction to extract proposal ID.

options:
  -h, --help     show this help message and exit
  --abi ABI      the ABI file of the Multisig Smart Contract
  --hash HASH    the transaction hash of the propose action
  --proxy PROXY  🔗 the URL of the proxy

```
## Group **Governance**


```
$ drtpy governance --help
usage: drtpy governance COMMAND [-h] ...

Propose, vote and interact with the governance contract.

COMMANDS:
  {propose,vote,close-proposal,clear-ended-proposals,claim-accumulated-fees,change-config,get-voting-power,get-config,get-proposal,get-delegated-vote-info}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
propose                        Create a new governance proposal.
vote                           Vote for a governance proposal.
close-proposal                 Close a governance proposal.
clear-ended-proposals          Clear ended proposals.
claim-accumulated-fees         Claim the accumulated fees.
change-config                  Change the config of the contract.
get-voting-power               Get the voting power of an user.
get-config                     Get the config of the governance contract.
get-proposal                   Get info about a proposal.
get-delegated-vote-info        Get info about a delegated vote.

```
### Governance.Propose


```
$ drtpy governance propose --help
usage: drtpy governance propose [-h] ...

Create a new governance proposal.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --commit-hash COMMIT_HASH                      the commit hash of the proposal
  --start-vote-epoch START_VOTE_EPOCH            the epoch in which the voting will start
  --end-vote-epoch END_VOTE_EPOCH                the epoch in which the voting will stop
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.Vote


```
$ drtpy governance vote --help
usage: drtpy governance vote [-h] ...

Vote for a governance proposal.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --proposal-nonce PROPOSAL_NONCE                the nonce of the proposal
  --vote {yes,no,veto,abstain}                   the type of vote
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.CloseProposal


```
$ drtpy governance close-proposal --help
usage: drtpy governance close-proposal [-h] ...

Close a governance proposal.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --proposal-nonce PROPOSAL_NONCE                the nonce of the proposal
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.ClearEndedProposals


```
$ drtpy governance clear-ended-proposals --help
usage: drtpy governance clear-ended-proposals [-h] ...

Clear ended proposals.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --proposers PROPOSERS [PROPOSERS ...]          a list of users who initiated the proposals (e.g. --proposers erd1...,
                                                 erd1...)
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.ClaimAccumulatedFees


```
$ drtpy governance claim-accumulated-fees --help
usage: drtpy governance claim-accumulated-fees [-h] ...

Claim the accumulated fees.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.ChangeConfig


```
$ drtpy governance change-config --help
usage: drtpy governance change-config [-h] ...

Change the config of the contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help                                     show this help message and exit
  --proposal-fee PROPOSAL_FEE                    the cost to create a new proposal
  --lost-proposal-fee LOST_PROPOSAL_FEE          the amount of native tokens the proposer loses if the proposal fails
  --min-quorum MIN_QUORUM                        the min quorum to be reached for the proposal to pass
  --min-veto-threshold MIN_VETO_THRESHOLD        the min veto threshold
  --min-pass-threshold MIN_PASS_THRESHOLD        the min pass threshold
  --sender SENDER                                the alias of the wallet set in the address config
  --pem PEM                                      🔑 the PEM file, if keyfile not provided
  --keyfile KEYFILE                              🔑 a JSON keyfile, if PEM not provided
  --passfile PASSFILE                            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --ledger                                       🔐 bool flag for signing transaction using ledger
  --sender-wallet-index SENDER_WALLET_INDEX      🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --sender-username SENDER_USERNAME              🖄 the username of the sender
  --hrp HRP                                      The hrp used to convert the address to its bech32 representation
  --proxy PROXY                                  🔗 the URL of the proxy
  --nonce NONCE                                  # the nonce for the transaction. If not provided, is fetched from the
                                                 network.
  --gas-price GAS_PRICE                          ⛽ the gas price (default: 1000000000)
  --gas-limit GAS_LIMIT                          ⛽ the gas limit
  --value VALUE                                  the value to transfer (default: 0)
  --chain CHAIN                                  the chain identifier
  --version VERSION                              the transaction version (default: 2)
  --options OPTIONS                              the transaction options (default: 0)
  --relayer RELAYER                              the bech32 address of the relayer
  --guardian GUARDIAN                            the bech32 address of the guardian
  --guardian-service-url GUARDIAN_SERVICE_URL    the url of the guardian service
  --guardian-2fa-code GUARDIAN_2FA_CODE          the 2fa code for the guardian
  --guardian-pem GUARDIAN_PEM                    🔑 the PEM file, if keyfile not provided
  --guardian-keyfile GUARDIAN_KEYFILE            🔑 a JSON keyfile, if PEM not provided
  --guardian-passfile GUARDIAN_PASSFILE          DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --guardian-ledger                              🔐 bool flag for signing transaction using ledger
  --guardian-wallet-index GUARDIAN_WALLET_INDEX  🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --relayer-pem RELAYER_PEM                      🔑 the PEM file, if keyfile not provided
  --relayer-keyfile RELAYER_KEYFILE              🔑 a JSON keyfile, if PEM not provided
  --relayer-passfile RELAYER_PASSFILE            DEPRECATED, do not use it anymore. Instead, you'll be prompted to enter
                                                 the password.
  --relayer-ledger                               🔐 bool flag for signing transaction using ledger
  --relayer-wallet-index RELAYER_WALLET_INDEX    🔑 the address index; can be used for PEM files, keyfiles of type
                                                 mnemonic or Ledger devices (default: 0)
  --outfile OUTFILE                              where to save the output (default: stdout)
  --send                                         ✓ whether to broadcast the transaction (default: False)
  --simulate                                     whether to simulate the transaction (default: False)
  --wait-result                                  signal to wait for the transaction result - only valid if --send is set
  --timeout TIMEOUT                              max num of seconds to wait for result - only valid if --wait-result is
                                                 set

```
### Governance.GetVotingPower


```
$ drtpy governance get-voting-power --help
usage: drtpy governance get-voting-power [-h] ...

Get the voting power of an user.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help     show this help message and exit
  --user USER    the bech32 address of the user
  --proxy PROXY  🔗 the URL of the proxy

```
### Governance.GetConfig


```
$ drtpy governance get-config --help
usage: drtpy governance get-config [-h] ...

Get the config of the governance contract.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help     show this help message and exit
  --proxy PROXY  🔗 the URL of the proxy

```
### Governance.GetDelegatedVoteInfo


```
$ drtpy governance get-delegated-vote-info --help
usage: drtpy governance get-delegated-vote-info [-h] ...

Get info about a delegated vote.

Output example:
===============
{
    "emittedTransaction": {
        "nonce": 42,
        "sender": "alice",
        "receiver": "bob",
        "...": "..."
    },
    "emittedTransactionData": "the transaction data, not encoded",
    "emittedTransactionHash": "the transaction hash"
}

options:
  -h, --help           show this help message and exit
  --contract CONTRACT  the bech32 address of the contract
  --user USER          the bech32 address of the user
  --proxy PROXY        🔗 the URL of the proxy

```
## Group **ConfigEnv**


```
$ drtpy config-env --help
usage: drtpy config-env COMMAND [-h] ...

Configure Dharitri CLI to use specific environment values.

COMMANDS:
  {new,get,set,dump,delete,switch,list,remove,reset}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
new                            Creates a new environment and sets it as the active environment.
get                            Gets an env value from the specified environment.
set                            Sets an env value for the specified environment.
dump                           Dumps the active environment.
delete                         Deletes an env value from the specified environment.
switch                         Switch to a different environment.
list                           List available environments
remove                         Deletes an environment from the env file. Use `drtpy config-env switch` to move to another env.
reset                          Deletes the environment file. Default env will be used.

```
### ConfigEnv.New


```
$ drtpy config-env new --help
usage: drtpy config-env new [-h] ...

Creates a new environment and sets it as the active environment.

positional arguments:
  name                 the name of the new environment

options:
  -h, --help           show this help message and exit
  --template TEMPLATE  an environment from which to create the new environment

```
### ConfigEnv.Set


```
$ drtpy config-env set --help
usage: drtpy config-env set [-h] ...

Sets an env value for the specified environment.

positional arguments:
  name        the name of the configuration entry
  value       the new value

options:
  -h, --help  show this help message and exit
  --env ENV   the name of the environment to operate on

```
### ConfigEnv.Get


```
$ drtpy config-env get --help
usage: drtpy config-env get [-h] ...

Gets an env value from the specified environment.

positional arguments:
  name        the name of the configuration entry

options:
  -h, --help  show this help message and exit
  --env ENV   the name of the environment to operate on

```
### ConfigEnv.Dump


```
$ drtpy config-env dump --help
usage: drtpy config-env dump [-h] ...

Dumps the active environment.

options:
  -h, --help  show this help message and exit
  --default   dumps the default environment instead of the active one.

```
### ConfigEnv.Switch


```
$ drtpy config-env switch --help
usage: drtpy config-env switch [-h] ...

Switch to a different environment.

options:
  -h, --help  show this help message and exit
  --env ENV   the name of the environment to operate on

```
### ConfigEnv.List


```
$ drtpy config-env list --help
usage: drtpy config-env list [-h] ...

List available environments

options:
  -h, --help  show this help message and exit

```
### ConfigEnv.Remove


```
$ drtpy config-env remove --help
usage: drtpy config-env remove [-h] ...

Deletes an environment from the env file. Use `drtpy config-env switch` to move to another env.

options:
  -h, --help  show this help message and exit
  --env ENV   the name of the environment to operate on

```
### ConfigEnv.Reset


```
$ drtpy config-env reset --help
usage: drtpy config-env reset [-h] ...

Deletes the environment file. Default env will be used.

options:
  -h, --help  show this help message and exit

```
## Group **ConfigWallet**


```
$ drtpy config-wallet --help
usage: drtpy config-wallet COMMAND [-h] ...

Configure Dharitri CLI to use a default wallet.

COMMANDS:
  {new,list,dump,get,set,delete,switch,remove,reset}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
new                            Creates a new wallet config and sets it as the active wallet.
list                           List configured wallets
dump                           Dumps the active wallet.
get                            Gets a config value from the specified wallet.
set                            Sets a config value for the specified wallet.
delete                         Deletes a config value from the specified wallet.
switch                         Switch to a different wallet.
remove                         Removes a wallet from the config using the alias. No default wallet will be set. Use `config-wallet switch` to set a new wallet.
reset                          Deletes the config file. No default wallet will be set.

```
### ConfigWallet.New


```
$ drtpy config-wallet new --help
usage: drtpy config-wallet new [-h] ...

Creates a new wallet config and sets it as the active wallet.

positional arguments:
  alias        the alias of the wallet

options:
  -h, --help   show this help message and exit
  --path PATH  the absolute path to the wallet file

```
### ConfigWallet.List


```
$ drtpy config-wallet list --help
usage: drtpy config-wallet list [-h] ...

List configured wallets

options:
  -h, --help  show this help message and exit

```
### ConfigWallet.Dump


```
$ drtpy config-wallet dump --help
usage: drtpy config-wallet dump [-h] ...

Dumps the active wallet.

options:
  -h, --help  show this help message and exit

```
### ConfigWallet.Get


```
$ drtpy config-wallet get --help
usage: drtpy config-wallet get [-h] ...

Gets a config value from the specified wallet.

positional arguments:
  value          the value to get from the specified wallet (e.g. path)

options:
  -h, --help     show this help message and exit
  --alias ALIAS  the alias of the wallet

```
### ConfigWallet.Set


```
$ drtpy config-wallet set --help
usage: drtpy config-wallet set [-h] ...

Sets a config value for the specified wallet.

positional arguments:
  key            the key to set for the specified wallet (e.g. index)
  value          the value to set for the specified key

options:
  -h, --help     show this help message and exit
  --alias ALIAS  the alias of the wallet

```
### ConfigWallet.Switch


```
$ drtpy config-wallet switch --help
usage: drtpy config-wallet switch [-h] ...

Switch to a different wallet.

options:
  -h, --help     show this help message and exit
  --alias ALIAS  the alias of the wallet

```
### ConfigWallet.Delete


```
$ drtpy config-wallet delete --help
usage: drtpy config-wallet delete [-h] ...

Deletes a config value from the specified wallet.

positional arguments:
  value          the value to delete for the specified address

options:
  -h, --help     show this help message and exit
  --alias ALIAS  the alias of the wallet

```
### ConfigWallet.Remove


```
$ drtpy config-wallet remove --help
usage: drtpy config-wallet remove [-h] ...

Removes a wallet from the config using the alias. No default wallet will be set. Use `config-wallet switch` to set a new wallet.

options:
  -h, --help     show this help message and exit
  --alias ALIAS  the alias of the wallet

```
### ConfigWallet.Reset


```
$ drtpy config-wallet reset --help
usage: drtpy config-wallet reset [-h] ...

Deletes the config file. No default wallet will be set.

options:
  -h, --help  show this help message and exit

```
## Group **Get**


```
$ drtpy get --help
usage: drtpy get COMMAND [-h] ...

Get info from the network.

COMMANDS:
  {account,storage,storage-entry,token,transaction,network-config,network-status}

OPTIONS:
  -h, --help            show this help message and exit

----------------
COMMANDS summary
----------------
account                        Get info about an account.
storage                        Get the storage (key-value pairs) of an account.
storage-entry                  Get a specific storage entry (key-value pair) of an account.
token                          Get a token of an account.
transaction                    Get a transaction from the network.
network-config                 Get the network configuration.
network-status                 Get the network status.

```
### Get.Account


```
$ drtpy get account --help
usage: drtpy get account [-h] ...

Get info about an account.

options:
  -h, --help         show this help message and exit
  --alias ALIAS      the alias of the wallet if configured in address config
  --address ADDRESS  the bech32 address
  --proxy PROXY      the proxy url
  --balance          whether to only fetch the balance of the address

```
### Get.Storage


```
$ drtpy get storage --help
usage: drtpy get storage [-h] ...

Get the storage (key-value pairs) of an account.

options:
  -h, --help         show this help message and exit
  --alias ALIAS      the alias of the wallet if configured in address config
  --address ADDRESS  the bech32 address
  --proxy PROXY      the proxy url

```
### Get.StorageEntry


```
$ drtpy get storage-entry --help
usage: drtpy get storage-entry [-h] ...

Get a specific storage entry (key-value pair) of an account.

options:
  -h, --help         show this help message and exit
  --alias ALIAS      the alias of the wallet if configured in address config
  --address ADDRESS  the bech32 address
  --proxy PROXY      the proxy url
  --key KEY          the storage key to read from

```
### Get.Token


```
$ drtpy get token --help
usage: drtpy get token [-h] ...

Get a token of an account.

options:
  -h, --help               show this help message and exit
  --alias ALIAS            the alias of the wallet if configured in address config
  --address ADDRESS        the bech32 address
  --proxy PROXY            the proxy url
  --identifier IDENTIFIER  the token identifier. Works for DCDT and NFT. (e.g. FNG-123456, NFT-987654-0a)

```
### Get.Transaction


```
$ drtpy get transaction --help
usage: drtpy get transaction [-h] ...

Get a transaction from the network.

options:
  -h, --help     show this help message and exit
  --proxy PROXY  the proxy url
  --hash HASH    the transaction hash

```

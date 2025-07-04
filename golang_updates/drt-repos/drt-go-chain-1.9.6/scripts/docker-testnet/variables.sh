# These paths must be absolute

########################################################################
# Docker network configuration

# Don't change the subnet, unless you know what you are doing. Prone to errors.
export DOCKER_NETWORK_SUBNET="172.18.0.0/24"
export DOCKER_NETWORK_NAME="local-testnet"

# By default ports won't be published. If set to 1, all containers will port-forward to host network.
export DOCKER_PUBLISH_PORTS=1

if [[ "$DOCKER_PUBLISH_PORTS" -gt 0 ]]; then
  export DOCKER_PUBLISH_PORT_RANGE=30000
fi

########################################################################


# METASHARD_ID will be used to identify a shard ID as metachain
export METASHARD_ID=4294967295

# Path to drt-go-chain. Determined automatically. Do not change.
export DHARITRIDIR=$(dirname $(dirname $DHARITRITESTNETSCRIPTSDIR))

# Enable the Dharitri Proxy. Note that this is a private repository
# (drt-go-chain-proxy).
export USE_PROXY=1

# Enable the Dharitri Transaction Generator. Note that this is a private
# repository (drt-go-chain-txgen).
export USE_TXGEN=0

# Path where the testnet will be instantiated. This folder is assumed to not
# exist, but it doesn't matter if it already does. It will be created if not,
# anyway.
export TESTNETDIR="$HOME/Dharitri/testnet"

# Path to drt-go-chain-deploy, branch: master. Default: near drt-go-chain.
export CONFIGGENERATORDIR="$(dirname $DHARITRIDIR)/drt-go-chain-deploy/cmd/filegen"

export CONFIGGENERATOR="$CONFIGGENERATORDIR/filegen"    # Leave unchanged.
export CONFIGGENERATOROUTPUTDIR="output"

# Path to the executable node. Leave unchanged unless well justified.
export NODEDIR="$DHARITRIDIR/cmd/node"
export NODE="$NODEDIR/node"     # Leave unchanged

# Path to the executable seednode. Leave unchanged unless well justified.
export SEEDNODEDIR="$DHARITRIDIR/cmd/seednode"
export SEEDNODE="$SEEDNODEDIR/seednode"   # Leave unchanged.

# Niceness value of the Seednode, Observer Nodes and Validator Nodes. Leave
# blank to not adjust niceness.
export NODE_NICENESS=10

# Start a watcher daemon for each validator node, which restarts the node if it
# is suffled out of its shard.
export NODE_WATCHER=0

# Delays after running executables.
export SEEDNODE_DELAY=5
export GENESIS_DELAY=30
export HARDFORK_DELAY=900 #15 minutes enough to take export and gracefully close
export NODE_DELAY=60

export GENESIS_STAKE_TYPE="direct" #'delegated' or 'direct' as in direct stake

#if set to 1, each observer will turn off the antiflooding capability, allowing spam in our network
export OBSERVERS_ANTIFLOOD_DISABLE=0

# Shard structure
export SHARDCOUNT=2
export SHARD_VALIDATORCOUNT=3
export SHARD_OBSERVERCOUNT=1
export SHARD_CONSENSUS_SIZE=3

# Metashard structure
export META_VALIDATORCOUNT=3
export META_OBSERVERCOUNT=1
export META_CONSENSUS_SIZE=$META_VALIDATORCOUNT

# MULTI_KEY_NODES if set to 1, one observer will be generated on each shard that will handle all generated keys
export MULTI_KEY_NODES=0

# EXTRA_KEYS if set to 1, extra keys will be added to the generated keys
export EXTRA_KEYS=1

# ALWAYS_NEW_CHAINID will generate a fresh new chain ID each time start.sh/config.sh is called
export ALWAYS_NEW_CHAINID=1

# ROUNDS_PER_EPOCH represents the number of rounds per epoch. If set to 0, it won't override the node's config
export ROUNDS_PER_EPOCH=0

# HYSTERESIS defines the hysteresis value for number of nodes in shard
export HYSTERESIS=0.0

# ALWAYS_NEW_APP_VERSION will set a new version each time the node will be compiled
export ALWAYS_NEW_APP_VERSION=0

# ALWAYS_UPDATE_CONFIGS will re-generate configs (toml + json) each time ./start.sh
# Set this variable to 0 when testing bootstrap from storage or other edge cases where you do not want a fresh new config
# each time.
export ALWAYS_UPDATE_CONFIGS=1

# IP of the seednode. This should be the first IP allocated in the local testnet network. If you modify the default
# DOCKER_NETWORK_SUBNET, you will need to edit this one accordingly too.
export SEEDNODE_IP="$(echo "$DOCKER_NETWORK_SUBNET" | rev | cut -d. -f2- | rev).2"

# Ports used by the Nodes
export PORT_SEEDNODE="9999"
export PORT_ORIGIN_OBSERVER="21100"
export PORT_ORIGIN_OBSERVER_REST="10000"
export PORT_ORIGIN_VALIDATOR="21500"
export PORT_ORIGIN_VALIDATOR_REST="9500"

########################################################################
# Proxy configuration

# Path to drt-go-chain-proxy, branch: master. Default: near drt-go-chain.
export PROXYDIR="$(dirname $DHARITRIDIR)/drt-go-chain-proxy/cmd/proxy"
export PROXY=$PROXYDIR/proxy    # Leave unchanged.

export PORT_PROXY="7950"
export PROXY_DELAY=10

########################################################################
# TxGen configuration

# Path to drt-go-chain-txgen. Default: near drt-go-chain.
export TXGENDIR="$(dirname $DHARITRIDIR)/drt-go-chain-txgen/cmd/txgen"
export TXGEN=$TXGENDIR/txgen    # Leave unchanged.

export PORT_TXGEN="7951"

export TXGEN_SCENARIOS_LINE='Scenarios = ["basic", "erc20", "dcdt"]'

# Number of accounts to be generated by txgen
export NUMACCOUNTS="250"

# Whether txgen should regenerate its accounts when starting, or not.
# Recommended value is 1, but 0 is useful to run the txgen a second time, to
# continue a testing session on the same accounts.
export TXGEN_REGENERATE_ACCOUNTS=0

# COPY_BACK_CONFIGS when set to 1 will copy back the configs and keys to the ./cmd/node/config directory
# in order to have a node in the IDE that can run a node in debug mode but in the same network with the rest of the nodes
# this option greatly helps the debugging process when running a small system test
export COPY_BACK_CONFIGS=0
# SKIP_VALIDATOR_IDX when setting a value greater than -1 will not launch the validator with the provided index
export SKIP_VALIDATOR_IDX=-1
# SKIP_OBSERVER_IDX when setting a value greater than -1 will not launch the observer with the provided index
export SKIP_OBSERVER_IDX=-1

# USE_HARDFORK will prepare the nodes to run the hardfork process, if needed
export USE_HARDFORK=1

# Load local overrides, .gitignored
LOCAL_OVERRIDES="$DHARITRITESTNETSCRIPTSDIR/local.sh"
if [ -f "$LOCAL_OVERRIDES" ]; then
  source "$DHARITRITESTNETSCRIPTSDIR/local.sh"
fi

# Leave unchanged.
let "total_observer_count = $SHARD_OBSERVERCOUNT * $SHARDCOUNT + $META_OBSERVERCOUNT"
export TOTAL_OBSERVERCOUNT=$total_observer_count

# to enable the full archive feature on the observers, please use the --full-archive flag
export EXTRA_OBSERVERS_FLAGS="-operation-mode db-lookup-extension"

if [[ $MULTI_KEY_NODES -eq 1 ]]; then
  EXTRA_OBSERVERS_FLAGS="--no-key"
fi

# Leave unchanged.
let "total_node_count = $SHARD_VALIDATORCOUNT * $SHARDCOUNT + $META_VALIDATORCOUNT + $TOTAL_OBSERVERCOUNT"
export TOTAL_NODECOUNT=$total_node_count

# VALIDATOR_KEY_PEM_FILE is the pem file name when running single key mode, with all nodes' keys
export VALIDATOR_KEY_PEM_FILE="validatorKey.pem"

# MULTI_KEY_PEM_FILE is the pem file name when running multi key mode, with all managed
export MULTI_KEY_PEM_FILE="allValidatorsKeys.pem"

# EXTRA_KEY_PEM_FILE is the pem file name when running multi key mode, with all extra managed
export EXTRA_KEY_PEM_FILE="extraValidatorsKeys.pem"

# WALLET_PEM="~/Numbat/MySandbox/testnet/wallets/users/alice.pem"
WALLET_PEM="~/Documents/shared_folder/numbat_testnet_wallet.pem"
DEPLOY_TRANSACTION=$(drtpy data load --key=deployTransaction-devnet)
DEPLOY_GAS="1000000000"
PROXY="https://testnet-gateway.dharitri.org"
CHAIN_ID="T"
# PROXY="http://localhost:7950"
# CHAIN_ID="local-testnet"

DCDT_ISSUE_ADDRESS="drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"
ROUTE_ADDRESS="drt1qqqqqqqqqqqqqpgq9s7ft5fj72zyt8qn8yd24xcl8rxu9v4m0n4s0rqhfe"
WREWA_WRAP_ADDRESS="drt1qqqqqqqqqqqqqpgq4axqc749vuqr27snr8d8qgvlmz44chsr0n4sxf2afk"
PAIR_ADDRESS="drt1qqqqqqqqqqqqqpgqr23zlc896w6qc2hw3evmmdmppw6jaucv0n4s36jp5d"
DEFAULT_GAS_LIMIT=50000000

##### ENDPOINTS #####

# params:
#   $1 = Token Name
#   $2 = Token Ticker
issueToken() {
    token_name="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    token_ticker="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    initial_supply=0xFFFFFFFFFFFFFFFFFFFF
    token_decimals=0x12

    drtpy --verbose contract call ${DCDT_ISSUE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --value=5000000000000000000 \
        --function="issue" \
        --arguments ${token_name} ${token_ticker} ${initial_supply} ${token_decimals} \
        --send || return
}

# params:
#   $1 = Token Identifier
changeTokenProperties() {
    true="0x$(echo -n 'true' | xxd -p -u | tr -d '\n')"
    canMint="0x$(echo -n 'canMint' | xxd -p -u | tr -d '\n')"
    token_identifier="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call ${DCDT_ISSUE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="controlChanges" \
        --arguments ${token_identifier} ${canMint} ${true} \
        --send || return
}

# params:
#   $1 = Token Identifier
#   $2 = Amount in hex
mintToken() {
    token_identifier="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call ${DCDT_ISSUE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="mint" \
        --arguments ${token_identifier} $2 \
        --send || return
}

# params:
#   $1 = Token Identifier
#   $2 = Address to assign the role
#   $3 = Special Role
setSpecialRoleToken() {
    token_identifier="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    address="0x$(drtpy wallet bech32 --decode $2)"
    special_role="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call ${DCDT_ISSUE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="setSpecialRole" \
        --arguments $token_identifier $address $special_role \
        --send || return
}


#### ROUTER ####

deployRouterContract() {
    drtpy --verbose contract deploy --recall-nonce \
          --pem=${WALLET_PEM} \
          --gas-price=1499999999 \
          --gas-limit=1499999999 \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --bytecode="../router/output/router.wasm" \
          --outfile="deploy-route-internal.interaction.json" --send || return
    
    ADDRESS=$(drtpy data parse --file="deploy-route-internal.interaction.json" --expression="data['contractAddress']")

    drtpy data store --key=router-address --value=${ADDRESS}

    echo ""
    echo "Route Smart contract address: ${ADDRESS}"
}

upgradeRouterContract() {
    drtpy --verbose contract upgrade ${ROUTE_ADDRESS} --recall-nonce \
          --pem=${WALLET_PEM} \
          --gas-limit=${DEPLOY_GAS} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --bytecode="../router/output/router.wasm" \
          --outfile="upgrade-route-internal.interaction.json" --send || return

    echo ""
    echo "Route Smart contract upgraded"
}

uploadPairContractCode() {
    echo "STARTING TO PUSH NEW PAIR CONTRACT"
    PAIR_CODE_HEX="$(xxd -p ../pair/output/pair.wasm | tr -d '\n')"
    PAIR_CODE_HEX="${PAIR_CODE_HEX::-4}"
    PAIR_CODE_HEX1="0x$(split -n1/3 <<<$PAIR_CODE_HEX)"
    PAIR_CODE_HEX2="0x$(split -n2/3 <<<$PAIR_CODE_HEX)"
    PAIR_CODE_HEX3="0x$(split -n3/3 <<<$PAIR_CODE_HEX)7575"

    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-price=1400000000 \
          --gas-limit=1400000000 \
          --function=startPairCodeConstruction \
          --send
    sleep 6

    echo "SENDING BATCH 1"
    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-price=1400000000 \
        --gas-limit=1400000000 \
        --function=appendPairCode \
        --arguments $PAIR_CODE_HEX1 \
        --send
    sleep 10

    echo "SENDING BATCH 2"
    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-price=1400000000 \
        --gas-limit=1400000000 \
        --function=appendPairCode \
        --arguments $PAIR_CODE_HEX2 \
        --send
    sleep 10

    echo "SENDING BATCH 3"
    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-price=1400000000 \
        --gas-limit=1400000000 \
        --function=appendPairCode \
        --arguments $PAIR_CODE_HEX3 \
        --send
    sleep 10

    echo "ENDING TO CREATE NEW PAIR"
    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-price=1400000000 \
          --gas-limit=1400000000 \
          --function=endPairCodeConstruction \
          --send
    sleep 6
}

# params:
#   $1 = First Token Identifier,
#   $2 = Second Token Identifier,
createPair() {
    first_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-price=1400000000 \
          --gas-limit=1400000000 \
          --function=createPair \
          --arguments $first_token $second_token 300 100\
          --send
    sleep 6

    echo "NEW PAIR CONTRACT ADDRESS:"
    drtpy --verbose contract query ${ROUTE_ADDRESS} \
    --proxy=${PROXY} \
    --function=getPair \
    --arguments $first_token $second_token || return
}

# params:
#   $1 = Pair Address,
#   $2 = LP Token Name,
#   $3 = LP Token Ticker
issueLpToken() {
    pair_address="0x$(drtpy wallet bech32 --decode $1)"
    lp_token_name="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    lp_token_ticker="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-limit=200000000 \
          --value=5000000000000000000 \
          --function=issueLpToken \
	      --arguments $pair_address $lp_token_name $lp_token_ticker \
	      --send || return
}

# params:
#   $1 = Pair Address,
setLpTokenLocalRoles() {
    pair_address="0x$(drtpy wallet bech32 --decode $1)"

    drtpy --verbose contract call ${ROUTE_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=200000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="setLocalRoles" \
        --arguments $pair_address \
        --send || return
}

#params:
#   $1 = pair contract to send fees,
#   $2 = farm contract to receive fees,
#   $3 = farm contract expected token
setFeeOn() {
    pair_address="0x$(drtpy wallet bech32 --decode $1)"
    farm_contract="0x$(drtpy wallet bech32 --decode $2)"
    farm_token="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $ROUTE_ADDRESS --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=200000000 \
        --function=setFeeOn \
        --arguments $pair_address $farm_contract $farm_token \
        --send || return
}

#params:
#   $1 = pair contract to send fees,
#   $2 = farm contract to receive fees,
#   $3 = farm contract expected token
setFeeOff() {
    pair_address="0x$(drtpy wallet bech32 --decode $1)"
    farm_contract="0x$(drtpy wallet bech32 --decode $2)"
    farm_token="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $ROUTE_ADDRESS --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=200000000 \
        --function=setFeeOff \
        --arguments $pair_address $farm_contract $farm_token \
        --send || return
}

# params
#   $1 = Token Identifier
#   $2 = Address
setLocalRolesOwner() {
    token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    address="0x$(drtpy wallet bech32 --decode $2)"

    drtpy --verbose contract call $ROUTE_ADDRESS --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=200000000 \
        --function=setLocalRolesOwner \
        --arguments $token $address 0x02 \
        --send || return
}

#### PAIRS ####

# params
#   $1 = First Token ID
#   $2 = Second Token ID
deployPairContract() {
    first_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    router_address="0x$(drtpy wallet bech32 --decode $ROUTE_ADDRESS)"
    user_address="$(drtpy wallet pem-address $WALLET_PEM)"
    user_address_decode="0x$(drtpy wallet bech32 --decode $user_address)"

    drtpy --verbose contract deploy --recall-nonce \
          --pem=${WALLET_PEM} \
          --gas-limit=250000000 \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --bytecode="../pair/output/pair.wasm" \
          --arguments $first_token $second_token $router_address $user_address_decode 0x000000000000012C 0x0000000000000032 \
          --outfile="deploy-pair-internal.interaction.json" --send || return
    
    ADDRESS=$(drtpy data parse --file="deploy-pair-internal.interaction.json" --expression="data['contractAddress']")

    echo ""
    echo "Pair Smart contract address: ${ADDRESS}"
}

# params
#   $1 = First Token ID
#   $2 = Second Token ID
#   $3 = Pair Address
upgradePairContract() {
    first_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    router_address="0x$(drtpy wallet bech32 --decode $ROUTE_ADDRESS)"
    user_address="$(drtpy wallet pem-address $WALLET_PEM)"
    user_address_decode="0x$(drtpy wallet bech32 --decode $user_address)"

    drtpy --verbose contract upgrade $3 --recall-nonce \
          --pem=${WALLET_PEM} \
          --gas-limit=250000000 \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --bytecode="../pair/output/pair.wasm" \
          --arguments $first_token $second_token $router_address $user_address_decode 0x000000000000012C 0x0000000000000032 \
          --outfile="upgrade-pair-internal.interaction.json" --send || return
    
    echo ""
    echo "Pair Smart contract upgraded"
}

# params
#   $1 = Pair address
#   $2 = Gas limit in hex
setTransferExecGasLimit() {
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=20000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function=set_transfer_exec_gas_limit \
        --arguments $2 \
        --send || return
}

# params
#   $1 = Pair Address
#   $2 = Token Identifier
setLpTokenIdentifier() {
    token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=20000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function=setLpTokenIdentifier \
        --arguments $token \
        --send || return
}

# params:
#   $1 = Pair Address,
#   $2 = First Token Identifier,
#   $3 = First Token Amount in hex,
#   $4 = Second Token Identifier,
#   $5 = Second Token Amount in hex,
#   $6 = Minimum First Token Amount in hex,
#   $7 = Minimum Second Token Amount in hex
addLiquidity() {
    method_name="0x$(echo -n 'acceptDcdtPayment' | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=30000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="multiDCDTNFTTransfer" 2 $2 0 $3 $4 0 $5 ${method_name} $6 $7 \
        --arguments $2 $3 $4 $5 \
        --send || return
}

# params:
#   $1 = Pair Address,
#   $2 = LP Token Identifier,
#   $3 = LP Token Amount in hex,
#   $4 = First Token Amount min in hex,
#   $5 = Second Token Amount min in hex
removeLiquidity() {
    method_name="0x$(echo -n 'removeLiquidity' | xxd -p -u | tr -d '\n')"
    lp_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    
    drtpy --verbose contract call $1 --recall-nonce \
      --pem=${WALLET_PEM} \
      --gas-limit=25000000 \
      --proxy=${PROXY} --chain=${CHAIN_ID} \
      --function="DCDTTransfer" \
      --arguments $lp_token $3 $method_name $4 $5 \
      --send || return
}

# params:
#   $1 = Pair Address,
#   $2 = Token In Identifier,
#   $3 = Amount In in hex,
#   $4 = Token Out Identifier,
#   $5 = Amount Out min in hex
swapFixedInput() {
    method_name="0x$(echo -n 'swapTokensFixedInput' | xxd -p -u | tr -d '\n')"
    token_in="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    token_out="0x$(echo -n $4 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
      --pem=${WALLET_PEM} \
      --gas-limit=100000000 \
      --proxy=${PROXY} --chain=${CHAIN_ID} \
      --function="DCDTTransfer" \
      --arguments $token_in $3 $method_name $token_out $5 \
      --send || return
}

# params:
#   $1 = Pair Address,
#   $2 = Token In max Identifier,
#   $3 = Amount In in hex,
#   $4 = Token Out Identifier,
#   $5 = Amount Out in hex
swapFixedOutput() {
    method_name="0x$(echo -n 'swapTokensFixedOutput' | xxd -p -u | tr -d '\n')"
    token_in="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    token_out="0x$(echo -n $4 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
      --pem=${WALLET_PEM} \
      --gas-limit=100000000 \
      --proxy=${PROXY} --chain=${CHAIN_ID} \
      --function="DCDTTransfer" \
      --arguments $token_in $3 $method_name $token_out $5 \
      --send || return
}

# params
#   $1 = destination pair contract,
#   $2 = pair contract to be whitelisted.
whitelist() {
    pair_address="0x$(drtpy wallet bech32 --decode $2)"
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEFAULT_GAS_LIMIT} \
        --function=whitelist \
        --arguments $pair_address \
        --send || return
}

#params
#   $1 = destination pair contract,
#   $2 = pair contract which will be added as trusted swap pair,
#   $3 = Trusted Pair First Token Identifier,
#   $4 = Trusted Pair Second Token Identifier,
addTrustedSwapPair() {
    pair_address="0x$(drtpy wallet bech32 --decode $2)"
    first_token="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $4 | xxd -p -u | tr -d '\n')"
    
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEFAULT_GAS_LIMIT} \
        --function=addTrustedSwapPair \
        --arguments $pair_address $first_token $second_token \
        --send || return
}

#### FARM ####

# params:
#   $1 = Farmed Pool Token Identifier (Rewards)
#   $2 = Farming Token Identifier (Farming Token)
#   $3 = Locked Asset Factory Address
deployFarmContract() {
    router_address="0x$(drtpy wallet bech32 --decode $ROUTE_ADDRESS)"
    farmed_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    farming_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    locked_asset_factory_address="0x$(drtpy wallet bech32 --decode $3)"

    drtpy --verbose contract deploy --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-price=1499999999 \
        --gas-limit=1499999999 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --metadata-payable \
        --bytecode="../farm/output/farm.wasm" \
        --arguments $router_address $farmed_token $farming_token $locked_asset_factory_address 0xE8D4A51000 \
        --outfile="deploy-farm-internal.interaction.json" --send || return

    ADDRESS=$(drtpy data parse --file="deploy-farm-internal.interaction.json" --expression="data['contractAddress']")

    drtpy data store --key=address-devnet --value=${ADDRESS}

    echo ""
    echo "Farm Smart Contract address: ${ADDRESS}"
}

# params:
#   $1 = Farm Address to upgrade
#   $2 = Farm Pool Token Identifier (Rewards)
#   $3 = Accepted Farming Token Identifier (Farming Token)
#   $4 = Locked Asset Factory Address
upgradeFarmContract() {
    router_address="0x$(drtpy wallet bech32 --decode $ROUTE_ADDRESS)"
    farmed_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    farming_token="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"
    locked_asset_factory_address="0x$(drtpy wallet bech32 --decode $4)"

    drtpy --verbose contract upgrade $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-price=1499999999 \
        --gas-limit=1499999999 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --metadata-payable \
        --bytecode="../farm/output/farm.wasm" \
        --arguments $router_address $farmed_token $farming_token $locked_asset_factory_address 0xE8D4A51000 \
        --outfile="upgrade-farm-internal.interaction.json" --send || return

    echo ""
    echo "Farm Smart Contract upgraded"
}

# params:
#   $1 = farm contract,
#   $2 = farm token name,
#   $3 = farm token ticker,
#   $3 = num decimals,
registerFarmToken() {
    farm_token_name="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    farm_token_ticker="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --value=5000000000000000000 \
        --function=registerFarmToken \
        --arguments $farm_token_name $farm_token_ticker $4 \
        --send || return
}

# params:
#   $1 = farm contract
setLocalRolesFarmToken() {
    drtpy --verbose contract call $1 --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-limit=${DEPLOY_GAS} \
          --function=setLocalRolesFarmToken \
          --send || return
}

#params:
#   $1 = farm contract,
#   $2 = lp token id,
#   $3 = lp token amount in hex
enterFarm() {
    method_name="0x$(echo -n 'enterFarm' | xxd -p -u | tr -d '\n')"
    lp_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=100000000 \
        --function=DCDTTransfer \
        --arguments $lp_token $3 $method_name \
        --send || return
}

#params:
#   $1 = farm contract,
#   $2 = lp token id,
#   $3 = lp token amount in hex
enterFarmAndLockRewards() {
    method_name="0x$(echo -n 'enterFarmAndLockRewards' | xxd -p -u | tr -d '\n')"
    lp_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"

    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=100000000 \
        --function=DCDTTransfer \
        --arguments $lp_token $3 $method_name \
        --send || return
}

#params:
#   $1 = farm token id,
#   $2 = farm token nonce in hex,
#   $3 = farm token amount in hex,
#   $4 = address of staking contract
exitFarm() {
    method_name="0x$(echo -n 'exitFarm' | xxd -p -u | tr -d '\n')"
    user_address="$(drtpy wallet pem-address $WALLET_PEM)"
    stake_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    farm_contract="0x$(drtpy wallet bech32 --decode $4)"

    drtpy --verbose contract call $user_address --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEFAULT_GAS_LIMIT} \
        --function=DCDTNFTTransfer \
        --arguments $stake_token $2 $3 $farm_contract $method_name \
        --send || return
}

#params:
#   $1 = farm token id,
#   $2 = farm token nonce in hex,
#   $3 = farm token amount in hex,
#   $4 = address of staking contract
claimRewards() {
    method_name="0x$(echo -n 'claimRewards' | xxd -p -u | tr -d '\n')"
    user_address="$(drtpy wallet pem-address $WALLET_PEM)"
    stake_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    farm_contract="0x$(drtpy wallet bech32 --decode $4)"

    drtpy --verbose contract call $user_address --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=500000000 \
        --function=DCDTNFTTransfer \
        --arguments $stake_token $2 $3 $farm_contract $method_name \
        --send || return
}


# params
#   $1 = Farm Address
#   $2 = PerBlockRewards in hex
setPerBlockRewardAmount() {
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEFAULT_GAS_LIMIT} \
        --function=setPerBlockRewardAmount \
        --arguments $2 \
        --send || return
}

# params
#   $1 = Farm Address
startProduceRewards() {
    drtpy --verbose contract call $1 --recall-nonce \
          --pem=${WALLET_PEM} \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --gas-limit=${DEPLOY_GAS} \
          --function=start_produce_rewards \
          --send || return
}

##### VIEW FUNCTIONS #####

# params:
#   $1 = First Token Identifier,
#   $2 = Second Token Identifier
getPairAddress() {
    first_token="0x$(echo -n $1 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract query ${ROUTE_ADDRESS} \
        --proxy=${PROXY} \
        --function=getPair \
        --arguments $first_token $second_token || return 
}

# params:
#   $1 = Pair Address
getLpTokenIdentifier() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getLpTokenIdentifier || return
}

# params:
#   $1 = farm Contract Address
getFarmTokenIdentifier() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getFarmTokenId || return
}

# params:
#   $1 = Pair Address,
#   $2 = First Token Identifier,
#   $3 = Second Token Identifier
getReserves() {
    first_token="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    second_token="0x$(echo -n $3 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getReserve \
        --arguments $first_token
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getReserve \
        --arguments $second_token || return
}

# params:
#   $1 = Pair Address,
#   $2 = Token In Identifier,
#   $3 = Token In Amount in hex
getEquivalent() {
    token_in="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getEquivalent \
        --arguments $token_in $3 || return
}

# params:
#   $1 = Pair Address,
#   $2 = Token In Identifier,
#   $3 = Token In Amount in hex
getAmountOut() {
    token_in="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getAmountOut \
        --arguments $token_in $3 || return
}

# params:
#   $1 = Pair Address,
#   $2 = Token Out Identifier,
#   $3 = Token Out Amount in hex
getAmountIn() {
    token_in="0x$(echo -n $2 | xxd -p -u | tr -d '\n')"
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getAmountIn \
        --arguments $token_in $3 || return
}


# params:
#   $1 = Pair Address,
#   $2 = Liquidity amount in hex
getTokensForGivenPosition() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getTokensForGivenPosition \
        --arguments $2 || return
}

# params
#   $1 = Farm Address
getPerBlockRewardAmount() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getPerBlockRewardAmount || return
}

# params
#   $1 = Farm Address
getFarmTokenSupply() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getFarmTokenSupply || return
}


# params:
#   $1 = farm contract,
#   $2 = farm token amount in hex
#   $3 = farm attributes in hex
calculateRewardsForGivenPosition() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=calculateRewardsForGivenPosition \
        --arguments $2 $3 || return
}

getState() {
    drtpy --verbose contract query $1 \
        --proxy=${PROXY} \
        --function=getState || return
}

##### UTILS #####

deployWREWAContract() {
    drtpy --verbose contract deploy --recall-nonce \
          --pem=${WALLET_PEM} \
          --gas-limit=100000000 \
          --proxy=${PROXY} --chain=${CHAIN_ID} \
          --metadata-payable \
          --bytecode="/home/numbat/Numbat/sc-bridge-numbat/rewa-dcdt-swap/output/rewa-dcdt-swap.wasm" \
          --outfile="deploy-wrewa-internal.interaction.json" --wait-result --send || return
    
    ADDRESS=$(drtpy data parse --file="deploy-wrewa-internal.interaction.json" --expression="data['contractAddress']")

    drtpy data store --key=router-address --value=${ADDRESS}

    echo ""
    echo "WREWA Smart contract address: ${ADDRESS}"
}

issueWrappedRewa() {
    drtpy --verbose contract call ${WREWA_WRAP_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --value=5000000000000000000 \
        --function=issueWrappedRewa \
        --arguments 0x5772617070656452455741 0x5752455741 10000000 \
        --send || return   
}

setWREWALocalRole() {
    drtpy --verbose contract call ${WREWA_WRAP_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --function=setLocalMintRole \
        --send || return   
}

wrapRewa() {
    drtpy --verbose contract call ${WREWA_WRAP_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --value=$1 \
        --function=wrapRewa \
        --send || return   

}

unwrapRewa() {
    drtpy --verbose contract call ${WREWA_WRAP_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --value=$1 \
        --function=wrapRewa \
        --send || return   

}

mintWrappedRewa() {
    drtpy --verbose contract call ${WREWA_WRAP_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --function=mintWrappedRewa \
        --arguments $1 \
        --send || return
}

# params:
#   $1 = Contract Address
pauseContract() {
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --function=pause \
        --send || return
}

# params:
#   $1 = Contract Address
resumeContract() {
    drtpy --verbose contract call $1 --recall-nonce \
        --pem=${WALLET_PEM} \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --gas-limit=${DEPLOY_GAS} \
        --function=resume \
        --send || return
}

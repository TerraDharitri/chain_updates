WALLET_PEM="/Users/sorinpetreasca/DevKitt/walletKey.pem"
WALLET_PEM2="/Users/sorinpetreasca/DevKitt/walletKey2.pem"
WALLET_PEM3="/Users/sorinpetreasca/DevKitt/walletKey3.pem"
PROXY="https://testnet-gateway.numbat.com"
CHAIN_ID="T"

LIQUID_STAKING_WASM_PATH="/Users/sorinpetreasca/Numbat/sc-liquid-staking-rs/liquid-staking/output/liquid-staking.wasm"

OWNER_ADDRESS="drt14nw9pukqyqu75gj0shm8upsegjft8l0awjefp877phfx74775dsqge8dz0"
CONTRACT_ADDRESS="drt1qqqqqqqqqqqqqpgq4dzfldya86ht366xrd4w78809rlxcfzn5dsqs072pq"

deploySC() {
    drtpy --verbose contract deploy --recall-nonce \
        --bytecode=${LIQUID_STAKING_WASM_PATH} \
        --pem=${WALLET_PEM} \
        --gas-limit=100000000 \
        --metadata-payable-by-sc \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --send || return
}

upgradeSC() {
    drtpy --verbose contract upgrade ${CONTRACT_ADDRESS} --recall-nonce \
        --bytecode=${LIQUID_STAKING_WASM_PATH} \
        --pem=${WALLET_PEM} \
        --gas-limit=100000000 \
        --metadata-payable-by-sc \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --send || return
}

TOKEN_NAME=0x4c53544f4b454e #LSTOKEN
TOKEN_TICKER=0x4c5354 #LST
TOKEN_DECIMALS=18
registerLsToken() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=100000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --value=50000000000000000 \
        --function="registerLsToken" \
        --arguments ${TOKEN_NAME} ${TOKEN_TICKER} ${TOKEN_DECIMALS} \
        --send || return
}

UNSTAKE_TOKEN_NAME=0x4c53554e5354414b45 #LSUNSTAKE
UNSTAKE_TOKEN_TICKER=0x4c5355 #LSU
registerUnstakeToken() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=100000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --value=50000000000000000 \
        --function="registerUnstakeToken" \
        --arguments ${UNSTAKE_TOKEN_NAME} ${UNSTAKE_TOKEN_TICKER} ${TOKEN_DECIMALS} \
        --send || return
}

setStateActive() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=6000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="setStateActive" \
        --send || return
}

setStateInactive() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=6000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="setStateInactive" \
        --send || return
}

###PARAMS 
### Contracts - drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqq80llllsfs8hsn drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqxllllls3s3f78
DELEGATION_ADDRESS="drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqq80llllsfs8hsn"
TOTAL_STAKED=1000000000000000000000
DELEGATION_CAP=3000000000000000000000
NR_NODES=3
APY=1974
whitelistDelegationContract() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=10000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="whitelistDelegationContract" \
        --arguments ${DELEGATION_ADDRESS} ${OWNER_ADDRESS} ${TOTAL_STAKED} ${DELEGATION_CAP} ${NR_NODES} ${APY}\
        --send || return
}

NEW_TOTAL_STAKED=1500000000000000000000
NEW_DELEGATION_CAP=5000000000000000000000
NEW_APY=18830
changeDelegationContractParams() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="changeDelegationContractParams" \
        --arguments ${DELEGATION_ADDRESS} ${NEW_TOTAL_STAKED} ${NEW_DELEGATION_CAP} ${NR_NODES} ${NEW_APY}\
        --send || return
}

###PARAMS
#1 - Amount
addLiquidity() {
        drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --value=$1 \
        --function="addLiquidity" \
        --send || return
}

###PARAMS
#1 - Amount
LS_TOKEN=str:LST-05ef8b
REMOVE_LIQUIDITY_METHOD=str:removeLiquidity
removeLiquidity() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="DCDTTransfer" \
        --arguments ${LS_TOKEN} $1 ${REMOVE_LIQUIDITY_METHOD} \
        --send || return
}

###PARAMS
## TESTING: 1437 unstake epoch
#1 - Nonce
unbondTokens() {
    user_address="$(drtpy wallet pem-address $WALLET_PEM2)"
    method_name=str:unbondTokens
    unbond_token=str:LSU-1cd7b4
    unbond_token_nonce=$1
    UNBOND_TOKEN_AMOUNT=1
    drtpy --verbose contract call $user_address --recall-nonce \
        --pem=${WALLET_PEM2} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="DCDTNFTTransfer" \
        --arguments $unbond_token $1 ${UNBOND_TOKEN_AMOUNT} ${CONTRACT_ADDRESS} $method_name \
        --send || return
}

claimRewards() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="claimRewards" \
        --send || return
}

recomputeTokenReserve() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=6000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="recomputeTokenReserve" \
        --send || return
}

delegateRewards() {
    drtpy --verbose contract call ${CONTRACT_ADDRESS} --recall-nonce \
        --pem=${WALLET_PEM} \
        --gas-limit=60000000 \
        --proxy=${PROXY} --chain=${CHAIN_ID} \
        --function="delegateRewards" \
        --send || return
}

# VIEWS

getLsTokenId() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getLsTokenId" \
}

###PARAMS
#1 - Amount
getLsValueForPosition() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --arguments $1 \
        --function="getLsValueForPosition" \    
}

getUnstakeTokenId() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getUnstakeTokenId" \
}

getUnstakeTokenSupply() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getUnstakeTokenSupply" \
}

getVirtualRewaReserve() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getVirtualRewaReserve" \
}

getRewardsReserve() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getRewardsReserve" \
}

getTotalWithdrawnRewa() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getTotalWithdrawnRewa" \
}

getDelegationAddressesList() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getDelegationAddressesList" \
}

###PARAMS
#1 - Address
getDelegationContractData() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --arguments $1 \
        --function="getDelegationContractData" \
}

getDelegationStatus() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getDelegationStatus" \
}

getDelegationClaimStatus() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --function="getDelegationClaimStatus" \
}

getDelegationContractStakedAmount() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --arguments ${DELEGATION_ADDRESS} \
        --function="getDelegationContractStakedAmount" \
}

getDelegationContractUnstakedAmount() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --arguments ${DELEGATION_ADDRESS} \
        --function="getDelegationContractUnstakedAmount" \
}

getDelegationContractUnbondedAmount() {
    drtpy --verbose contract query ${CONTRACT_ADDRESS} \
        --proxy=${PROXY} \
        --arguments ${DELEGATION_ADDRESS} \
        --function="getDelegationContractUnbondedAmount" \
}

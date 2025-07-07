ALICE="/home/numbat/dharitri-sdk/testwallets/latest/users/alice.pem"
ADDRESS=$(drtpy data load --key=address-testnet)
DEPLOY_TRANSACTION=$(drtpy data load --key=deployTransaction-testnet)
PROXY=http://localhost:7950
CHAIN_ID=local-testnet

CHILD_CODE=0x"$(xxd -p ../child/output/child.wasm | tr -d '\n')"
DCDT_ISSUE_COST=5000000000000000000

TOKEN_DISPLAY_NAME=0x5772617070656452657761  # "WrappedRewa"
TOKEN_TICKER=0x5752455741  # "WREWA"
INITIAL_SUPPLY=0x03e8 # 1000

deployParent() {
    drtpy --verbose contract deploy --project=${PROJECT} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --outfile="deploy-testnet.interaction.json" --send --proxy=${PROXY} --chain=${CHAIN_ID} || return

    TRANSACTION=$(drtpy data parse --file="deploy-testnet.interaction.json" --expression="data['emittedTransactionHash']")
    ADDRESS=$(drtpy data parse --file="deploy-testnet.interaction.json" --expression="data['contractAddress']")

    drtpy data store --key=address-testnet --value=${ADDRESS}
    drtpy data store --key=deployTransaction-testnet --value=${TRANSACTION}

    echo ""
    echo "Smart contract address: ${ADDRESS}"
}

deployChildThroughParent() {
    drtpy --verbose contract call ${ADDRESS} --recall-nonce --pem=${ALICE} --gas-limit=400000000 --function="deployChildContract" --arguments ${CHILD_CODE} --send --outfile="deploy-child-sc-spam.json" --proxy=${PROXY} --chain=${CHAIN_ID}
}

executeOnDestIssueToken() {
    drtpy --verbose contract call ${ADDRESS} --recall-nonce --pem=${ALICE} --gas-limit=200000000 --value=${DCDT_ISSUE_COST} --function="executeOnDestIssueToken" --arguments ${TOKEN_DISPLAY_NAME} ${TOKEN_TICKER} ${INITIAL_SUPPLY} --send --proxy=${PROXY} --chain=${CHAIN_ID}
}

getChildContractAddress() {
    local QUERY_OUTPUT=$(drtpy --verbose contract query ${ADDRESS} --function="getChildContractAddress" --proxy=${PROXY})
    parseQueryOutput
    parsedAddressToBech32

    CHILD_ADDRESS=${ADDRESS_BECH32}
    echo "Child address: ${CHILD_ADDRESS}"
}

getWrappedRewaTokenIdentifier() {
    getChildContractAddress
    drtpy --verbose contract call ${CHILD_ADDRESS} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --function="getWrappedRewaTokenIdentifier" --send --proxy=${PROXY} --chain=${CHAIN_ID}
}

# helpers

parseQueryOutput() {
    PARSED=$(jq -r '.[0].hex' <<< "${QUERY_OUTPUT}")
}

parsedAddressToBech32() {
    ADDRESS_BECH32=$(drtpy wallet bech32 --encode ${PARSED})
}

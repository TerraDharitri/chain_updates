generateConfig() {
  echo "Generating configuration using values from scripts/variables.sh..."

  TMP_SHARD_OBSERVERCOUNT=$SHARD_OBSERVERCOUNT
  TMP_META_OBSERVERCOUNT=$META_OBSERVERCOUNT
  # set num of observers to 0, they will start with generated keys
  if [[ $MULTI_KEY_NODES -eq 1 ]]; then
    TMP_SHARD_OBSERVERCOUNT=0
    TMP_META_OBSERVERCOUNT=0
  fi

  pushd $TESTNETDIR/filegen
  ./filegen \
    -output-directory $CONFIGGENERATOROUTPUTDIR               \
    -num-of-shards $SHARDCOUNT                                \
    -num-of-nodes-in-each-shard $SHARD_VALIDATORCOUNT         \
    -num-of-observers-in-each-shard $TMP_SHARD_OBSERVERCOUNT  \
    -consensus-group-size $SHARD_CONSENSUS_SIZE               \
    -num-of-metachain-nodes $META_VALIDATORCOUNT              \
    -num-of-observers-in-metachain $TMP_META_OBSERVERCOUNT    \
    -metachain-consensus-group-size $META_CONSENSUS_SIZE      \
    -stake-type $GENESIS_STAKE_TYPE \
    -hysteresis $HYSTERESIS \
    -round-duration $ROUND_DURATION_IN_MS
  popd
}

copyConfig() {
  pushd $TESTNETDIR

  cp ./filegen/"$CONFIGGENERATOROUTPUTDIR"/genesis.json ./node/config
  cp ./filegen/"$CONFIGGENERATOROUTPUTDIR"/nodesSetup.json ./node/config
  cp ./filegen/"$CONFIGGENERATOROUTPUTDIR"/*.pem ./node/config #there might be more .pem files there
  if [[ $MULTI_KEY_NODES -eq 1 ]]; then
      mv ./node/config/"$VALIDATOR_KEY_PEM_FILE" ./node/config/"$MULTI_KEY_PEM_FILE"
      if [[ $EXTRA_KEYS -eq 1 ]]; then
        cat $NODEDIR/config/testKeys/"${EXTRA_KEY_PEM_FILE}" >> ./node/config/"$MULTI_KEY_PEM_FILE"
      fi
  fi
  echo "Configuration files copied from the configuration generator to the working directories of the executables."
  popd
}

copySeednodeConfig() {
  pushd $TESTNETDIR
  cp $SEEDNODEDIR/config/* ./seednode/config
  popd

  pushd $DHARITRIDIR/cmd/keygenerator

  if [[ ! -f "p2pKey.pem" ]]; then
      go build
      ./keygenerator --key-type p2p
  fi

  cp p2pKey.pem $TESTNETDIR/seednode/config
  GENERATED_P2P_PUB_KEY=$(grep for  ./p2pKey.pem | head -1 | grep -oP '[^[:blank:]-]*' | tail -1)
  export P2P_SEEDNODE_ADDRESS="/ip4/$SEEDNODE_IP/tcp/$PORT_SEEDNODE/p2p/$GENERATED_P2P_PUB_KEY"

  popd
}

updateSeednodeConfig() {
  pushd $TESTNETDIR/seednode/config
  cp p2p.toml p2p_edit.toml

  updateTOMLValue p2p_edit.toml "Port" "\"$PORT_SEEDNODE\""

  cp p2p_edit.toml p2p.toml
  rm p2p_edit.toml

  echo "Updated configuration for the Seednode."
  popd
}

copyNodeConfig() {
  pushd $TESTNETDIR
  cp $NODEDIR/config/api.toml ./node/config
  cp $NODEDIR/config/config.toml ./node/config/config_validator.toml
  cp $NODEDIR/config/config.toml ./node/config/config_observer.toml
  cp $NODEDIR/config/economics.toml ./node/config
  cp $NODEDIR/config/ratings.toml ./node/config
  cp $NODEDIR/config/prefs.toml ./node/config
  cp $NODEDIR/config/external.toml ./node/config
  cp $NODEDIR/config/p2p.toml ./node/config
  cp $NODEDIR/config/fullArchiveP2P.toml ./node/config
  cp $NODEDIR/config/enableEpochs.toml ./node/config
  cp $NODEDIR/config/enableRounds.toml ./node/config
  cp $NODEDIR/config/systemSmartContractsConfig.toml ./node/config
  cp $NODEDIR/config/genesisSmartContracts.json ./node/config
  mkdir ./node/config/genesisContracts -p
  cp $NODEDIR/config/genesisContracts/*.* ./node/config/genesisContracts
  mkdir ./node/config/gasSchedules -p
  cp $NODEDIR/config/gasSchedules/*.* ./node/config/gasSchedules

  echo "Configuration files copied from the Node to the working directories of the executables."
  popd
}

updateNodeConfig() {
  pushd $TESTNETDIR/node/config
  cp p2p.toml p2p_edit.toml

  updateTOMLValue p2p_edit.toml "InitialPeerList" "[\"$P2P_SEEDNODE_ADDRESS\"]"

  cp p2p_edit.toml p2p.toml
  rm p2p_edit.toml

  cp nodesSetup.json nodesSetup_edit.json
  
  let startTime="$(date +%s) + $GENESIS_DELAY"
  updateJSONValue nodesSetup_edit.json "startTime" "$startTime"

  updateJSONValue nodesSetup_edit.json "minTransactionVersion" "1"

	if [ $ALWAYS_NEW_CHAINID -eq 1 ]; then
		updateTOMLValue config_validator.toml "ChainID" "\"local-testnet"\"
		updateTOMLValue config_observer.toml "ChainID" "\"local-testnet"\"
	fi

	if [ $ROUNDS_PER_EPOCH -ne 0 ]; then
    sed -i "s,RoundsPerEpoch.*$,RoundsPerEpoch = $ROUNDS_PER_EPOCH," config_observer.toml
    sed -i "s,MinRoundsBetweenEpochs.*$,MinRoundsBetweenEpochs = $ROUNDS_PER_EPOCH," config_observer.toml
	  sed -i "s,RoundsPerEpoch.*$,RoundsPerEpoch = $ROUNDS_PER_EPOCH," config_validator.toml
    sed -i "s,MinRoundsBetweenEpochs.*$,MinRoundsBetweenEpochs = $ROUNDS_PER_EPOCH," config_validator.toml
	fi

  cp nodesSetup_edit.json nodesSetup.json
  rm nodesSetup_edit.json

  if [ $OBSERVERS_ANTIFLOOD_DISABLE -eq 1 ]
  then
     sed -i '/\[Antiflood\]/,/\[Logger\]/ s/true/false/' config_observer.toml
  fi

  updateConfigsForStakingV4

  # Update chain parameters
  updateChainParameters config_observer.toml
  updateChainParameters config_validator.toml

  echo "Updated configuration for Nodes."
  popd
}

updateChainParameters() {
  tomlFile=$1

  sed -i "s,ShardConsensusGroupSize\([^,]*\),ShardConsensusGroupSize = $SHARD_CONSENSUS_SIZE," $tomlFile
  sed -i "s,ShardMinNumNodes\([^,]*\),ShardMinNumNodes = $SHARD_CONSENSUS_SIZE," $tomlFile
  sed -i "s,MetachainConsensusGroupSize\([^,]*\),MetachainConsensusGroupSize = $META_CONSENSUS_SIZE," $tomlFile
  sed -i "s,MetachainMinNumNodes\([^,]*\),MetachainMinNumNodes = $META_CONSENSUS_SIZE," $tomlFile
  sed -i "s,RoundDuration\([^,]*\),RoundDuration = $ROUND_DURATION_IN_MS," $tomlFile
  sed -i "s,Hysteresis\([^,]*\),Hysteresis = $HYSTERESIS," $tomlFile
}

updateConfigsForStakingV4() {
  config=$(cat enableEpochs.toml)

  echo "Updating staking v4 configs"

  # Get the StakingV4Step3EnableEpoch value
  staking_enable_epoch=$(echo "$config" | awk -F '=' '/  StakingV4Step3EnableEpoch/{gsub(/^[ \t]+|[ \t]+$/,"", $2); print $2; exit}')
  # Count the number of entries in MaxNodesChangeEnableEpoch
  entry_count=$(echo "$config" | awk '/MaxNodesChangeEnableEpoch/,/\]/{if ($0 ~ /\{/) {count++}} END {print count}')

  # Check if entry_count is less than 2
  if [ "$entry_count" -lt 2 ]; then
    echo "Not enough entries found to update"
  else
    # Find all entries in MaxNodesChangeEnableEpoch
    all_entries=$(awk '/MaxNodesChangeEnableEpoch/,/\]/{if ($0 ~ /^[[:space:]]*\{/) {p=1}; if (p) print; if ($0 ~ /\]/) p=0}' enableEpochs.toml | grep -vE '^\s*#' | sed '/^\s*$/d')

    # Get the index of the entry with EpochEnable equal to StakingV4Step3EnableEpoch
    index=$(echo "$all_entries" | grep -n "EpochEnable = $staking_enable_epoch" | cut -d: -f1)

    if [[ -z "${index// }" ]]; then
      echo -e "\033[1;33mWarning: MaxNodesChangeEnableEpoch does not contain an entry enable epoch for StakingV4Step3EnableEpoch, nodes might fail to start...\033[0m"
    else
      prev_entry=$(echo "$all_entries" | sed -n "$((index-1))p")
      curr_entry=$(echo "$all_entries" | sed -n "$((index))p")

      # Extract the value of MaxNumNodes & NodesToShufflePerShard from prev_entry
      max_nodes_from_prev_epoch=$(echo "$prev_entry" | awk -F 'MaxNumNodes = ' '{print $2}' | cut -d ',' -f1)
      nodes_to_shuffle_per_shard=$(echo "$prev_entry" | awk -F 'NodesToShufflePerShard = ' '{gsub(/[^0-9]+/, "", $2); print $2}')

      # Calculate the new MaxNumNodes value based on the formula
      new_max_nodes=$((max_nodes_from_prev_epoch - (SHARDCOUNT + 1) * nodes_to_shuffle_per_shard))
      curr_entry_updated=$(echo "$curr_entry" | awk -v new_max_nodes="$new_max_nodes" '{gsub(/MaxNumNodes = [0-9]+,/, "MaxNumNodes = " new_max_nodes ",")}1')

      echo "Updating entry in MaxNodesChangeEnableEpoch from $curr_entry to $curr_entry_updated"

      sed -i "/$staking_enable_epoch/,/$staking_enable_epoch/ s|.*$curr_entry.*|$curr_entry_updated|" enableEpochs.toml
    fi
  fi
}

copyProxyConfig() {
  pushd $TESTNETDIR

  cp -r $PROXYDIR/config/apiConfig ./proxy/config/
  cp $PROXYDIR/config/config.toml ./proxy/config/
  cp -r $PROXYDIR/config/apiConfig ./proxy/config

  cp ./node/config/economics.toml ./proxy/config/
  cp ./node/config/external.toml ./proxy/config/
  cp ./node/config/walletKey.pem ./proxy/config

  echo "Copied configuration for the Proxy."
  popd
}

updateProxyConfig() {
  pushd $TESTNETDIR/proxy/config
  cp config.toml config_edit.toml

  # Truncate config.toml before the [[Observers]] list
  sed -i -n '/\[\[Observers\]\]/q;p' config_edit.toml
  
  updateTOMLValue config_edit.toml "ServerPort" $PORT_PROXY
  generateProxyObserverList config_edit.toml

  cp config_edit.toml config.toml
  rm config_edit.toml

  echo "Updated configuration for the Proxy."
  popd
}

copyTxGenConfig() {
  pushd $TESTNETDIR

  cp $TXGENDIR/config/config.toml ./txgen/config/

  cp $TXGENDIR/config/sc.toml ./txgen/config/
  cp $TXGENDIR/config/*.wasm ./txgen/config/

  cp ./node/config/economics.toml ./txgen/config/
  cp ./node/config/walletKey.pem ./txgen/config
  cp ./node/config/enableEpochs.toml ./txgen/config/nodeConfig/config

  echo "Copied configuration for the TxGen."
  popd
}

updateTxGenConfig() {
  pushd $TESTNETDIR/txgen/config
  cp config.toml config_edit.toml

  updateTOMLValue config_edit.toml "ServerPort" $PORT_TXGEN
  updateTOMLValue config_edit.toml "ProxyServerURL" "\"http://127.0.0.1:$PORT_PROXY\""
  sed -i "/Scenarios = \[/c ${TXGEN_SCENARIOS_LINE}" config_edit.toml

  cp config_edit.toml config.toml
  rm config_edit.toml

  echo "Updated configuration for the TxGen."
  popd
}


generateProxyObserverList() {
  OBSERVER_INDEX=0
  OUTPUTFILE=$!
  # Start Shard Observers
  (( max_shard_id=$SHARDCOUNT - 1 ))
  for SHARD in $(seq 0 1 $max_shard_id); do
    for _ in $(seq $SHARD_OBSERVERCOUNT); do
      (( PORT=$PORT_ORIGIN_OBSERVER_REST+$OBSERVER_INDEX))

      echo "[[Observers]]" >> config_edit.toml
      echo "   ShardId = $SHARD" >> config_edit.toml
      echo "   Address = \"http://127.0.0.1:$PORT\"" >> config_edit.toml
      echo ""$'\n' >> config_edit.toml

      (( OBSERVER_INDEX++ ))
    done
  done
  # Start Meta Observers
  for META_OBSERVER in $(seq $META_OBSERVERCOUNT); do
    (( PORT=$PORT_ORIGIN_OBSERVER_REST+$OBSERVER_INDEX ))

      echo "[[Observers]]" >> config_edit.toml
      echo "   ShardId = $METASHARD_ID" >> config_edit.toml
      echo "   Address = \"http://127.0.0.1:$PORT\"" >> config_edit.toml
      echo ""$'\n' >> config_edit.toml

      (( OBSERVER_INDEX++ ))
    done
}

updateTOMLValue() {
  local filename=$1
  local key=$2
  local value=$3

  escaped_value=$(printf "%q" $value)

  sed -i "s,$key = .*\$,$key = $escaped_value," $filename
}


updateJSONValue() {
  local filename=$1
  local key=$2
  local value=$3

  escaped_value=$(printf "%q" $value)

  sed -i "s,\"$key\": .*\$,\"$key\": $escaped_value\,," $filename
}

changeConfigForHardfork(){
  pushd $TESTNETDIR/node/config

  export FIRST_PUBKEY=$(cat nodesSetup.json | grep pubkey -m 1 | sed -E 's/^.*"([0-9a-f]+)".*$/\1/g')
  updateTOMLValue config_observer.toml "PublicKeyToListenFrom" "\"$FIRST_PUBKEY\""
  updateTOMLValue config_validator.toml "PublicKeyToListenFrom" "\"$FIRST_PUBKEY\""

  popd
}

copyBackConfigs(){
  pushd $TESTNETDIR

  echo "trying to copy-back the configs"
  cp ./node/config/*.* $NODEDIR/config
  cp $NODEDIR/config/config_validator.toml $NODEDIR/config/config.toml

  popd
}

package configs

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/TerraDharitri/drt-go-chain/common/factory"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/genesis/data"
	"github.com/TerraDharitri/drt-go-chain/node"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	"github.com/TerraDharitri/drt-go-chain/testscommon"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	shardingCore "github.com/TerraDharitri/drt-go-chain-core/core/sharding"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/mcl"
)

var oneRewa = big.NewInt(1000000000000000000)
var initialStakedRewaPerNode = big.NewInt(0).Mul(oneRewa, big.NewInt(2500))
var initialSupply = big.NewInt(0).Mul(oneRewa, big.NewInt(20000000)) // 20 million REWA
const (
	// ChainID contains the chain id
	ChainID = "chain"

	// ChainSimulatorConsensusGroupSize defines the size of the consensus group for chain simulator
	ChainSimulatorConsensusGroupSize = 1
	allValidatorsPemFileName         = "allValidatorsKeys.pem"
)

// ArgsChainSimulatorConfigs holds all the components needed to create the chain simulator configs
type ArgsChainSimulatorConfigs struct {
	NumOfShards                 uint32
	OriginalConfigsPath         string
	GenesisTimeStamp            int64
	RoundDurationInMillis       uint64
	TempDir                     string
	MinNodesPerShard            uint32
	ConsensusGroupSize          uint32
	MetaChainMinNodes           uint32
	MetaChainConsensusGroupSize uint32
	Hysteresis                  float32
	InitialEpoch                uint32
	RoundsPerEpoch              core.OptionalUint64
	NumNodesWaitingListShard    uint32
	NumNodesWaitingListMeta     uint32
	AlterConfigsFunction        func(cfg *config.Configs)
}

// ArgsConfigsSimulator holds the configs for the chain simulator
type ArgsConfigsSimulator struct {
	GasScheduleFilename   string
	Configs               config.Configs
	ValidatorsPrivateKeys []crypto.PrivateKey
	InitialWallets        *dtos.InitialWalletKeys
}

// CreateChainSimulatorConfigs will create the chain simulator configs
func CreateChainSimulatorConfigs(args ArgsChainSimulatorConfigs) (*ArgsConfigsSimulator, error) {
	configs, err := testscommon.CreateTestConfigs(args.TempDir, args.OriginalConfigsPath)
	if err != nil {
		return nil, err
	}

	configs.GeneralConfig.GeneralSettings.ChainID = ChainID

	// empty genesis smart contracts file
	err = os.WriteFile(configs.ConfigurationPathsHolder.SmartContracts, []byte("[]"), os.ModePerm)
	if err != nil {
		return nil, err
	}

	// update genesis.json
	initialWallets, err := generateGenesisFile(args, configs)
	if err != nil {
		return nil, err
	}

	// generate validators key and nodesSetup.json
	privateKeys, publicKeys, err := generateValidatorsKeyAndUpdateFiles(
		configs,
		initialWallets.StakeWallets,
		args,
	)
	if err != nil {
		return nil, err
	}

	configs.ConfigurationPathsHolder.AllValidatorKeys = path.Join(args.TempDir, allValidatorsPemFileName)
	err = generateValidatorsPem(configs.ConfigurationPathsHolder.AllValidatorKeys, publicKeys, privateKeys)
	if err != nil {
		return nil, err
	}

	configs.GeneralConfig.SmartContractsStorage.DB.Type = string(storageunit.MemoryDB)
	configs.GeneralConfig.SmartContractsStorageForSCQuery.DB.Type = string(storageunit.MemoryDB)
	configs.GeneralConfig.SmartContractsStorageSimulate.DB.Type = string(storageunit.MemoryDB)

	eligibleNodes := args.MinNodesPerShard*args.NumOfShards + args.MetaChainMinNodes
	waitingNodes := args.NumNodesWaitingListShard*args.NumOfShards + args.NumNodesWaitingListMeta

	SetMaxNumberOfNodesInConfigs(configs, eligibleNodes, waitingNodes, args.NumOfShards)

	// set compatible trie configs
	configs.GeneralConfig.StateTriesConfig.SnapshotsEnabled = false
	// enable db lookup extension
	configs.GeneralConfig.DbLookupExtensions.Enabled = true

	configs.GeneralConfig.EpochStartConfig.ExtraDelayForRequestBlockInfoInMilliseconds = 1
	configs.GeneralConfig.EpochStartConfig.GenesisEpoch = args.InitialEpoch
	configs.GeneralConfig.EpochStartConfig.MinRoundsBetweenEpochs = 1

	if args.RoundsPerEpoch.HasValue {
		configs.GeneralConfig.EpochStartConfig.RoundsPerEpoch = int64(args.RoundsPerEpoch.Value)
	}

	gasScheduleName, err := GetLatestGasScheduleFilename(configs.ConfigurationPathsHolder.GasScheduleDirectoryName)
	if err != nil {
		return nil, err
	}

	updateConfigsChainParameters(args, configs)
	node.ApplyArchCustomConfigs(configs)

	if args.AlterConfigsFunction != nil {
		args.AlterConfigsFunction(configs)
		// this is needed to keep in sync Andromeda flag and the second entry from chain parameters
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[1].EnableEpoch = configs.EpochConfig.EnableEpochs.AndromedaEnableEpoch
	}

	return &ArgsConfigsSimulator{
		Configs:               *configs,
		ValidatorsPrivateKeys: privateKeys,
		GasScheduleFilename:   gasScheduleName,
		InitialWallets:        initialWallets,
	}, nil
}

func updateConfigsChainParameters(args ArgsChainSimulatorConfigs, configs *config.Configs) {
	for idx := 0; idx < len(configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch); idx++ {
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].ShardMinNumNodes = args.MinNodesPerShard
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].MetachainMinNumNodes = args.MetaChainMinNodes
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].MetachainConsensusGroupSize = args.MetaChainConsensusGroupSize
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].ShardConsensusGroupSize = args.ConsensusGroupSize
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].RoundDuration = args.RoundDurationInMillis
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[idx].Hysteresis = args.Hysteresis
	}

	if len(configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch) > 1 {
		configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch[1].ShardConsensusGroupSize = args.MinNodesPerShard
	}
}

// SetMaxNumberOfNodesInConfigs will correctly set the max number of nodes in configs
func SetMaxNumberOfNodesInConfigs(cfg *config.Configs, eligibleNodes uint32, waitingNodes uint32, numOfShards uint32) {
	cfg.SystemSCConfig.StakingSystemSCConfig.MaxNumberOfNodesForStake = uint64(eligibleNodes + waitingNodes)
	numMaxNumNodesEnableEpochs := len(cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch)
	for idx := 0; idx < numMaxNumNodesEnableEpochs-1; idx++ {
		cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[idx].MaxNumNodes = eligibleNodes + waitingNodes
	}

	cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[numMaxNumNodesEnableEpochs-1].EpochEnable = cfg.EpochConfig.EnableEpochs.StakingV4Step3EnableEpoch
	prevEntry := cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[numMaxNumNodesEnableEpochs-2]
	cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[numMaxNumNodesEnableEpochs-1].NodesToShufflePerShard = prevEntry.NodesToShufflePerShard

	stakingV4NumNodes := eligibleNodes + waitingNodes
	if stakingV4NumNodes-(numOfShards+1)*prevEntry.NodesToShufflePerShard >= eligibleNodes {
		// prevent the case in which we are decreasing the eligible number of nodes because we are working with 0 waiting list size
		stakingV4NumNodes -= (numOfShards + 1) * prevEntry.NodesToShufflePerShard
	}
	cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[numMaxNumNodesEnableEpochs-1].MaxNumNodes = stakingV4NumNodes
}

// SetQuickJailRatingConfig will set the rating config in a way that leads to rapid jailing of a node
func SetQuickJailRatingConfig(cfg *config.Configs) {
	cfg.RatingsConfig.ShardChain.RatingStepsByEpoch[0].ConsecutiveMissedBlocksPenalty = 100
	cfg.RatingsConfig.ShardChain.RatingStepsByEpoch[0].HoursToMaxRatingFromStartRating = 1
	cfg.RatingsConfig.MetaChain.RatingStepsByEpoch[0].ConsecutiveMissedBlocksPenalty = 100
	cfg.RatingsConfig.MetaChain.RatingStepsByEpoch[0].HoursToMaxRatingFromStartRating = 1
}

// SetStakingV4ActivationEpochs configures activation epochs for Staking V4.
// It takes an initial epoch and sets three consecutive steps for enabling Staking V4 features:
//   - Step 1 activation epoch
//   - Step 2 activation epoch
//   - Step 3 activation epoch
func SetStakingV4ActivationEpochs(cfg *config.Configs, initialEpoch uint32) {
	cfg.EpochConfig.EnableEpochs.StakeLimitsEnableEpoch = initialEpoch
	cfg.EpochConfig.EnableEpochs.StakingV4Step1EnableEpoch = initialEpoch
	cfg.EpochConfig.EnableEpochs.StakingV4Step2EnableEpoch = initialEpoch + 1
	cfg.EpochConfig.EnableEpochs.StakingV4Step3EnableEpoch = initialEpoch + 2

	// Set the MaxNodesChange enable epoch for index 2
	cfg.EpochConfig.EnableEpochs.MaxNodesChangeEnableEpoch[2].EpochEnable = initialEpoch + 2
	cfg.SystemSCConfig.StakingSystemSCConfig.NodeLimitPercentage = 1
}

func generateGenesisFile(args ArgsChainSimulatorConfigs, configs *config.Configs) (*dtos.InitialWalletKeys, error) {
	addressConverter, err := factory.NewPubkeyConverter(configs.GeneralConfig.AddressPubkeyConverter)
	if err != nil {
		return nil, err
	}

	initialWalletKeys := &dtos.InitialWalletKeys{
		BalanceWallets: make(map[uint32]*dtos.WalletKey),
		StakeWallets:   make([]*dtos.WalletKey, 0),
	}

	addresses := make([]data.InitialAccount, 0)
	numOfNodes := int((args.NumNodesWaitingListShard+args.MinNodesPerShard)*args.NumOfShards + args.NumNodesWaitingListMeta + args.MetaChainMinNodes)
	for i := 0; i < numOfNodes; i++ {
		wallet, errGenerate := generateWalletKey(addressConverter)
		if errGenerate != nil {
			return nil, errGenerate
		}

		stakedValue := big.NewInt(0).Set(initialStakedRewaPerNode)
		addresses = append(addresses, data.InitialAccount{
			Address:      wallet.Address.Bech32,
			StakingValue: stakedValue,
			Supply:       stakedValue,
		})

		initialWalletKeys.StakeWallets = append(initialWalletKeys.StakeWallets, wallet)
	}

	// generate an address for every shard
	initialBalance := big.NewInt(0).Set(initialSupply)
	totalStakedValue := big.NewInt(int64(numOfNodes))
	totalStakedValue = totalStakedValue.Mul(totalStakedValue, big.NewInt(0).Set(initialStakedRewaPerNode))
	initialBalance = initialBalance.Sub(initialBalance, totalStakedValue)

	walletBalance := big.NewInt(0).Set(initialBalance)
	walletBalance.Div(walletBalance, big.NewInt(int64(args.NumOfShards)))

	// remainder = balance % numTotalWalletKeys
	remainder := big.NewInt(0).Set(initialBalance)
	remainder.Mod(remainder, big.NewInt(int64(args.NumOfShards)))

	for shardID := uint32(0); shardID < args.NumOfShards; shardID++ {
		walletKey, errG := generateWalletKeyForShard(shardID, args.NumOfShards, addressConverter)
		if errG != nil {
			return nil, errG
		}

		addresses = append(addresses, data.InitialAccount{
			Address: walletKey.Address.Bech32,
			Balance: big.NewInt(0).Set(walletBalance),
			Supply:  big.NewInt(0).Set(walletBalance),
		})

		initialWalletKeys.BalanceWallets[shardID] = walletKey
	}

	addresses[len(addresses)-1].Balance.Add(walletBalance, remainder)
	addresses[len(addresses)-1].Supply.Add(walletBalance, remainder)

	addressesBytes, errM := json.Marshal(addresses)
	if errM != nil {
		return nil, errM
	}

	err = os.WriteFile(configs.ConfigurationPathsHolder.Genesis, addressesBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return initialWalletKeys, nil
}

func generateValidatorsKeyAndUpdateFiles(
	configs *config.Configs,
	stakeWallets []*dtos.WalletKey,
	args ArgsChainSimulatorConfigs,
) ([]crypto.PrivateKey, []crypto.PublicKey, error) {
	blockSigningGenerator := signing.NewKeyGenerator(mcl.NewSuiteBLS12())

	nodesSetupFile := configs.ConfigurationPathsHolder.Nodes
	nodes := &sharding.NodesSetup{}
	err := core.LoadJsonFile(nodes, nodesSetupFile)
	if err != nil {
		return nil, nil, err
	}

	nodes.RoundDuration = args.RoundDurationInMillis
	nodes.StartTime = args.GenesisTimeStamp

	nodes.Hysteresis = 0

	nodes.InitialNodes = make([]*sharding.InitialNode, 0)
	configs.NodesConfig.InitialNodes = make([]*config.InitialNodeConfig, 0)
	privateKeys := make([]crypto.PrivateKey, 0)
	publicKeys := make([]crypto.PublicKey, 0)
	walletIndex := 0
	// generate meta keys
	for idx := uint32(0); idx < args.NumNodesWaitingListMeta+args.MetaChainMinNodes; idx++ {
		sk, pk := blockSigningGenerator.GeneratePair()
		privateKeys = append(privateKeys, sk)
		publicKeys = append(publicKeys, pk)

		pkBytes, errB := pk.ToByteArray()
		if errB != nil {
			return nil, nil, errB
		}

		nodes.InitialNodes = append(nodes.InitialNodes, &sharding.InitialNode{
			PubKey:  hex.EncodeToString(pkBytes),
			Address: stakeWallets[walletIndex].Address.Bech32,
		})

		configs.NodesConfig.InitialNodes = append(configs.NodesConfig.InitialNodes, &config.InitialNodeConfig{
			PubKey:        hex.EncodeToString(pkBytes),
			Address:       stakeWallets[walletIndex].Address.Bech32,
			InitialRating: 5000001,
		})

		walletIndex++
	}

	// generate shard keys
	for idx1 := uint32(0); idx1 < args.NumOfShards; idx1++ {
		for idx2 := uint32(0); idx2 < args.NumNodesWaitingListShard+args.MinNodesPerShard; idx2++ {
			sk, pk := blockSigningGenerator.GeneratePair()
			privateKeys = append(privateKeys, sk)
			publicKeys = append(publicKeys, pk)

			pkBytes, errB := pk.ToByteArray()
			if errB != nil {
				return nil, nil, errB
			}

			nodes.InitialNodes = append(nodes.InitialNodes, &sharding.InitialNode{
				PubKey:  hex.EncodeToString(pkBytes),
				Address: stakeWallets[walletIndex].Address.Bech32,
			})

			configs.NodesConfig.InitialNodes = append(configs.NodesConfig.InitialNodes, &config.InitialNodeConfig{
				PubKey:        hex.EncodeToString(pkBytes),
				Address:       stakeWallets[walletIndex].Address.Bech32,
				InitialRating: 5000001,
			})

			walletIndex++
		}
	}

	marshaledNodes, err := json.Marshal(nodes)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(nodesSetupFile, marshaledNodes, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	return privateKeys, publicKeys, nil
}

func generateValidatorsPem(validatorsFile string, publicKeys []crypto.PublicKey, privateKey []crypto.PrivateKey) error {
	validatorPubKeyConverter, err := pubkeyConverter.NewHexPubkeyConverter(96)
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	for idx := 0; idx < len(publicKeys); idx++ {
		publicKeyBytes, errA := publicKeys[idx].ToByteArray()
		if errA != nil {
			return errA
		}

		pkString, errE := validatorPubKeyConverter.Encode(publicKeyBytes)
		if errE != nil {
			return errE
		}

		privateKeyBytes, errP := privateKey[idx].ToByteArray()
		if errP != nil {
			return errP
		}

		blk := pem.Block{
			Type:  "PRIVATE KEY for " + pkString,
			Bytes: []byte(hex.EncodeToString(privateKeyBytes)),
		}

		err = pem.Encode(&buff, &blk)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(validatorsFile, buff.Bytes(), 0644)
}

// GetLatestGasScheduleFilename will parse the provided path and get the latest gas schedule filename
func GetLatestGasScheduleFilename(directory string) (string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	extension := ".toml"
	versionMarker := "V"

	highestVersion := 0
	filename := ""
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		splt := strings.Split(name, versionMarker)
		if len(splt) != 2 {
			continue
		}

		versionAsString := splt[1][:len(splt[1])-len(extension)]
		number, errConversion := strconv.Atoi(versionAsString)
		if errConversion != nil {
			continue
		}

		if number > highestVersion {
			highestVersion = number
			filename = name
		}
	}

	return path.Join(directory, filename), nil
}

func generateWalletKeyForShard(shardID, numOfShards uint32, converter core.PubkeyConverter) (*dtos.WalletKey, error) {
	for {
		walletKey, err := generateWalletKey(converter)
		if err != nil {
			return nil, err
		}

		addressShardID := shardingCore.ComputeShardID(walletKey.Address.Bytes, numOfShards)
		if addressShardID != shardID {
			continue
		}

		return walletKey, nil
	}
}

func generateWalletKey(converter core.PubkeyConverter) (*dtos.WalletKey, error) {
	walletSuite := ed25519.NewEd25519()
	walletKeyGenerator := signing.NewKeyGenerator(walletSuite)

	sk, pk := walletKeyGenerator.GeneratePair()
	pubKeyBytes, err := pk.ToByteArray()
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := sk.ToByteArray()
	if err != nil {
		return nil, err
	}

	bech32Address, err := converter.Encode(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return &dtos.WalletKey{
		Address: dtos.WalletAddress{
			Bech32: bech32Address,
			Bytes:  pubKeyBytes,
		},
		PrivateKeyHex: hex.EncodeToString(privateKeyBytes),
	}, nil
}

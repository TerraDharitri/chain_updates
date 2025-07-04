package genesis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/parsers"
	"github.com/TerraDharitri/drt-go-chain/trie/keyBuilder"
	"github.com/TerraDharitri/drt-go-chain/update"
)

var _ update.ExportHandler = (*stateExport)(nil)

// ArgsNewStateExporter defines the arguments needed to create new state exporter
type ArgsNewStateExporter struct {
	ShardCoordinator         sharding.Coordinator
	StateSyncer              update.StateSyncer
	Marshalizer              marshal.Marshalizer
	Hasher                   hashing.Hasher
	HardforkStorer           update.HardforkStorer
	ExportFolder             string
	AddressPubKeyConverter   core.PubkeyConverter
	ValidatorPubKeyConverter core.PubkeyConverter
	GenesisNodesSetupHandler update.GenesisNodesSetupHandler
}

type stateExport struct {
	stateSyncer              update.StateSyncer
	shardCoordinator         sharding.Coordinator
	marshalizer              marshal.Marshalizer
	hasher                   hashing.Hasher
	hardforkStorer           update.HardforkStorer
	exportFolder             string
	addressPubKeyConverter   core.PubkeyConverter
	validatorPubKeyConverter core.PubkeyConverter
	genesisNodesSetupHandler update.GenesisNodesSetupHandler
}

var log = logger.GetOrCreate("update/genesis")

// NewStateExporter exports all the data at a specific moment to a hardfork storer
func NewStateExporter(args ArgsNewStateExporter) (*stateExport, error) {
	if check.IfNil(args.ShardCoordinator) {
		return nil, data.ErrNilShardCoordinator
	}
	if check.IfNil(args.StateSyncer) {
		return nil, update.ErrNilStateSyncer
	}
	if check.IfNil(args.Marshalizer) {
		return nil, data.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return nil, update.ErrNilHasher
	}
	if check.IfNil(args.HardforkStorer) {
		return nil, update.ErrNilHardforkStorer
	}
	if len(args.ExportFolder) == 0 {
		return nil, update.ErrEmptyExportFolderPath
	}
	if check.IfNil(args.AddressPubKeyConverter) {
		return nil, fmt.Errorf("%w for address", update.ErrNilPubKeyConverter)
	}
	if check.IfNil(args.ValidatorPubKeyConverter) {
		return nil, fmt.Errorf("%w for validators", update.ErrNilPubKeyConverter)
	}
	if check.IfNil(args.GenesisNodesSetupHandler) {
		return nil, update.ErrNilGenesisNodesSetupHandler
	}

	se := &stateExport{
		stateSyncer:              args.StateSyncer,
		shardCoordinator:         args.ShardCoordinator,
		marshalizer:              args.Marshalizer,
		hasher:                   args.Hasher,
		hardforkStorer:           args.HardforkStorer,
		exportFolder:             args.ExportFolder,
		addressPubKeyConverter:   args.AddressPubKeyConverter,
		validatorPubKeyConverter: args.ValidatorPubKeyConverter,
		genesisNodesSetupHandler: args.GenesisNodesSetupHandler,
	}

	return se, nil
}

// ExportAll syncs and exports all the data from every shard for a certain epoch start block
func (se *stateExport) ExportAll(epoch uint32) error {
	defer func() {
		errClose := se.hardforkStorer.Close()
		log.LogIfError(errClose)
	}()

	err := se.stateSyncer.SyncAllState(epoch)
	if err != nil {
		return err
	}

	err = se.exportEpochStartMetaBlock()
	if err != nil {
		return err
	}

	err = se.exportUnFinishedMetaBlocks()
	if err != nil {
		return err
	}

	err = se.exportAllTries()
	if err != nil {
		return err
	}

	err = se.exportAllMiniBlocks()
	if err != nil {
		return err
	}

	err = se.exportAllTransactions()
	if err != nil {
		return err
	}

	err = se.exportAllValidatorsInfo()
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportAllTransactions() error {
	toExportTransactions, err := se.stateSyncer.GetAllTransactions()
	if err != nil {
		return err
	}

	log.Debug("Starting export for transactions", "len", len(toExportTransactions))
	for key, tx := range toExportTransactions {
		errExport := se.exportTx(key, tx)
		if errExport != nil {
			return errExport
		}
	}

	return se.hardforkStorer.FinishedIdentifier(TransactionsIdentifier)
}

func (se *stateExport) exportAllValidatorsInfo() error {
	toExportValidatorsInfo, err := se.stateSyncer.GetAllValidatorsInfo()
	if err != nil {
		return err
	}

	log.Debug("Starting export for validators info", "len", len(toExportValidatorsInfo))
	for key, validatorInfo := range toExportValidatorsInfo {
		errExport := se.exportValidatorInfo(key, validatorInfo)
		if errExport != nil {
			return errExport
		}
	}

	return se.hardforkStorer.FinishedIdentifier(ValidatorsInfoIdentifier)
}

func (se *stateExport) exportAllMiniBlocks() error {
	toExportMBs, err := se.stateSyncer.GetAllMiniBlocks()
	if err != nil {
		return err
	}

	log.Debug("Starting export for miniBlocks", "len", len(toExportMBs))
	for key, mb := range toExportMBs {
		errExport := se.exportMBs(key, mb)
		if errExport != nil {
			return errExport
		}
	}

	return se.hardforkStorer.FinishedIdentifier(MiniBlocksIdentifier)
}

func (se *stateExport) exportAllTries() error {
	toExportTries, err := se.stateSyncer.GetAllTries()
	if err != nil {
		return err
	}

	log.Debug("Starting export for tries", "len", len(toExportTries))
	for key, trie := range toExportTries {
		err = se.exportTrie(key, trie)
		if err != nil {
			return err
		}
	}

	return nil
}

func (se *stateExport) exportEpochStartMetaBlock() error {
	metaBlock, err := se.stateSyncer.GetEpochStartMetaBlock()
	if err != nil {
		return err
	}

	log.Debug("Starting export for epoch start metaBlock")
	err = se.exportMetaBlock(metaBlock, EpochStartMetaBlockIdentifier)
	if err != nil {
		return err
	}

	err = se.hardforkStorer.FinishedIdentifier(EpochStartMetaBlockIdentifier)
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportUnFinishedMetaBlocks() error {
	unFinishedMetaBlocks, err := se.stateSyncer.GetUnFinishedMetaBlocks()
	if err != nil {
		return err
	}

	log.Debug("Starting export for unFinished metaBlocks", "len", len(unFinishedMetaBlocks))
	for _, metaBlock := range unFinishedMetaBlocks {
		errExportMetaBlock := se.exportMetaBlock(metaBlock, UnFinishedMetaBlocksIdentifier)
		if errExportMetaBlock != nil {
			return errExportMetaBlock
		}
	}

	err = se.hardforkStorer.FinishedIdentifier(UnFinishedMetaBlocksIdentifier)
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportMetaBlock(metaBlock data.HeaderHandler, identifier string) error {
	jsonData, err := json.Marshal(metaBlock)
	if err != nil {
		return err
	}

	metaHash := se.hasher.Compute(string(jsonData))
	versionKey := CreateVersionKey(metaBlock, metaHash)
	err = se.hardforkStorer.Write(identifier, []byte(versionKey), jsonData)
	if err != nil {
		return err
	}

	log.Debug("Exported metaBlock",
		"identifier", identifier,
		"version key", versionKey,
		"hash", metaHash,
		"epoch", metaBlock.GetEpoch(),
		"round", metaBlock.GetRound(),
		"nonce", metaBlock.GetNonce(),
		"start of epoch block", metaBlock.GetNonce() == 0 || metaBlock.IsStartOfEpochBlock(),
		"rootHash", metaBlock.GetRootHash(),
	)

	return nil
}

func (se *stateExport) exportTrie(key string, trie common.Trie) error {
	identifier := TrieIdentifier + atSep + key

	accType, shId, err := GetTrieTypeAndShId(identifier)
	if err != nil {
		return err
	}

	rootHash, err := trie.RootHash()
	if err != nil {
		return err
	}

	leavesChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err = trie.GetAllLeavesOnChannel(
		leavesChannels,
		context.Background(),
		rootHash,
		keyBuilder.NewKeyBuilder(),
		parsers.NewMainTrieLeafParser(),
	)
	if err != nil {
		return err
	}

	if accType == ValidatorAccount {
		validatorData, err := getValidatorDataFromLeaves(leavesChannels, se.marshalizer)
		if err != nil {
			return err
		}

		nodesSetupFilePath := filepath.Join(se.exportFolder, common.NodesSetupJsonFileName)
		err = se.exportNodesSetupJson(validatorData)
		if err == nil {
			log.Debug("hardfork nodesSetup.json exported successfully", "file path", nodesSetupFilePath)
		} else {
			log.Warn("hardfork nodesSetup.json not exported", "file path", nodesSetupFilePath, "error", err)
		}

		return err
	}

	if shId > se.shardCoordinator.NumberOfShards() && shId != core.MetachainShardId {
		return nodesCoordinator.ErrInvalidShardId
	}

	rootHashKey := CreateRootHashKey(key)

	err = se.hardforkStorer.Write(identifier, []byte(rootHashKey), rootHash)
	if err != nil {
		return err
	}

	if accType == DataTrie {
		return se.exportDataTries(leavesChannels, accType, shId, identifier)
	}

	log.Debug("exporting trie",
		"identifier", identifier,
		"root hash", rootHash,
	)

	return se.exportAccountLeaves(leavesChannels, accType, shId, identifier)
}

func (se *stateExport) exportDataTries(
	leavesChannels *common.TrieIteratorChannels,
	accType Type,
	shId uint32,
	identifier string,
) error {
	for leaf := range leavesChannels.LeavesChan {
		keyToExport := CreateAccountKey(accType, shId, leaf.Key())
		err := se.hardforkStorer.Write(identifier, []byte(keyToExport), leaf.Value())
		if err != nil {
			return err
		}
	}

	err := leavesChannels.ErrChan.ReadFromChanNonBlocking()
	if err != nil {
		return err
	}

	return se.hardforkStorer.FinishedIdentifier(identifier)
}

func (se *stateExport) exportAccountLeaves(
	leavesChannels *common.TrieIteratorChannels,
	accType Type,
	shId uint32,
	identifier string,
) error {
	for leaf := range leavesChannels.LeavesChan {
		keyToExport := CreateAccountKey(accType, shId, leaf.Key())
		err := se.hardforkStorer.Write(identifier, []byte(keyToExport), leaf.Value())
		if err != nil {
			return err
		}
	}

	err := leavesChannels.ErrChan.ReadFromChanNonBlocking()
	if err != nil {
		return err
	}

	return se.hardforkStorer.FinishedIdentifier(identifier)
}

func (se *stateExport) exportMBs(key string, mb *block.MiniBlock) error {
	marshaledData, err := json.Marshal(mb)
	if err != nil {
		return err
	}

	keyToSave := CreateMiniBlockKey(key)

	err = se.hardforkStorer.Write(MiniBlocksIdentifier, []byte(keyToSave), marshaledData)
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportTx(key string, tx data.TransactionHandler) error {
	marshaledData, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	keyToSave := CreateTransactionKey(key, tx)

	err = se.hardforkStorer.Write(TransactionsIdentifier, []byte(keyToSave), marshaledData)
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportValidatorInfo(key string, validatorInfo *state.ShardValidatorInfo) error {
	marshaledData, err := json.Marshal(validatorInfo)
	if err != nil {
		return err
	}

	keyToSave := CreateValidatorInfoKey(key)

	err = se.hardforkStorer.Write(ValidatorsInfoIdentifier, []byte(keyToSave), marshaledData)
	if err != nil {
		return err
	}

	return nil
}

func (se *stateExport) exportNodesSetupJson(validators state.ShardValidatorsInfoMapHandler) error {
	acceptedListsForExport := []common.PeerType{common.EligibleList, common.WaitingList, common.JailedList}
	initialNodes := make([]*config.InitialNodeConfig, 0)

	for _, validator := range validators.GetAllValidatorsInfo() {
		if shouldExportValidator(validator, acceptedListsForExport) {

			pubKey, err := se.validatorPubKeyConverter.Encode(validator.GetPublicKey())
			if err != nil {
				return nil
			}

			rewardAddress, err := se.addressPubKeyConverter.Encode(validator.GetRewardAddress())
			if err != nil {
				return nil
			}

			initialNodes = append(initialNodes, &config.InitialNodeConfig{
				PubKey:        pubKey,
				Address:       rewardAddress,
				InitialRating: validator.GetRating(),
			})
		}
	}

	sort.SliceStable(initialNodes, func(i, j int) bool {
		return strings.Compare(initialNodes[i].PubKey, initialNodes[j].PubKey) < 0
	})

	exportedNodesConfig := se.genesisNodesSetupHandler.ExportNodesConfig()
	exportedNodesConfig.InitialNodes = initialNodes

	nodesSetupBytes, err := json.MarshalIndent(exportedNodesConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(se.exportFolder, common.NodesSetupJsonFileName), nodesSetupBytes, 0664)
}

// IsInterfaceNil returns true if underlying object is nil
func (se *stateExport) IsInterfaceNil() bool {
	return se == nil
}

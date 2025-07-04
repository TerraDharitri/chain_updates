package factory

import (
	"bytes"
	"fmt"
	"path"
	"runtime/pprof"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/p2p"
)

var log = logger.GetOrCreate("integrationtests")

// PrintStack -
func PrintStack() {
	buffer := new(bytes.Buffer)
	err := pprof.Lookup("goroutine").WriteTo(buffer, 2)
	if err != nil {
		log.Debug("could not dump goroutines")
	}

	log.Debug(fmt.Sprintf("\n%s", buffer.String()))
}

// CreateDefaultConfig -
func CreateDefaultConfig(tb testing.TB) *config.Configs {
	configPathsHolder := createConfigurationsPathsHolder()

	generalConfig, _ := common.LoadMainConfig(configPathsHolder.MainConfig)
	ratingsConfig, _ := common.LoadRatingsConfig(configPathsHolder.Ratings)
	economicsConfig, _ := common.LoadEconomicsConfig(configPathsHolder.Economics)
	prefsConfig, _ := common.LoadPreferencesConfig(configPathsHolder.Preferences)
	mainP2PConfig, _ := common.LoadP2PConfig(configPathsHolder.MainP2p)
	fullArchiveP2PConfig, _ := common.LoadP2PConfig(configPathsHolder.FullArchiveP2p)
	externalConfig, _ := common.LoadExternalConfig(configPathsHolder.External)
	systemSCConfig, _ := common.LoadSystemSmartContractsConfig(configPathsHolder.SystemSC)
	epochConfig, _ := common.LoadEpochConfig(configPathsHolder.Epoch)
	roundConfig, _ := common.LoadRoundConfig(configPathsHolder.RoundActivation)
	var nodesConfig config.NodesConfig
	_ = core.LoadJsonFile(&nodesConfig, NodesSetupPath)

	mainP2PConfig.KadDhtPeerDiscovery.Enabled = false
	prefsConfig.Preferences.DestinationShardAsObserver = "0"
	prefsConfig.Preferences.ConnectionWatcherType = p2p.ConnectionWatcherTypePrint

	configs := &config.Configs{}
	configs.GeneralConfig = generalConfig
	configs.RatingsConfig = ratingsConfig
	configs.EconomicsConfig = economicsConfig
	configs.SystemSCConfig = systemSCConfig
	configs.PreferencesConfig = prefsConfig
	configs.MainP2pConfig = mainP2PConfig
	configs.FullArchiveP2pConfig = fullArchiveP2PConfig
	configs.ExternalConfig = externalConfig
	configs.EpochConfig = epochConfig
	configs.RoundConfig = roundConfig
	workingDir := tb.TempDir()
	dbDir := tb.TempDir()
	logsDir := tb.TempDir()
	configs.FlagsConfig = &config.ContextFlagsConfig{
		WorkingDir:  workingDir,
		DbDir:       dbDir,
		LogsDir:     logsDir,
		UseLogView:  true,
		BaseVersion: BaseVersion,
		Version:     Version,
	}
	configs.ConfigurationPathsHolder = configPathsHolder
	configs.ImportDbConfig = &config.ImportDbConfig{}
	configs.NodesConfig = &nodesConfig

	configs.GeneralConfig.GeneralSettings.ChainParametersByEpoch = computeChainParameters(uint32(len(configs.NodesConfig.InitialNodes)), configs.GeneralConfig.GeneralSettings.GenesisMaxNumberOfShards)

	return configs
}

func computeChainParameters(numInitialNodes uint32, numShardsWithoutMeta uint32) []config.ChainParametersByEpochConfig {
	numShardsWithMeta := numShardsWithoutMeta + 1
	nodesPerShards := numInitialNodes / numShardsWithMeta
	shardCnsGroupSize := nodesPerShards
	if shardCnsGroupSize > 1 {
		shardCnsGroupSize--
	}
	diff := numInitialNodes - nodesPerShards*numShardsWithMeta
	return []config.ChainParametersByEpochConfig{
		{
			ShardConsensusGroupSize:     shardCnsGroupSize,
			ShardMinNumNodes:            nodesPerShards,
			MetachainConsensusGroupSize: nodesPerShards,
			MetachainMinNumNodes:        nodesPerShards + diff,
			RoundDuration:               2000,
		},
	}
}

func createConfigurationsPathsHolder() *config.ConfigurationPathsHolder {
	var concatPath = func(filename string) string {
		return path.Join(BaseNodeConfigPath, filename)
	}

	return &config.ConfigurationPathsHolder{
		MainConfig:               concatPath(ConfigPath),
		Ratings:                  concatPath(RatingsPath),
		Economics:                concatPath(EconomicsPath),
		Preferences:              concatPath(PrefsPath),
		External:                 concatPath(ExternalPath),
		MainP2p:                  concatPath(MainP2pPath),
		FullArchiveP2p:           concatPath(FullArchiveP2pPath),
		Epoch:                    concatPath(EpochPath),
		SystemSC:                 concatPath(SystemSCConfigPath),
		GasScheduleDirectoryName: concatPath(GasSchedule),
		RoundActivation:          concatPath(RoundActivationPath),
		Nodes:                    NodesSetupPath,
		Genesis:                  GenesisPath,
		SmartContracts:           GenesisSmartContracts,
		ValidatorKey:             ValidatorKeyPemPath,
		ApiRoutes:                "",
		P2pKey:                   P2pKeyPath,
	}
}
